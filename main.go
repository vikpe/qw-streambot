package main

import (
	"os"

	"github.com/joho/godotenv"
	ezquake2 "github.com/vikpe/streambot/pkg/ezquake"
	zeromq2 "github.com/vikpe/streambot/pkg/zeromq"
	"github.com/vikpe/streambot/third_party/twitch"
)

func main() {
	godotenv.Load()

	playerName := os.Getenv("EZQUAKE_PLAYER_NAME")
	process := ezquake2.NewProcess(os.Getenv("EZQUAKE_BIN_PATH"))
	pipe := ezquake2.NewPipeWriter(os.Getenv("EZQUAKE_PROCESS_USERNAME"))
	twitchClient := twitch.NewClient(
		os.Getenv("TWITCH_BOT_CLIENT_ID"),
		os.Getenv("TWITCH_CHANNEL_TITLE_ACCESS_TOKEN"),
		os.Getenv("TWITCH_CHANNEL_BROADCASTER_ID"),
	)
	publisher := zeromq2.NewPublisher(os.Getenv("ZMQ_PUBLISHER_ADDRESS"))
	subscriber := zeromq2.NewSubscriber(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"), zeromq2.TopicsAll)

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
