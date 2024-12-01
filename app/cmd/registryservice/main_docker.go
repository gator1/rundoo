// +build docker

package main

import (
    "app/registry"
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
)

func main() {
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
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        fmt.Println("Received shutdown signal")
        srv.Shutdown(ctx)
        cancel()
    }()

    <-ctx.Done()
    fmt.Println("Shutting down registry service")
}