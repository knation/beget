// Configuration data and functions
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package util

import (
	"bytes"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
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
	Kafka KafkaWriterConfig
}

type KafkaWriterConfig struct {
	Brokers []string
	Topics  []string

	//
	// Below is a subset of initiation options:
	// https://pkg.go.dev/github.com/segmentio/kafka-go#Writer
	//

	// Limit on how many attempts will be made to deliver a message.
	//
	// The default is to try at most 10 times.
	MaxAttempts int `mapstructure:"max_attempts"`

	// WriteBackoffMin optionally sets the smallest amount of time the writer waits before
	// it attempts to write a batch of messages
	//
	// Default: 100ms
	WriteBackoffMin time.Duration `mapstructure:"write_backoff_min"`

	// WriteBackoffMax optionally sets the maximum amount of time the writer waits before
	// it attempts to write a batch of messages
	//
	// Default: 1s
	WriteBackoffMax time.Duration `mapstructure:"write_backoff_max"`

	// Limit on how many messages will be buffered before being sent to a
	// partition.
	//
	// The default is to use a target batch size of 100 messages.
	BatchSize int `mapstructure:"batch_size"`

	// Limit the maximum size of a request in bytes before being sent to
	// a partition.
	//
	// The default is to use a kafka default value of 1048576.
	BatchBytes int64 `mapstructure:"batch_bytes"`

	// Time limit on how often incomplete message batches will be flushed to
	// kafka.
	//
	// The default is to flush at least every second.
	BatchTimeout time.Duration `mapstructure:"batch_timeout"`

	// Timeout for read operations performed by the Writer.
	//
	// Defaults to 10 seconds.
	ReadTimeout time.Duration `mapstructure:"read_timeout"`

	// Timeout for write operation performed by the Writer.
	//
	// Defaults to 10 seconds.
	WriteTimeout time.Duration `mapstructure:"write_timeout"`

	// Number of acknowledges from partition replicas required before receiving
	// a response to a produce request, the following values are supported:
	//
	//  RequireNone (0)  fire-and-forget, do not wait for acknowledgements from the
	//  RequireOne  (1)  wait for the leader to acknowledge the writes
	//  RequireAll  (-1) wait for the full ISR to acknowledge the writes
	//
	// Defaults to RequireNone.
	RequiredAcks kafka.RequiredAcks `mapstructure:"required_acks"`

	// Setting this flag to true causes the WriteMessages method to never block.
	// It also means that errors are ignored since the caller will not receive
	// the returned value. Use this only if you don't care about guarantees of
	// whether the messages were written to kafka.
	//
	// Defaults to false.
	Async bool `mapstructure:"async"`

	// AllowAutoTopicCreation notifies writer to create topic if missing.
	AllowAutoTopicCreation bool `mapstructure:"allow_auto_topic_creation"`
}

var Config Configuration

// InitConfig load configuration from a `config.yaml` in the same directory
// as the executable or by parsing provided flags and ENV variables.
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

	return finalizeConfig()
}

// InitConfig load configuration from the given YAML as a string.
func InitConfigFromYaml(s string) error {
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer([]byte(s))); err != nil {
		return fmt.Errorf("invalid configuration")
	}

	return finalizeConfig()
}

// finalizeConfig finalizes configuration before use, including setting
// defaults and unmarshalling to a struct
func finalizeConfig() error {
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
