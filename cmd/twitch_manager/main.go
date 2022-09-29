package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/vikpe/prettyfmt"
	"github.com/vikpe/streambot/internal/app/twitch_manager"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("unable to load environment variables", err)
		return
	}

	manager, err := twitch_manager.New(
		os.Getenv("TWITCH_MANAGER_APP_CLIENT_ID"),
		os.Getenv("TWITCH_MANAGER_APP_ACCESS_TOKEN"),
		os.Getenv("TWITCH_CHANNEL_BROADCASTER_ID"),
		os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"),
	)

	if err != nil {
		fmt.Println("unable to create twitch manager:", err)
		return
	}

	var pfmt = prettyfmt.New("twitch_manager", color.FgHiRed, "15:04:05", color.FgWhite)
	manager.OnStarted = func() { pfmt.Println("started") }
	manager.OnStopped = func(signal os.Signal) { pfmt.Println("stopped", signal) }
	manager.OnError = func(err error) { pfmt.Println("error", err) }
	manager.Start()
}
