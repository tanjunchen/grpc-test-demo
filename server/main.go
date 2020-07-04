package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc-test-demo/service"
	 "grpc-test-demo/src/prod"
	"net"
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
	s := grpc.NewServer()

	// 注册HelloService
	prod.RegisterProductServiceServer(s, testService)

	fmt.Println("Listen on " + Address)

	s.Serve(listen)
}
