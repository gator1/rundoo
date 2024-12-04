package main

import (
	"app/registry"
	"app/service"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
    "os/signal"
    "syscall"
)

var isLocalhost bool

func main() {
	localhost := flag.Bool("localhost", false, "Run the application in localhost mode")
    flag.Parse()

    // Set the global variable
    isLocalhost = *localhost
	service.IsLocalhost = isLocalhost

	registry.SetupRegistryService()
	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = ":3000"

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	// Set up signal handling to gracefully shut down the server
    sigChan := make(chan os.Signal, 1)
    	
	if !isLocalhost {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	}
    go func() {
		if isLocalhost {
			fmt.Println("Registry service started. Press any key to stop.")
			var s string
			fmt.Scanln(&s)
			srv.Shutdown(ctx)
			cancel()
		} else {
			<-sigChan
			fmt.Println("Received shutdown signal")
			srv.Shutdown(ctx)
			cancel()
		}
    }()


	<-ctx.Done()
	fmt.Println("Shutting down registry service")
}
