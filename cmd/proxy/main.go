package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/vikpe/prettyfmt"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("unable to load environment variables", err)
		return
	}

	proxy := zeromq.NewProxyService(
		os.Getenv("ZMQ_PROXY_FRONTEND_ADDRESS"),
		os.Getenv("ZMQ_PROXY_BACKEND_ADDRESS"),
	)
	pfmt := prettyfmt.New("proxy", color.FgHiGreen, "15:04:05", color.FgWhite)
	proxy.OnStarted = func() { pfmt.Println("start") }
	proxy.OnStopped = func(sig os.Signal) { pfmt.Printfln("stop (%s)", sig) }
	proxy.OnError = func(err error) { pfmt.Println("error", err) }
	proxy.Start()
}
