package main

import (
	"fmt"
	"log"
	"net/http"

	"app/registry"
)

func (app *application) routes(r *registry.ServiceConfig, serviceAddress string) *http.ServeMux {
	r.Name = registry.RundooPortal
	r.URL = serviceAddress
	r.HeartbeatURL = r.URL + "/heartbeat"
	r.RequiredServices = []registry.ServiceName{
		registry.LogService,
		registry.RundooService,
	}
	r.UpdateURL = r.URL + "/services"
	log.Printf("rundoo-web (portal) Service url %s %v\n", r.URL, r)
	fmt.Printf("rundoo-web (portal) Service url %s %v\n", r.URL, r)
	
	r.Mux = http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	r.Mux.HandleFunc("/", app.home)
	r.Mux.HandleFunc("/product/view/", app.productView)
	r.Mux.HandleFunc("/product/create", app.productCreate)
	r.Mux.HandleFunc("/products/search", app.productsSearch)

	return r.Mux
}