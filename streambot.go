package main

import (
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/streambot/ezquake"
	"github.com/vikpe/streambot/qws"
	"github.com/vikpe/streambot/task"
	"github.com/vikpe/streambot/topics"
	"github.com/vikpe/streambot/util/term"
	"github.com/vikpe/streambot/util/twitch"
	"github.com/vikpe/streambot/zeromq"
)

var pp = term.NewPrettyPrinter("brain", color.FgHiMagenta)

type Streambot struct {
	clientPlayerName string
	pipe             ezquake.PipeWriter
	process          ezquake.Process
	serverMonitor    task.ServerMonitor
	twitch           twitch.Client
	publisher        zeromq.Publisher
	subscriber       zeromq.Subscriber
	AutoMode         bool
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
		AutoMode:         true,
	}
}

func (s *Streambot) Start() {
	// event listeners
	s.subscriber.Start(s.OnMessage)
	zeromq.WaitForConnection()

	// event dispatchers
	processMonitor := task.NewProcessMonitor(&s.process, s.publisher.SendMessage)
	processMonitor.Start(3 * time.Second)
	s.serverMonitor.Start(5 * time.Second)

	// evaluate every x seconds
	ev := task.NewPeriodicalTask(func() { s.publisher.SendMessage(topics.Evaluate, "") })
	ev.Start(10 * time.Second)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func (s *Streambot) OnMessage(msg zeromq.Message) {
	handlers := map[string]zeromq.MessageDataHandler{
		// commands
		topics.EnableAuto:      s.OnEnableAuto,
		topics.DisableAuto:     s.OnDisableAuto,
		topics.ConnectToServer: s.OnConnectToServer,
		topics.SuggestServer:   s.OnSuggestServer,
		topics.ClientCommand:   s.OnClientCommand,
		topics.StopClient:      s.OnStopClient,
		topics.SystemUpdate:    s.OnSystemUpdate,
		topics.Evaluate:        s.OnSystemEvaluate,

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

func (s *Streambot) OnEnableAuto(data zeromq.MessageData) {
	s.AutoMode = true
	s.publisher.SendMessage(topics.EnableAuto, "")
}

func (s *Streambot) OnDisableAuto(data zeromq.MessageData) {
	s.AutoMode = false
}

func (s *Streambot) OnSystemEvaluate(data zeromq.MessageData) {
	pp.Print("OnSystemEvaluate")

	if !s.process.IsStarted() {
		pp.Print("not started: do nothing (wait until started)")
		return
	}

	currentServer := s.CurrentServer()

	if s.AutoMode {
		shouldConsiderChange := 0 == currentServer.Score || currentServer.Mode.IsCustom() || currentServer.Status.IsStandby()

		if !shouldConsiderChange {
			pp.Print("should not consider change: do nothing")
			return
		}

		bestServer, err := qws.GetBestServer()

		if err != nil {
			pp.Print("no server found..")
			return
		}

		isAtBestServer := currentServer.Score >= bestServer.Score

		if isAtBestServer {
			pp.Print("at best server: do nothing")
			return
		}

		s.publisher.SendMessage(topics.ConnectToServer, bestServer)

	} else {
		const MinScore = 30
		isCrapServer := currentServer.Score < MinScore

		if !isCrapServer {
			pp.Print("server is ok: do nothing")
			return
		}
		fmt.Print("server is shit: enable auto")

		s.publisher.SendMessage(topics.EnableAuto, "")
	}
}

func (s *Streambot) CurrentServer() mvdsv.Mvdsv {
	return GetServer(s.serverMonitor.GetAddress())
}

func (s *Streambot) OnSuggestServer(data zeromq.MessageData) {
	s.publisher.SendMessage(topics.DisableAuto, "")
	s.publisher.SendMessage(topics.ConnectToServer, data)
}

func (s *Streambot) OnConnectToServer(data zeromq.MessageData) {
	var server mvdsv.Mvdsv
	data.To(&server)

	fmt.Print("OnConnectToServer", server.Address, data)

	if s.serverMonitor.GetAddress() == server.Address {
		pp.Print(" .. already connected to server")
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

	pp.Print(" .. new server!", server.Address)
	s.serverMonitor.SetAddress(server.Address)

	// validate that we connected
	time.AfterFunc(8*time.Second, func() {
		pp.Print("VALIDATE THAT WE ARE ON SERVER")
		genericServer, _ := serverstat.GetInfo(server.Address)
		server := convert.ToMvdsv(genericServer)

		if analyze.HasSpectator(server, s.clientPlayerName) {
			pp.Print(" - oooh yes. ggggggggggggggggggg")
		} else {
			pp.Print(" - NIET!")
		}
	})
}

func (s *Streambot) ClientCommand(command string) {
	s.publisher.SendMessage(topics.ClientCommand, command)
}

func (s *Streambot) OnClientCommand(data zeromq.MessageData) {
	pp.Print("OnClientCommand", data.ToString())

	if s.process.IsStarted() {
		s.pipe.Write(data.ToString())
	}
}

func (s *Streambot) OnClientStarted(data zeromq.MessageData) {
	pp.Print("OnClientStarted", data.ToString())

	time.AfterFunc(4*time.Second, func() {
		s.publisher.SendMessage(topics.ClientCommand, "toggleconsole")
	})
}

func (s *Streambot) OnStopClient(data zeromq.MessageData) {
	pp.Print("OnStopClient", data.ToString())
	s.process.Stop(syscall.SIGTERM)
}

func (s *Streambot) OnClientStopped(data zeromq.MessageData) {
	pp.Print("OnClientStopped", data.ToString())
}

func (s *Streambot) OnClientConnected(data zeromq.MessageData) {
	pp.Print("OnClientConnected", data.ToString())
}

func (s *Streambot) OnClientDisconnected(data zeromq.MessageData) {
	pp.Print("OnClientDisconnected", data.ToString())
}

func (s *Streambot) OnSystemUpdate(data zeromq.MessageData) {
	pp.Print("OnSystemUpdate", data.ToString())
}

func (s *Streambot) OnServerMapChanged(data zeromq.MessageData) {
	pp.Print("OnServerMapChanged", data.ToString())
}

func (s *Streambot) OnServerScoreChanged(data zeromq.MessageData) {
	pp.Print("OnServerScoreChanged", data.ToInt())
}

func (s *Streambot) OnServerStatusChanged(data zeromq.MessageData) {
	pp.Print("OnServerStatusChanged", data.ToString())
}

func (s *Streambot) OnServerTitleChanged(data zeromq.MessageData) {
	pp.Print("OnServerTitleChanged", data.ToString())
	//s.twitch.SetTitle(data.ToString())
}

func GetServer(address string) mvdsv.Mvdsv {
	nullResult := mvdsv.Mvdsv{}

	if "" == address {
		return nullResult
	}

	genericServer, err := serverstat.GetInfo(address)

	if err != nil {
		return nullResult
	}

	return convert.ToMvdsv(genericServer)
}
