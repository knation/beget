package handler

import (
	"beget/downstream"
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestInitRouter(t *testing.T) {
	r := InitRouter()
	assert.NotNil(t, r)
}

func TestProduceHandlerFailure(t *testing.T) {

	results := make([]kafka.Message, 0)
	stubKafkaProduce := downstream.KafkaProduce
	downstream.KafkaProduce = func(ctx context.Context, m kafka.Message) error {
		results = append(results, m)
		return nil
	}

	t.Run("failure", func(t *testing.T) {

		w := httptest.NewRecorder()
		requestBody := ioutil.NopCloser(bytes.NewReader([]byte(``)))
		req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)

		topicProduceHandler(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		assert.Equal(t, 415, res.StatusCode)
		assert.Equal(t, "missing Content-Type header\n", string(data))
		assert.Empty(t, results)
	})

	t.Run("success", func(t *testing.T) {
		downstream.KafkaTopics = make(map[string]struct{})
		downstream.KafkaTopics["foo"] = struct{}{}
		downstream.KafkaTopics["bar"] = struct{}{}

		tests := []string{
			`{"topic":"foo","value":{"foo":1}}`,
			`{"topic":"bar","value":"foobar"}`,
			`{"topic":"bar","value":"foobar","key":"somekey"}`,
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

		for _, test := range tests {
			w := httptest.NewRecorder()
			requestBody := ioutil.NopCloser(bytes.NewReader([]byte(test)))
			req, _ := http.NewRequest(http.MethodGet, "/produce", requestBody)
			req.Header.Add("Content-Type", "application/json")

			topicProduceHandler(w, req)

			///
			res := w.Result()
			defer res.Body.Close()
			_, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			// fmt.Println(string(data))
			///
		}

		for i := range tests {
			if len(results) >= (i + 1) {
				assert.Equal(t, expected[i], results[i])
			} else {
				assert.FailNow(t, "invalid result length")
			}
		}

		downstream.KafkaTopics = make(map[string]struct{})
	})

	// Restore stubs
	downstream.KafkaProduce = stubKafkaProduce
}
