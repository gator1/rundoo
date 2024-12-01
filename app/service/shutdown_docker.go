// +build docker

package service

import (
	"app/registry"
	"context"
	"fmt"
	"net/http"
	"os"
    "os/signal"
    "syscall"	
)

func waitForShutdown(ctx context.Context, srv *http.Server, config registry.ServiceConfig, cancel context.CancelFunc) {
	fmt.Printf("%v started. Press any key to stop.\n", config.Name)
	// Set up signal handling to gracefully shut down the server
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        fmt.Println("Received shutdown signal")
        srv.Shutdown(ctx)
        cancel()
    }()
}

