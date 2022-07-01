package main

import (
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/streambot/ezquake"
	"github.com/vikpe/streambot/qws"
	"github.com/vikpe/streambot/task"
	"github.com/vikpe/streambot/topics"
	"github.com/vikpe/streambot/util/sstat"
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
	evaluateTask     task.PeriodicalTask
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
		evaluateTask:     task.NewPeriodicalTask(func() { publisher.SendMessage(topics.StreambotEvaluate, "") }),
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

	if s.process.IsStarted() {
		s.evaluateTask.Start(10 * time.Second)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func (s *Streambot) OnMessage(msg zeromq.Message) {
	handlers := map[string]zeromq.MessageDataHandler{
		// commands
		topics.StreambotEnableAuto:      s.OnStreambotEnableAuto,
		topics.StreambotDisableAuto:     s.OnStreambotDisableAuto,
		topics.StreambotConnectToServer: s.OnStreambotConnectToServer,
		topics.StreambotSuggestServer:   s.OnStreambotSuggestServer,
		topics.ClientCommand:            s.OnClientCommand,
		topics.StopClient:               s.OnStopClient,
		topics.StreambotSystemUpdate:    s.OnStreambotSystemUpdate,
		topics.StreambotEvaluate:        s.OnStreambotEvaluate,

		// client events
		topics.ClientStarted: s.OnClientStarted,
		topics.ClientStopped: s.OnClientStopped,

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

func (s *Streambot) OnStreambotEnableAuto(data zeromq.MessageData) {
	s.AutoMode = true
	s.publisher.SendMessage(topics.StreambotEnableAuto, "")
}

func (s *Streambot) OnStreambotDisableAuto(data zeromq.MessageData) {
	s.AutoMode = false
}

func (s *Streambot) ValidateCurrentServer() {
	if "" == s.serverMonitor.GetAddress() {
		return
	}

	secondsConnected := time.Now().Sub(s.serverMonitor.GetAddressTimestamp()).Seconds()
	connectionGracePeriod := 10.0

	if secondsConnected <= connectionGracePeriod {
		return
	}

	currentServer := sstat.GetMvdsvServer(s.serverMonitor.GetAddress())

	if analyze.HasSpectator(currentServer, s.clientPlayerName) {
		return
	}

	altName := fmt.Sprintf("%s(1)", s.clientPlayerName)

	if analyze.HasSpectator(currentServer, altName) {
		s.publisher.SendMessage(topics.ClientCommand, fmt.Sprintf("name %s", s.clientPlayerName))
		return
	}

	fmt.Println("not connected to current server (reset server address)", currentServer.SpectatorNames, currentServer.QtvStream.SpectatorNames)
	s.serverMonitor.SetAddress("")
}

func (s *Streambot) OnStreambotEvaluate(data zeromq.MessageData) {
	fmt.Println()
	pp.Print("OnStreambotEvaluate - ")

	// check process
	if !s.process.IsStarted() {
		fmt.Println("not started: do nothing (wait until started)")
		return
	}

	// validate current server
	s.ValidateCurrentServer()

	// check server
	currentServer := sstat.GetMvdsvServer(s.serverMonitor.GetAddress())

	// validate based on auto mode enabled/disabled
	if s.AutoMode {
		s.evaluateAutoModeEnabled(currentServer)
	} else {
		s.evaluateAutoModeDisabled(currentServer)
	}
}

func (s *Streambot) evaluateAutoModeEnabled(currentServer mvdsv.Mvdsv) {
	shouldConsiderChange := 0 == currentServer.Score || currentServer.Mode.IsCustom() || currentServer.Status.IsStandby()

	if !shouldConsiderChange {
		return
	}

	bestServer, err := qws.GetBestServer()

	if err != nil {
		return
	}

	shouldStay := currentServer.Score >= bestServer.Score || currentServer.Address == bestServer.Address

	if shouldStay {
		return
	}

	s.publisher.SendMessage(topics.StreambotConnectToServer, bestServer)
}

func (s *Streambot) evaluateAutoModeDisabled(currentServer mvdsv.Mvdsv) {
	const MinScore = 30
	isOkServer := currentServer.Score >= MinScore

	if !isOkServer {
		fmt.Println("server is ok: do nothing")
		return
	}
	fmt.Println("server is shit: enable auto")

	s.publisher.SendMessage(topics.StreambotEnableAuto, "")
}

func (s *Streambot) OnStreambotSuggestServer(data zeromq.MessageData) {
	s.publisher.SendMessage(topics.StreambotDisableAuto, "")
	s.publisher.SendMessage(topics.StreambotConnectToServer, data)
}

func (s *Streambot) OnStreambotConnectToServer(data zeromq.MessageData) {
	var server mvdsv.Mvdsv
	data.To(&server)

	pp.Print("OnStreambotConnectToServer", server.Address, data)

	if s.serverMonitor.GetAddress() == server.Address {
		fmt.Println(" .. already connected to server")
		return
	}

	if len(server.QtvStream.Url) > 0 {
		s.ClientCommand(fmt.Sprintf("qtvplay %s", server.QtvStream.Url))
	} else {
		s.ClientCommand(fmt.Sprintf("connect %s", server.Address))
	}

	time.AfterFunc(4*time.Second, func() {
		s.ClientCommand("bot_track")
	})

	fmt.Println(" .. new server!", server.Address)
	s.serverMonitor.SetAddress(server.Address)
}

func (s *Streambot) ClientCommand(command string) {
	s.publisher.SendMessage(topics.ClientCommand, command)
}

func (s *Streambot) OnClientCommand(data zeromq.MessageData) {
	pp.Println("OnClientCommand", data.ToString())

	if s.process.IsStarted() {
		s.pipe.Write(data.ToString())
	}
}

func (s *Streambot) OnClientStarted(data zeromq.MessageData) {
	pp.Println("OnClientStarted", data.ToString())

	s.evaluateTask.Start(10 * time.Second)

	time.AfterFunc(5*time.Second, func() {
		s.publisher.SendMessage(topics.ClientCommand, "toggleconsole")
	})
}

func (s *Streambot) OnStopClient(data zeromq.MessageData) {
	pp.Println("OnStopClient", data.ToString())
	s.process.Stop(syscall.SIGTERM)
}

func (s *Streambot) OnClientStopped(data zeromq.MessageData) {
	pp.Println("OnClientStopped", data.ToString())
	s.evaluateTask.Stop()
}

func (s *Streambot) OnStreambotSystemUpdate(data zeromq.MessageData) {
	pp.Println("OnStreambotSystemUpdate", data.ToString())
}

func (s *Streambot) OnServerMapChanged(data zeromq.MessageData) {
	pp.Println("OnServerMapChanged", data.ToString())
}

func (s *Streambot) OnServerScoreChanged(data zeromq.MessageData) {
	pp.Println("OnServerScoreChanged", data.ToInt())
}

func (s *Streambot) OnServerStatusChanged(data zeromq.MessageData) {
	pp.Println("OnServerStatusChanged", data.ToString())
}

func (s *Streambot) OnServerTitleChanged(data zeromq.MessageData) {
	pp.Println("OnServerTitleChanged", data.ToString())
	//s.twitch.SetTitle(data.ToString())
}
