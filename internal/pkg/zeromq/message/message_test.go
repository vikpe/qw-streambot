package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	message2 "github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

func TestNewMessageFromParts(t *testing.T) {
	t.Run("0 args", func(t *testing.T) {
		expect := message2.Message{}
		msg, err := message2.NewMessageFromFrames([]string{})
		assert.Equal(t, expect, msg)
		assert.EqualError(t, err, "expected 3 message frames, got 0")
	})

	t.Run("1 arg", func(t *testing.T) {
		expect := message2.Message{}
		msg, err := message2.NewMessageFromFrames([]string{"HELLO"})
		assert.Equal(t, expect, msg)
		assert.EqualError(t, err, "expected 3 message frames, got 1")
	})

	t.Run("2 args", func(t *testing.T) {
		expect := message2.Message{}
		msg, err := message2.NewMessageFromFrames([]string{"HELLO", "WORLD"})
		assert.Equal(t, expect, msg)
		assert.EqualError(t, err, "expected 3 message frames, got 2")
	})

	t.Run("3 args", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			serializedContent := message2.Serialize("WORLD")
			expect := message2.Message{
				Topic:       "HELLO",
				Content:     serializedContent,
				ContentType: "string",
			}
			msg, err := message2.NewMessageFromFrames([]string{"HELLO", "string", string(serializedContent)})
			assert.Equal(t, expect, msg)
			assert.Nil(t, err)
			assert.Equal(t, "WORLD", msg.Content.ToString())
		})

		t.Run("int", func(t *testing.T) {
			serializedContent := message2.Serialize(33)
			expect := message2.Message{
				Topic:       "HELLO",
				Content:     serializedContent,
				ContentType: "int",
			}
			msg, err := message2.NewMessageFromFrames([]string{"HELLO", "int", string(serializedContent)})
			assert.Equal(t, expect, msg)
			assert.Nil(t, err)
			assert.Equal(t, 33, msg.Content.ToInt())
		})

		t.Run("[]string", func(t *testing.T) {
			serializedContent := message2.Serialize([]string{"a", "b"})
			expect := message2.Message{
				Topic:       "HELLO",
				Content:     serializedContent,
				ContentType: "int",
			}
			msg, err := message2.NewMessageFromFrames([]string{"HELLO", "int", string(serializedContent)})

			assert.Equal(t, expect, msg)
			assert.Nil(t, err)

			var unserializedContent []string
			msg.Content.To(&unserializedContent)
			assert.Equal(t, []string{"a", "b"}, unserializedContent)
		})
	})

	t.Run("3+ args", func(t *testing.T) {
		expect := message2.Message{}
		msg, err := message2.NewMessageFromFrames([]string{"a", "b", "c", "d"})
		assert.Equal(t, expect, msg)
		assert.EqualError(t, err, "expected 3 message frames, got 4")
	})
}

func TestNewMessage(t *testing.T) {
	expect := message2.Message{
		Topic:       "foo",
		ContentType: "string",
		Content:     message2.Serialize("bar"),
	}
	assert.Equal(t, expect, message2.NewMessage("foo", "bar"))
}
