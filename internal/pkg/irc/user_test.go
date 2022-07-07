package irc_test

import (
	"testing"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/internal/pkg/irc"
)

func TestIsBroadcaster(t *testing.T) {
	t.Run("undefined value", func(t *testing.T) {
		user := twitch.User{Badges: map[string]int{}}
		assert.False(t, irc.IsBroadcaster(user))
	})

	t.Run("is not a broadcaster", func(t *testing.T) {
		user := twitch.User{Badges: map[string]int{"broadcaster": 0}}
		assert.False(t, irc.IsBroadcaster(user))
	})

	t.Run("is a broadcaster", func(t *testing.T) {
		user := twitch.User{Badges: map[string]int{"broadcaster": 1}}
		assert.True(t, irc.IsBroadcaster(user))
	})
}
