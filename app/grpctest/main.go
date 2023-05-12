package main

import (
	"context"
    "log"

    "google.golang.org/grpc"
    rundoogrpc "app/api/v1"
    rundoogrpcclient "app/api/v1/client"
)

func main() {
	var err error
	
	rpcPort := ":6001"

	grpcService, err := rundoogrpcclient.NewGRPCService(rpcPort)
	if err != nil {
		log.Printf("error instantiating gRPC service: %v\n", err)

	}
	
    response, err := grpcService.GetProducts()
	if err != nil {
        log.Fatalf("failed to get products: %v", err)
    }


	if err != nil {
		log.Printf("grpcService.GetProducts() returned an error: %v\n", err)
	} else {
		log.Printf("grpcService.GetProducts() returned:  %+v\n", response)
		
	}

    for _, product := range response {
        log.Printf("Product: %v\n", product)
    }
	
}


func test() {
    conn, err := grpc.Dial("localhost:6001", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("failed to dial: %v", err)
    }
    defer conn.Close()

    client := rundoogrpc.NewProductServiceClient(conn)

    response, err := client.GetProducts(context.Background(), &rundoogrpc.GetProductsRequest{})
    if err != nil {
        log.Fatalf("failed to get products: %v", err)
    }

    for _, product := range response.GetProducts() {
        log.Printf("Product: %v\n", product)
    }
}

