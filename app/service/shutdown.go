package service

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

var IsLocalhost bool

func waitForShutdown(ctx context.Context, srv *http.Server, config registry.ServiceConfig, cancel context.CancelFunc) {
	
	if IsLocalhost {
		fmt.Printf("%v started. Press any key to stop.\n", config.Name)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%v:%v", config.Host, config.Port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx)
		cancel()
	} else {
		fmt.Printf("%v started. Wait for Sigterm, docker.\n", config.Name)
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
	

}