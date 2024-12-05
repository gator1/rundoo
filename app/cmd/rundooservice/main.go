package main

import (
	
	"context"
	"database/sql"
	"flag"
	"fmt"
	stlog "log"
	"time"
	

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
	stlog.Printf("before open db dsn: %v\n", dsn)	
	var db *sql.DB
	var err error

	for i := 0; i < 3; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			break
		}
		stlog.Printf("Attempt %d: open db failed dsn: %v, error: %v\n", i+1, dsn, err)
		time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
	}

	if err != nil {
		stlog.Printf("Failed to open db after 3 attempts dsn: %v\n", dsn)
		stlog.Fatal(err)
	}
	stlog.Printf("after open db success dsn: %v\n", dsn)	

	defer db.Close()

	for i := 0; i < 3; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		stlog.Printf("Attempt %d: ping db failed dsn: %v, error: %v\n", i+1, dsn, err)
		time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
	}

	if err != nil {
		stlog.Printf("Failed to ping db after 3 attempts dsn: %v\n", dsn)
		stlog.Fatal(err)
	}
	
	stlog.Printf("database connection pool established")
	
	r. Name = registry.RundooService
	r.Host = host
	r.Port = port
	
	app := &application{
		models: data.NewModels(db),
		handler: &rundoo.ProductsHandler{},	
	}
	// configure our service
	productService := rundoo.NewService(&app.models)


	app.routes(&r, serviceAddress, productService)

	stlog.Printf("rundooService rundoo-api runs on a %s\n", host)

	ctx, err := service.Start(context.Background(), r)
	if err != nil {
		stlog.Printf("dservice.Start failed: %v\n", err)
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}

	<-ctx.Done()
	stlog.Println("Shutting down product service")
}
