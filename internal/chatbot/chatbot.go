package chatbot

import (
	"os"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/prettyfmt"
	"github.com/vikpe/streambot/internal/pkg/irc"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/third_party/qws"
	"github.com/vikpe/streambot/pkg/commander"
	"golang.org/x/exp/slices"
)

func New(username string, accessToken string, channel string, publisherAddress string) *irc.Bot {
	var pfmt = prettyfmt.New("chatbot", color.FgHiBlue, "15:04:05", color.FgWhite)
	cmder := commander.NewCommander(zeromq.NewPublisher(publisherAddress).SendMessage)

	chatbot := irc.NewBot(username, accessToken, channel, '!')

	chatbot.OnConnected = func() {
		pfmt.Println("connected as", username)
	}

	chatbot.OnStarted = func() {
		pfmt.Println("start")
	}

	chatbot.OnStopped = func(sig os.Signal) {
		pfmt.Printfln("stop (%s)", sig)
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
		if !irc.IsBroadcaster(msg.User) {
			chatbot.Reply(msg, "cmd is a mod-only command.")
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
