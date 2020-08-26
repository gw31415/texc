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

func main() {
	lis, _ := net.Listen("unix", "/home/ama/texcd.sock")
	defer lis.Close()
	opts := []grpc.ServerOption{}
	wd, _ := os.Getwd()
	grpcServer := grpc.NewServer(opts...)
	wd += "/temp"
	if f, err := os.Stat(wd); os.IsNotExist(err) {
		os.Mkdir(wd, 0755)
	} else if !f.IsDir() {
		panic(err)
	}
	sv, err := server.NewTexcServiceServer(wd)
	if err != nil {
		panic(err)
	}
	proto.RegisterTexcServiceServer(grpcServer, sv)
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()
	// SIGINT待ち
	fmt.Println("Press Ctrl-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
