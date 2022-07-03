package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/zeromq/message"
)

func TestSerialization(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		var valueBefore = "hello"
		var valueAfter string
		message.Unserialize(message.Serialize(valueBefore), &valueAfter)
		assert.Equal(t, valueBefore, valueAfter)
	})

	t.Run("int", func(t *testing.T) {
		var valueBefore = 5
		var valueAfter int
		message.Unserialize(message.Serialize(valueBefore), &valueAfter)
		assert.Equal(t, valueBefore, valueAfter)
	})

	t.Run("float64", func(t *testing.T) {
		var valueBefore = 5.5
		var valueAfter float64
		message.Unserialize(message.Serialize(valueBefore), &valueAfter)
		assert.Equal(t, valueBefore, valueAfter)
	})

	t.Run("[]string", func(t *testing.T) {
		var valueBefore = []string{"a", "b"}
		var valueAfter []string
		message.Unserialize(message.Serialize(valueBefore), &valueAfter)
		assert.Equal(t, valueBefore, valueAfter)
	})

	t.Run("map", func(t *testing.T) {
		var valueBefore = map[string]int{"foo": 2}
		var valueAfter map[string]int
		message.Unserialize(message.Serialize(valueBefore), &valueAfter)
		assert.Equal(t, valueBefore, valueAfter)
	})
}
