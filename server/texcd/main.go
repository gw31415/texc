package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gw31415/texc/proto"
	"github.com/gw31415/texc/server"
	"google.golang.org/grpc"
)

const (
	port = ":3475"
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
		}
		os.Exit(1)
	}()
	lis, _ := net.Listen("tcp", port)
	defer lis.Close()
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	if len(os.Args) != 2 {
		panic("please specify the temp directory.")
	}
	wd := os.Args[1]
	if f, err := os.Stat(wd); os.IsNotExist(err) || !f.IsDir() {
		panic("directory not found.")
	}
	sv, err := server.NewTexcServiceServer(&server.TexcServiceServerConfig{
		CacheDir: wd,
		CmdWhiteList: []string{
			"latexmk",
			"pdftoppm",
			"pdflatex",
			"pdfcrop",
		},
	})
	if err != nil {
		panic(err)
	}
	proto.RegisterTexcServiceServer(grpcServer, sv)
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()
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
		}
	}()
	// SIGINT待ち
	fmt.Println("Press Ctrl-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
