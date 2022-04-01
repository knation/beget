// Main HTTP handler
//
// Author: Kirk Morales
// Copyright 2022. All Rights Reserved.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

// Expected request body for producing to a topic.
type ProduceBody struct {
	Topic    string      // The topic to write the message to (required)
	Key      string      // The key of the message (optional)
	Value    interface{} // The message value as JSON (required if no ValueStr)
	ValueStr []byte      // The message value as a byte slice (required if no Value)
}

// Validates the given request context. Returns true if valid with the
// message body, false otherwise.
func validateRequest(c *gin.Context) (bool, *ProduceBody) {
	// Bind and validate JSON payload
	var body ProduceBody
	if err := c.BindJSON(&body); err != nil {
		c.String(400, "Invalid Body")
		c.Abort()
		return false, nil
	}

	// Look for required "topic" value
	if body.Topic == "" {
		c.String(400, "Missing `topic` parameter")
		c.Abort()
		return false, nil
	}
	fmt.Printf("\n\n%v\n\n", body)

	// Make sure we have a body (either as JSON or a string)
	if body.Value == nil && len(body.ValueStr) == 0 {
		c.String(400, "Missing message value")
		c.Abort()
		return false, nil
	} else if body.Value != nil {
		// Test if string
		if reflect.TypeOf(body.Value) == reflect.TypeOf("") {
			body.ValueStr = []byte(body.Value.(string))
		} else {
			// JSON provided. Encode it to a string. No need to capture error --
			// since the body was decoded to begin with, we know this is valid JSON
			str, _ := json.Marshal(body.Value)
			body.ValueStr = str
		}
	}

	// Make sure topic is allowed
	if !topics[body.Topic] {
		c.String(404, "Invalid topic specified")
		c.Abort()
		return false, nil
	}

	return true, &body
}

// Handles a request. This is called after `validateRequest`, so you can
// ensure that `body` is valid and contains the required fields (and that
// the given topic is allowed)
func requestHandler(c *gin.Context, body *ProduceBody) {
	m := kafka.Message{
		Topic: body.Topic,
		Value: body.ValueStr,
	}

	if body.Key != "" {
		m.Key = []byte(body.Key)
	}

	writeMessage(m)

	c.Status(200)
}

// Writes the given message to Kafka. This syntax allows us to stub the function
// for testing.
var writeMessage = func(m kafka.Message) {
	// If not in release mode, log message
	if serviceMode == "debug" {
		sugar.Debugf("PRODUCE: %v", m)
		return
	}

	if err := kafkaWriter.WriteMessages(context.Background(), m); err != nil {
		sugar.Error("failed to write kafka messages:", err)
	}
}
