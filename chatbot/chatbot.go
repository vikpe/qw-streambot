package chatbot

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/streambot/chatbot/ircbot"
	"github.com/vikpe/streambot/chatbot/ircbot/command"
	"github.com/vikpe/streambot/third_party/qws"
	"github.com/vikpe/streambot/util/term"
	"github.com/vikpe/streambot/zeromq"
	"github.com/vikpe/streambot/zeromq/commander"
	"golang.org/x/exp/slices"
)

func New(username string, accessToken string, channel string, publisherAddress string) *ircbot.Bot {
	pp := term.NewPrettyPrinter("chatbot", color.FgHiBlue)
	cmder := commander.NewCommander(zeromq.NewPublisher(publisherAddress).SendMessage)

	chatbot := ircbot.New(username, accessToken, channel, '!')

	chatbot.OnConnected = func() {
		pp.Println("connected as", username)
	}

	chatbot.OnStarted = func() {
		pp.Println("start")
	}

	chatbot.OnStopped = func(sig os.Signal) {
		pp.Println(fmt.Sprintf("stop (%s)", sig))
	}

	chatbot.AddCommand("auto", func(cmd command.Command, msg twitch.PrivateMessage) {
		shouldDisable := slices.Contains([]string{"0", "off"}, cmd.ArgsAsString())

		if shouldDisable {
			cmder.DisableAuto()
		} else {
			cmder.EnableAuto()
		}
	})

	chatbot.AddCommand("autotrack", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Autotrack()
	})

	chatbot.AddCommand("cfg_load", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Command("cfg_load")
	})

	chatbot.AddCommand("cmd", func(cmd command.Command, msg twitch.PrivateMessage) {
		if !ircbot.IsBroadcaster(msg.User) {
			chatbot.Reply(msg, "cmd is a mod-only command.")
			return
		}

		cmder.Command(cmd.ArgsAsString())
	})

	chatbot.AddCommand("console", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Command("toggleconsole")
	})

	chatbot.AddCommand("find", func(cmd command.Command, msg twitch.PrivateMessage) {
		server, err := qws.FindPlayer(cmd.ArgsAsString())
		if err != nil {
			chatbot.Reply(msg, err.Error())
			return
		}
		cmder.SuggestServer(server)
	})

	chatbot.AddCommand("lastscores", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Lastscores()
	})

	chatbot.AddCommand("restart", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.StopEzquake()
	})

	chatbot.AddCommand("showscores", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Showscores()
	})

	chatbot.AddCommand("track", func(cmd command.Command, msg twitch.PrivateMessage) {
		cmder.Track(cmd.ArgsAsString())
	})

	return chatbot
}
