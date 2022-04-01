// Entrypoint for beget web service.
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var serviceMode string
var logger *zap.Logger
var sugar *zap.SugaredLogger
var kafkaWriter *kafka.Writer
var topics map[string]bool

func main() {

	// Initialize logger
	logger, _ = zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()

	// Get mode to launch in: debug|release
	serviceMode = os.Getenv("MODE")
	if serviceMode == "" {
		serviceMode = "release"
	}

	// Get web server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Check for Kafka host or brokers
	var kafkaHosts net.Addr
	if b := os.Getenv("KAFKA_BROKERS"); b != "" {
		kafkaBrokers := strings.Split(b, ",")
		if len(kafkaBrokers) == 0 {
			sugar.Panic("No `KAFKA_BROKERS` provided")
		}
		kafkaHosts = kafka.TCP(kafkaBrokers...)

	} else if h := os.Getenv("KAFKA_HOST"); h != "" {
		kafkaHosts = kafka.TCP(h)

	} else {
		sugar.Panic("Must provide either `KAFKA_BROKERS` or `KAFKA_HOST`")
	}

	// Parse topics
	topics = make(map[string]bool)
	if t := strings.Split(os.Getenv("TOPICS"), ","); len(t) > 0 {
		for _, topic := range t {
			topics[topic] = true
		}

	} else {
		sugar.Panic("No topics specified")
	}

	sugar.Infof("Starting service in '%s' mode...", serviceMode)

	if serviceMode == "release" {
		kafkaWriter = &kafka.Writer{
			Addr:     kafkaHosts,
			Balancer: &kafka.LeastBytes{},
		}
	}

	// Configure Gin
	gin.SetMode(serviceMode) // Set the run mode (release/debug)

	router := initRouter()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start webserver in background to allow for graceful shutdown code below
	go func() {
		sugar.Infof("Listening on port %v", port)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			sugar.Info(err.Error())
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
	sugar.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalf("server forced to shutdown: %s", err.Error())
	}

	// Close Kafka writer
	if err := kafkaWriter.Close(); err != nil {
		sugar.Fatalf("failed to close writer: %s", err.Error())
	}

	sugar.Info("Server exiting")
}

// Initializes the gin engine
func initRouter() *gin.Engine {
	router := gin.New()        // Create the router
	router.Use(gin.Recovery()) // Recovery middleware recovers from any panics and writes a 500 if there was one

	// Log requests
	if serviceMode == "release" {
		router.Use(releaseGinLogger)
	} else {
		router.Use(debugGinLogger)
	}

	router.Any("/healthz", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.POST("/produce", func(c *gin.Context) {
		ok, body := validateRequest(c)
		if ok {
			requestHandler(c, body)
		}
	})

	return router
}

// Middleware for logging requests from gin in "release" mode
func releaseGinLogger(c *gin.Context) {
	// Don't log health checks
	if c.Request.URL.Path == "/healthz" {
		c.Next()
		return
	}

	start := time.Now()

	c.Next()

	duration := time.Since(start)

	sugar.Infow(c.Request.URL.Path,
		"timestamp", start.Format(time.RFC3339),
		"status", c.Writer.Status(),
		"method", c.Request.Method,
		"query", c.Request.URL.RawQuery,
		"ip", c.ClientIP(),
		"user-agent", c.Request.UserAgent(),
		"errors", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		"duration", duration,
	)
}

// Middleware for logging requests from gin in "debug" mode
func debugGinLogger(c *gin.Context) {
	start := time.Now()

	c.Next()

	sugar.Infow(c.Request.URL.Path,
		"timestamp", start.Format(time.RFC3339),
		"status", c.Writer.Status(),
		"method", c.Request.Method,
		"query", c.Request.URL.RawQuery,
		"errors", c.Errors.ByType(gin.ErrorTypePrivate).String(),
	)
}
