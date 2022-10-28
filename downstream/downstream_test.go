package downstream_test

import (
	"beget/downstream"
	"beget/util"
	"os"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestInitDebug(t *testing.T) {

	t.Run("no topics", func(t *testing.T) {
		err := downstream.Init(util.DebugMode)

		assert.EqualError(t, err, `no topics specified`)
	})

	t.Run("empty topics", func(t *testing.T) {
		os.Setenv("KAFKA_TOPICS", "")
		err := downstream.Init(util.DebugMode)

		assert.EqualError(t, err, `no topics specified`)
	})

	t.Run("success in debug mode", func(t *testing.T) {
		os.Setenv("KAFKA_TOPICS", "foo")
		err := downstream.Init(util.DebugMode)

		assert.Nil(t, err)
		assert.EqualValues(t, map[string]bool{"foo": true}, downstream.KafkaTopics)

		// Make sure close doesn't break
		err = downstream.Close()
		assert.Nil(t, err)

		os.Setenv("KAFKA_TOPICS", "")
	})

	t.Run("no brokers", func(t *testing.T) {
		os.Setenv("KAFKA_TOPICS", "foo")

		err := downstream.Init(util.ReleaseMode)

		assert.EqualError(t, err, `must provide either "KAFKA_BROKERS"`)
	})

	t.Run("empty brokers", func(t *testing.T) {
		os.Setenv("KAFKA_TOPICS", "foo")
		os.Setenv("KAFKA_BROKERS", "")

		err := downstream.Init(util.ReleaseMode)

		assert.EqualError(t, err, `must provide either "KAFKA_BROKERS"`)

		os.Setenv("KAFKA_TOPICS", "")
	})

	t.Run("single broker", func(t *testing.T) {
		os.Setenv("KAFKA_TOPICS", "foo")
		os.Setenv("KAFKA_BROKERS", "broker.foo.com")

		err := downstream.Init(util.ReleaseMode)

		assert.Nil(t, err)
		assert.NotNil(t, downstream.KafkaWriter)
		assert.EqualValues(t, kafka.TCP("broker.foo.com"), downstream.KafkaWriter.Addr)

		// Make sure close doesn't break
		err = downstream.Close()
		assert.Nil(t, err)

		os.Setenv("KAFKA_TOPICS", "")
		os.Setenv("KAFKA_BROKERS", "")
	})

	t.Run("multiple brokers", func(t *testing.T) {
		os.Setenv("KAFKA_TOPICS", "foo")
		os.Setenv("KAFKA_BROKERS", "broker.foo.com,broker.bar.com")

		err := downstream.Init(util.ReleaseMode)

		assert.Nil(t, err)
		assert.NotNil(t, downstream.KafkaWriter)
		assert.EqualValues(t, kafka.TCP("broker.foo.com", "broker.bar.com"), downstream.KafkaWriter.Addr)

		// Make sure close doesn't break
		err = downstream.Close()
		assert.Nil(t, err)

		os.Setenv("KAFKA_TOPICS", "")
		os.Setenv("KAFKA_BROKERS", "")
	})
}

func TestKafkaProduce(t *testing.T) {
	t.Run("debug", func(t *testing.T) {

	})
}
