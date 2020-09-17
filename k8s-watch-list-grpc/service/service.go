package service

import k8s "grpc-test-demo/k8s-watch-list-grpc/proto"

type K8sServiceService struct{}

var MsgChan = make(chan k8s.SyncServiceResponse, 1024)

func (k8sServiceService K8sServiceService) SyncServiceWatchListService(send k8s.ServiceService_SyncServiceWatchListServiceServer) error {
	for {
		m := <-MsgChan
		// 消息处理函数
		send.Send(&m)
	}
	return nil
}
