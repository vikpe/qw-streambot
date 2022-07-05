package chatbot

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/streambot/chatbot/irc"
	"github.com/vikpe/streambot/third_party/qws"
	"github.com/vikpe/streambot/util/term"
	"github.com/vikpe/streambot/zeromq"
	"github.com/vikpe/streambot/zeromq/commander"
	"golang.org/x/exp/slices"
)

func New(username string, accessToken string, channel string, publisherAddress string) *irc.Bot {
	pp := term.NewPrettyPrinter("chatbot", color.FgHiBlue)
	cmder := commander.NewCommander(zeromq.NewPublisher(publisherAddress).SendMessage)

	chatbot := irc.NewBot(username, accessToken, channel, '!')

	chatbot.OnConnected = func() {
		pp.Println("connected as", username)
	}

	chatbot.OnStarted = func() {
		pp.Println("start")
	}

	chatbot.OnStopped = func(sig os.Signal) {
		pp.Println(fmt.Sprintf("stop (%s)", sig))
	}

	chatbot.AddCommand("auto", func(cmd irc.Command, msg twitch.PrivateMessage) {
		shouldDisable := slices.Contains([]string{"0", "off"}, cmd.ArgsToString())

		if shouldDisable {
			cmder.DisableAuto()
		} else {
			cmder.EnableAuto()
		}
	})

	chatbot.AddCommand("autotrack", func(cmd irc.Command, msg twitch.PrivateMessage) {
		cmder.Autotrack()
	})

	chatbot.AddCommand("cfg_load", func(cmd irc.Command, msg twitch.PrivateMessage) {
		cmder.Command("cfg_load")
	})

	chatbot.AddCommand("cmd", func(cmd irc.Command, msg twitch.PrivateMessage) {
		if !irc.UserIsBroadcaster(msg.User) {
			chatbot.Reply(msg, "cmd is a mod-only irccommand.")
			return
		}

		cmder.Command(cmd.ArgsToString())
	})

	chatbot.AddCommand("console", func(cmd irc.Command, msg twitch.PrivateMessage) {
		cmder.Command("toggleconsole")
	})

	chatbot.AddCommand("find", func(cmd irc.Command, msg twitch.PrivateMessage) {
		server, err := qws.FindPlayer(cmd.ArgsToString())
		if err != nil {
			chatbot.Reply(msg, err.Error())
			return
		}
		cmder.SuggestServer(server)
	})

	chatbot.AddCommand("lastscores", func(cmd irc.Command, msg twitch.PrivateMessage) {
		cmder.Lastscores()
	})

	chatbot.AddCommand("restart", func(cmd irc.Command, msg twitch.PrivateMessage) {
		cmder.StopEzquake()
	})

	chatbot.AddCommand("showscores", func(cmd irc.Command, msg twitch.PrivateMessage) {
		cmder.Showscores()
	})

	chatbot.AddCommand("track", func(cmd irc.Command, msg twitch.PrivateMessage) {
		cmder.Track(cmd.ArgsToString())
	})

	return chatbot
}
