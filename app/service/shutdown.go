// +build !docker

package service

import (
	"app/registry"
	"context"
	"fmt"
	"log"
	
	"net/http"
)

func waitForShutdown(ctx context.Context, srv *http.Server, config registry.ServiceConfig, cancel context.CancelFunc) {
	
	
	fmt.Printf("%v started. Press any key to stop.\n", config.Name)
	var s string
	fmt.Scanln(&s)
	err := registry.ShutdownService(fmt.Sprintf("http://%v:%v", config.Host, config.Port))
	if err != nil {
		log.Println(err)
	}
	srv.Shutdown(ctx)
	cancel()

}