package main

import (
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/streambot/ezquake"
	"github.com/vikpe/streambot/qws"
	"github.com/vikpe/streambot/task"
	"github.com/vikpe/streambot/topics"
	"github.com/vikpe/streambot/util/twitch"
	"github.com/vikpe/streambot/zeromq"
)

type Streambot struct {
	clientPlayerName string
	pipe             ezquake.PipeWriter
	process          ezquake.Process
	serverMonitor    task.ServerMonitor
	twitch           twitch.Client
	publisher        zeromq.Publisher
	subscriber       zeromq.Subscriber
}

func NewStreambot(
	clientPlayerName string,
	process ezquake.Process,
	pipe ezquake.PipeWriter,
	twitchClient twitch.Client,
	publisher zeromq.Publisher,
	subscriber zeromq.Subscriber,
) Streambot {
	return Streambot{
		clientPlayerName: clientPlayerName,
		pipe:             pipe,
		process:          process,
		serverMonitor:    task.NewServerMonitor(publisher.SendMessage),
		twitch:           twitchClient,
		publisher:        publisher,
		subscriber:       subscriber,
	}
}

func (s *Streambot) Start() {
	// event listeners
	s.subscriber.Start(s.OnMessage)

	// event dispatchers
	processMonitor := task.NewProcessMonitor(&s.process, s.publisher.SendMessage)
	processMonitor.Start(3 * time.Second)
	s.serverMonitor.Start(5 * time.Second)

	// health check
	go func() {
		ticker := time.NewTicker(10 * time.Second)

		for ; true; <-ticker.C {
			s.publisher.SendMessage(topics.SystemHealthCheck, "")
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func (s *Streambot) OnMessage(msg zeromq.Message) {
	handlers := map[string]zeromq.MessageDataHandler{
		// commands
		topics.ConnectToServer:   s.OnConnectToServer,
		topics.ClientCommand:     s.OnClientCommand,
		topics.StopClient:        s.OnStopClient,
		topics.SystemUpdate:      s.OnSystemUpdate,
		topics.SystemHealthCheck: s.OnSystemHealthCheck,

		// client events
		topics.ClientStarted:      s.OnClientStarted,
		topics.ClientStopped:      s.OnClientStopped,
		topics.ClientConnected:    s.OnClientConnected,
		topics.ClientDisconnected: s.OnClientDisconnected,

		// server events
		topics.ServerMapChanged:    s.OnServerMapChanged,
		topics.ServerScoreChanged:  s.OnServerScoreChanged,
		topics.ServerStatusChanged: s.OnServerStatusChanged,
		topics.ServerTitleChanged:  s.OnServerTitleChanged,
	}

	if handler, ok := handlers[msg.Topic]; ok {
		handler(msg.Data)
	} else {
		fmt.Println("no handler defined for", msg.Topic, fmt.Sprintf("%T", msg.Topic), msg.Data)
	}
}

func (s *Streambot) Evaluate() {
	fmt.Println("Evaluate", s.serverMonitor.GetAddress())

	if "" == s.serverMonitor.GetAddress() {
		server, _ := qws.GetBestServer()
		s.publisher.SendMessage(topics.ConnectToServer, server)
	}
}

func (s *Streambot) OnStart() {
	fmt.Println("OnStart")
	s.Evaluate()
}

func (s *Streambot) OnConnectToServer(data zeromq.MessageData) {
	var server mvdsv.Mvdsv
	data.To(&server)

	fmt.Print("OnConnectToServer", server.Address, data)

	if s.serverMonitor.GetAddress() == server.Address {
		fmt.Println(" .. already connected to server")
		return
	}

	if len(server.QtvStream.Url) > 0 {
		s.ClientCommand(fmt.Sprintf("qtvplay %s", server.QtvStream.Url))
		s.ClientCommand("bot_track")
	} else {
		s.ClientCommand(fmt.Sprintf("connect %s", server.Address))

		time.AfterFunc(4*time.Second, func() {
			s.ClientCommand("bot_track")
		})
	}

	fmt.Println(" .. new server!", server.Address)
	s.serverMonitor.SetAddress(server.Address)

	// validate that we connected
	time.AfterFunc(8*time.Second, func() {
		fmt.Println("VALIDATE THAT WE ARE ON SERVER")
		genericServer, _ := serverstat.GetInfo(server.Address)
		server := convert.ToMvdsv(genericServer)

		if analyze.HasSpectator(server, s.clientPlayerName) {
			fmt.Println(" - oooh yes. ggggggggggggggggggg")
		} else {
			fmt.Println(" - NIET!")
		}
	})
}

func (s *Streambot) ClientCommand(command string) {
	s.publisher.SendMessage(topics.ClientCommand, command)
}

func (s *Streambot) OnClientCommand(data zeromq.MessageData) {
	fmt.Println("OnClientCommand", data.ToString())

	if s.process.IsStarted() {
		s.pipe.Write(data.ToString())
	}
}

func (s *Streambot) OnClientStarted(data zeromq.MessageData) {
	fmt.Println("OnClientStarted", data.ToString())

	time.AfterFunc(4*time.Second, func() {
		s.publisher.SendMessage(topics.ClientCommand, "toggleconsole")
	})
}

func (s *Streambot) OnStopClient(data zeromq.MessageData) {
	fmt.Println("OnStopClient", data.ToString())
	s.process.Stop(syscall.SIGTERM)
}

func (s *Streambot) OnClientStopped(data zeromq.MessageData) {
	fmt.Println("OnClientStopped", data.ToString())
}

func (s *Streambot) OnClientConnected(data zeromq.MessageData) {
	fmt.Println("OnClientConnected", data.ToString())
}

func (s *Streambot) OnClientDisconnected(data zeromq.MessageData) {
	fmt.Println("OnClientDisconnected", data.ToString())
}

func (s *Streambot) OnSystemHealthCheck(data zeromq.MessageData) {
	fmt.Println("OnSystemHealthCheck", data.ToString())
	s.Evaluate()
}

func (s *Streambot) OnSystemUpdate(data zeromq.MessageData) {
	fmt.Println("OnSystemUpdate", data.ToString())
}

func (s *Streambot) OnServerMapChanged(data zeromq.MessageData) {
	fmt.Println("OnServerMapChanged", data.ToString())
}

func (s *Streambot) OnServerScoreChanged(data zeromq.MessageData) {
	fmt.Println("OnServerScoreChanged", data.ToInt())
}

func (s *Streambot) OnServerStatusChanged(data zeromq.MessageData) {
	fmt.Println("OnServerStatusChanged", data.ToString())
}

func (s *Streambot) OnServerTitleChanged(data zeromq.MessageData) {
	fmt.Println("OnServerTitleChanged", data.ToString())
	//s.twitch.SetTitle(data.ToString())
}
