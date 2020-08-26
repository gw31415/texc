package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gw31415/texc/proto"
	"google.golang.org/grpc"
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
			return
		}
	}()
	if len(os.Args) != 2 {
		panic("Args length does not match.")
	}
	dialer := func(a string, t time.Duration) (net.Conn, error) {
		return net.Dial("unix", a)
	}
	conn, err := grpc.Dial(os.Args[1], grpc.WithInsecure(), grpc.WithDialer(dialer))
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
	for _, path := range paths {
		fmt.Println(path)
		stat, _ := os.Stat(path)
		f, err := os.OpenFile(path, os.O_RDONLY, 0755)
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
	}
	in_pb := new(proto.Input)
	in_pb.Data = make([]byte, 0xffff)
	for {
		_, err := tar_data.Read(in_pb.Data)
		if err == io.EOF {
			break
		}
		stream.Send(in_pb)
	}
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		fields := strings.Fields(scan.Text())
		stream.Send(&proto.Input{
			Exec: fields,
		})
	}
	stream.CloseSend()
	out_ch := make(chan int)
	go func() {
		err_count := 0
		for {
			out, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				err_count++
				fmt.Println(err.Error())
			}
			if out.Stdout != nil {
				os.Stdout.Write(out.Stdout)
			}
		}
		out_ch <- err_count
	}()
	<-out_ch
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
