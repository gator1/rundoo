package main

import (
	"app/log"
	"app/registry"
	"app/rundooportal"
	"app/service"
	"context"
	"fmt"
	stlog "log"
	"net/http"
)

func main() {
	err := rundooportal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}

	host, port := "localhost", "5050"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.ServiceConfig
	handler := &rundooportal.RundooHandler{}

	r.Name = registry.RundooPortal
	r.Host = host
	r.Port = port
	r.URL = serviceAddress
	r.HeartbeatURL = r.URL + "/heartbeat"
	r.RequiredServices = []registry.ServiceName{
		registry.LogService,
		registry.ProductService,
	}
	r.UpdateURL = r.URL + "/services"
	r.HttpHandler = handler
	r.Mux = http.NewServeMux()
	r.Mux.Handle("/products", handler)
	r.Mux.Handle("/products/", handler)
	ctx, err := service.Start(context.Background(), r)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}

	<-ctx.Done()
	fmt.Println("Shutting down rundoo portal")

}
