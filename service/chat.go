package service

import (
	"fmt"
	"io"
	"strconv"

	"grpc-test-demo/src/chat"
)

type Streamer struct{}

func (s *Streamer) BindStream(stream chat.Chat_BindStreamServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("收到客户端通过 context 发出的停止信号")
			return ctx.Err()
		default:
			input, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("客户端发送的数据流结束")
				return nil
			}
			if err != nil {
				fmt.Println("接收数据出错:", err)
				return err
			}
			switch input.In {
			case "结束对话\n":
				fmt.Println("收到'结束对话'指令")
				if err := stream.Send(&chat.Response{Out: "收到结束指令"}); err != nil {
					return err
				}
				return nil

			case "返回数据流\n":
				fmt.Println("收到'返回数据流'指令")
				for i := 0; i < 10; i++ {
					if err := stream.Send(&chat.Response{Out: "数据流 #" + strconv.Itoa(i)}); err != nil {
						return err
					}
				}
			default:
				fmt.Printf("[收到消息]: %s", input.In)
				if err := stream.Send(&chat.Response{Out: "服务端返回: " + input.In}); err != nil {
					return err
				}
			}
		}
	}
}
