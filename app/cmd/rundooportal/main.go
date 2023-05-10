package main

import (
	"app/log"
	"app/registry"
	"app/service"
	"app/rundooportal"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	err := rundooportal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}

	host, port := "localhost", "5050"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.ServiceConfig
	r.Name = registry.RundooPortal
	r.URL = serviceAddress
	r.HeartbeatURL = r.URL + "/heartbeat"
	r.RequiredServices = []registry.ServiceName{
		registry.LogService,
		registry.ProductService,
	}
	r.UpdateURL = r.URL + "/services"

	ctx, err := service.Start(context.Background(),
		host,
		port,
		r,
		rundooportal.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}

	<-ctx.Done()
	fmt.Println("Shutting down rundoo portal")

}
