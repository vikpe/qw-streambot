package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/internal/app/quakebot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("unable to load environment variables", err)
		return
	}

	brn := quakebot.New(
		os.Getenv("EZQUAKE_PLAYER_NAME"),
		os.Getenv("EZQUAKE_BIN_PATH"),
		os.Getenv("EZQUAKE_PROCESS_USERNAME"),
		os.Getenv("ZMQ_PUBLISHER_ADDRESS"),
		os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"),
	)
	brn.Start()
}
