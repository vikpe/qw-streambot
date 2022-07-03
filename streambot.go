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
	"github.com/vikpe/streambot/monitor"
	"github.com/vikpe/streambot/third_party/qws"
	"github.com/vikpe/streambot/third_party/sstat"
	"github.com/vikpe/streambot/third_party/twitch"
	"github.com/vikpe/streambot/util/task"
	"github.com/vikpe/streambot/util/term"
	"github.com/vikpe/streambot/zeromq"
	"github.com/vikpe/streambot/zeromq/message"
	"github.com/vikpe/streambot/zeromq/topic"
)

var pp = term.NewPrettyPrinter("brain", color.FgHiMagenta)

type Streambot struct {
	clientPlayerName string
	pipe             ezquake.PipeWriter
	process          ezquake.Process
	serverMonitor    monitor.ServerMonitor
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
		serverMonitor:    monitor.NewServerMonitor(publisher.SendMessage),
		evaluateTask:     task.NewPeriodicalTask(func() { publisher.SendMessage(topic.StreambotEvaluate) }),
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
	processMonitor := monitor.NewProcessMonitor(&s.process, s.publisher.SendMessage)
	processMonitor.Start(3 * time.Second)
	s.serverMonitor.Start(5 * time.Second)

	if s.process.IsStarted() {
		s.evaluateTask.Start(10 * time.Second)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func (s *Streambot) OnMessage(msg message.Message) {
	handlers := map[string]message.Handler{
		// commands
		topic.StreambotEnableAuto:      s.OnStreambotEnableAuto,
		topic.StreambotDisableAuto:     s.OnStreambotDisableAuto,
		topic.StreambotConnectToServer: s.OnStreambotConnectToServer,
		topic.StreambotSuggestServer:   s.OnStreambotSuggestServer,
		topic.EzquakeCommand:           s.OnEzquakeCommand,
		topic.EzquakeScript:            s.OnEzquakeScript,
		topic.StopEzquake:              s.OnStopEzquake,
		topic.StreambotSystemUpdate:    s.OnStreambotSystemUpdate,
		topic.StreambotEvaluate:        s.OnStreambotEvaluate,

		// ezquake events
		topic.EzquakeStarted: s.OnEzquakeStarted,
		topic.EzquakeStopped: s.OnEzquakeStopped,

		// server events
		topic.ServerTitleChanged: s.OnServerTitleChanged,
	}

	if handler, ok := handlers[msg.Topic]; ok {
		handler(msg)
	} else {
		fmt.Println("no handler defined for", msg.Topic, fmt.Sprintf("%T", msg.Topic), msg.Content)
	}
}

func (s *Streambot) OnStreambotEnableAuto(msg message.Message) {
	s.AutoMode = true
	s.publisher.SendMessage(topic.StreambotEvaluate)
}

func (s *Streambot) OnStreambotDisableAuto(msg message.Message) {
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
		s.ClientCommand(fmt.Sprintf("name %s", s.clientPlayerName))
		return
	}

	fmt.Println("not connected to current server (reset server address)", currentServer.SpectatorNames, currentServer.QtvStream.SpectatorNames)
	s.serverMonitor.SetAddress("")
}

func (s *Streambot) OnStreambotEvaluate(msg message.Message) {
	// check process
	if !s.process.IsStarted() {
		return
	}

	// validate current server
	s.ValidateCurrentServer()

	// validate based on auto mode
	if s.AutoMode {
		s.evaluateAutoModeEnabled()
	} else {
		s.evaluateAutoModeDisabled()
	}
}

func (s *Streambot) evaluateAutoModeEnabled() {
	currentServer := sstat.GetMvdsvServer(s.serverMonitor.GetAddress())
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

	s.publisher.SendMessage(topic.StreambotConnectToServer, bestServer)
}

func (s *Streambot) evaluateAutoModeDisabled() {
	currentServer := sstat.GetMvdsvServer(s.serverMonitor.GetAddress())
	const MinScore = 30
	isOkServer := currentServer.Score >= MinScore

	if !isOkServer {
		fmt.Println("server is ok: do nothing")
		return
	}

	secondsConnected := time.Now().Sub(s.serverMonitor.GetAddressTimestamp()).Seconds()
	gracePeriod := 60.0 * 5 // 5 minutes

	if secondsConnected < gracePeriod {
		return
	}

	fmt.Println("server is shit: enable auto")

	s.publisher.SendMessage(topic.StreambotEnableAuto)
}

func (s *Streambot) OnStreambotSuggestServer(msg message.Message) {
	var server mvdsv.Mvdsv
	msg.Content.To(&server)

	s.publisher.SendMessage(topic.StreambotDisableAuto)
	s.publisher.SendMessage(topic.StreambotConnectToServer, server)
}

func (s *Streambot) OnStreambotConnectToServer(msg message.Message) {
	var server mvdsv.Mvdsv
	msg.Content.To(&server)

	pp.Print("OnStreambotConnectToServer", server.Address, msg.Content)

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
	s.publisher.SendMessage(topic.EzquakeCommand, command)
}

func (s *Streambot) OnEzquakeCommand(msg message.Message) {
	pp.Println("OnEzquakeCommand", msg.Content.ToString())

	if !s.process.IsStarted() {
		return
	}

	s.pipe.Write(msg.Content.ToString())
}

func (s *Streambot) OnEzquakeScript(msg message.Message) {
	script := msg.Content.ToString()

	switch script {
	case "lastscores":
		s.ClientCommand("toggleconsole;lastscores")
		time.AfterFunc(8*time.Second, func() { s.ClientCommand("toggleconsole") })
	case "showscores":
		s.ClientCommand("+showscores")
		time.AfterFunc(8*time.Second, func() { s.ClientCommand("-showscores") })
	}
}

func (s *Streambot) OnEzquakeStarted(msg message.Message) {
	pp.Println("OnEzquakeStarted")

	s.evaluateTask.Start(10 * time.Second)

	time.AfterFunc(5*time.Second, func() { s.ClientCommand("toggleconsole") })
}

func (s *Streambot) OnStopEzquake(msg message.Message) {
	pp.Println("OnStopEzquake")
	s.process.Stop(syscall.SIGTERM)
}

func (s *Streambot) OnEzquakeStopped(msg message.Message) {
	pp.Println("OnEzquakeStopped")
	s.evaluateTask.Stop()
}

func (s *Streambot) OnStreambotSystemUpdate(msg message.Message) {
	pp.Println("OnStreambotSystemUpdate")
}

func (s *Streambot) OnServerTitleChanged(msg message.Message) {
	pp.Println("OnServerTitleChanged", msg.Content.ToString())
	s.twitch.SetTitle(msg.Content.ToString())
}
