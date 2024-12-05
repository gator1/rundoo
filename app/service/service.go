package service

import (
	"app/registry"
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
	
)

func Start(ctx context.Context, config registry.ServiceConfig) (context.Context, error) {
	
	//var err error

	ctx = startService(ctx, config)

	// if config.Mux != nil {
	// 	err = registry.RegisterServiceMux(config)
	// } else {
	// 	err = registry.RegisterService(config)	
	// }
	
	// if err != nil {
	// 	return ctx, err
	// }

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

	if config.GrpcServer != nil {
		portInt, _ := strconv.Atoi(config.Port)
		rpcPort := ":"+strconv.Itoa(portInt + 1)

		con, err := net.Listen("tcp", rpcPort)
		if err != nil {
			log.Printf("Starting gRPC user service listen  error on %s...\n", con.Addr().String())
			panic(err)
		}

		go func() {
			log.Printf("Starting gRPC user service on %s...\n", con.Addr().String())
			err = config.GrpcServer.Serve(con)
			if err != nil {
				log.Printf("Starting gRPC user service serve error on %s...\n", con.Addr().String())
				panic(err)
			}
			cancel()
		}()

	}

	go waitForShutdown(ctx, &srv, config, cancel)
	

	return ctx
}
