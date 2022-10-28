// Route handling
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package handler

import (
	"beget/downstream"
	"beget/util"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/segmentio/kafka-go"
)

var mode util.ServiceMode

// Initializes the gin engine
func InitRouter(m util.ServiceMode) http.Handler {
	mode = m

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/healthz", healthCheckHandler)
	r.Post("/produce", topicProduceHandler)

	return r
}

// Handler for health checks
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// Handles a request for topic production
func topicProduceHandler(w http.ResponseWriter, r *http.Request) {
	var body *RequestBody
	var ok bool

	if body, ok = validate(w, r); !ok {
		return
	}

	message := kafka.Message{
		Topic: body.Topic,
		Value: body.valueStr,
	}

	if body.Key != "" {
		message.Key = []byte(body.Key)
	}

	if err := downstream.KafkaProduce(mode, message); err != nil {
		util.Sugar.Error("failed to write kafka messages:", err)
	}

	w.Write([]byte("OK"))
}
