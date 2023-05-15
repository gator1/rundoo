package main

import (
	"net/http"

	rundooportal "app/portal"
	"app/registry"
)

func (app *application) routes(r *registry.ServiceConfig, serviceAddress string) *http.ServeMux {
	handler := &rundooportal.RundooHandler{}
	app.handler = handler
	r.Name = registry.RundooPortal
	r.URL = serviceAddress
	r.HeartbeatURL = r.URL + "/heartbeat"
	r.RequiredServices = []registry.ServiceName{
		registry.LogService,
		registry.RundooService,
	}
	r.UpdateURL = r.URL + "/services"
	r.HttpHandler = handler

	r.Mux = http.NewServeMux()
	//r.Mux.Handle("/products", handler)
	//r.Mux.Handle("/products/", handler)
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	r.Mux.HandleFunc("/", app.home)
	r.Mux.HandleFunc("/product/view/", app.productView)
	r.Mux.HandleFunc("/product/create", app.productCreate)

	return r.Mux
}