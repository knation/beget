package handler

import (
	"beget/downstream"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidBody(t *testing.T) {
	t.Run("missing content-type header", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(``)))
		req, _ := http.NewRequest(http.MethodGet, "/healthz", requestBody)

		body, ok := validate(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 415, res.StatusCode)
		assert.Equal(t, "missing Content-Type header\n", string(data))
	})

	t.Run("invalid content-type header", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(``)))
		req, _ := http.NewRequest(http.MethodGet, "/healthz", requestBody)
		req.Header.Add("Content-Type", "text/plain")

		body, ok := validate(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 415, res.StatusCode)
		assert.Equal(t, "Content-Type header is not application/json\n", string(data))
	})

	t.Run("empty body", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(``)))
		req, _ := http.NewRequest(http.MethodGet, "/healthz", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "Request body must not be empty\n", string(data))
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`<>`)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "Request body contains badly-formed JSON (at position 1)\n", string(data))
	})

	t.Run("empty JSON body (missing `topic`)", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{}`)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "missing topic\n", string(data))
	})

	t.Run("invalid JSON value", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo","value":"badjson}`)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "Request body contains badly-formed JSON\n", string(data))
	})

	t.Run("invalid topic", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo","value":{"foo":1}}`)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "invalid topic\n", string(data))
	})

	t.Run("invalid data type", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":false,"value":{"foo":1}}`)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "Request body contains an invalid value for the \"Topic\" field (at position 14)\n", string(data))
	})

	t.Run("invalid field", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"foobar":false,"value":{"foo":1}}`)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "Request body contains unknown field \"foobar\"\n", string(data))
	})

	t.Run("more than just the json", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo","value":{"foo":1}}something else`)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "Request body must only contain a single JSON object\n", string(data))
	})

	t.Run("more than just the json", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		// Make a request larger than 1MB
		requestBodyRaw := []byte(`{"topic":"foo","value":{"foo":"`)
		for count := 0; count <= 1048576; count++ {
			requestBodyRaw = append(requestBodyRaw, " "...)
		}
		requestBodyRaw = append(requestBodyRaw, `"}}`...)

		requestBody := ioutil.NopCloser(bytes.NewReader(requestBodyRaw))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 413, res.StatusCode)
		assert.Equal(t, "Request body must not be larger than 1MB\n", string(data))
	})

	t.Run("missing message value", func(t *testing.T) {
		downstream.KafkaTopics = make(map[string]bool)
		downstream.KafkaTopics["foo"] = true

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo"}`)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.False(t, ok)
		assert.Nil(t, body)
		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "missing message value\n", string(data))

		downstream.KafkaTopics = make(map[string]bool)
	})
}

func TestValidRequest(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		downstream.KafkaTopics = make(map[string]bool)
		downstream.KafkaTopics["foo"] = true

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo","value":{"foo":1}}`)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)
		downstream.KafkaTopics["foo"] = false

		assert.True(t, ok)

		expected := &RequestBody{
			Topic: "foo",
			Value: map[string]interface{}{
				"foo": float64(1),
			},
			valueStr: []byte(`{"foo":1}`),
		}

		assert.Equal(t, expected, body)

		downstream.KafkaTopics = make(map[string]bool)
	})

	t.Run("string", func(t *testing.T) {
		downstream.KafkaTopics = make(map[string]bool)
		downstream.KafkaTopics["foo"] = true

		w := httptest.NewRecorder()

		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"topic":"foo","value":"foobar"}`)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
		req.Header.Add("Content-Type", "application/json")

		body, ok := validate(w, req)
		downstream.KafkaTopics["foo"] = false

		assert.True(t, ok)

		expected := &RequestBody{
			Topic:    "foo",
			Value:    "foobar",
			valueStr: []byte("foobar"),
		}

		assert.Equal(t, expected, body)
	})
}
