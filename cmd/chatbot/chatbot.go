package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/joho/godotenv"
)

type Chatbot struct {
	client  *twitch.Client
	channel string
}

func NewChatbot(username string, oath string, channel string) Chatbot {
	client := twitch.NewClient(username, oath)

	client.OnConnect(func() {
		fmt.Println("connected as ", username)
	})

	client.OnRoomStateMessage(func(message twitch.RoomStateMessage) {
		fmt.Println("RoomStateMessage", message)
	})

	client.OnUserStateMessage(func(message twitch.UserStateMessage) {
		fmt.Println("UserStateMessage", message)
	})

	client.OnWhisperMessage(func(message twitch.WhisperMessage) {
		fmt.Println("WhisperMessage", message)
	})

	client.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		fmt.Println("UserNoticeMessage", message)
	})

	client.OnUnsetMessage(func(message twitch.RawMessage) {
		fmt.Println("OnUnsetMessage", message)
	})
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println("OnPrivateMessage", message)
	})

	return Chatbot{
		client:  client,
		channel: channel,
	}
}

func (c Chatbot) Start() {
	err := c.client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.client.Join(c.channel)
}

func (c Chatbot) Stop() {
	err := c.client.Disconnect()
	if err != nil {
		return
	}
}

func main() {
	godotenv.Load("../../.env")

	bot := NewChatbot(
		os.Getenv("TWITCH_CHATBOT_USERNAME"),
		os.Getenv("TWITCH_CHATBOT_OATH"),
		os.Getenv("TWITCH_CHATBOT_CHANNEL"),
	)
	bot.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
