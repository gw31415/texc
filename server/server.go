package server

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

func (sv *TexcServiceServer) Sync(stream pb.TexcService_SyncServer) error {
	b := new(bytes.Buffer)
	dl_err := make(chan error)
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				dl_err <- err
				return
			}
			b.Write(in.Data)
		}
		dl_err <- nil
	}()
	tr := tar.NewReader(b)
	tar_err := make(chan error)
	new_dir := "temp"
	cache_dir := fmt.Sprintf("%s/%s/", sv.cache_dir, new_dir)
	go func() {
		for {
			h, err := tr.Next()
			if err == io.EOF {
				break
			}
			path := cache_dir + h.Name
			dir := filepath.Dir(path)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.MkdirAll(dir, 0755)
			}
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				tar_err <- err
				return
			}
			io.Copy(file, tr)
		}
		tar_err <- nil
	}()
	if err := <-dl_err; err != nil {
		return err
	}
	if err := <-tar_err; err != nil {
		return err
	}
	return nil
}
