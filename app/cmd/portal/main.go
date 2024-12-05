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
	
	stlog.Printf("rundoo portal runs on a %s\n", host)

	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)


	app := &application{
		productlist: &models.RundooModel{Endpoint: fmt.Sprintf("%s/products", serviceAddress)},
	}


	var r registry.ServiceConfig
	r.Name = registry.RundooPortal
	if isLocalhost {
		r.Host = "localhost"
		stlog.Println("registry runs on localhost")
	} else {
		stlog.Println("registry runs on a container")
		r.Host = "portal"
	}

	
	r.Port = port
	
	app.routes(&r, serviceAddress)

	ctx, err := service.Start(context.Background(), r)
	if err != nil {
		stlog.Println("Portal can't start", err)
		stlog.Fatal(err)
	}

	fmt.Printf("rundoo portal before GetProviderfor LogService %s\n", host)
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}
	stlog.Printf("rundoo portal before GetProviderfor RundooService %s\n", host)
	fmt.Printf("rundoo portal before GetProviderfor RundooService %s\n", host)
	if _, err := registry.GetProvider(registry.RundooService); err != nil {
		stlog.Println("rundoo-api is not avilable in the registry for portal")
		fmt.Println("rundoo-api is not avilable in the registry for portal")
	}
	<-ctx.Done()
	stlog.Println("Shutting down rundoo portal")


}
