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

	playerName := os.Getenv("EZQUAKE_PLAYER_NAME")
	process := ezquake.NewProcess(os.Getenv("EZQUAKE_BIN_PATH"))
	pipe := ezquake.NewPipeWriter(os.Getenv("EZQUAKE_PROCESS_USERNAME"))
	twitchClient := twitch.NewClient(
		os.Getenv("TWITCH_BOT_CLIENT_ID"),
		os.Getenv("TWITCH_CHANNEL_ACCESS_TOKEN"),
		os.Getenv("TWITCH_CHANNEL_BROADCASTER_ID"),
	)
	publisher := zeromq.NewPublisher(os.Getenv("ZMQ_PUBLISHER_ADDRESS"))
	subscriber := zeromq.NewSubscriber(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"), zeromq.TopicsAll)

	bot := NewStreambot(
		playerName,
		process,
		pipe,
		twitchClient,
		publisher,
		subscriber,
	)
	bot.Start()
}
