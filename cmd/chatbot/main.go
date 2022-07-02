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
		os.Getenv("TWITCH_BOT_CHAT_ACCESS_TOKEN"),
		os.Getenv("TWITCH_CHANNEL_USERNAME"),
	)
	pp := term.NewPrettyPrinter("chatbot", color.FgHiBlue)
	bot.OnConnect = func() { pp.Println("connected as", os.Getenv("TWITCH_BOT_USERNAME")) }
	bot.OnStart = func() { pp.Println("start") }
	bot.OnStop = func(sig os.Signal) { pp.Println(fmt.Sprintf("stop (%s)", sig)) }
	bot.Start()
}
