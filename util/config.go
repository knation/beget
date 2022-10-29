// Configuration data and functions
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package util

import (
	"fmt"

	"github.com/spf13/viper"
)

// Defines the `serviceMode` type used for an enum of acceptable modes to run the app in.
type ServiceMode string

const (
	ReleaseMode ServiceMode = "release"
	DebugMode   ServiceMode = "debug"
)

type Configuration struct {
	App struct {
		Mode ServiceMode
	}
	Server struct {
		Port int
	}
	Kafka struct {
		Brokers []string
		Topics  []string
	}
}

var Config Configuration

func InitConfig() error {
	// Set the file name and type of the configurations file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return fmt.Errorf("no configuration found")
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("error parsing configuration")
		}
	}

	// Set configuration defaults
	viper.SetDefault("app.mode", "debug")
	viper.SetDefault("server.port", 8080)

	// Get configuration into our `Config` variable
	err := viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	return nil
}
