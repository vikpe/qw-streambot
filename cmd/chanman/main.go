package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
	"github.com/vikpe/streambot/internal/third_party/twitch"
	"github.com/vikpe/streambot/pkg/topic"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("unable to load environment variables", err)
		return
	}

	twitchClient := twitch.NewClient(
		os.Getenv("TWITCH_BOT_CLIENT_ID"),
		os.Getenv("TWITCH_CHANNEL_TITLE_ACCESS_TOKEN"),
		os.Getenv("TWITCH_CHANNEL_BROADCASTER_ID"),
	)

	subscriber := zeromq.NewSubscriber(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"), topic.ServerTitleChanged)
	subscriber.Start(func(msg message.Message) {
		switch msg.Topic {
		case topic.ServerTitleChanged:
			twitchClient.SetTitle(msg.Content.ToString())
		}
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
