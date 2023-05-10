package main

import (
	"app/log"
	"app/registry"
	"app/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./app.log")

	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.ServiceConfig
	r.Name = registry.LogService
	r.URL = serviceAddress
	r.HeartbeatURL = r.URL + "/heartbeat"
	r.RequiredServices = make([]registry.ServiceName, 0)
	r.UpdateURL = r.URL + "/services"

	ctx, err := service.Start(context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down log service")

}
