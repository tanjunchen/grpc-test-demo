package main

import (
	"fmt"
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	k8s "grpc-test-demo/k8s-watch-list-grpc/proto"
)

const (
	// Address gRPC服务地址
	Address = "localhost:9999"
)

// customCredential 自定义认证
type customCredential struct{}

func (customCredential customCredential) RequireTransportSecurity() bool {
	return false
}

// GetRequestMetadata 实现自定义认证接口
func (customCredential customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "admin",
		"appkey": "admin",
	}, nil
}

func main() {
	var err error
	var opts []grpc.DialOption

	// insecure way
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))

	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	client := k8s.NewServiceServiceClient(conn)

	stream, err := client.SyncServiceWatchListService(context.TODO())

	if err != nil {
		fmt.Println(err)
	}
	for {
		// 接收从 服务端返回的数据流
		response, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("收到服务端的结束信号")
			break // 如果收到结束信号, 则退出接收循环, 结束客户端程序
		}

		if err != nil {
			fmt.Println("接收数据出错:", err)
			break
		}

		// 没有错误的情况下，打印来自服务端的消息
		fmt.Printf("[客户端收到]: %s \n", response)
	}
}
