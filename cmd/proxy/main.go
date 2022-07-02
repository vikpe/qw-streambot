package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/util/term"
	"github.com/vikpe/streambot/zeromq"
)

func main() {
	godotenv.Load("../../.env")

	proxy := zeromq.NewProxy(
		os.Getenv("ZMQ_PROXY_FRONTEND_ADDRESS"),
		os.Getenv("ZMQ_PROXY_BACKEND_ADDRESS"),
	)
	pp := term.NewPrettyPrinter("proxy", color.FgHiGreen)
	proxy.OnStart = func() { pp.Println("start") }
	proxy.OnStop = func(sig os.Signal) { pp.Println(fmt.Sprintf("stop (%s)", sig)) }

	err := proxy.Start()
	if err != nil {
		pp.Println("error", err)
	}
}
