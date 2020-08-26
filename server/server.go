package server

import (
	"archive/tar"
	"bytes"
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
	execs := make([][]string, 0)
	b := new(bytes.Buffer)
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if in.Data != nil {
			b.Write(in.Data)
		}
		if in.Exec != nil {
			execs = append(execs, in.Exec)
		}
	}
	tr := tar.NewReader(b)
	new_dir := strconv.FormatInt(randSrc.Int63(), 16)
	fmt.Printf("Login: %s\n", new_dir)
	cache_dir := fmt.Sprintf("%s/%s/", sv.cache_dir, new_dir)
	defer os.RemoveAll(cache_dir)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		path := cache_dir + h.Name
		dir := filepath.Dir(path)
		if !h.FileInfo().IsDir() {
			os.MkdirAll(dir, 0755)
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			io.Copy(file, tr)
			fmt.Printf(" -- %s\n", h.Name)
		}
	}
	for _, exe := range execs {
		cmd := exec.Command(exe[0], exe[1:]...)
		cmd.Dir = cache_dir
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		cmd.Start()
		buf := make([]byte, 0xffff)
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
				return err
			}
		}
		if err := cmd.Wait(); err != nil {
			return err
		}
	}
	fmt.Printf("Logout: %s\n", new_dir)
	return nil
}
