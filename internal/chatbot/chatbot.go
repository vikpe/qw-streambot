package chatbot

import (
	"os"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/prettyfmt"
	"github.com/vikpe/streambot/internal/pkg/irc"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
	"github.com/vikpe/streambot/internal/third_party/qws"
	"github.com/vikpe/streambot/pkg/commander"
	"github.com/vikpe/streambot/pkg/topic"
	"golang.org/x/exp/slices"
)

type Chatbot struct {
	*irc.Bot
	subscriber zeromq.Subscriber
}

func New(username, accessToken, channel, subscriberAddress, publisherAddress string) *Chatbot {
	var pfmt = prettyfmt.New("chatbot", color.FgHiMagenta, "15:04:05", color.FgWhite)

	chatbot := Chatbot{
		Bot:        irc.NewBot(username, accessToken, channel, '!'),
		subscriber: zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll),
	}

	// zmq messages
	onZmqMessage := func(message message.Message) {
		switch message.Topic {
		case topic.ChatbotSay:
			chatbot.Say(message.Content.ToString())
		}
	}

	// bot events
	chatbot.OnConnected = func() {
		pfmt.Println("connected as", username)
		chatbot.subscriber.Start(onZmqMessage)
	}

	chatbot.OnStarted = func() {
		pfmt.Println("start")
	}

	chatbot.OnStopped = func(sig os.Signal) {
		pfmt.Printfln("stop (%s)", sig)
	}

	// channel commands
	cmder := commander.NewCommander(zeromq.NewPublisher(publisherAddress).SendMessage)

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

	return &chatbot
}
