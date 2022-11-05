// Functions related to logging
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package util

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func InitLogging() {

	config := zap.NewProductionConfig()

	// Default timestamp is epoch. Use ISO8601 instead
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Shorten message key to save bytes
	config.EncoderConfig.MessageKey = "m"

	Logger, _ = config.Build()

	defer Logger.Sync() // flushes buffer, if any

	Sugar = Logger.Sugar()
}

// HttpLogger returns a new go-chi logging middleware configured
// using the options defined in `Config.Server.HttpLogging`.
func HttpLogger(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Log if we're not skipping the health check or it's not a request to `/healthz`
		if !Config.Server.HttpLogging.SkipHealthCheck || r.URL.Path != "/healthz" {
			start := time.Now()
			defer func() {

				ua := "-"
				if !Config.Server.HttpLogging.SkipUserAgent {
					ua = r.UserAgent()
				}

				Sugar.Infow("http",
					"path", r.URL.Path,
					"status", ww.Status(),
					"method", r.Method,
					"query", r.URL.Query(),
					"size", ww.BytesWritten(),
					"user-agent", ua,
					"duration", time.Since(start),
				)
			}()
		}
		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}
