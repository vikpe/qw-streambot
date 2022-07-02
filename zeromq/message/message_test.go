package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/zeromq/message"
)

func TestNewFromMultipart(t *testing.T) {
	t.Run("1 arg", func(t *testing.T) {
		expect := message.Message{
			Topic:    "HELLO",
			Data:     "",
			DataType: "string",
		}
		assert.Equal(t, expect, message.NewFromMultipart([]string{"HELLO"}))
	})

	t.Run("2 args", func(t *testing.T) {
		expect := message.Message{
			Topic:    "HELLO",
			Data:     "WORLD",
			DataType: "string",
		}
		assert.Equal(t, expect, message.NewFromMultipart([]string{"HELLO", "WORLD"}))
	})

	t.Run("3 args", func(t *testing.T) {
		expect := message.Message{
			Topic:    "SCORE",
			Data:     "6",
			DataType: "int",
		}
		assert.Equal(t, expect, message.NewFromMultipart([]string{"SCORE", "6", "int"}))
	})

	t.Run("3+ args", func(t *testing.T) {
		expect := message.Message{
			Topic:    "SCORE",
			Data:     "6",
			DataType: "int",
		}
		assert.Equal(t, expect, message.NewFromMultipart([]string{"SCORE", "6", "int", "foo", "bar"}))
	})
}
