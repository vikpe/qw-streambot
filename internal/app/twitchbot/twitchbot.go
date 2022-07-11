package twitchbot

import (
	"os"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/prettyfmt"
	"github.com/vikpe/streambot/internal/com/commander"
	"github.com/vikpe/streambot/internal/com/topic"
	"github.com/vikpe/streambot/internal/pkg/qws"
	zeromq2 "github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
	chatbot "github.com/vikpe/twitch-chatbot"
	"golang.org/x/exp/slices"
)

type Twitchbot struct {
	*chatbot.Chatbot
	subscriber *zeromq2.Subscriber
}

func New(botUsername, botAccessToken, channelName, subscriberAddress, publisherAddress string) *Twitchbot {
	var pfmt = prettyfmt.New("twitchbot", color.FgHiMagenta, "15:04:05", color.FgWhite)

	bot := Twitchbot{
		Chatbot:    chatbot.NewChatbot(botUsername, botAccessToken, channelName, '!'),
		subscriber: zeromq2.NewSubscriber(subscriberAddress, zeromq2.TopicsAll),
	}

	// zmq messages
	onZmqMessage := func(message message.Message) {
		switch message.Topic {
		case topic.TwitchbotSay:
			bot.Say(message.Content.ToString())
		}
	}

	// bot events
	bot.OnConnected = func() {
		pfmt.Println("connected as", botUsername)
		bot.subscriber.Start(onZmqMessage)
	}

	bot.OnStarted = func() {
		pfmt.Println("started")
	}

	bot.OnStopped = func(sig os.Signal) {
		pfmt.Printfln("stopped (%s)", sig)
	}

	// channel commands
	cmder := commander.NewCommander(zeromq2.NewPublisher(publisherAddress).SendMessage)

	bot.AddCommand("auto", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		shouldDisable := slices.Contains([]string{"0", "off"}, cmd.ArgsToString())

		if shouldDisable {
			cmder.DisableAuto()
		} else {
			cmder.EnableAuto()
		}
	})

	bot.AddCommand("autotrack", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Autotrack()
	})

	bot.AddCommand("cfg_load", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Command("cfg_load")
	})

	bot.AddCommand("cmd", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		if !chatbot.IsBroadcaster(msg.User) {
			bot.Reply(msg, "cmd is a mod-only chatbot.")
			return
		}

		cmder.Command(cmd.ArgsToString())
	})

	bot.AddCommand("console", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Command("toggleconsole")
	})

	bot.AddCommand("find", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		server, err := qws.FindPlayer(cmd.ArgsToString())
		if err != nil {
			bot.Reply(msg, err.Error())
			return
		}
		cmder.SuggestServer(server)
	})

	bot.AddCommand("lastscores", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Lastscores()
	})

	bot.AddCommand("restart", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.StopEzquake()
	})

	bot.AddCommand("showscores", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Showscores()
	})

	bot.AddCommand("track", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Track(cmd.ArgsToString())
	})

	return &bot
}
