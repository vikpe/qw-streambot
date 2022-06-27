package zeromq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/zeromq"
)

func TestParseEvent(t *testing.T) {
	t.Run("1 arg", func(t *testing.T) {
		expect := zeromq.Event{
			Topic:    "HELLO",
			Data:     "",
			DataType: "string",
		}
		assert.Equal(t, expect, zeromq.ParseEvent([]string{"HELLO"}))
	})

	t.Run("2 args", func(t *testing.T) {
		expect := zeromq.Event{
			Topic:    "HELLO",
			Data:     "WORLD",
			DataType: "string",
		}
		assert.Equal(t, expect, zeromq.ParseEvent([]string{"HELLO", "WORLD"}))
	})

	t.Run("3 args", func(t *testing.T) {
		expect := zeromq.Event{
			Topic:    "SCORE",
			Data:     "6",
			DataType: "int",
		}
		assert.Equal(t, expect, zeromq.ParseEvent([]string{"SCORE", "6", "int"}))
	})

	t.Run("3+ args", func(t *testing.T) {
		expect := zeromq.Event{
			Topic:    "SCORE",
			Data:     "6",
			DataType: "int",
		}
		assert.Equal(t, expect, zeromq.ParseEvent([]string{"SCORE", "6", "int", "foo", "bar"}))
	})
}
