package chatbot

import (
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
)

const (
	CommandPrefix      = "!"
	CommandFind        = "find"
	CommandDisableAuto = "manual"
	CommandEnableAuto  = "auto"
)

type Command struct {
	Name    string
	Args    []string
	Message twitch.PrivateMessage
}

type CommandHandler func(Command)

func IsCommand(messageText string) bool {
	return strings.HasPrefix(strings.TrimSpace(messageText), CommandPrefix)
}

func NewCommand(msg twitch.PrivateMessage) Command {
	parts := strings.SplitN(strings.ToLower(msg.Message[1:]), " ", 2)
	name := parts[0]
	args := parts[1:]

	return Command{
		Name:    name,
		Args:    args,
		Message: msg,
	}
}
