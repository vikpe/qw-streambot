package chatbot

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/streambot/irc/bot"
	"github.com/vikpe/streambot/irc/bot/command"
	"github.com/vikpe/streambot/third_party/qws"
	"github.com/vikpe/streambot/util/term"
	"github.com/vikpe/streambot/zeromq"
	"github.com/vikpe/streambot/zeromq/commander"
	"golang.org/x/exp/slices"
)

func New(username string, accessToken string, channel string, publisherAddress string) *bot.Bot {
	pp := term.NewPrettyPrinter("chatbot", color.FgHiBlue)
	cmder := commander.NewCommander(zeromq.NewPublisher(publisherAddress).SendMessage)

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

	chatbot.OnCommand("find", func(call command.Command, msg twitch.PrivateMessage) {
		pp.Println("find player", call.Args)
		server, err := qws.FindPlayer(call.ArgsAsString())
		if err != nil {
			chatbot.Reply(msg, err.Error())
		}
		cmder.SuggestServer(server)
	})

	chatbot.OnCommand("auto", func(call command.Command, msg twitch.PrivateMessage) {
		pp.Println("auto", call.Args)

		if slices.Contains([]string{"0", "off"}, call.ArgsAsString()) {
			cmder.DisableAuto()
		} else {
			cmder.EnableAuto()
		}
	})

	chatbot.OnCommand("track", func(call command.Command, msg twitch.PrivateMessage) {
		pp.Println("track", call.Args)
		cmder.Track(call.ArgsAsString())
	})

	chatbot.OnCommand("autotrack", func(call command.Command, msg twitch.PrivateMessage) {
		pp.Println("autotrack", call.Args)
		cmder.Autotrack()
	})

	return chatbot
}
