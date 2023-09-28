package twitchbot

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/samber/lo"
	"github.com/vikpe/go-qwhub"
	"github.com/vikpe/prettyfmt"
	hubTwitch "github.com/vikpe/qw-hub-api/pkg/twitch"
	"github.com/vikpe/streambot/internal/app/twitchbot/monitor"
	"github.com/vikpe/streambot/internal/comms/commander"
	"github.com/vikpe/streambot/internal/comms/topic"
	"github.com/vikpe/streambot/internal/pkg/qws"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
	chatbot "github.com/vikpe/twitch-chatbot"
)

func New(botUsername, botAccessToken, channelName, subscriberAddress, publisherAddress string) *chatbot.Chatbot {
	var pfmt = prettyfmt.New("twitchbot", color.FgHiMagenta, "15:04:05", color.FgWhite)

	bot := chatbot.NewChatbot(botUsername, botAccessToken, channelName, '!')

	// zmq messages
	subscriber := zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll)
	subscriber.OnMessage = func(message message.Message) {
		switch message.Topic {
		case topic.TwitchbotSay:
			bot.Say(message.Content.ToString())
		}
	}

	// announce when streamers go live
	streamsMonitor := monitor.NewStreamsMonitor(qwhub.NewClient().Streams, func(stream hubTwitch.Stream) {
		bot.Say(fmt.Sprintf("%s is now streaming @ %s - %s", stream.ClientName, stream.Url, stream.Title))
	})

	// bot events
	bot.OnConnected = func() {
		pfmt.Println("connected as", botUsername)
	}

	bot.OnStarted = func() {
		pfmt.Println("started")
		go subscriber.Start()
		go streamsMonitor.Start(15 * time.Second)
	}

	bot.OnStopped = func(sig os.Signal) {
		subscriber.Stop()
		streamsMonitor.Stop()
		pfmt.Printfln("stopped (%s)", sig)
	}

	// channel commands
	cmder := commander.NewCommander(zeromq.NewPublisher(publisherAddress).SendMessage)

	bot.AddCommand("attack", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Attack()
	})

	bot.AddCommand("auto", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		shouldDisable := lo.Contains([]string{"0", "off"}, cmd.ArgsToString())

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
		cmder.LoadConfig()
	})

	bot.AddCommand("cmd", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		if !chatbot.IsModerator(msg.User) {
			bot.Reply(msg, "cmd is a mod-only command.")
			return
		}

		cmder.Command(cmd.ArgsToString())
	})

	bot.AddCommand("console", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Command("toggleconsole")
	})

	bot.AddCommand("help", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		bot.Reply(msg, "see the channel description for info/commands.")
	})

	bot.AddCommand("commands", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		replyMessage := fmt.Sprintf(`available commands: %s`, bot.GetCommands(", "))
		bot.Reply(msg, replyMessage)
	})

	bot.AddCommand("find", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		server, err := qws.FindPlayer(cmd.ArgsToString())
		if err != nil {
			bot.Reply(msg, err.Error())
			return
		}
		cmder.SuggestServer(server)
	})

	bot.AddCommand("jump", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Jump()
	})

	bot.AddCommand("lastscores", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Lastscores()
	})

	bot.AddCommand("restart", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.StopEzquake()
		time.AfterFunc(1250*time.Millisecond, func() {
			cmder.StopQuakeManager()
		})
	})

	bot.AddCommand("say", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Commandf("bot_say %s: %s", msg.User.DisplayName, cmd.ArgsToString())
	})

	bot.AddCommand("showscores", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Showscores()
	})

	bot.AddCommand("track", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Track(cmd.ArgsToString())
	})

	return bot
}
