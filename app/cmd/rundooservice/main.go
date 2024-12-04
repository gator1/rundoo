package main

import (
	
	"context"
	"database/sql"
	"flag"
	"fmt"
	stlog "log"
	

	_ "github.com/lib/pq"

	"app/internal/data"
	"app/log"
	"app/rundoo"
	"app/registry"
	"app/service"
)

var isLocalhost bool

type application struct {
	models data.Models
	handler *rundoo.ProductsHandler
}

func main() {
	localhost := flag.Bool("localhost", false, "Run the application in localhost mode")
    flag.Parse()

    // Set the global variable
    isLocalhost = *localhost
	service.IsLocalhost = isLocalhost

	host, port := "rundoo-api", "6000"
	registry.ServicesURL = "http://registryservice:3000/services"
	if isLocalhost {
		registry.ServicesURL = "http://localhost:3000/services"
		host = "localhost"
	}
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.ServiceConfig

	dsn := "postgres://postgres:uber@rundoo-db/postgres?sslmode=disable"
	if isLocalhost {
		dsn = "postgres://postgres:mysecretpassword@localhost/rundoo?sslmode=disable"
	}
	stlog.Printf("before open db dsn: %v", dsn)	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		stlog.Printf(" open db failed dsn: %v", dsn)	
		stlog.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		stlog.Fatal(err)
	}

	stlog.Printf("database connection pool established")

	r.Host = host
	r.Port = port
	
	app := &application{
		models: data.NewModels(db),
		handler: &rundoo.ProductsHandler{},	
	}
	// configure our service
	productService := rundoo.NewService(&app.models)


	app.routes(&r, serviceAddress, productService)


	ctx, err := service.Start(context.Background(), r)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}

	<-ctx.Done()
	fmt.Println("Shutting down product service")
}
