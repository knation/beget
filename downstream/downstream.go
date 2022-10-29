// Functions associated with Kafka production
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package downstream

import (
	"beget/util"
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer
var KafkaTopics map[string]struct{} = make(map[string]struct{})

// Initializes the Kafka connection given env variables provided
func Init() error {

	// Parse topics
	if len(util.Config.Kafka.Topics) > 0 {
		for _, topic := range util.Config.Kafka.Topics {
			KafkaTopics[topic] = struct{}{}
		}
	} else {
		return fmt.Errorf("no topics provided")
	}

	if util.Config.App.Mode == util.ReleaseMode {
		// Check for Kafka host or brokers
		if len(util.Config.Kafka.Brokers) == 0 {
			return fmt.Errorf("no brokers provided")
		}

		// All options can be found here: https://pkg.go.dev/github.com/segmentio/kafka-go?utm_source=godoc#Writer
		KafkaWriter = &kafka.Writer{
			Addr:     kafka.TCP(util.Config.Kafka.Brokers...),
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
var KafkaProduce = func(m kafka.Message) error {
	// If not in release mode, log message
	if util.Config.App.Mode == util.DebugMode {
		util.Sugar.Debugf("PRODUCE: %v", m)
		return nil
	}

	return KafkaWriter.WriteMessages(context.Background(), m)
}
