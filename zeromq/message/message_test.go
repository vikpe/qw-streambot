package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/zeromq/message"
)

func TestNewMessageFromParts(t *testing.T) {
	t.Run("0 args", func(t *testing.T) {
		expect := message.Message{}
		msg, err := message.NewMessageFromFrames([]string{})
		assert.Equal(t, expect, msg)
		assert.EqualError(t, err, "expected 3 message frames, got 0")
	})

	t.Run("1 arg", func(t *testing.T) {
		expect := message.Message{}
		msg, err := message.NewMessageFromFrames([]string{"HELLO"})
		assert.Equal(t, expect, msg)
		assert.EqualError(t, err, "expected 3 message frames, got 1")
	})

	t.Run("2 args", func(t *testing.T) {
		expect := message.Message{}
		msg, err := message.NewMessageFromFrames([]string{"HELLO", "WORLD"})
		assert.Equal(t, expect, msg)
		assert.EqualError(t, err, "expected 3 message frames, got 2")
	})

	t.Run("3 args", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			expect := message.Message{
				Topic:       "HELLO",
				Content:     message.SerializedValue(`"WORLD"`),
				ContentType: "string",
			}
			msg, err := message.NewMessageFromFrames([]string{"HELLO", "string", `"WORLD"`})
			assert.Equal(t, expect, msg)
			assert.Nil(t, err)
			assert.Equal(t, "WORLD", msg.Content.ToString())
		})

		t.Run("int", func(t *testing.T) {
			expect := message.Message{
				Topic:       "HELLO",
				Content:     message.SerializedValue("33"),
				ContentType: "int",
			}
			msg, err := message.NewMessageFromFrames([]string{"HELLO", "int", "33"})
			assert.Equal(t, expect, msg)
			assert.Nil(t, err)
			assert.Equal(t, 33, msg.Content.ToInt())
		})

		t.Run("slice of strings", func(t *testing.T) {
			expect := message.Message{
				Topic:       "HELLO",
				Content:     message.SerializedValue(`["a","b"]`),
				ContentType: "int",
			}
			msg, err := message.NewMessageFromFrames([]string{"HELLO", "int", `["a","b"]`})

			assert.Equal(t, expect, msg)
			assert.Nil(t, err)

			var unserializedContent []string
			msg.Content.To(&unserializedContent)
			assert.Equal(t, []string{"a", "b"}, unserializedContent)
		})
	})

	t.Run("3+ args", func(t *testing.T) {
		expect := message.Message{}
		msg, err := message.NewMessageFromFrames([]string{"a", "b", "c", "d"})
		assert.Equal(t, expect, msg)
		assert.EqualError(t, err, "expected 3 message frames, got 4")
	})
}

func TestSerializedValue_ToString(t *testing.T) {
	assert.Equal(t, "abc", message.NewSerializedValue("abc").ToString())
}

func TestSerializedValue_ToInt(t *testing.T) {
	assert.Equal(t, 123, message.NewSerializedValue(123).ToInt())
}

func TestSerializedValue_To(t *testing.T) {
	valueBefore := []string{"a", "b", "c"}
	var valueAfter []string
	message.NewSerializedValue(valueBefore).To(&valueAfter)
	assert.Equal(t, valueBefore, valueAfter)
}
