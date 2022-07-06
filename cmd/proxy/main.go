package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/pkg/prettyprint"
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
	pp := prettyprint.New("proxy", color.FgHiGreen)
	proxy.OnStart = func() { pp.Println("start") }
	proxy.OnStop = func(sig os.Signal) { pp.Println(fmt.Sprintf("stop (%s)", sig)) }
	proxy.OnError = func(err error) { pp.Println("error", err) }
	proxy.Start()
}
