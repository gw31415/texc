package main

import (
	"net"
	"os"

	"github.com/gw31415/texc/proto"
	"github.com/gw31415/texc/server"
	"google.golang.org/grpc"
)
func main() {
	lis, _ := net.FileListener(os.Stdin)
	wd, _ := os.Getwd()
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterTexcServiceServer(grpcServer, server.NewTexcServiceServer(wd))
	grpcServer.Serve(lis)
}
