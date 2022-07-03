package chatbot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/vikpe/streambot/zeromq"
	"github.com/vikpe/streambot/zeromq/topic"
)

type Chatbot struct {
	client      *twitch.Client
	publisher   zeromq.Publisher
	stopChan    chan os.Signal
	OnStarted   func()
	OnConnected func()
	OnStopped   func(sig os.Signal)
}

func New(username string, accessToken string, channel string, publisherAddress string) *Chatbot {
	client := twitch.NewClient(username, fmt.Sprintf("oauth:%s", accessToken))
	client.Join(channel)

	handler := NewMessageHandler(client, publisherAddress)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.Channel == channel {
			handler.OnPrivateMessage(message)
		}
	})

	return &Chatbot{
		client:      client,
		publisher:   zeromq.NewPublisher(publisherAddress),
		OnStarted:   func() {},
		OnConnected: func() {},
		OnStopped:   func(sig os.Signal) {},
	}
}

func (c *Chatbot) Start() {
	c.OnStarted()

	c.stopChan = make(chan os.Signal, 1)
	signal.Notify(c.stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		c.client.OnConnect(func() {
			c.publisher.SendMessage(topic.ChatbotConnected)
			c.OnConnected()
		})
		c.client.Connect()
		defer c.client.Disconnect()
	}()
	sig := <-c.stopChan

	c.publisher.SendMessage(topic.ChatbotDisconnected)
	c.OnStopped(sig)
}

func (c *Chatbot) Stop() {
	if c.stopChan == nil {
		return
	}
	c.stopChan <- syscall.SIGINT
	time.Sleep(50 * time.Millisecond)
}
