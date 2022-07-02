package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/zeromq"
)

func main() {
	godotenv.Load("../../.env")

	zeromq.NewProxy(
		os.Getenv("ZMQ_PROXY_FRONTEND_ADDRESS"),
		os.Getenv("ZMQ_PROXY_BACKEND_ADDRESS"),
	).Start()
}
