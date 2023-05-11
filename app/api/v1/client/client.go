package client

import (
	"context"
	"time"

	"app/products"
	rundoogrpc "app/api/v1"
	"google.golang.org/grpc"
)

var defaultRequestTimeout = time.Second * 10

type grpcService struct {
	grpcClient rundoogrpc.ProductServiceClient
}

// NewGRPCService creates a new gRPC user service connection using the specified connection string.
func NewGRPCService(connString string) (products.ServiceInterface, error) {
	conn, err := grpc.Dial(connString, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &grpcService{grpcClient: rundoogrpc.NewProductServiceClient(conn)}, nil
}

func (s *grpcService) GetProducts() (result products.Products, err error) {
	req := &rundoogrpc.GetProductsRequest{}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancelFunc()
	// resp, err := s.grpcClient.GetProducts(ctx, req)
	s.grpcClient.GetProducts(ctx, req)
	if err != nil {
		return
	}
/*
	for _, grpcProduct:= range resp.GetProducts() {
		//u := unmarshalProduct(*grpcProduct)
		//result[u.ID] = u
	}
	*/
	return
}

func unmarshalProduct(grpcProduct *products.Product) (result products.Product) {
	//result.ID = grpcUser.Id
	result.Name = grpcProduct.Name
	return
}
