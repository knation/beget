// Functions associated with Kafka production
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package downstream

import (
	"beget/util"
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer
var KafkaTopics map[string]bool

// Initializes the Kafka connection given env variables provided
func Init(mode util.ServiceMode) error {

	// Parse topics
	KafkaTopics = make(map[string]bool)
	if t := os.Getenv("KAFKA_TOPICS"); t != "" {
		for _, topic := range strings.Split(t, ",") {
			KafkaTopics[topic] = true
		}

	} else {
		return fmt.Errorf("no topics specified")
	}

	if mode == util.ReleaseMode {
		// Check for Kafka host or brokers
		var kafkaHosts net.Addr
		if b := os.Getenv("KAFKA_BROKERS"); b != "" {
			kafkaHosts = kafka.TCP(strings.Split(b, ",")...)

		} else {
			return fmt.Errorf(`must provide either "KAFKA_BROKERS"`)
		}

		// All options can be found here: https://pkg.go.dev/github.com/segmentio/kafka-go?utm_source=godoc#Writer
		KafkaWriter = &kafka.Writer{
			Addr:     kafkaHosts,
			Balancer: &kafka.LeastBytes{},
		}
	}

	return nil
}

// Closes active downstream connections
func Close() error {
	if KafkaWriter != nil {
		return KafkaWriter.Close()
	}
	return nil
}

// Writes the given message to Kafka. This syntax allows us to stub the function
// for testing.
var KafkaProduce = func(mode util.ServiceMode, m kafka.Message) error {
	// If not in release mode, log message
	if mode == util.DebugMode {
		util.Sugar.Debugf("PRODUCE: %v", m)
		return nil
	}

	return KafkaWriter.WriteMessages(context.Background(), m)
}
