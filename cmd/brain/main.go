package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/internal/brain"
	"github.com/vikpe/streambot/internal/brain/util/proc"
	"github.com/vikpe/streambot/internal/pkg/ezquake"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/third_party/twitch"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("unable to load environment variables", err)
		return
	}

	playerName := os.Getenv("EZQUAKE_PLAYER_NAME")
	process := proc.NewProcessController(os.Getenv("EZQUAKE_BIN_PATH"))
	pipe := ezquake.NewPipeWriter(os.Getenv("EZQUAKE_PROCESS_USERNAME"))
	twitchClient := twitch.NewClient(
		os.Getenv("TWITCH_BOT_CLIENT_ID"),
		os.Getenv("TWITCH_CHANNEL_TITLE_ACCESS_TOKEN"),
		os.Getenv("TWITCH_CHANNEL_BROADCASTER_ID"),
	)
	publisher := zeromq.NewPublisher(os.Getenv("ZMQ_PUBLISHER_ADDRESS"))
	subscriber := zeromq.NewSubscriber(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"), zeromq.TopicsAll)

	brn := brain.NewBrain(
		playerName,
		process,
		pipe,
		twitchClient,
		publisher,
		subscriber,
	)
	brn.Start()
}
