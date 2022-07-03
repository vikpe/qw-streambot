package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/zeromq/message"
)

func TestSerialization(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		var value = "hello"
		serialized := message.Serialize(value)

		var unserialized string
		message.Unserialize(serialized, &unserialized)
		assert.Equal(t, value, unserialized)
	})

	t.Run("int", func(t *testing.T) {
		var value = 5
		serialized := message.Serialize(value)

		var unserialized int
		message.Unserialize(serialized, &unserialized)
		assert.Equal(t, value, unserialized)
	})

	t.Run("float64", func(t *testing.T) {
		var value = 5.5
		serialized := message.Serialize(value)

		var unserialized float64
		message.Unserialize(serialized, &unserialized)
		assert.Equal(t, value, unserialized)
	})

	t.Run("[]string", func(t *testing.T) {
		var value = []string{"a", "b"}
		serialized := message.Serialize(value)

		var unserialized []string
		message.Unserialize(serialized, &unserialized)
		assert.Equal(t, value, unserialized)
	})

	t.Run("map", func(t *testing.T) {
		var value = map[string]int{"foo": 2}
		serialized := message.Serialize(value)

		var unserialized map[string]int
		message.Unserialize(serialized, &unserialized)
		assert.Equal(t, value, unserialized)
	})
}
