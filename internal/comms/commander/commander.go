// Package commander sends command messages
package commander

import (
	"fmt"

	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/internal/comms/topic"
)

type Commander struct {
	sendMessage func(string, ...any)
}

func NewCommander(sendMessage func(topic string, data ...any)) Commander {
	return Commander{
		sendMessage: sendMessage,
	}
}

func (c Commander) Autotrack() {
	c.sendMessage(topic.EzquakeCommand, "bot_track")
}

func (c Commander) TwitchbotSay(text string) {
	c.sendMessage(topic.TwitchbotSay, text)
}

func (c Commander) Command(cmd string) {
	c.sendMessage(topic.EzquakeCommand, cmd)
}

func (c Commander) Commandf(format string, args ...any) {
	c.Command(fmt.Sprintf(format, args...))
}

func (c Commander) EnableAuto() {
	c.sendMessage(topic.StreambotEnableAuto)
}

func (c Commander) Evaluate() {
	c.sendMessage(topic.StreambotEvaluate)
}

func (c Commander) DisableAuto() {
	c.sendMessage(topic.StreambotDisableAuto)
}

func (c Commander) Lastscores() {
	c.sendMessage(topic.EzquakeScript, "lastscores")
}

func (c Commander) Showscores() {
	c.sendMessage(topic.EzquakeScript, "showscores")
}

func (c Commander) SuggestServer(server mvdsv.Mvdsv) {
	c.sendMessage(topic.StreambotSuggestServer, server)
}

func (c Commander) StopEzquake() {
	c.sendMessage(topic.EzquakeStop)
}

func (c Commander) Track(target string) {
	c.sendMessage(topic.EzquakeCommand, fmt.Sprintf("bot_track %s", target))
}
