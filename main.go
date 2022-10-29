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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var mode util.ServiceMode

func main() {

	// Initialize logger
	util.InitLogging()

	// Read config
	// if err := viper.ReadInConfig(); err != nil {
	// 	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	// 		// Config file not found; ignore error if desired
	// 		util.Sugar.Panic(errors.New("no configuration found"))
	// 	} else {
	// 		// Config file was found but another error was produced
	// 		util.Sugar.Panic(errors.New("error parsing configuration"))
	// 	}
	// }

	// Get mode to launch in: debug|release
	switch os.Getenv("MODE") {
	case "debug":
		mode = util.DebugMode
	case "release":
		mode = util.ReleaseMode
	default:
		mode = util.DebugMode
	}

	// Get web server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	util.Sugar.Infof("Starting service in '%s' mode on port %s...", mode, port)

	// Initialize kafka or panic if there was a problem
	if err := downstream.Init(mode); err != nil {
		util.Sugar.Panic(err)
	}

	router := handler.InitRouter(mode)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start webserver in background to allow for graceful shutdown code below
	go func() {
		util.Sugar.Infof("Listening on port %v", port)
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
