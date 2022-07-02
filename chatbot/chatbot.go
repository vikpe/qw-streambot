package chatbot

import (
	"fmt"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

type Chatbot struct {
	client  *twitch.Client
	channel string
}

func New(username string, oath string, channel string) *Chatbot {
	client := twitch.NewClient(username, oath)

	client.OnConnect(func() {
		fmt.Println("connected as", username)
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println("OnPrivateMessage", message)
	})

	return &Chatbot{
		client:  client,
		channel: channel,
	}
}

func (c *Chatbot) Start() {
	c.client.Join(c.channel)

	go func() {
		err := c.client.Connect()
		if err != nil {
			fmt.Println("chatbot connect error", err)
			return
		}
	}()
	time.Sleep(time.Second)
}

func (c *Chatbot) Stop() {
	err := c.client.Disconnect()
	if err != nil {
		fmt.Println("chatbot disconnect error", err)
		return
	}
}
