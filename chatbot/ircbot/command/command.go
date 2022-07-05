package command

import (
	"errors"
	"strings"
	"unicode"

	"github.com/gempir/go-twitch-irc/v3"
)

type Command struct {
	Name string
	Args []string
}

func New(name string, args ...string) Command {
	return Command{
		Name: name,
		Args: args,
	}
}

func (c Command) ArgsAsString() string {
	return strings.Join(c.Args, " ")
}

type Handler func(cmd Command, msg twitch.PrivateMessage)

func IsCommand(prefix rune, text string) bool {
	txt := strings.TrimLeft(text, " ")

	if len(txt) < 2 || 0 != strings.IndexRune(txt, prefix) {
		return false
	}

	parts := strings.FieldsFunc(txt[1:], unicode.IsSpace)

	if 0 == len(parts) {
		return false
	}

	firstRune := rune(parts[0][0])
	return unicode.IsLetter(firstRune)
}

func NewFromText(prefix rune, text string) (Command, error) {
	if !IsCommand(prefix, text) {
		return Command{}, errors.New("unable to parse command call")
	}
	txt := strings.TrimLeft(text, " ")
	txt = strings.ToLower(txt[1:])
	parts := strings.FieldsFunc(txt, unicode.IsSpace)
	name := parts[0]
	args := parts[1:]
	return New(name, args...), nil
}
