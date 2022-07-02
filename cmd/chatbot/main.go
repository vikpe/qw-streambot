package main

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/chatbot"
)

func main() {
	godotenv.Load("../../.env")

	bot := chatbot.New(
		os.Getenv("TWITCH_BOT_USERNAME"),
		os.Getenv("TWITCH_BOT_CHAT_ACCESS_TOKEN"),
		os.Getenv("TWITCH_CHANNEL_NAME"),
	)
	bot.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
