package main

import (
	"context"
    "log"

    "google.golang.org/grpc"
    rundoogrpc "app/api/v1"
	
)


func main() {
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
