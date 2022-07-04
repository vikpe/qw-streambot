package bot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/irc/bot"
)

func TestIsCommand(t *testing.T) {
	testCases := map[string]bool{
		"":       false,
		" ":      false,
		"find":   false,
		"#":      false,
		"# ":     false,
		"# #":    false,
		"#.":     false,
		"#.find": false,
		"##":     false,

		"#find":       true,
		" #find":      true,
		" #find ":     true,
		"  #  find  ": true,
		"# find":      true,
		"# find ":     true,
	}

	for text, expect := range testCases {
		t.Run(text, func(t *testing.T) {
			assert.Equal(t, expect, bot.IsCommand(text))
		})
	}
}

func TestParseCommandCall(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		command, err := bot.ParseCommandCall("")
		assert.Equal(t, command, bot.CommandCall{})
		assert.EqualError(t, err, "unable to parse command")
	})

	t.Run("valid", func(t *testing.T) {
		testCases := map[string]bot.CommandCall{
			"#find":         {Name: "find", Args: []string{}},
			" #find XantoM": {Name: "find", Args: []string{"xantom"}},
		}

		for text, expect := range testCases {
			t.Run(text, func(t *testing.T) {
				foo, err := bot.ParseCommandCall(text)
				assert.Equal(t, expect, foo)
				assert.Nil(t, err)
			})
		}
	})
}
