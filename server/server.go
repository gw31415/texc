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

const (
	block_size = 0xffff
)

type TexcServiceServer struct {
	pb.TexcServiceServer
	cache_dir string
}

func NewTexcServiceServer(cache_dir string) (*TexcServiceServer, error) {
	if f, err := os.Stat(cache_dir); os.IsNotExist(err) || !f.IsDir() {
		return nil, err
	}
	abs_path, err := filepath.Abs(cache_dir)
	if err != nil {
		return nil, err
	}
	return &TexcServiceServer{
		cache_dir: abs_path,
	}, nil
}

var randSrc = rand.NewSource(time.Now().UnixNano())

func (sv *TexcServiceServer) Sync(stream pb.TexcService_SyncServer) error {
	done := make(chan error, 1)
	go func() {
		execs := make([][]string, 0)
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
				execs = append(execs, in.Exec)
			}
			if in.Dl != "" {
				dl_filelist = append(dl_filelist, in.Dl)
			}
		}
		tr := tar.NewReader(b)
		new_dir := strconv.FormatInt(randSrc.Int63(), 16)
		fmt.Printf("Login: %s\n", new_dir)
		cache_dir := fmt.Sprintf("%s/%s/", sv.cache_dir, new_dir)
		defer func() {
			os.RemoveAll(cache_dir)
			fmt.Printf("Logout: %s\n", new_dir)
		}()
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
		for _, exe := range execs {
			cmd := exec.Command(exe[0], exe[1:]...)
			cmd.Dir = cache_dir
			stdout, err := cmd.StdoutPipe()
			stderr := bytes.NewBuffer([]byte{})
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
		out_pb := new(pb.Output)
		out_pb.Data = make([]byte, block_size)
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
	case <-time.After(time.Second * 30):
		return errors.New("timeout.")
	case err := <-done:
		return err
	}
}
