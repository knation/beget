package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestInvalidBody(t *testing.T) {
	t.Run("empty body", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(``)))
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/produce", requestBody)

		ok, body := validateRequest(ctx)
		assert.False(t, ok)
		assert.Nil(t, body)
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`<>`)))
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/produce", requestBody)

		ok, body := validateRequest(ctx)
		assert.False(t, ok)
		assert.Nil(t, body)
	})

	t.Run("empty JSON body (missing `topic`)", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{}`)))
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/produce", requestBody)

		ok, body := validateRequest(ctx)
		assert.False(t, ok)
		assert.Nil(t, body)
	})

	t.Run("missing message value", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo"}`)))
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/produce", requestBody)

		ok, body := validateRequest(ctx)
		assert.False(t, ok)
		assert.Nil(t, body)
	})

	t.Run("invalid JSON value", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo","value":"badjson"}`)))
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/produce", requestBody)

		ok, body := validateRequest(ctx)
		assert.False(t, ok)
		assert.Nil(t, body)
	})

	t.Run("invalid topic", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo","value":{"foo":1}}`)))
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/produce", requestBody)

		ok, body := validateRequest(ctx)
		assert.False(t, ok)
		assert.Nil(t, body)
	})
}

func TestValidBody(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		topics = make(map[string]bool)
		topics["foo"] = true

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo","value":{"foo":1}}`)))
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/produce", requestBody)

		ok, body := validateRequest(ctx)
		topics["foo"] = false

		assert.True(t, ok)

		expected := &ProduceBody{
			Topic: "foo",
			Value: map[string]interface{}{
				"foo": float64(1),
			},
			ValueStr: []byte(`{"foo":1}`),
		}

		assert.Equal(t, expected, body)
	})

	t.Run("string", func(t *testing.T) {
		topics = make(map[string]bool)
		topics["foo"] = true

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo","value":"foobar"}`)))
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/produce", requestBody)

		ok, body := validateRequest(ctx)
		topics["foo"] = false

		assert.True(t, ok)

		expected := &ProduceBody{
			Topic:    "foo",
			Value:    "foobar",
			ValueStr: []byte("foobar"),
		}

		assert.Equal(t, expected, body)
	})
}

func TestRequestHandler(t *testing.T) {
	tests := []*ProduceBody{
		{
			Topic: "foo",
			Value: map[string]interface{}{
				"foo": float64(1),
			},
			ValueStr: []byte(`{"foo":1}`),
		},
		{
			Topic:    "bar",
			Value:    "foobar",
			ValueStr: []byte("foobar"),
		},
		{
			Topic:    "bar",
			Key:      "somekey",
			Value:    "foobar",
			ValueStr: []byte("foobar"),
		},
	}

	expected := []kafka.Message{
		{
			Topic: "foo",
			Value: []byte(`{"foo":1}`),
		},
		{
			Topic: "bar",
			Value: []byte("foobar"),
		},
		{
			Topic: "bar",
			Value: []byte("foobar"),
			Key:   []byte("somekey"),
		},
	}

	results := make([]kafka.Message, 0)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Stub `writeMessage` -- store all messages passes to the function in order
	stubWriteMessage := writeMessage
	writeMessage = func(m kafka.Message) {
		results = append(results, m)
	}

	for _, test := range tests {
		requestHandler(ctx, test)
	}

	// Test all results
	for i, _ := range tests {
		assert.Equal(t, expected[i], results[i])
	}

	// Restore stubs
	writeMessage = stubWriteMessage
}
