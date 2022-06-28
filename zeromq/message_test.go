package zeromq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/zeromq"
)

func TestEventFromMessage(t *testing.T) {
	t.Run("1 arg", func(t *testing.T) {
		expect := zeromq.Message{
			Topic:    "HELLO",
			Data:     "",
			DataType: "string",
		}
		assert.Equal(t, expect, zeromq.NewMessage([]string{"HELLO"}))
	})

	t.Run("2 args", func(t *testing.T) {
		expect := zeromq.Message{
			Topic:    "HELLO",
			Data:     "WORLD",
			DataType: "string",
		}
		assert.Equal(t, expect, zeromq.NewMessage([]string{"HELLO", "WORLD"}))
	})

	t.Run("3 args", func(t *testing.T) {
		expect := zeromq.Message{
			Topic:    "SCORE",
			Data:     "6",
			DataType: "int",
		}
		assert.Equal(t, expect, zeromq.NewMessage([]string{"SCORE", "6", "int"}))
	})

	t.Run("3+ args", func(t *testing.T) {
		expect := zeromq.Message{
			Topic:    "SCORE",
			Data:     "6",
			DataType: "int",
		}
		assert.Equal(t, expect, zeromq.NewMessage([]string{"SCORE", "6", "int", "foo", "bar"}))
	})
}
