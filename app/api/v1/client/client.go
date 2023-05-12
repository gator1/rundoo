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
	GrpcClient rundoogrpc.ProductServiceClient
}

// NewGRPCService creates a new gRPC user service connection using the specified connection string.
func NewGRPCService(connString string) (products.ServiceInterface, error) {
	conn, err := grpc.Dial(connString, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &grpcService{GrpcClient: rundoogrpc.NewProductServiceClient(conn)}, nil
}

func (s *grpcService) GetProducts() (result products.Products, err error) {
	req := &rundoogrpc.GetProductsRequest{}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancelFunc()
	// resp, err := s.GetProducts(req)
	resp, err := s.GrpcClient.GetProducts(ctx, req)
	if err != nil {
		return
	}


	for _, grpcProduct:= range resp.GetProducts() {
		result = append(result, products.Product{
			Name:     grpcProduct.Name,
			Category: products.CategoryType(grpcProduct.Category),
			Sku:      products.SKU(grpcProduct.Sku),
		})
	}
	
	return
	
}

/*

func (s *grpcService) SearchProducts(filters []rundoogrpc.Filter) (result products.Products, err error) {
	req := &rundoogrpc.SearchProductsRequest{
		Filters: filters,
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancelFunc()

	resp, err := s.GrpcClient.SearchProducts(ctx, req)
	if err != nil {
		return
	}

	for _, grpcProduct := range resp.GetProducts() {
		u := unmarshalProduct(grpcProduct)
		result = append(result, u) 
	}

	return
}*/

func (s *grpcService) SearchProducts(filters []rundoogrpc.Filter) (result products.Products, err error) {
	// Convert the []products.Filter slice to []*rundoogrpc.Filter slice
	grpcFilters := make([]*rundoogrpc.Filter, len(filters))
	for i, f := range filters {
		grpcFilters[i] = &rundoogrpc.Filter{
			Field: f.Field,
			Value: f.Value,
		}
	}

	req := &rundoogrpc.SearchProductsRequest{
		Filters: grpcFilters,
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancelFunc()

	resp, err := s.GrpcClient.SearchProducts(ctx, req)
	if err != nil {
		return
	}

	for _, grpcProduct := range resp.GetProducts() {
		u := unmarshalProduct(grpcProduct)
		result = append(result, *u) 
	}

	return
}


func (s *grpcService) AddProduct(product products.Product) (ok bool, err error) {
	req := &rundoogrpc.AddProductRequest{
		Product: &rundoogrpc.Product{
			Name:  product.Name,
		},
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancelFunc()

	resp, err := s.GrpcClient.AddProduct(ctx, req)
	if err != nil {
		return
	}
	return resp.Ok, nil
}


func unmarshalProduct(grpcProduct *rundoogrpc.Product) *products.Product {
	p := &products.Product{
		Name:  grpcProduct.Name,
		Category:  products.CategoryType(grpcProduct.Category),
		Sku: products.SKU(grpcProduct.Sku),
	}
	return p
}