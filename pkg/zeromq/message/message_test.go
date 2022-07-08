package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/pkg/zeromq/message"
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
			serializedContent := message.Serialize("WORLD")
			expect := message.Message{
				Topic:       "HELLO",
				Content:     serializedContent,
				ContentType: "string",
			}
			msg, err := message.NewMessageFromFrames([]string{"HELLO", "string", string(serializedContent)})
			assert.Equal(t, expect, msg)
			assert.Nil(t, err)
			assert.Equal(t, "WORLD", msg.Content.ToString())
		})

		t.Run("int", func(t *testing.T) {
			serializedContent := message.Serialize(33)
			expect := message.Message{
				Topic:       "HELLO",
				Content:     serializedContent,
				ContentType: "int",
			}
			msg, err := message.NewMessageFromFrames([]string{"HELLO", "int", string(serializedContent)})
			assert.Equal(t, expect, msg)
			assert.Nil(t, err)
			assert.Equal(t, 33, msg.Content.ToInt())
		})

		t.Run("[]string", func(t *testing.T) {
			serializedContent := message.Serialize([]string{"a", "b"})
			expect := message.Message{
				Topic:       "HELLO",
				Content:     serializedContent,
				ContentType: "int",
			}
			msg, err := message.NewMessageFromFrames([]string{"HELLO", "int", string(serializedContent)})

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
