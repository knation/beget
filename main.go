// Entrypoint for beget web service.
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package main

import (
	"beget/downstream"
	"beget/handler"
	"beget/util"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Initialize logger
	util.InitLogging()

	// Load configuration
	if err := util.InitConfig(); err != nil {
		util.Sugar.Panic(err)
	}

	// Get web server port
	port := util.Config.Server.Port
	if port <= 0 {
		port = 8080
	}

	util.Sugar.Infof("Starting service in '%s' mode on port %d...", util.Config.App.Mode, port)

	// Initialize kafka or panic if there was a problem
	if err := downstream.Init(); err != nil {
		util.Sugar.Panic(err)
	}

	router := handler.InitRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// Start webserver in background to allow for graceful shutdown code below
	go func() {
		util.Sugar.Infof("Listening on port %d", port)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			util.Sugar.Info(err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	util.Sugar.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		util.Sugar.Fatalf("server forced to shutdown: %s", err.Error())
	}

	// Close Kafka writer
	if err := downstream.Close(); err != nil {
		util.Sugar.Fatalf("failed to close writer: %s", err.Error())
	}

	util.Sugar.Info("Server exiting")
}
