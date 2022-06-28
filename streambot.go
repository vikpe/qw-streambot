package main

import (
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/ezquake"
	"github.com/vikpe/streambot/qws"
	"github.com/vikpe/streambot/topics"
	"github.com/vikpe/streambot/zeromq"
)

type Streambot struct {
	pipe       ezquake.PipeWriter
	process    ezquake.Process
	publisher  zeromq.Publisher
	subscriber zeromq.Subscriber
}

func NewStreambot(
	ezquakeUsername string,
	ezquakePath string,
	publisherAddress string,
	subscriberAddress string,
) Streambot {
	bot := Streambot{
		pipe:      ezquake.NewPipeWriter(ezquakeUsername),
		process:   ezquake.NewProcess(ezquakePath),
		publisher: zeromq.NewPublisher(publisherAddress),
	}
	bot.subscriber = zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll, bot.OnMessage)

	return bot
}

func (s Streambot) Start() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		s.subscriber.Start()
	}()

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(4 * time.Second)
		for ; true; <-ticker.C {
			bestServer, err := qws.GetBestServer()

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(bestServer.Address, bestServer.Score)
		}
	}()

	wg.Wait()
}

func (s Streambot) OnMessage(msg zeromq.Message) {
	handlers := map[string]zeromq.MessageDataHandler{
		// client
		topics.ClientConnect:      s.OnClientConnect,
		topics.ClientCommand:      s.OnClientCommand,
		topics.ClientStarted:      s.OnClientStarted,
		topics.ClientStopped:      s.OnClientStopped,
		topics.ClientConnected:    s.OnClientConnected,
		topics.ClientDisconnected: s.OnClientDisconnected,
		topics.StopClient:         s.OnStopClient,

		// server
		topics.ServerMapChanged:    s.OnServerMapChanged,
		topics.ServerScoreChanged:  s.OnServerScoreChanged,
		topics.ServerStatusChanged: s.OnServerStatusChanged,
		topics.ServerTitleChanged:  s.OnServerTitleChanged,

		// system
		topics.SystemHealthCheck: s.OnSystemHealthCheck,
		topics.SystemUpdate:      s.OnSystemUpdate,
	}

	if handler, ok := handlers[msg.Topic]; ok {
		handler(msg.Data)
	} else {
		fmt.Println("no handler defined for", msg.Topic, fmt.Sprintf("%T", msg.Topic), msg.Data)
	}
}

func (s Streambot) Evaluate() {
	fmt.Println("Evaluate")
}

func (s Streambot) OnStart() {
	fmt.Println("OnStart")
	s.Evaluate()
}

func (s Streambot) OnClientConnect(data zeromq.MessageData) {
	var server mvdsv.Mvdsv
	data.To(&server)

	fmt.Println("OnClientConnect", server.Address, data)

	time.AfterFunc(4*time.Second, func() {
		s.publisher.SendMessage(topics.ClientCommand, "bot_track")
	})
}

func (s Streambot) OnClientCommand(data zeromq.MessageData) {
	fmt.Println("OnClientCommand", data.ToString())
	s.pipe.Write(data.ToString())
}

func (s Streambot) OnClientStarted(data zeromq.MessageData) {
	fmt.Println("OnClientStarted", data.ToString())

	time.AfterFunc(4*time.Second, func() {
		s.publisher.SendMessage(topics.ClientCommand, "toggleconsole")
	})
}

func (s Streambot) OnStopClient(data zeromq.MessageData) {
	fmt.Println("OnStopClient", data.ToString())
	s.process.Stop(syscall.SIGTERM)
}

func (s Streambot) OnClientStopped(data zeromq.MessageData) {
	fmt.Println("OnClientStopped", data.ToString())
}

func (s Streambot) OnClientConnected(data zeromq.MessageData) {
	fmt.Println("OnClientConnected", data.ToString())
}

func (s Streambot) OnClientDisconnected(data zeromq.MessageData) {
	fmt.Println("OnClientDisconnected", data.ToString())
}

func (s Streambot) OnSystemHealthCheck(data zeromq.MessageData) {
	fmt.Println("OnSystemHealthCheck", data.ToString())
}

func (s Streambot) OnSystemUpdate(data zeromq.MessageData) {
	fmt.Println("OnSystemUpdate", data.ToString())
}

func (s Streambot) OnServerMapChanged(data zeromq.MessageData) {
	fmt.Println("OnServerMapChanged", data.ToString())
}

func (s Streambot) OnServerScoreChanged(data zeromq.MessageData) {
	fmt.Println("OnServerScoreChanged", data.ToInt())
}

func (s Streambot) OnServerStatusChanged(data zeromq.MessageData) {
	fmt.Println("OnServerStatusChanged", data.ToString())
}

func (s Streambot) OnServerTitleChanged(data zeromq.MessageData) {
	fmt.Println("OnServerTitleChanged", data.ToString())
}
