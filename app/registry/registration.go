package registry

import (
	"log"
	"net"
	"net/http"
	"sync"

	"google.golang.org/grpc"
)

type ServiceConfig struct {
	Name             ServiceName
	URL              string
	HeartbeatURL     string
	UpdateURL        string
	Host             string
	Port             string
	Mux              *http.ServeMux
	HttpHandler      http.Handler `json:"-"`
	GrpcServer       *grpc.Server `json:"-"`
	RequiredServices []ServiceName
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")
	RundooPortal  = ServiceName("RundooPortal")
	ProductService = ServiceName("ProductService")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}

func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]ServiceConfig),
		mutex: new(sync.RWMutex),
		

	}
}

func (r *Registry) RegisterService(config ServiceConfig) {
	r.services[string(config.Name)] = config
}

func (r *Registry) StartServices() error {
	for _, config := range r.services {
		if config.HttpHandler != nil {
			go func() {
				log.Printf("Starting HTTP service '%s'", config.Name)
				err := http.ListenAndServe(":8080", config.HttpHandler)
				if err != nil {
					log.Fatalf("Failed to start HTTP service '%s': %s", config.Name, err)
				}
			}()
		}
		config := config

		if config.GrpcServer != nil {
			go func() {
				log.Printf("Starting gRPC service '%s'", config.Name)
				lis, err := net.Listen("tcp", ":9090")
				if err != nil {
					log.Fatalf("Failed to start gRPC service '%s': %s", config.Name, err)
				}
				if err := config.GrpcServer.Serve(lis); err != nil {
					log.Fatalf("Failed to serve gRPC service '%s': %s", config.Name, err)
				}
			}()
		}
	}

	return nil
}

func (r *Registry) StopServices() error {
	// TODO: Stop services gracefully
	return nil
}
