package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/ezquake"
	"github.com/vikpe/streambot/util/twitch"
	"github.com/vikpe/streambot/zeromq"
)

func main() {
	godotenv.Load()

	proxy := zeromq.NewProxy(
		os.Getenv("ZMQ_PROXY_FRONTEND_ADDRESS"),
		os.Getenv("ZMQ_PROXY_BACKEND_ADDRESS"),
	)

	go func() {
		err := proxy.Start()
		if err != nil {
			fmt.Println("PROXY START FAIL", err.Error())
			return
		}
	}()
	zeromq.WaitForConnection()

	process := ezquake.NewProcess(os.Getenv("EZQUAKE_BIN_PATH"))
	pipe := ezquake.NewPipeWriter(os.Getenv("EZQUAKE_USERNAME"))
	twitchClient := twitch.NewClient(
		os.Getenv("TWITCH_API_CLIENT_ID"),
		os.Getenv("TWITCH_API_ACCESS_TOKEN"),
		os.Getenv("TWITCH_API_BROADCASTER_ID"),
	)
	publisher := zeromq.NewPublisher(os.Getenv("ZMQ_PUBLISHER_ADDRESS"))

	bot := NewStreambot(
		process,
		pipe,
		twitchClient,
		publisher,
		os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"),
	)
	bot.Start()
}
