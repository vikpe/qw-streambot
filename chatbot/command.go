package chatbot

import (
	"errors"
	"strings"
	"unicode"

	"github.com/gempir/go-twitch-irc/v3"
)

const (
	CommandPrefix      = "#"
	CommandFind        = "find"
	CommandDisableAuto = "manual"
	CommandEnableAuto  = "auto"
)

type Command struct {
	Name    string
	Args    []string
	Message twitch.PrivateMessage
}

type CommandHandler func(Command)

func IsCommand(text string) bool {
	txt := strings.TrimLeft(text, " ")

	if !strings.HasPrefix(txt, "#") {
		return false
	}

	parts := strings.FieldsFunc(txt[1:], unicode.IsSpace)

	if 0 == len(parts) {
		return false
	}

	firstRune := rune(parts[0][0])
	return unicode.IsLetter(firstRune) || unicode.IsDigit(firstRune)
}

type Foo struct {
	Name string
	Args []string
}

func Parse(text string) (Foo, error) {
	if !IsCommand(text) {
		return Foo{}, errors.New("unable to parse command")
	}

	s := strings.TrimLeft(text, " ")[1:]
	parts := strings.FieldsFunc(strings.ToLower(s), unicode.IsSpace)

	return Foo{
		Name: parts[0],
		Args: parts[1:],
	}, nil
}

func NewCommand(msg twitch.PrivateMessage) Command {
	parts := strings.SplitN(strings.ToLower(msg.Message[1:]), " ", 2)
	name := parts[0]
	args := parts[1:]

	return Command{
		Name:    name,
		Args:    args,
		Message: msg,
	}
}
