package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gopkg.in/olahol/melody.v1"
	"grpc-test-demo/src/chat"
)

func main() {

	// 设置 log
	log.SetFlags(log.LstdFlags)

	// 创建 gRPC 连接
	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		log.Printf("连接失败: [%v]\n", err)
		return
	}
	defer conn.Close()

	// 声明客户端
	client := chat.NewChatClient(conn)

	r := gin.Default()

	// 封装了 gorrila websocket 的框架 melody
	m := melody.New()

	// 静态页面路由
	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "websocket-demo/html/index.html")
	})

	// websocket 路由
	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	// 处理 websocket 客户端新连接, 并为每一个新连接创建一个双向数据流
	m.HandleConnect(func(s *melody.Session) {
		log.Println("有新用户接入")
		// 给每个连入的新用户创建一个数据流
		// 声明 context
		ctx := context.Background()

		// 创建双向数据流
		stream, err := client.BindStream(ctx)
		if err != nil {
			log.Printf("创建数据流失败: [%v]\n", err)
			// 如果创建数据流失败，向客户端发送失败信息 同时 关闭websocket连接
			s.CloseWithMsg([]byte("创建数据流失败:" + err.Error()))
			return
		}

		// 如果创建成功，将数据流保存在 session 中
		s.Set("stream", stream)

		// 启动一个 goroutine 用于接收从服务端返回的消息
		go func() {
			for {
				// 接收从服务端返回的数据流
				response, err := stream.Recv()

				if err == io.EOF {
					log.Println("⚠️ 收到服务端的结束信号")
					s.CloseWithMsg([]byte("⚠️ 收到服务端的结束信号"))
					return
				}

				if err != nil {
					// TODO: 处理接收错误
					log.Println("接收数据出错:", err)
					s.CloseWithMsg([]byte("接收数据出错" + err.Error()))
					return
				}

				log.Printf("[客户端收到]: %s", response.Out)
				// 如果成功收到从服务端返回的消息, 将消息序列化后返回给 websocket 客户端
				data, _ := json.Marshal(response)
				s.Write(data)
			}
		}()
	})

	// 处理用户发来的消息
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		log.Println("收到消息:", msg)
		// 把用户输入的信息原样返回 websocket 客户端
		s.Write(msg)

		// 将 []byte 类型的 msg 解析为 proto.Request
		var in chat.Request
		if err := json.Unmarshal(msg, &in); err != nil {
			log.Println("解析输入信息失败:", err)
			s.CloseWithMsg([]byte("输入信息解析失败"))
			return
		}

		// 从 session 中取出 stream
		saveData, ok := s.Get("stream")
		if !ok {
			s.CloseWithMsg([]byte("没有找到stream!"))
			return
		}

		// 断言 stream
		stream, ok := saveData.(chat.Chat_BindStreamClient)
		if !ok {
			s.CloseWithMsg([]byte("被保存的数据流不是Chat_BindStreamClient!"))
			return
		}

		if err := stream.Send(&in); err != nil {
			s.CloseWithMsg([]byte("向gRPC服务端发送消息失败:" + err.Error()))
			return
		}
	})

	// 处理 websocket 连接断开事件，并关闭 session 中 stream 的连接
	m.HandleDisconnect(func(s *melody.Session) {
		log.Println("websocket客户端断开连接")
		// 从 session中取出 stream
		saveData, ok := s.Get("stream")
		if !ok {
			log.Println("没有找到stream!")
			return
		}

		// 断言stream
		stream, ok := saveData.(chat.Chat_BindStreamClient)
		if !ok {
			log.Println("被保存的数据流不是Chat_BindStreamClient!")
			return
		}

		if err := stream.CloseSend(); err != nil {
			log.Println("断开stream连接出错:", err)
		}
	})

	r.Run(":8989")
}
