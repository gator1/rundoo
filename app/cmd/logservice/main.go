package main

import (
	"app/log"
	"app/registry"
	"app/service"
	"context"
	"fmt"
	stlog "log"
	"net/http"

)

func main() {
	log.Run("./app.log")

	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.ServiceConfig

	handler := &log.LogHandler{}
	
	r.Name = registry.LogService
	r.Host = host
	r.Port = port
	r.URL = serviceAddress
	r.HeartbeatURL = r.URL + "/heartbeat"
	r.RequiredServices = make([]registry.ServiceName, 0)
	r.UpdateURL = r.URL + "/services"
	//r.HttpHandler = http.HandlerFunc((&log.LogHandler{}).ServeHTTP)
	r.HttpHandler = handler
	r.Mux = http.NewServeMux()
	r.Mux.Handle("/log", handler)
	r.Mux.Handle("/log/", handler)

	ctx, err := service.Start(context.Background(), r)
	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down log service")

}
