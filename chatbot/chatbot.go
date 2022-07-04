package chatbot

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/streambot/irc/bot"
	"github.com/vikpe/streambot/third_party/qws"
	"github.com/vikpe/streambot/util/term"
	"github.com/vikpe/streambot/zeromq"
	"github.com/vikpe/streambot/zeromq/topic"
)

const (
	CommandFind        = "find"
	CommandEnableAuto  = "auto"
	CommandDisableAuto = "manual"
)

func New(username string, accessToken string, channel string, publisherAddress string) *bot.Bot {
	publisher := zeromq.NewPublisher(publisherAddress)
	pp := term.NewPrettyPrinter("chatbot", color.FgHiBlue)

	chatbot := bot.New(username, accessToken, channel, '#')
	chatbot.OnConnected = func() {
		pp.Println("connected as", username)
	}
	chatbot.OnStarted = func() {
		pp.Println("start")
	}
	chatbot.OnStopped = func(sig os.Signal) {
		pp.Println(fmt.Sprintf("stop (%s)", sig))
	}

	chatbot.OnCommand(CommandFind, func(call bot.CommandCall, msg twitch.PrivateMessage) {
		pp.Println("find", call.Args)

		playerName := strings.TrimSpace(strings.Join(call.Args, " "))
		const minFindLength = 2

		if len(playerName) < minFindLength {
			chatbot.Reply(msg, fmt.Sprintf(`Provide at least %d characters.`, minFindLength))
		}

		servers := qws.GetMvdsvServersByQueryParams(map[string]string{
			"has_player": playerName,
		})

		if len(servers) > 0 {
			publisher.SendMessage(topic.StreambotSuggestServer, servers[0])
			fmt.Println("found player", servers[0].Address)

		} else {
			chatbot.Reply(msg, fmt.Sprintf(`"%s" not found.`, playerName))
		}
	})

	chatbot.OnCommand(CommandEnableAuto, func(call bot.CommandCall, msg twitch.PrivateMessage) {
		pp.Println("enable auto", call.Args)
	})

	chatbot.OnCommand(CommandDisableAuto, func(call bot.CommandCall, msg twitch.PrivateMessage) {
		pp.Println("disable auto", call.Args)
	})

	return chatbot
}
