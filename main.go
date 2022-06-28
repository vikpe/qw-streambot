package main

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	bot := NewStreambot(
		os.Getenv("EZQUAKE_USERNAME"),
		os.Getenv("EZQUAKE_BIN_PATH"),
		os.Getenv("ZMQ_PUBLISHER_ADDRESS"),
		os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"),
	)
	bot.Start()
}
