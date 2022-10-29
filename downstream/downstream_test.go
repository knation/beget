package downstream_test

import (
	"beget/downstream"
	"beget/util"
	"testing"

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
	})

	// Reset config
	util.Config.Kafka.Topics = []string{}
	util.Config.Kafka.Brokers = []string{}
}

func TestKafkaProduce(t *testing.T) {
	t.Run("debug", func(t *testing.T) {

	})
}
