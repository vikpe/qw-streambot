package chatbot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

type Chatbot struct {
	client    *twitch.Client
	stopChan  chan os.Signal
	OnStart   func()
	OnConnect func()
	OnStop    func(sig os.Signal)
}

func New(username string, accessToken string, channel string) *Chatbot {
	client := twitch.NewClient(username, fmt.Sprintf("oauth:%s", accessToken))

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println("OnPrivateMessage", message)
	})

	client.Join(channel)

	return &Chatbot{
		client:    client,
		OnStart:   func() {},
		OnConnect: func() {},
		OnStop:    func(sig os.Signal) {},
	}
}

func (c *Chatbot) Start() {
	c.OnStart()

	c.stopChan = make(chan os.Signal, 1)
	signal.Notify(c.stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		c.client.OnConnect(c.OnConnect)
		c.client.Connect()
		defer c.client.Disconnect()
	}()
	sig := <-c.stopChan
	c.OnStop(sig)
}

func (c *Chatbot) Stop() {
	if c.stopChan == nil {
		return
	}
	c.stopChan <- syscall.SIGINT
	time.Sleep(10 * time.Millisecond)
}
