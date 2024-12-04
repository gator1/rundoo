package main

import (
	"app/log"
	"app/registry"
	"app/service"
	"context"
	"flag"
	"fmt"
	stlog "log"
	"net/http"

)

var isLocalhost bool

func main() {
	log.Run("./app.log")
	localhost := flag.Bool("localhost", false, "Run the application in localhost mode")
    flag.Parse()

    // Set the global variable
    isLocalhost = *localhost
	registry.ServicesURL = "http://registryservice:3000/services"
	service.IsLocalhost = isLocalhost
	host, port := "logservice", "4000"
	if isLocalhost {
		host = "localhost"
		registry.ServicesURL = "http://localhost:3000/services"
	} 
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
	r.HttpHandler = handler
	r.Mux = http.NewServeMux()
	r.Mux.Handle("/log", handler)
	r.Mux.Handle("/log/", handler)

	fmt.Printf("starting log service on %s with registry at %s Host = %s port=%s \n", serviceAddress, registry.ServicesURL,	r.Host, r.Port)

	ctx, err := service.Start(context.Background(), r)
	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down log service")

}
