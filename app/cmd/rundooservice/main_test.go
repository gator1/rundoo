package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jarcoal/httpmock"

	"app/internal/data"
	"app/log"
	"app/registry"
	"app/rundoo"
	"app/service"

)

func TestService(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T, db *sql.DB,
	){
		"Test Main": testMain,
		
	} {
		t.Run(scenario, func(t *testing.T) {
			// Create a mock DB and a mock connection
			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %s", err)
			}
			defer db.Close()

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			expectedResponse := "Mocked response"
			httpmock.RegisterResponder("POST", "", httpmock.NewStringResponder(200, expectedResponse))
			httpmock.RegisterResponder("POST", "http://localhost:3000/services", httpmock.NewStringResponder(200, expectedResponse))


			fn(t, db)
		})
	}
}

func testMain(t *testing.T, db *sql.DB) {

	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)


	app := &application{
		models: data.NewModels(db),
		handler: &rundoo.ProductsHandler{},	
	}
	var r registry.ServiceConfig

	// configure our service
	productService := rundoo.NewService(&app.models)


	app.routes(&r, serviceAddress, productService)


	_, err := service.Start(context.Background(), r)
	require.NoError(t, err)

	
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.Name)
	}
	
	
}

