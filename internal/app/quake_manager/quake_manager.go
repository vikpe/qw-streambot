package quake_manager

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/vikpe/prettyfmt"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/streambot/internal/app/quake_manager/monitor"
	"github.com/vikpe/streambot/internal/comms/commander"
	"github.com/vikpe/streambot/internal/comms/topic"
	"github.com/vikpe/streambot/internal/pkg/calc"
	"github.com/vikpe/streambot/internal/pkg/ezquake"
	"github.com/vikpe/streambot/internal/pkg/proc"
	"github.com/vikpe/streambot/internal/pkg/qws"
	"github.com/vikpe/streambot/internal/pkg/sstat"
	"github.com/vikpe/streambot/internal/pkg/task"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

var pfmt = prettyfmt.New("quakemanager", color.FgHiCyan, "15:04:05", color.FgWhite)

type QuakeManager struct {
	clientPlayerName string
	pipe             *ezquake.PipeWriter
	process          *proc.ProcessController
	processMonitor   *monitor.ProcessMonitor
	serverMonitor    *monitor.ServerMonitor
	evaluateTask     *task.PeriodicalTask
	publisher        *zeromq.Publisher
	subscriber       *zeromq.Subscriber
	commander        *commander.Commander
	stopChan         chan os.Signal
	AutoMode         bool
}

func New(
	clientPlayerName string,
	ezquakeBinPath string,
	ezquakeProcessUsername string,
	publisherAddress string,
	subscriberAddress string,
) *QuakeManager {
	publisher := zeromq.NewPublisher(publisherAddress)
	subscriber := zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll)

	process := proc.NewProcessController(ezquakeBinPath)
	manager := QuakeManager{
		clientPlayerName: clientPlayerName,
		pipe:             ezquake.NewPipeWriter(ezquakeProcessUsername),
		process:          process,
		processMonitor:   monitor.NewProcessMonitor(process.IsStarted, publisher.SendMessage),
		serverMonitor:    monitor.NewServerMonitor(sstat.GetMvdsvServer, publisher.SendMessage),
		evaluateTask:     task.NewPeriodicalTask(func() { publisher.SendMessage(topic.StreambotEvaluate) }),
		subscriber:       subscriber,
		publisher:        publisher,
		commander:        commander.NewCommander(publisher.SendMessage),
		AutoMode:         true,
	}
	subscriber.OnMessage = manager.OnMessage

	return &manager
}

func (b *QuakeManager) Start() {
	b.stopChan = make(chan os.Signal, 1)
	signal.Notify(b.stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		// event listeners
		go b.subscriber.Start()
		zeromq.WaitForConnection()

		// event dispatchers
		go b.processMonitor.Start(3 * time.Second)
		go b.serverMonitor.Start(5 * time.Second)

		if b.process.IsStarted() {
			go b.evaluateTask.Start(10 * time.Second)
		}
	}()
	<-b.stopChan
}

func (b *QuakeManager) Stop() {
	if b.stopChan == nil {
		return
	}
	b.stopChan <- syscall.SIGINT
	time.Sleep(30 * time.Millisecond)
}

func (b *QuakeManager) OnMessage(msg message.Message) {
	handlers := map[string]message.Handler{
		// commands
		topic.StreambotDisableAuto:   b.OnStreambotDisableAuto,
		topic.StreambotEnableAuto:    b.OnStreambotEnableAuto,
		topic.StreambotEvaluate:      b.OnStreambotEvaluate,
		topic.StreambotSuggestServer: b.OnStreambotSuggestServer,
		topic.EzquakeCommand:         b.OnEzquakeCommand,
		topic.EzquakeScript:          b.OnEzquakeScript,
		topic.EzquakeStop:            b.OnStopEzquake,

		// ezquake events
		topic.EzquakeStarted: b.OnEzquakeStarted,
		topic.EzquakeStopped: b.OnEzquakeStopped,

		// server events
		topic.ServerMatchtagChanged: b.OnServerMatchtagChanged,
	}

	if handler, ok := handlers[msg.Topic]; ok {
		handler(msg)
	}
}

func (b *QuakeManager) OnStreambotEnableAuto(msg message.Message) {
	b.AutoMode = true
	b.commander.Evaluate()
}

func (b *QuakeManager) OnStreambotDisableAuto(msg message.Message) {
	b.AutoMode = false
}

func (b *QuakeManager) ValidateCurrentServer() {
	if "" == b.serverMonitor.GetAddress() {
		return
	}

	secondsConnected := time.Now().Sub(b.serverMonitor.GetAddressTimestamp()).Seconds()
	connectionGracePeriod := 10.0
	if secondsConnected <= connectionGracePeriod {
		return
	}

	currentServer := sstat.GetMvdsvServer(b.serverMonitor.GetAddress())
	if analyze.HasSpectator(currentServer, b.clientPlayerName) {
		return
	}

	altName := fmt.Sprintf("%s(1)", b.clientPlayerName)
	if analyze.HasSpectator(currentServer, altName) {
		b.commander.Commandf("name %s", b.clientPlayerName)
		return
	}

	pfmt.Println("not connected to current server (reset server address)", currentServer.SpectatorNames, currentServer.QtvStream.SpectatorNames)
	b.serverMonitor.SetAddress("")
}

func (b *QuakeManager) OnStreambotEvaluate(msg message.Message) {
	if !b.process.IsStarted() {
		return
	}

	b.ValidateCurrentServer()

	if b.AutoMode {
		b.evaluateAutoModeEnabled()
	} else {
		b.evaluateAutoModeDisabled()
	}
}

func (b *QuakeManager) evaluateAutoModeEnabled() {
	currentServer := sstat.GetMvdsvServer(b.serverMonitor.GetAddress())
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

	b.connectToServer(bestServer)
}

func (b *QuakeManager) evaluateAutoModeDisabled() {
	currentServer := sstat.GetMvdsvServer(b.serverMonitor.GetAddress())
	const MinScore = 25
	isOkServer := currentServer.Score >= MinScore

	if !isOkServer {
		return
	}

	secondsConnected := time.Now().Sub(b.serverMonitor.GetAddressTimestamp()).Seconds()
	gracePeriod := 60.0 * 5 // 5 minutes

	if secondsConnected < gracePeriod {
		return
	}

	b.commander.EnableAuto()
}

func (b *QuakeManager) OnStreambotSuggestServer(msg message.Message) {
	var server mvdsv.Mvdsv
	msg.Content.To(&server)

	b.commander.DisableAuto()
	b.connectToServer(server)
}

func (b *QuakeManager) connectToServer(server mvdsv.Mvdsv) {
	if b.serverMonitor.GetAddress() == server.Address {
		return
	}

	if len(server.QtvStream.Url) > 0 {
		b.commander.Commandf("qtvplay %s", server.QtvStream.Url)
	} else {
		b.commander.Commandf("connect %s", server.Address)
	}

	time.AfterFunc(4*time.Second, func() {
		b.commander.Autotrack()
	})

	b.serverMonitor.SetAddress(server.Address)
}

func (b *QuakeManager) OnEzquakeCommand(msg message.Message) {
	if !b.process.IsStarted() {
		return
	}

	b.pipe.Write(msg.Content.ToString())
}

func (b *QuakeManager) OnEzquakeScript(msg message.Message) {
	script := msg.Content.ToString()

	switch script {
	case "lastscores":
		b.commander.Command("toggleconsole;lastscores")
		time.AfterFunc(8*time.Second, func() { b.commander.Command("toggleconsole") })
	case "showscores":
		b.commander.Command("+showscores")
		time.AfterFunc(8*time.Second, func() { b.commander.Command("-showscores") })
	}
}

func (b *QuakeManager) OnEzquakeStarted(msg message.Message) {
	pfmt.Println("OnEzquakeStarted")
	go b.evaluateTask.Start(10 * time.Second)
	time.AfterFunc(5*time.Second, func() { b.commander.Command("toggleconsole") })
}

func (b *QuakeManager) OnStopEzquake(msg message.Message) {
	pfmt.Println("OnStopEzquake")
	b.process.Stop(syscall.SIGTERM)

	time.AfterFunc(2*time.Second, func() {
		if b.process.IsStarted() {
			b.process.Stop(syscall.SIGKILL)
		}
	})
}

func (b *QuakeManager) OnEzquakeStopped(msg message.Message) {
	pfmt.Println("OnEzquakeStopped")
	b.serverMonitor.SetAddress("")
	b.evaluateTask.Stop()
}

func (b *QuakeManager) OnServerMatchtagChanged(msg message.Message) {
	matchtag := msg.Content.ToString()

	pfmt.Println("OnServerMatchtagChanged", matchtag)
	textScale := calc.StaticTextScale(matchtag)
	b.commander.Commandf("hud_static_text_scale %f;bot_set_statictext %s", textScale, matchtag)
}
