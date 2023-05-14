package main

import (
	
	"context"
	"fmt"
	stlog "log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	rundoogrpc "app/api/v1"
	"app/log"
	"app/rundoo"
	"app/registry"
	"app/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.ServiceConfig

	// configure our service
	productService := rundoo.NewService()


	handler := &rundoo.ProductsHandler{}
	r.Name = registry.RundooService
	r.Host = host
	r.Port = port
	r.URL = serviceAddress
	r.HeartbeatURL = r.URL + "/heartbeat"
	r.RequiredServices = []registry.ServiceName{
		registry.LogService,
	}
	r.UpdateURL = r.URL + "/services"
	r.HttpHandler = handler

	r.HttpHandler = handler
	r.Mux = http.NewServeMux()
	r.Mux.Handle("/products", handler)
	r.Mux.Handle("/products/", handler)

	
	// configure our gRPC service controller
	productServiceController := NewProductsServiceController(productService)

	// start a gRPC server
	server := grpc.NewServer()
	rundoogrpc.RegisterProductServiceServer(server, productServiceController)
	reflection.Register(server)
	r.GrpcServer = server

	ctx, err := service.Start(context.Background(), r)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}

	<-ctx.Done()
	fmt.Println("Shutting down product service")
}
