/*
Texcのサーバーインスタンスです.
*/
package server

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	pb "github.com/gw31415/texc/proto"
)

// Texcのサーバーインスタンス
type TexcServiceServer struct {
	pb.TexcServiceServer
	cache_dir      string
	cmd_white_list []string
	Timeout        time.Duration // タイムアウト時間
	BlockSize      int64         // 送信時のブロックサイズ
}

type TexcServiceServerConfig struct {
	CacheDir     string   // 一時ファイル用ディレクトリ
	CmdWhiteList []string // 許可するコマンド一覧
}

// Texcのサーバーインスタンスを作成します
func NewTexcServiceServer(config *TexcServiceServerConfig) (*TexcServiceServer, error) {
	if f, err := os.Stat(config.CacheDir); os.IsNotExist(err) || !f.IsDir() {
		return nil, err
	}
	abs_path, err := filepath.Abs(config.CacheDir)
	if err != nil {
		return nil, err
	}
	return &TexcServiceServer{
		cmd_white_list: config.CmdWhiteList,
		cache_dir:      abs_path,
		Timeout:        time.Second * 30,
		BlockSize:      0xffff,
	}, nil
}

var randSrc = rand.NewSource(time.Now().UnixNano())

type exec_unit struct {
	exec          []string
	no_out_stream bool
}

// 呼びだされるSync関数
func (sv *TexcServiceServer) Sync(stream pb.TexcService_SyncServer) error {
	done := make(chan error, 1)
	go func() {
		// データの受けとり
		var execs []*exec_unit
		b := new(bytes.Buffer)
		dl_filelist := make([]string, 0)
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				done <- err
				return
			}
			if in.Data != nil {
				b.Write(in.Data)
			}
			if in.Exec != nil {
				execs = append(execs, &exec_unit{
					exec:          in.Exec,
					no_out_stream: in.NoOutStream,
				})
			}
			if in.Dl != "" {
				dl_filelist = append(dl_filelist, in.Dl)
			}
		}
		// 一時ファイル用のディレクトリのセットアップ
		new_dir := strconv.FormatInt(randSrc.Int63(), 16)
		fmt.Printf("Login: %s\n", new_dir)
		cache_dir := fmt.Sprintf("%s/%s/", sv.cache_dir, new_dir)
		// 後処理の設定
		defer func() {
			os.RemoveAll(cache_dir)
			fmt.Printf("Logout: %s\n", new_dir)
		}()

		// Tarの展開
		tr := tar.NewReader(b)
		for {
			h, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err.Error())
				done <- err
				return
			}
			path := cache_dir + h.Name
			dir := filepath.Dir(path)
			if !h.FileInfo().IsDir() {
				os.MkdirAll(dir, 0755)
				file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					fmt.Println(err.Error())
					done <- err
					return
				}
				io.Copy(file, tr)
				fmt.Printf(" -- %s\n", h.Name)
			}
		}

		// コマンドの処理
		for _, exe := range execs {
			// ホワイトリスト処理
			is_include := false
			for _, e := range sv.cmd_white_list {
				if exe.exec[0] == e {
					is_include = true
					break
				}
			}
			if !is_include {
				done <- errors.New("this command is not allowed.")
				return
			}
			// cmdインスタンス
			cmd := exec.Command(exe.exec[0], exe.exec[1:]...)
			stderr := bytes.NewBuffer([]byte{})
			cmd.Dir = cache_dir
			// リアルタイム出力なし
			if exe.no_out_stream {
				cmd.Run()
				if cmd.ProcessState.ExitCode() != 0 {
					out, err := cmd.CombinedOutput()
					if err != nil {
						done <- err
						return
					}
					stream.Send(
						&pb.Output{
							Stdout: out,
							Stderr: stderr.Bytes(),
						},
					)
					done <- errors.New(fmt.Sprintf("process exited with error code %d", cmd.ProcessState.ExitCode()))
					return
				}
				out, err := cmd.CombinedOutput()
				if err != nil {
					done <- err
					return
				}
				stream.Send(
					&pb.Output{
						Stdout: out,
					},
				)
				done <- nil
				return
			} else {
				// リアルタイム出力
				stdout, err := cmd.StdoutPipe()
				cmd.Stderr = stderr
				if err != nil {
					done <- err
					return
				}
				if err := cmd.Start(); err != nil {
					done <- err
					return
				}
				buf := make([]byte, 0xff)
				for {
					i, err := stdout.Read(buf)
					if i > 0 {
						stream.Send(
							&pb.Output{
								Stdout: buf,
							},
						)
					}
					if err == io.EOF {
						break
					}
					if err != nil {
						done <- err
						return
					}
				}
				if cmd.Wait(); cmd.ProcessState.ExitCode() != 0 {
					stream.Send(
						&pb.Output{
							Stderr: stderr.Bytes(),
						},
					)
					done <- errors.New(fmt.Sprintf("process exited with error code %d", cmd.ProcessState.ExitCode()))
					return
				}
			}
		}

		// Tarの作成
		tar_data := bytes.NewBuffer([]byte{})
		tar_w := tar.NewWriter(tar_data)
		for _, file_name := range dl_filelist {
			path := cache_dir + file_name
			if !filepath.IsAbs(path) {
				done <- errors.New("path is not abs.")
			}
			f, err := os.OpenFile(path, os.O_RDONLY, 0644)
			if err != nil {
				done <- err
				return
			}
			stat, err := f.Stat()
			if err != nil {
				done <- err
				return
			}
			tar_w.WriteHeader(&tar.Header{
				Name:    file_name,
				Mode:    int64(stat.Mode()),
				ModTime: stat.ModTime(),
				Size:    stat.Size(),
			})
			if err != nil {
				done <- err
				return
			}
			io.Copy(tar_w, f)
			fmt.Printf("Send: %s\n", file_name)
		}
		tar_w.Close()

		// データの送信
		out_pb := new(pb.Output)
		out_pb.Data = make([]byte, sv.BlockSize)
		var sent_size int = 0
		for {
			i, err := tar_data.Read(out_pb.Data)
			sent_size += i
			if err == io.EOF {
				break
			}
			stream.Send(out_pb)
		}
		done <- nil
	}()
	select {
	case <-time.After(sv.Timeout):
		return errors.New("timeout.")
	case err := <-done:
		return err
	}
}
