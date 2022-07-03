package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/chatbot"
	"github.com/vikpe/streambot/util/term"
)

func main() {
	godotenv.Load("../../.env")

	bot := chatbot.New(
		os.Getenv("TWITCH_BOT_USERNAME"),
		os.Getenv("TWITCH_BOT_ACCESS_TOKEN"),
		os.Getenv("TWITCH_CHANNEL_USERNAME"),
		os.Getenv("ZMQ_PUBLISHER_ADDRESS"),
	)
	pp := term.NewPrettyPrinter("chatbot", color.FgHiBlue)
	bot.OnConnected = func() { pp.Println("connected as", os.Getenv("TWITCH_BOT_USERNAME")) }
	bot.OnStarted = func() { pp.Println("start") }
	bot.OnStopped = func(sig os.Signal) { pp.Println(fmt.Sprintf("stop (%s)", sig)) }
	bot.Start()
}
