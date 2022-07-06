package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/internal/chatbot"
)

func main() {
	godotenv.Load("../../.env")

	bot := chatbot.New(
		os.Getenv("TWITCH_BOT_USERNAME"),
		os.Getenv("TWITCH_BOT_CHAT_ACCESS_TOKEN"),
		os.Getenv("TWITCH_CHANNEL_USERNAME"),
		os.Getenv("ZMQ_PUBLISHER_ADDRESS"),
	)
	bot.Start()
}
