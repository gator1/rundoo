package main

import (
	"context"
	"fmt"
	stlog "log"
	
	"app/internal/models"
	"app/log"
	rundooportal "app/portal"
	"app/registry"
	"app/service"
)

type application struct {
	productlist *models.RundooModel
	handler *rundooportal.RundooHandler

}

func main() {
	err := rundooportal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}
	host, port := "localhost", "5050"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)


	app := &application{
		productlist: &models.RundooModel{Endpoint: fmt.Sprintf("%s/products", serviceAddress)},
	}


	var r registry.ServiceConfig
	r.Host = host
	r.Port = port
	
	app.routes(&r, serviceAddress)

	ctx, err := service.Start(context.Background(), r)
	if err != nil {
		stlog.Fatal(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}

	<-ctx.Done()
	fmt.Println("Shutting down rundoo portal")
	stlog.Println("Shutting down rundoo portal")


}
