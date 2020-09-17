package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-test-demo/service"
	"grpc-test-demo/src/router"
)

const (
	// Address gRPC服务地址
	Address = "0.0.0.0:9999"
)

var testService = service.TestService{}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}

	// 实例化 grpc Server
	server := grpc.NewServer()

	router.Init(server)

	reflection.Register(server)

	fmt.Println("Listen on " + Address)

	server.Serve(listen)
}
