package chatbot

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/streambot/chatbot/irc/bot"
	"github.com/vikpe/streambot/chatbot/irc/bot/command"
	"github.com/vikpe/streambot/third_party/qws"
	"github.com/vikpe/streambot/util/term"
	"github.com/vikpe/streambot/zeromq"
	"github.com/vikpe/streambot/zeromq/commander"
	"golang.org/x/exp/slices"
)

func IsBroadcaster(user twitch.User) bool {
	if broadcasterValue, ok := user.Badges["broadcaster"]; ok {
		return 1 == broadcasterValue
	}

	return false
}

func New(username string, accessToken string, channel string, publisherAddress string) *bot.Bot {
	pp := term.NewPrettyPrinter("chatbot", color.FgHiBlue)
	cmder := commander.NewCommander(zeromq.NewPublisher(publisherAddress).SendMessage)

	chatbot := bot.New(username, accessToken, channel, '!')

	chatbot.OnConnected = func() {
		pp.Println("connected as", username)
	}

	chatbot.OnStarted = func() {
		pp.Println("start")
	}

	chatbot.OnStopped = func(sig os.Signal) {
		pp.Println(fmt.Sprintf("stop (%s)", sig))
	}

	chatbot.OnCommand("auto", func(cmd command.Command, msg twitch.PrivateMessage) {
		shouldDisable := slices.Contains([]string{"0", "off"}, cmd.ArgsAsString())

		if shouldDisable {
			cmder.DisableAuto()
		} else {
			cmder.EnableAuto()
		}
	})

	chatbot.OnCommand("autotrack", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Autotrack()
	})

	chatbot.OnCommand("cfg_load", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Command("cfg_load")
	})

	chatbot.OnCommand("cmd", func(cmd command.Command, msg twitch.PrivateMessage) {
		if !IsBroadcaster(msg.User) {
			chatbot.Reply(msg, "cmd is a mod-only command.")
			return
		}

		cmder.Command(cmd.ArgsAsString())
	})

	chatbot.OnCommand("console", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Command("toggleconsole")
	})

	chatbot.OnCommand("find", func(cmd command.Command, msg twitch.PrivateMessage) {
		server, err := qws.FindPlayer(cmd.ArgsAsString())
		if err != nil {
			chatbot.Reply(msg, err.Error())
			return
		}
		cmder.SuggestServer(server)
	})

	chatbot.OnCommand("lastscores", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Lastscores()
	})

	chatbot.OnCommand("restart", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.StopEzquake()
	})

	chatbot.OnCommand("showscores", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Showscores()
	})

	chatbot.OnCommand("track", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Track(cmd.ArgsAsString())
	})

	return chatbot
}
