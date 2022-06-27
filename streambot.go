package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/events"
	"github.com/vikpe/streambot/zeromq"
)

type MessageDataHandler func(data zeromq.EventData)

func OnMessage(msg zeromq.Event) {
	handlers := map[string]MessageDataHandler{
		// ezquak events
		events.EzquakeProcessStart:  OnEzquakeProcessStart,
		events.EzquakeProcessStop:   OnEzquakeProcessStop,
		events.EzquakeServerConnect: OnEzquakeServerConnect,

		// server events
		events.ServerMapChange:        OnServerMapChange,
		events.ServerScoreChange:      OnServerScoreChange,
		events.ServerStatusNameChange: OnServerStatusNameChange,
		events.ServerTitleChange:      OnServerTitleChange,

		// streambot events
		events.StreambotHealthCheck:         OnStreambotHealthCheck,
		events.StreambotActionSuggestServer: OnStreambotActionSuggestServer,
	}

	if handler, ok := handlers[msg.Topic]; ok {
		handler(msg.Data)
	} else {
		fmt.Println("no handler defined for", msg.Topic, fmt.Sprintf("%T", msg.Topic), msg.Data)
	}
}

func OnStreambotActionSuggestServer(data zeromq.EventData) {
	var server mvdsv.Mvdsv
	data.To(&server)
	fmt.Println("StreambotActionSuggestServer", server.Address, data)
}

func OnEzquakeProcessStart(data zeromq.EventData) {
	fmt.Println("OnEzquakeProcessStart", data.ToString())
}

func OnEzquakeProcessStop(data zeromq.EventData) {
	fmt.Println("OnEzquakeProcessStop", data.ToString())
}

func OnEzquakeServerConnect(data zeromq.EventData) {
	fmt.Println("OnEzquakeServerConnect", data.ToString())
}

func OnStreambotHealthCheck(data zeromq.EventData) {
	fmt.Println("OnStreambotHealthCheck", data.ToString())
}

func OnServerMapChange(data zeromq.EventData) {
	fmt.Println("OnServerMapChange", data.ToString())
}

func OnServerScoreChange(data zeromq.EventData) {
	fmt.Println("OnServerScoreChange", data.ToInt())
}

func OnServerStatusNameChange(data zeromq.EventData) {
	fmt.Println("OnServerStatusNameChange", data.ToString())
}

func OnServerTitleChange(data zeromq.EventData) {
	fmt.Println("OnServerTitleChange", data.ToString())
}

func main() {
	godotenv.Load()
	wg := sync.WaitGroup{}

	/*wg.Add(1)
	go func() {
		proxy := zeromq.NewProxy(
			os.Getenv("ZMQ_PROXY_FRONTEND_ADDRESS"),
			os.Getenv("ZMQ_PROXY_BACKEND_ADDRESS"),
		)
		proxy.Start()
	}()
	zeromq.WaitForConnection()*/

	subscriber := zeromq.NewSubscriber(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"), zeromq.TopicsAll, OnMessage)
	wg.Add(1)
	go func() {
		subscriber.Start()
	}()
	zeromq.WaitForConnection()

	wg.Wait()

	/*


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

	/*cfg, _ := config.NewFromFile(os.Getenv("STREAMBOT_CONFIG_PATH"))
	fmt.Println("cfg.Mode", cfg.Mode)
	fmt.Println("cfg.MapChangeTimestamp", cfg.MapChangeTimestamp)
	fmt.Println("cfg.SpecAddress", cfg.SpecAddress)
	fmt.Println("cfg.ServerAddress", cfg.ServerAddress, "\n")

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

	wg.Wait()*/
}
