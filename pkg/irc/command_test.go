package irc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/pkg/irc"
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
		"!123":   false,

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
			assert.Equal(t, expect, irc.IsCommand(prefix, text))
		})
	}
}

func BenchmarkIsCommand(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		irc.IsCommand('!', "!find xantom")
	}
}

func TestNewCommandFromText(t *testing.T) {
	t.Run("invalid prefix", func(t *testing.T) {
		cmd, err := irc.NewCommandFromText('!', "#bar")
		assert.Equal(t, cmd, irc.Command{})
		assert.EqualError(t, err, "unable to parse irccommand call")
	})

	t.Run("invalid irccommand", func(t *testing.T) {
		cmd, err := irc.NewCommandFromText('#', "##")
		assert.Equal(t, cmd, irc.Command{})
		assert.EqualError(t, err, "unable to parse irccommand call")
	})

	t.Run("valid", func(t *testing.T) {
		testCases := map[string]irc.Command{
			"!find":       {Name: "find", Args: []string{}},
			" !find arg1": {Name: "find", Args: []string{"arg1"}},
			"!find a b c": {Name: "find", Args: []string{"a", "b", "c"}},
		}

		const prefix = '!'

		for text, expect := range testCases {
			t.Run(text, func(t *testing.T) {
				foo, err := irc.NewCommandFromText(prefix, text)
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
		irc.NewCommandFromText('!', "!find xantom")
	}
}

func TestArgsToString(t *testing.T) {
	assert.Equal(t, "foo bar", irc.NewCommand("find", "foo", "bar").ArgsToString())
}
