package router

import (
	"google.golang.org/grpc"
	"grpc-test-demo/service"
	"grpc-test-demo/src/chat"
	"grpc-test-demo/src/prod"
)


var streamer = service.Streamer{}

func Init(s *grpc.Server) {

	prod.RegisterProductServiceServer(s, service.TestService{})

	chat.RegisterChatServer(s, &streamer)
}
