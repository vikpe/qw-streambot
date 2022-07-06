package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/pkg/prettyfmt"
	"github.com/vikpe/streambot/pkg/zeromq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("unable to load environment variables", err)
		return
	}

	proxy := zeromq.NewProxy(
		os.Getenv("ZMQ_PROXY_FRONTEND_ADDRESS"),
		os.Getenv("ZMQ_PROXY_BACKEND_ADDRESS"),
	)
	pfmt := prettyfmt.New("proxy", color.FgHiGreen)
	proxy.OnStart = func() { pfmt.Println("start") }
	proxy.OnStop = func(sig os.Signal) { pfmt.Printfln("stop (%s)", sig) }
	proxy.OnError = func(err error) { pfmt.Println("error", err) }
	proxy.Start()
}
