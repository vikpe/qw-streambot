package ircbot

import "github.com/gempir/go-twitch-irc/v3"

func IsBroadcaster(user twitch.User) bool {
	if broadcasterValue, ok := user.Badges["broadcaster"]; ok {
		return 1 == broadcasterValue
	}

	return false
}
