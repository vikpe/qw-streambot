package chatbot

import (
	"fmt"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/streambot/third_party/qws"
	"github.com/vikpe/streambot/zeromq"
	"github.com/vikpe/streambot/zeromq/topic"
)

/*func IsModerator(name string) bool {
	mods := []string{"vikpe", "circle1", "hangtime_of_qw", "vikpebot", "wimpeeh"}
	return slices.Contains(mods, name)
}*/

type MessageHandler struct {
	client    *twitch.Client
	publisher zeromq.Publisher
}

func NewMessageHandler(client *twitch.Client, publisherAddress string) *MessageHandler {
	return &MessageHandler{
		client:    client,
		publisher: zeromq.NewPublisher(publisherAddress),
	}
}

func (c *MessageHandler) OnPrivateMessage(message twitch.PrivateMessage) {
	fmt.Println(message.Channel, message.Message)

	if IsCommand(message.Message) {
		c.OnCommand(NewCommand(message))
	}
}

func (c *MessageHandler) OnCommand(cmd Command) {
	handlers := map[string]CommandHandler{
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
	fmt.Println("find", cmd.Args)

	playerName := strings.TrimSpace(strings.Join(cmd.Args, " "))
	const minFindLength = 2

	if len(playerName) < minFindLength {
		c.ReplyToPrivateMessage(cmd.Message, fmt.Sprintf(`Provide at least %d characters.`, minFindLength))
	}

	servers := qws.GetMvdsvServersByQueryParams(map[string]string{
		"has_player": playerName,
	})

	if len(servers) > 0 {
		c.publisher.SendMessage(topic.StreambotSuggestServer, servers[0])
		fmt.Println("found player", servers[0].Address)

	} else {
		c.ReplyToPrivateMessage(cmd.Message, fmt.Sprintf(`"%s" not found.`, playerName))
	}
}

func (c *MessageHandler) ReplyToPrivateMessage(msg twitch.PrivateMessage, replyText string) {
	c.client.Reply(msg.Channel, msg.ID, replyText)
}

func (c *MessageHandler) OnCommandEnableAuto(cmd Command) {

}

func (c *MessageHandler) OnCommandDisableAuto(cmd Command) {

}
