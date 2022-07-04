package command_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/irc/bot/command"
)

func TestIsCommand(t *testing.T) {
	testCases := map[string]bool{
		"":       false,
		" ":      false,
		"find":   false,
		"#find":  false,
		"!":      false,
		"! ":     false,
		"! !":    false,
		"!.":     false,
		"!.find": false,
		"!!":     false,

		"!find":       true,
		" !find":      true,
		" !find ":     true,
		"  !  find  ": true,
		"! find":      true,
		"! find ":     true,
	}

	const prefix = '!'

	for text, expect := range testCases {
		t.Run(text, func(t *testing.T) {
			assert.Equal(t, expect, command.IsCommand(prefix, text))
		})
	}
}

func BenchmarkIsCommand(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		command.IsCommand('!', "!find xantom")
	}
}

func TestNewCommandFromText(t *testing.T) {
	t.Run("invalid prefix", func(t *testing.T) {
		cmd, err := command.NewFromText('!', "#bar")
		assert.Equal(t, cmd, command.Command{})
		assert.EqualError(t, err, "unable to parse command call")
	})

	t.Run("invalid command", func(t *testing.T) {
		cmd, err := command.NewFromText('#', "##")
		assert.Equal(t, cmd, command.Command{})
		assert.EqualError(t, err, "unable to parse command call")
	})

	t.Run("valid", func(t *testing.T) {
		testCases := map[string]command.Command{
			"!find":       {Name: "find", Args: []string{}},
			" !find arg1": {Name: "find", Args: []string{"arg1"}},
			"!find a b c": {Name: "find", Args: []string{"a", "b", "c"}},
		}

		const prefix = '!'

		for text, expect := range testCases {
			t.Run(text, func(t *testing.T) {
				foo, err := command.NewFromText(prefix, text)
				assert.Equal(t, expect, foo)
				assert.Nil(t, err)
			})
		}
	})
}

func BenchmarkNewCommandFromText(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		command.NewFromText('!', "!find xantom")
	}
}

func TestArgsAsString(t *testing.T) {
	assert.Equal(t, "foo bar", command.New("find", "foo", "bar").ArgsAsString())
}
