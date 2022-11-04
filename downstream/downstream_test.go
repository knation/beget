package downstream_test

import (
	"beget/downstream"
	"beget/util"
	"fmt"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestInitDebug(t *testing.T) {

	t.Run("no topics", func(t *testing.T) {
		err := downstream.Init()

		assert.EqualError(t, err, "no topics provided")
	})

	t.Run("empty topics", func(t *testing.T) {
		err := downstream.Init()

		assert.EqualError(t, err, "no topics provided")
	})

	// Set to debug mode for next tests
	util.Config.App.Mode = util.DebugMode

	// Set sample topics
	util.Config.Kafka.Topics = []string{"foo"}

	t.Run("success in debug mode", func(t *testing.T) {
		err := downstream.Init()

		assert.Nil(t, err)
		assert.EqualValues(t, map[string]struct{}{"foo": {}}, downstream.KafkaTopics)

		// Make sure close doesn't break
		err = downstream.Close()
		assert.Nil(t, err)
	})

	// Set to release mode for next tests
	util.Config.App.Mode = util.ReleaseMode

	t.Run("no brokers", func(t *testing.T) {
		err := downstream.Init()

		assert.EqualError(t, err, "no brokers provided")
	})

	t.Run("empty brokers", func(t *testing.T) {
		err := downstream.Init()

		assert.EqualError(t, err, "no brokers provided")
	})

	t.Run("single broker", func(t *testing.T) {
		util.Config.Kafka.Brokers = []string{"broker.foo.com"}

		err := downstream.Init()

		assert.Nil(t, err)
		assert.NotNil(t, downstream.KafkaWriter)
		assert.EqualValues(t, kafka.TCP("broker.foo.com"), downstream.KafkaWriter.Addr)

		// Make sure close doesn't break
		err = downstream.Close()
		assert.Nil(t, err)
	})

	t.Run("multiple brokers", func(t *testing.T) {
		util.Config.Kafka.Brokers = []string{"broker.foo.com", "broker.bar.com"}

		err := downstream.Init()

		assert.Nil(t, err)
		assert.NotNil(t, downstream.KafkaWriter)
		assert.EqualValues(t, kafka.TCP("broker.foo.com", "broker.bar.com"), downstream.KafkaWriter.Addr)

		// Make sure close doesn't break
		err = downstream.Close()
		assert.Nil(t, err)

		// Test default writer options
		stats := downstream.KafkaWriter.Stats()
		assert.Equal(t, int64(0), stats.MaxAttempts)
		assert.Equal(t, time.Duration(0), stats.WriteBackoffMin)
		assert.Equal(t, time.Duration(0), stats.WriteBackoffMax)
		assert.Equal(t, 0, downstream.KafkaWriter.BatchSize)
		assert.Equal(t, int64(0), downstream.KafkaWriter.BatchBytes)
		assert.Equal(t, time.Duration(0), stats.BatchTimeout)
		assert.Equal(t, time.Duration(0), stats.ReadTimeout)
		assert.Equal(t, time.Duration(0), stats.WriteTimeout)
		assert.Equal(t, int64(0), stats.RequiredAcks)
		assert.Equal(t, false, stats.Async)
		assert.Equal(t, false, downstream.KafkaWriter.AllowAutoTopicCreation)
	})

	// Reset config
	util.Config.Kafka.Topics = []string{}
	util.Config.Kafka.Brokers = []string{}
}

// func TestKafkaProduce(t *testing.T) {
// 	t.Run("debug", func(t *testing.T) {

// 	})
// }

func TestKafkaOverrideOptions(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		config := `
app:
  mode: release

server:
  port: 8000

kafka:
  brokers:
    - foo.bar.com
  topics:
    - foo
  max_attempts: 11
  write_backoff_min: 12
  write_backoff_max: 113
  batch_size: 14
  batch_bytes: 15
  batch_timeout: 16
  read_timeout: 17
  write_timeout: 18
  required_acks: -1
  async: true
  allow_auto_topic_creation: true
`

		err := util.InitConfigFromYaml(config)
		assert.Nil(t, err)

		fmt.Printf("%v\n", util.Config)

		err = downstream.Init()
		assert.Nil(t, err)

		stats := downstream.KafkaWriter.Stats()

		assert.Equal(t, int64(11), stats.MaxAttempts)
		assert.Equal(t, time.Duration(12), stats.WriteBackoffMin)
		assert.Equal(t, time.Duration(113), stats.WriteBackoffMax)
		assert.Equal(t, 14, downstream.KafkaWriter.BatchSize)
		assert.Equal(t, int64(15), downstream.KafkaWriter.BatchBytes)
		assert.Equal(t, time.Duration(16), stats.BatchTimeout)
		assert.Equal(t, time.Duration(17), stats.ReadTimeout)
		assert.Equal(t, time.Duration(18), stats.WriteTimeout)
		assert.Equal(t, int64(-1), stats.RequiredAcks)
		assert.Equal(t, true, stats.Async)
		assert.Equal(t, true, downstream.KafkaWriter.AllowAutoTopicCreation)
	})
}
