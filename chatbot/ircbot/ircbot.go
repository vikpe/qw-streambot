package ircbot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/streambot/chatbot/ircbot/command"
)

type Bot struct {
	client          *twitch.Client
	channel         string
	commandHandlers map[string]CommandHandler
	stopChan        chan os.Signal
	OnStarted       func()
	OnConnected     func()
	OnStopped       func(os.Signal)
}

type CommandHandler func(cmd command.Command, msg twitch.PrivateMessage)

func New(username string, oauth string, channel string, commandPrefix rune) *Bot {
	client := twitch.NewClient(username, oauth)
	client.Join(channel)

	bot := Bot{
		client:          client,
		channel:         channel,
		commandHandlers: make(map[string]CommandHandler, 0),
		OnStarted:       func() {},
		OnConnected:     func() {},
		OnStopped:       func(os.Signal) {},
	}

	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		if msg.Channel != channel {
			return
		}

		cmd, err := command.NewFromText(commandPrefix, msg.Message)

		if err != nil {
			return
		}

		if cmdHandler, ok := bot.commandHandlers[cmd.Name]; ok {
			cmdHandler(cmd, msg)
		} else {
			bot.Reply(msg, fmt.Sprintf(`unknown command "%s".`, cmd.Name))
		}
	})

	return &bot
}

func (b *Bot) AddCommand(name string, handler CommandHandler) {
	b.commandHandlers[name] = handler
}

func (b *Bot) Reply(msg twitch.PrivateMessage, replyText string) {
	b.client.Reply(msg.Channel, msg.ID, replyText)
}

func (b *Bot) Say(text string) {
	b.client.Say(b.channel, text)
}

func (b *Bot) Start() {
	b.OnStarted()

	b.stopChan = make(chan os.Signal, 1)
	signal.Notify(b.stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		b.client.OnConnect(func() {
			b.OnConnected()
		})
		b.client.Connect()
		defer b.client.Disconnect()
	}()
	sig := <-b.stopChan

	b.OnStopped(sig)
}

func (b *Bot) Stop() {
	if b.stopChan == nil {
		return
	}
	b.stopChan <- syscall.SIGINT
	time.Sleep(50 * time.Millisecond)
}
