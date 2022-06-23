package main

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/zeromq"
)

func main() {
	godotenv.Load()
	/*
		godotenv.Load()
		wg := sync.WaitGroup{}

		wg.Add(1)
		go func() {
			proxy := zeromq.NewProxy(
				os.Getenv("ZMQ_PROXY_FRONTEND_ADDRESS"),
				os.Getenv("ZMQ_PROXY_BACKEND_ADDRESS"),
			)
			proxy.Start()
		}()
		subscriber := zeromq.NewSubscriber(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"), "")
		time.Sleep(time.Millisecond * 20)

		wg.Add(1)
		go func() {
			subscriber.Start()
		}()
		time.Sleep(time.Millisecond * 20)

		ticker := time.NewTicker(time.Duration(5) * time.Second)
		//process := ezquake.NewProcess("/home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64")

		wg.Add(1)

		publisher := zeromq.NewPublisher(os.Getenv("ZMQ_PUBLISHER_ADDRESS"))

		go func() {
			for ; true; <-ticker.C {
				address := process.TcpAddress()
				fmt.Println("process.TcpAddress", address)
				info, _ := serverstat.GetInfo("troopers.fi:28001")
				fmt.Println(info)
				publisher.SendMessage("hello", "world")
			}
		}()

		wg.Wait()
	*/

	wg := sync.WaitGroup{}
	subscriber := zeromq.NewSubscriber(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"), "")

	wg.Add(1)
	go func() {
		subscriber.Start()
	}()

	publiser := zeromq.NewPublisher(os.Getenv("ZMQ_PUBLISHER_ADDRESS"))
	publiser.SendMessage("hehe", "")

	wg.Wait()
}
