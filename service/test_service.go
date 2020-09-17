package service

import (
	"context"
	"fmt"

	"grpc-test-demo/src/prod"
	"grpc-test-demo/src/status"
)

type TestService struct{}

func (testService TestService) GetProductStock(_ context.Context, in *prod.ProdRequest) (*prod.ProdResponse, error) {
	podId := in.ProdId
	fmt.Println(podId)
	response := prod.ProdResponse{
		ProdStock: in.ProdId,
		Status: &status.Status{
			Code: "200",
			Msg:  "success",
		},
	}
	return &response, nil
}
