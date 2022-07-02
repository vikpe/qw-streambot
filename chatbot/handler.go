package chatbot

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v3"
	"golang.org/x/exp/slices"
)

func IsModerator(name string) bool {
	mods := []string{"vikpe", "circle1", "hangtime_of_qw", "vikpebot", "wimpeeh"}
	return slices.Contains(mods, name)
}

type MessageHandler struct {
	client *twitch.Client
}

func NewMessageHandler(client *twitch.Client) MessageHandler {
	return MessageHandler{client: client}
}

func (c *MessageHandler) OnPrivateMessage(message twitch.PrivateMessage) {
	fmt.Println(message.Channel, message.Message)

	if IsCommand(message.Message) {
		c.OnCommand(NewCommand(message))
	}
}

func (c *MessageHandler) OnCommand(cmd Command) {
	handlers := map[string]CommandHandler{
		// commands
		CommandFind:        c.OnCommandFind,
		CommandEnableAuto:  c.OnCommandEnableAuto,
		CommandDisableAuto: c.OnCommandDisableAuto,
	}

	if handler, ok := handlers[cmd.Name]; ok {
		handler(cmd)
	} else {
		fmt.Println("no handler defined for", cmd.Name)
	}
}

func (c *MessageHandler) OnCommandFind(cmd Command) {

}

func (c *MessageHandler) OnCommandEnableAuto(cmd Command) {

}

func (c *MessageHandler) OnCommandDisableAuto(cmd Command) {

}
