// Global stuff that, frankly, doesn't have a better home.
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package util

// Defines the `serviceMode` type used for an enum of acceptable modes to run the app in.
type ServiceMode string

const (
	ReleaseMode ServiceMode = "release"
	DebugMode   ServiceMode = "debug"
)
