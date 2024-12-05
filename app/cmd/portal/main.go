package main

import (
	"context"
	"flag"
	"fmt"
	stlog "log"
	
	"app/internal/models"
	"app/log"
	"app/registry"
	"app/service"
)

var isLocalhost bool

type application struct {
	productlist *models.RundooModel
}

func main() {
	localhost := flag.Bool("localhost", false, "Run the application in localhost mode")
    flag.Parse()

    // Set the global variable
    isLocalhost = *localhost
	service.IsLocalhost = isLocalhost

	host, port := "portal", "5050"
	// Set the ServicesURL based on the flag
    registry.ServicesURL = "http://registryservice:3000/services"
    if isLocalhost {
		host = "localhost"
		registry.ServicesURL = "http://localhost:3000/services"
	} 
	
	fmt.Printf("We  runs on a %s\n", host)

	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)


	app := &application{
		productlist: &models.RundooModel{Endpoint: fmt.Sprintf("%s/products", serviceAddress)},
	}


	var r registry.ServiceConfig
	if isLocalhost {
		r.Host = "localhost"
		fmt.Println("registry runs on localhost")
	} else {
		fmt.Println("registry runs on a container")
		r.Host = "registryservice"
	}

	
	r.Port = port
	
	app.routes(&r, serviceAddress)

	ctx, err := service.Start(context.Background(), r)
	if err != nil {
		fmt.Println("Portal can't start", err)
		stlog.Fatal(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}

	<-ctx.Done()
	fmt.Println("Shutting down rundoo portal")
	stlog.Println("Shutting down rundoo portal")


}
