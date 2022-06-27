package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/vikpe/streambot/config"
	"github.com/vikpe/streambot/events"
	"github.com/vikpe/streambot/ezquake"
	"github.com/vikpe/streambot/task"
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
	/*
		wg := sync.WaitGroup{}
		subscriber := zeromq.NewSubscriber(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"), "")

		wg.Add(1)
		go func() {
			subscriber.Start()
		}()

		publiser := zeromq.NewPublisher(os.Getenv("ZMQ_PUBLISHER_ADDRESS"))
		publiser.SendMessage("hehe", "")*/

	cfg, _ := config.NewFromFile(os.Getenv("STREAMBOT_CONFIG_PATH"))
	fmt.Println("cfg.Mode", cfg.Mode)
	fmt.Println("cfg.MapChangeTimestamp", cfg.MapChangeTimestamp)
	fmt.Println("cfg.SpecAddress", cfg.SpecAddress)
	fmt.Println("cfg.ServerAddress", cfg.ServerAddress, "\n")

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		proc := ezquake.NewProcess(os.Getenv("EZQUAKE_BIN_PATH"))
		eventHandler := func(topic string, data any) {
			fmt.Println("got event", topic, data)
		}

		pmon := task.NewProcessMonitor(&proc, eventHandler)
		pmon.Start(4 * time.Second)

		foo := func() { fmt.Println(events.StreambotHealthCheck) }
		hth := task.NewPeriodicalTask(foo)
		hth.Start(10 * time.Second)
	}()

	wg.Wait()
}
