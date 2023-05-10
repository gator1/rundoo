package main

import (
	rundoogrpc "app/api/v1"
	"app/log"
	"app/products"
	"app/registry"
	"app/service"
	"context"
	"fmt"
	stlog "log"
	"net"
	"strconv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.ServiceConfig

	// configure our service
	productService := products.NewService()


	handler := &products.ProductsHandler{}
	r.Name = registry.ProductService
	r.Host = host
	r.Port = port
	r.URL = serviceAddress
	r.HeartbeatURL = r.URL + "/heartbeat"
	r.RequiredServices = []registry.ServiceName{
		registry.LogService,
	}
	r.UpdateURL = r.URL + "/services"
	r.HttpHandler = handler
	ctx, err := service.Start(context.Background(), r)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}

	
	// configure our gRPC service controller
	productServiceController := NewProductsServiceController(productService)

	// start a gRPC server
	server := grpc.NewServer()
	rundoogrpc.RegisterProductServiceServer(server, productServiceController)
	reflection.Register(server)

	portInt, _ := strconv.Atoi(port)
	rpcPort := ":"+strconv.Itoa(portInt + 1)

	con, err := net.Listen("tcp", rpcPort)
	if err != nil {
		stlog.Printf("Starting gRPC user service on %s...\n", con.Addr().String())
		panic(err)
	}

	stlog.Printf("Starting gRPC user service on %s...\n", con.Addr().String())
	err = server.Serve(con)
	if err != nil {
		stlog.Printf("Starting gRPC user service on %s...\n", con.Addr().String())
		panic(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down product service")
}
