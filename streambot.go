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
		topics.EzquakeCommand:           s.OnEzquakeCommand,
		topics.EzquakeLastscores:        s.OnEzquakeLastscores,
		topics.EzquakeShowscores:        s.OnEzquakeShowscores,
		topics.StopEzquake:              s.OnStopEzquake,
		topics.StreambotSystemUpdate:    s.OnStreambotSystemUpdate,
		topics.StreambotEvaluate:        s.OnStreambotEvaluate,

		// ezquake events
		topics.EzquakeStarted: s.OnEzquakeStarted,
		topics.EzquakeStopped: s.OnEzquakeStopped,

		// server events
		topics.ServerTitleChanged: s.OnServerTitleChanged,
	}

	if handler, ok := handlers[msg.Topic]; ok {
		handler(msg.Data)
	} else {
		fmt.Println("no handler defined for", msg.Topic, fmt.Sprintf("%T", msg.Topic), msg.Data)
	}
}

func (s *Streambot) OnStreambotEnableAuto(data zeromq.MessageData) {
	s.AutoMode = true
	s.publisher.SendMessage(topics.StreambotEvaluate, "")
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
		s.ClientCommand(fmt.Sprintf("name %s", s.clientPlayerName))
		return
	}

	fmt.Println("not connected to current server (reset server address)", currentServer.SpectatorNames, currentServer.QtvStream.SpectatorNames)
	s.serverMonitor.SetAddress("")
}

func (s *Streambot) OnStreambotEvaluate(data zeromq.MessageData) {
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

	s.publisher.SendMessage(topics.StreambotConnectToServer, bestServer)
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

	s.publisher.SendMessage(topics.StreambotEnableAuto, "")
}

func (s *Streambot) OnStreambotSuggestServer(data zeromq.MessageData) {
	var server mvdsv.Mvdsv
	data.To(&server)

	s.publisher.SendMessage(topics.StreambotDisableAuto, "")
	s.publisher.SendMessage(topics.StreambotConnectToServer, server)
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
	s.publisher.SendMessage(topics.EzquakeCommand, command)
}

func (s *Streambot) OnEzquakeCommand(data zeromq.MessageData) {
	pp.Println("OnEzquakeCommand", data.ToString())

	if s.process.IsStarted() {
		s.pipe.Write(data.ToString())
	}
}

func (s *Streambot) OnEzquakeLastscores(data zeromq.MessageData) {
	s.ClientCommand("toggleconsole;lastscores")

	time.AfterFunc(8*time.Second, func() {
		s.ClientCommand("toggleconsole")
	})
}

func (s *Streambot) OnEzquakeShowscores(data zeromq.MessageData) {
	s.ClientCommand("+showscores")

	time.AfterFunc(8*time.Second, func() {
		s.ClientCommand("-showscores")
	})
}

func (s *Streambot) OnEzquakeStarted(data zeromq.MessageData) {
	pp.Println("OnEzquakeStarted", data.ToString())

	s.evaluateTask.Start(10 * time.Second)

	time.AfterFunc(5*time.Second, func() {
		s.ClientCommand("toggleconsole")
	})
}

func (s *Streambot) OnStopEzquake(data zeromq.MessageData) {
	pp.Println("OnStopEzquake", data.ToString())
	s.process.Stop(syscall.SIGTERM)
}

func (s *Streambot) OnEzquakeStopped(data zeromq.MessageData) {
	pp.Println("OnEzquakeStopped", data.ToString())
	s.evaluateTask.Stop()
}

func (s *Streambot) OnStreambotSystemUpdate(data zeromq.MessageData) {
	pp.Println("OnStreambotSystemUpdate", data.ToString())
}

func (s *Streambot) OnServerTitleChanged(data zeromq.MessageData) {
	pp.Println("OnServerTitleChanged", data.ToString())
	s.twitch.SetTitle(data.ToString())
}
