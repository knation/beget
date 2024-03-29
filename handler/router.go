// Route handling
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package handler

import (
	"beget/downstream"
	"beget/util"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

var httpLogger zerolog.Logger

// Initializes the gin engine
func InitRouter() http.Handler {

	r := chi.NewRouter()
	r.Use(util.HttpLogger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Heartbeat("/healthz"))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(time.Duration(util.Config.Server.Timeout) * time.Second))

	r.Post("/produce", topicProduceHandler)

	return r
}

// Handles a request for topic production
func topicProduceHandler(w http.ResponseWriter, r *http.Request) {
	var body *RequestBody
	var ok bool

	// NOTE: We expect that `validate` is not making any external requests or otherwise
	// doing any long-running work. If that changes, it should become context aware
	// to gracefully stop what it's doing immediately upon timeout
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

	// NOTE: We intentionally do not pass `r.Context()` down the chain to the Kafka
	// writer because we don't believe an HTTP timeout should cause writing to cease immediately.
	// If you feel differently, can you change that below.
	if err := downstream.KafkaProduce(context.Background(), message); err != nil {
		util.Sugar.Error("failed to write kafka messages:", err)
	}

	w.Write([]byte("OK"))
}
