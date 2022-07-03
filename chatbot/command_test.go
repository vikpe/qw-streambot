package chatbot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/chatbot"
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
			assert.Equal(t, expect, chatbot.IsCommand(text))
		})
	}
}

func TestParse(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		command, err := chatbot.Parse("")
		assert.Equal(t, command, chatbot.Foo{})
		assert.EqualError(t, err, "unable to parse command")
	})

	t.Run("valid", func(t *testing.T) {
		testCases := map[string]chatbot.Foo{
			"#find":         {Name: "find", Args: []string{}},
			" #find XantoM": {Name: "find", Args: []string{"xantom"}},
		}

		for text, expect := range testCases {
			t.Run(text, func(t *testing.T) {
				foo, err := chatbot.Parse(text)
				assert.Equal(t, expect, foo)
				assert.Nil(t, err)
			})
		}
	})
}
