package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-test-demo/src/prod"
)

const (
	// Address gRPC服务地址
	Address = "0.0.0.0:9999"
)

func main() {
	// 连接
	conn, err := grpc.Dial(Address, grpc.WithInsecure())

	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	// 初始化客户端
	c := prod.NewProductServiceClient(conn)

	// 调用方法
	request := new(prod.ProdRequest)
	request.ProdId = 211

	r, err := c.GetProductStock(context.Background(), request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r)
}
