package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/ezquake"
	"github.com/vikpe/streambot/util/twitch"
	"github.com/vikpe/streambot/zeromq"
)

func main() {
	godotenv.Load()

	process := ezquake.NewProcess(os.Getenv("EZQUAKE_BIN_PATH"))
	pipe := ezquake.NewPipeWriter(os.Getenv("EZQUAKE_USERNAME"))
	twitchClient := twitch.NewClient(
		os.Getenv("TWITCH_API_CLIENT_ID"),
		os.Getenv("TWITCH_API_ACCESS_TOKEN"),
		os.Getenv("TWITCH_API_BROADCASTER_ID"),
	)
	publisher := zeromq.NewPublisher(os.Getenv("ZMQ_PUBLISHER_ADDRESS"))
	subscriber := zeromq.NewSubscriber(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"), zeromq.TopicsAll)

	bot := NewStreambot(
		process,
		pipe,
		twitchClient,
		publisher,
		subscriber,
	)
	bot.Start()
}
