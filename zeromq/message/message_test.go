package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/zeromq/message"
)

func TestNewMessageFromParts(t *testing.T) {
	t.Run("0 args", func(t *testing.T) {
		expect := message.Message{}
		msg, err := message.NewMessageFromParts([]string{})
		assert.Equal(t, expect, msg)
		assert.ErrorContains(t, err, "expected 1-3 message frames, got 0")
	})

	t.Run("1 arg", func(t *testing.T) {
		expect := message.Message{
			Topic:       "HELLO",
			Content:     `""`,
			ContentType: "string",
		}
		msg, err := message.NewMessageFromParts([]string{"HELLO"})
		assert.Equal(t, expect, msg)
		assert.Nil(t, err)
	})

	t.Run("2 args", func(t *testing.T) {
		expect := message.Message{
			Topic:       "HELLO",
			Content:     `"WORLD"`,
			ContentType: "string",
		}
		msg, err := message.NewMessageFromParts([]string{"HELLO", "WORLD"})
		assert.Equal(t, expect, msg)
		assert.Nil(t, err)
	})

	t.Run("3 args", func(t *testing.T) {

		t.Run("string", func(t *testing.T) {
			expect := message.Message{
				Topic:       "HELLO",
				Content:     `"WORLD"`,
				ContentType: "string",
			}
			msg, err := message.NewMessageFromParts([]string{"HELLO", "string", "WORLD"})
			assert.Equal(t, expect, msg)
			assert.Nil(t, err)
		})

		t.Run("int", func(t *testing.T) {
			expect := message.Message{
				Topic:       "HELLO",
				Content:     `"3"`,
				ContentType: "int",
			}
			msg, err := message.NewMessageFromParts([]string{"HELLO", "int", "3"})
			assert.Equal(t, expect, msg)
			assert.Nil(t, err)
		})
	})

	t.Run("3+ args", func(t *testing.T) {
		expect := message.Message{}
		msg, err := message.NewMessageFromParts([]string{"a", "b", "c", "d"})
		assert.Equal(t, expect, msg)
		assert.ErrorContains(t, err, "expected 1-3 message frames, got 4")
	})
}
