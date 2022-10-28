// Functions related to logging
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package util

import "go.uber.org/zap"

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func InitLogging() {
	Logger, _ = zap.NewProduction()
	defer Logger.Sync() // flushes buffer, if any
	Sugar = Logger.Sugar()
}
