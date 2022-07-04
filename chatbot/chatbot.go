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

	chatbot.OnCommand("auto", func(call command.Command, msg twitch.PrivateMessage) {
		shouldDisable := slices.Contains([]string{"0", "off"}, call.ArgsAsString())

		if shouldDisable {
			cmder.DisableAuto()
		} else {
			cmder.EnableAuto()
		}
	})

	chatbot.OnCommand("autotrack", func(call command.Command, msg twitch.PrivateMessage) {
		cmder.Autotrack()
	})

	chatbot.OnCommand("console", func(call command.Command, msg twitch.PrivateMessage) {
		cmder.Command("toggleconsole")
	})

	chatbot.OnCommand("find", func(call command.Command, msg twitch.PrivateMessage) {
		server, err := qws.FindPlayer(call.ArgsAsString())
		if err != nil {
			chatbot.Reply(msg, err.Error())
			return
		}
		cmder.SuggestServer(server)
	})

	chatbot.OnCommand("lastscores", func(call command.Command, msg twitch.PrivateMessage) {
		cmder.Lastscores()
	})

	chatbot.OnCommand("showscores", func(call command.Command, msg twitch.PrivateMessage) {
		cmder.Showscores()
	})

	chatbot.OnCommand("track", func(call command.Command, msg twitch.PrivateMessage) {
		cmder.Track(call.ArgsAsString())
	})

	return chatbot
}
