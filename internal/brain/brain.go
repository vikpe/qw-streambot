package brain

import (
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/vikpe/prettyfmt"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/streambot/com/topic"
	"github.com/vikpe/streambot/internal/brain/util/calc"
	"github.com/vikpe/streambot/internal/brain/util/proc"
	"github.com/vikpe/streambot/internal/brain/util/task"
	"github.com/vikpe/streambot/internal/monitor"
	"github.com/vikpe/streambot/pkg/ezquake"
	"github.com/vikpe/streambot/pkg/qws"
	"github.com/vikpe/streambot/pkg/zeromq"
	"github.com/vikpe/streambot/pkg/zeromq/message"
	"github.com/vikpe/streambot/third_party/sstat"
	"github.com/vikpe/streambot/third_party/twitch"
)

var pfmt = prettyfmt.New("brain", color.FgHiMagenta, "15:04:05", color.FgWhite)

type Brain struct {
	clientPlayerName string
	pipe             ezquake.PipeWriter
	process          proc.ProcessController
	serverMonitor    monitor.ServerMonitor
	evaluateTask     task.PeriodicalTask
	twitch           twitch.Client
	publisher        zeromq.Publisher
	subscriber       zeromq.Subscriber
	AutoMode         bool
}

func NewBrain(
	clientPlayerName string,
	process proc.ProcessController,
	pipe ezquake.PipeWriter,
	twitchClient twitch.Client,
	publisher zeromq.Publisher,
	subscriber zeromq.Subscriber,
) Brain {
	return Brain{
		clientPlayerName: clientPlayerName,
		pipe:             pipe,
		process:          process,
		serverMonitor:    monitor.NewServerMonitor(sstat.GetMvdsvServer, publisher.SendMessage),
		evaluateTask:     task.NewPeriodicalTask(func() { publisher.SendMessage(topic.StreambotEvaluate) }),
		twitch:           twitchClient,
		publisher:        publisher,
		subscriber:       subscriber,
		AutoMode:         true,
	}
}

func (s *Brain) Start() {
	// event listeners
	s.subscriber.Start(s.OnMessage)
	zeromq.WaitForConnection()

	// event dispatchers
	processMonitor := monitor.NewProcessMonitor(s.process.IsStarted, s.publisher.SendMessage)
	processMonitor.Start(3 * time.Second)
	s.serverMonitor.Start(5 * time.Second)

	if s.process.IsStarted() {
		s.evaluateTask.Start(10 * time.Second)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func (s *Brain) OnMessage(msg message.Message) {
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
		topic.ServerTitleChanged:    s.OnServerTitleChanged,
		topic.ServerMatchtagChanged: s.OnServerMatchtagChanged,
	}

	if handler, ok := handlers[msg.Topic]; ok {
		handler(msg)
	} else {
		fmt.Println("no handler defined for", msg.Topic, fmt.Sprintf("%T", msg.Topic), msg.Content)
	}
}

func (s *Brain) OnStreambotEnableAuto(msg message.Message) {
	s.AutoMode = true
	s.publisher.SendMessage(topic.StreambotEvaluate)
}

func (s *Brain) OnStreambotDisableAuto(msg message.Message) {
	s.AutoMode = false
}

func (s *Brain) ValidateCurrentServer() {
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

func (s *Brain) OnStreambotEvaluate(msg message.Message) {
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

func (s *Brain) evaluateAutoModeEnabled() {
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

func (s *Brain) evaluateAutoModeDisabled() {
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

func (s *Brain) OnStreambotSuggestServer(msg message.Message) {
	var server mvdsv.Mvdsv
	msg.Content.To(&server)

	s.publisher.SendMessage(topic.StreambotDisableAuto)
	s.publisher.SendMessage(topic.StreambotConnectToServer, server)
}

func (s *Brain) OnStreambotConnectToServer(msg message.Message) {
	var server mvdsv.Mvdsv
	msg.Content.To(&server)

	pfmt.Println("OnStreambotConnectToServer", server.Address, msg.Content)

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

	s.serverMonitor.SetAddress(server.Address)
}

func (s *Brain) ClientCommand(command string) {
	s.publisher.SendMessage(topic.EzquakeCommand, command)
}

func (s *Brain) OnEzquakeCommand(msg message.Message) {
	pfmt.Println("OnEzquakeCommand", msg.Content.ToString())

	if !s.process.IsStarted() {
		return
	}

	s.pipe.Write(msg.Content.ToString())
}

func (s *Brain) OnEzquakeScript(msg message.Message) {
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

func (s *Brain) OnEzquakeStarted(msg message.Message) {
	pfmt.Println("OnEzquakeStarted")
	s.evaluateTask.Start(10 * time.Second)

	time.AfterFunc(5*time.Second, func() { s.ClientCommand("toggleconsole") })
}

func (s *Brain) OnStopEzquake(msg message.Message) {
	pfmt.Println("OnStopEzquake")
	s.process.Stop(syscall.SIGTERM)

	time.AfterFunc(2*time.Second, func() {
		if s.process.IsStarted() {
			s.process.Stop(syscall.SIGKILL)
		}
	})
}

func (s *Brain) OnEzquakeStopped(msg message.Message) {
	pfmt.Println("OnEzquakeStopped")
	s.serverMonitor.SetAddress("")
	s.evaluateTask.Stop()
}

func (s *Brain) OnStreambotSystemUpdate(msg message.Message) {
	pfmt.Println("OnStreambotSystemUpdate")
}

func (s *Brain) OnServerTitleChanged(msg message.Message) {
	pfmt.Println("OnServerTitleChanged", msg.Content.ToString())
	s.twitch.SetTitle(msg.Content.ToString())
}

func (s *Brain) OnServerMatchtagChanged(msg message.Message) {
	matchtag := msg.Content.ToString()
	textScale := calc.StaticTextScale(matchtag)
	s.ClientCommand(fmt.Sprintf("hud_static_text_scale %f;bot_set_statictext %s", textScale, matchtag))
}
