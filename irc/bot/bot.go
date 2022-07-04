package bot

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/gempir/go-twitch-irc/v3"
)

type Bot struct {
	client          *twitch.Client
	channel         string
	channelCommands map[string]CommandCallHandler
	stopChan        chan os.Signal
	CommandPrefix   rune
	OnStarted       func()
	OnConnected     func()
	OnStopped       func(os.Signal)
}

func New(username string, oauth string, channel string, commandPrefix rune) *Bot {
	client := twitch.NewClient(username, oauth)
	client.Join(channel)

	bot := Bot{
		client:          client,
		channel:         channel,
		channelCommands: make(map[string]CommandCallHandler, 0),
		CommandPrefix:   commandPrefix,
		OnStarted:       func() {},
		OnConnected:     func() {},
		OnStopped:       func(os.Signal) {},
	}

	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		if msg.Channel != channel {
			return
		}

		cmdCall, err := ParseCommandCall(msg.Message)

		if err != nil {
			return
		}

		if handler, ok := bot.channelCommands[cmdCall.Name]; ok {
			handler(cmdCall, msg)
		}
	})

	return &bot
}

func (b *Bot) OnCommand(name string, handler CommandCallHandler) {
	b.channelCommands[name] = handler
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

const (
	CommandPrefix = "#"
)

type CommandCall struct {
	Name string
	Args []string
}

type CommandCallHandler func(CommandCall, twitch.PrivateMessage)

func IsCommand(text string) bool {
	txt := strings.TrimLeft(text, " ")

	if !strings.HasPrefix(txt, CommandPrefix) {
		return false
	}

	parts := strings.FieldsFunc(txt[1:], unicode.IsSpace)

	if 0 == len(parts) {
		return false
	}

	firstRune := rune(parts[0][0])
	return unicode.IsLetter(firstRune) || unicode.IsDigit(firstRune)
}

func ParseCommandCall(text string) (CommandCall, error) {
	if !IsCommand(text) {
		return CommandCall{}, errors.New("unable to parse command")
	}

	s := strings.TrimLeft(text, " ")[1:]
	parts := strings.FieldsFunc(strings.ToLower(s), unicode.IsSpace)

	return CommandCall{
		Name: parts[0],
		Args: parts[1:],
	}, nil
}
