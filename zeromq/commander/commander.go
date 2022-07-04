// Package commander sends commands over zeromq
package commander

import (
	"fmt"

	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/zeromq/topic"
)

type Commander struct {
	sendMessage func(string, ...any)
}

func NewCommander(sendMessage func(string, ...any)) Commander {
	return Commander{
		sendMessage: sendMessage,
	}
}

func (c Commander) Autotrack() {
	c.sendMessage(topic.EzquakeCommand, "bot_track")
}

func (c Commander) EnableAuto() {
	c.sendMessage(topic.StreambotEnableAuto)
}

func (c Commander) DisableAuto() {
	c.sendMessage(topic.StreambotDisableAuto)
}

func (c Commander) SuggestServer(server mvdsv.Mvdsv) {
	c.sendMessage(topic.StreambotSuggestServer, server)
}

func (c Commander) StopEzquake() {
	c.sendMessage(topic.StopEzquake)
}

func (c Commander) Track(target string) {
	c.sendMessage(topic.EzquakeCommand, fmt.Sprintf("bot_track %s", target))
}
