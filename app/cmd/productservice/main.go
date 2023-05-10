package main

import (
	"app/products"
	"app/log"
	"app/registry"
	"app/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.ServiceConfig
	r.Name = registry.ProductService
	r.URL = serviceAddress
	r.HeartbeatURL = r.URL + "/heartbeat"
	r.RequiredServices = make([]registry.ServiceName, 0)
	r.UpdateURL = r.URL + "/services"

	ctx, err := service.Start(context.Background(),
		host,
		port,
		r,
		products.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}

	<-ctx.Done()
	fmt.Println("Shutting down product service")
}
