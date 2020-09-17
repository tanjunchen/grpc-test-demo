package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/grpc"
	"grpc-test-demo/src/chat"
	"grpc-test-demo/src/prod"
)

const (
	// Address gRPC服务地址
	Address = "localhost:9999"
)

func main() {
	// 创建连接
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("连接失败: [%v]\n", err)
		return
	}
	defer conn.Close()

	// 声明客户端
	client := chat.NewChatClient(conn)

	// 声明 context
	ctx := context.Background()

	// 创建双向数据流
	stream, err := client.BindStream(ctx)
	if err != nil {
		fmt.Printf("创建数据流失败: [%v]\n", err)
	}

	// 启动一个 goroutine 接收命令行输入的指令
	go func() {
		fmt.Println("请输入消息...")
		in := bufio.NewReader(os.Stdin)
		for {
			// 获取命令行输入的字符串， 以回车 \n 作为结束标志
			str, _ := in.ReadString('\n')

			// 向服务端发送 指令
			if err := stream.Send(&chat.Request{In: str}); err != nil {
				return
			}
		}
	}()

	for {
		// 接收从 服务端返回的数据流
		response, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("⚠️收到服务端的结束信号")
			break // 如果收到结束信号, 则退出接收循环, 结束客户端程序
		}

		if err != nil {
			fmt.Println("接收数据出错:", err)
		}

		// 没有错误的情况下，打印来自服务端的消息
		fmt.Printf("[客户端收到]: %s \n", response.Out)
	}
}

func test() {
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
