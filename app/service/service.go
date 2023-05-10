package service

import (
	"app/registry"
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, config registry.ServiceConfig) (context.Context, error) {
	
	var err error

	ctx = startService(ctx, config)

	if config.Mux != nil {
		err = registry.RegisterServiceMux(config)
	} else {
		err = registry.RegisterService(config)	
	}
	
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, config registry.ServiceConfig) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = ":" + config.Port

	if config.Mux != nil {
		srv.Handler = config.Mux
	} else if config.HttpHandler != nil {
		if config.Name != registry.LogService {
			srv.Handler = config.HttpHandler
		}
	}

	go func() {
		log.Printf("Starting HTTP service '%s'", config.Name)
		log.Println(srv.ListenAndServe())
		cancel()
	}()


	go func() {
		fmt.Printf("%v started. Press any key to stop.\n", config.Name)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%v:%v", config.Host, config.Port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
