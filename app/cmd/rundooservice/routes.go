package main

import (
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	rundoogrpc "app/api/v1"
	"app/registry"
	"app/rundoo"
)

func (app *application) routes(r *registry.ServiceConfig, serviceAddress string, service rundoo.ServiceInterface) *http.ServeMux {
	
	r.Mux = http.NewServeMux()
	
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	
	r.Name = registry.RundooService
	r.URL = serviceAddress
	r.HeartbeatURL = r.URL + "/heartbeat"
	r.RequiredServices = []registry.ServiceName{
		registry.LogService,
	}
	r.UpdateURL = r.URL + "/services"
	
	r.Mux.HandleFunc("/", http.NotFound) // Catch-all route
	r.Mux.HandleFunc("/products", app.getCreateProductsHandler)
	r.Mux.HandleFunc("/products/", app.getUpdateDeleteProductsHandler)
	
	// configure our gRPC service controller
	productServiceController := NewProductsServiceController(service)

	// start a gRPC server
	server := grpc.NewServer()
	rundoogrpc.RegisterProductServiceServer(server, productServiceController)
	reflection.Register(server)
	r.GrpcServer = server

	return r.Mux
}