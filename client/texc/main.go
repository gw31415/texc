package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gw31415/texc/proto"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
)

const (
	block_size = 0xffff
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			switch e := e.(type) {
			case error:
				fmt.Println(e.Error())
			case string:
				fmt.Println(e)
			default:
				fmt.Println(e)
			}
			os.Exit(1)
		}
	}()
	if len(os.Args) != 2 {
		panic("please specify the main tex file.")
	}
	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		panic("file not found.")
	}
	conn, err := grpc.Dial("texc.amas.dev:3475", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := proto.NewTexcServiceClient(conn)
	stream, err := client.Sync(context.Background())
	if err != nil {
		panic(err)
	}
	tar_data := bytes.NewBuffer([]byte{})
	tar_w := tar.NewWriter(tar_data)
	paths, err := dirwalk(".")
	if err != nil {
		panic(err)
	}
	var total_size int64 = 0
	for _, path := range paths {
		f, err := os.OpenFile(path, os.O_RDONLY, 0755)
		if err != nil {
			panic(err)
		}
		stat, err := f.Stat()
		if err != nil {
			panic(err)
		}
		tar_w.WriteHeader(&tar.Header{
			Name:    path,
			Mode:    int64(stat.Mode()),
			ModTime: stat.ModTime(),
			Size:    stat.Size(),
		})
		if err != nil {
			panic(err)
		}
		io.Copy(tar_w, f)
		total_size += stat.Size()
		fmt.Printf("Add: %s\n", path)
	}
	tar_w.Close()
	in_pb := new(proto.Input)
	in_pb.Data = make([]byte, block_size)
	var sent_size int = 0
	send_status := "Send: %d%%\n"
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		send_status = "\rSend: %d%%"
	}
	for {
		i, err := tar_data.Read(in_pb.Data)
		sent_size += i
		if err == io.EOF {
			break
		}
		stream.Send(in_pb)
		percent := int64(sent_size) * 100 / total_size
		if percent > 99 {
			percent = 99
		}
		fmt.Printf(send_status, percent)
	}
	fmt.Printf(send_status, 100)
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println()
	}
	stream.Send(&proto.Input{
		Exec: []string{"latexmk", os.Args[1]},
	})
	stream.Send(&proto.Input{
		Dl: getFileNameWithoutExt(os.Args[1]) + ".pdf",
	})
	stream.CloseSend()
	dl := false
	b := new(bytes.Buffer)
	for {
		out, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if out.Stdout != nil {
			os.Stdout.Write(out.Stdout)
		}
		if out.Stderr != nil {
			os.Stderr.Write(out.Stderr)
		}
		if out.Data != nil {
			dl = true
			b.Write(out.Data)
		}
	}
	if dl {
		tr := tar.NewReader(b)
		for {
			h, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			wd, _ := os.Getwd()
			path := fmt.Sprintf("%s/%s", wd, h.Name)
			dir := filepath.Dir(path)
			if !h.FileInfo().IsDir() {
				os.MkdirAll(dir, 0755)
				file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0744)
				if err != nil {
					panic(err)
				}
				io.Copy(file, tr)
				fmt.Printf("Download: %s\n", h.Name)
			}
		}
	}
}

func dirwalk(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, file := range files {
		if file.IsDir() {
			f, err := dirwalk(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			paths = append(paths, f...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths, nil
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
