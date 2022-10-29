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
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

var httpLogger zerolog.Logger

// Initializes the gin engine
func InitRouter() http.Handler {

	// Logger
	conciseLogging := false

	if util.Config.App.Mode == util.DebugMode {
		conciseLogging = true
	}

	httpLogger = httplog.NewLogger("httplog-example", httplog.Options{
		JSON:    true,
		Concise: conciseLogging,
	})

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(httpLogger))
	r.Use(middleware.Recoverer)

	r.Use(middleware.Heartbeat("/healthz"))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(30 * time.Second))

	r.Post("/produce", topicProduceHandler)

	return r
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

	if err := downstream.KafkaProduce(message); err != nil {
		util.Sugar.Error("failed to write kafka messages:", err)
	}

	w.Write([]byte("OK"))
}
