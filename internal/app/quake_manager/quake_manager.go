package quake_manager

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/vikpe/go-ezquake"
	"github.com/vikpe/prettyfmt"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/streambot/internal/app/quake_manager/monitor"
	"github.com/vikpe/streambot/internal/comms/commander"
	"github.com/vikpe/streambot/internal/comms/topic"
	"github.com/vikpe/streambot/internal/pkg/calc"
	"github.com/vikpe/streambot/internal/pkg/qws"
	"github.com/vikpe/streambot/internal/pkg/sstat"
	"github.com/vikpe/streambot/internal/pkg/task"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

var pfmt = prettyfmt.New("quakemanager", color.FgHiCyan, "15:04:05", color.FgWhite)

type QuakeManager struct {
	clientPlayerName string
	controller       *ezquake.ClientController
	processMonitor   *monitor.ProcessMonitor
	serverMonitor    *monitor.ServerMonitor
	evaluateTask     *task.PeriodicalTask
	subscriber       *zeromq.Subscriber
	commander        *commander.Commander
	assetManager     *ezquake.AssetManager
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
	controller := ezquake.NewClientController(ezquakeProcessUsername, ezquakeBinPath)
	publisher := zeromq.NewPublisher(publisherAddress)
	subscriber := zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll)

	manager := QuakeManager{
		assetManager:     ezquake.NewAssetManager(filepath.Dir(ezquakeBinPath)),
		clientPlayerName: clientPlayerName,
		controller:       controller,
		processMonitor:   monitor.NewProcessMonitor(controller.Process.IsStarted, publisher.SendMessage),
		serverMonitor:    monitor.NewServerMonitor(sstat.GetMvdsvServer, publisher.SendMessage),
		evaluateTask:     task.NewPeriodicalTask(func() { publisher.SendMessage(topic.StreambotEvaluate) }),
		subscriber:       subscriber,
		commander:        commander.NewCommander(publisher.SendMessage),
		AutoMode:         true,
	}
	subscriber.OnMessage = manager.OnMessage

	return &manager
}

func (m *QuakeManager) Start() {
	m.stopChan = make(chan os.Signal, 1)
	signal.Notify(m.stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		// event listeners
		go m.subscriber.Start()
		zeromq.WaitForConnection()

		// event dispatchers
		go m.processMonitor.Start(3 * time.Second)
		go m.serverMonitor.Start(5 * time.Second)

		if m.controller.Process.IsStarted() {
			go m.evaluateTask.Start(10 * time.Second)
		}
	}()
	<-m.stopChan
}

func (m *QuakeManager) Stop() {
	if m.stopChan == nil {
		return
	}
	m.stopChan <- syscall.SIGINT
	time.Sleep(30 * time.Millisecond)
}

func (m *QuakeManager) OnMessage(msg message.Message) {
	handlers := map[string]message.Handler{
		// commands
		topic.StreambotDisableAuto:   m.OnStreambotDisableAuto,
		topic.StreambotEnableAuto:    m.OnStreambotEnableAuto,
		topic.StreambotEvaluate:      m.OnStreambotEvaluate,
		topic.StreambotSuggestServer: m.OnStreambotSuggestServer,
		topic.EzquakeCommand:         m.OnEzquakeCommand,
		topic.EzquakeScript:          m.OnEzquakeScript,
		topic.EzquakeStop:            m.OnStopEzquake,

		// ezquake events
		topic.EzquakeStarted: m.OnEzquakeStarted,
		topic.EzquakeStopped: m.OnEzquakeStopped,

		// server events
		topic.ServerMatchtagChanged: m.OnServerMatchtagChanged,
	}

	if handler, ok := handlers[msg.Topic]; ok {
		handler(msg)
	}
}

func (m *QuakeManager) OnStreambotEnableAuto(msg message.Message) {
	m.AutoMode = true
	m.commander.Evaluate()
}

func (m *QuakeManager) OnStreambotDisableAuto(msg message.Message) {
	m.AutoMode = false
}

func (m *QuakeManager) ValidateCurrentServer() {
	if !m.serverMonitor.IsConnected() {
		return
	}

	const connectionGraceDuration = 10 * time.Second
	if m.serverMonitor.GetConnectionDuration() <= connectionGraceDuration {
		return
	}

	currentServer := sstat.GetMvdsvServer(m.serverMonitor.GetAddress())
	if analyze.HasSpectator(currentServer, m.clientPlayerName) {
		return
	}

	//
	altName := fmt.Sprintf("%s(1)", m.clientPlayerName)
	if analyze.HasSpectator(currentServer, altName) {
		m.commander.Commandf("name %s", m.clientPlayerName)
		return
	}

	// download missing maps
	mapName := currentServer.Settings.Get("map", "")

	if len(mapName) > 0 && !m.assetManager.HasMap(mapName) {
		pfmt.Printfln("trying to download map %s", mapName)

		err := m.assetManager.DownloadMap(mapName)

		if err == nil {
			pfmt.Printfln("fail")
		} else {
			pfmt.Printfln("success")
			return
		}
	}

	pfmt.Println("not connected to current server (reset server address)", currentServer.SpectatorNames, currentServer.QtvStream.SpectatorNames)
	m.serverMonitor.ClearAddress()
}

func (m *QuakeManager) OnStreambotEvaluate(msg message.Message) {
	if !m.controller.Process.IsStarted() {
		return
	}

	m.ValidateCurrentServer()

	if m.AutoMode {
		m.evaluateAutoModeEnabled()
	} else {
		m.evaluateAutoModeDisabled()
	}
}

func (m *QuakeManager) evaluateAutoModeEnabled() {
	const idleGraceDuration = 20 * time.Second
	currentServer := sstat.GetMvdsvServer(m.serverMonitor.GetAddress())
	isAllowedIdle := m.serverMonitor.IsConnected() && m.serverMonitor.GetIdleDuration() <= idleGraceDuration && currentServer.Mode.IsXonX()

	if isAllowedIdle {
		return
	}

	shouldConsiderChange := (0 == currentServer.Score) || !currentServer.Mode.IsXonX() || currentServer.Status.IsStandby() || (currentServer.Status.Description == "Score screen")

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

	m.connectToServer(bestServer)
}

func (m *QuakeManager) evaluateAutoModeDisabled() {
	currentServer := sstat.GetMvdsvServer(m.serverMonitor.GetAddress())

	const minAcceptableScore = 4
	if currentServer.Score >= minAcceptableScore {
		return
	}

	var idleGraceDuration float64

	if currentServer.Score >= 30 {
		idleGraceDuration = 5
	} else {
		idleGraceDuration = 3
	}

	if m.serverMonitor.GetIdleDuration().Minutes() <= idleGraceDuration {
		return
	}

	m.commander.EnableAuto()
}

func (m *QuakeManager) OnStreambotSuggestServer(msg message.Message) {
	var server mvdsv.Mvdsv
	msg.Content.To(&server)

	m.commander.DisableAuto()
	m.connectToServer(server)
}

func (m *QuakeManager) connectToServer(server mvdsv.Mvdsv) {
	if m.serverMonitor.GetAddress() == server.Address {
		return
	}

	if len(server.QtvStream.Url) > 0 {
		m.commander.Commandf("qtvplay %s", server.QtvStream.Url)
	} else {
		m.commander.Commandf("connect %s", server.Address)
	}

	time.AfterFunc(4*time.Second, func() {
		m.commander.Autotrack()
	})

	m.serverMonitor.SetAddress(server.Address)
}

func (m *QuakeManager) OnEzquakeCommand(msg message.Message) {
	m.controller.Command(msg.Content.ToString())
}

func (m *QuakeManager) OnEzquakeScript(msg message.Message) {
	script := msg.Content.ToString()

	switch script {
	case "lastscores":
		m.commander.Command("toggleconsole;lastscores")
		time.AfterFunc(8*time.Second, func() { m.commander.Command("toggleconsole") })
	case "showscores":
		m.commander.Command("+showscores")
		time.AfterFunc(8*time.Second, func() { m.commander.Command("-showscores") })
	}
}

func (m *QuakeManager) OnEzquakeStarted(msg message.Message) {
	pfmt.Println("OnEzquakeStarted")
	go m.evaluateTask.Start(10 * time.Second)
	time.AfterFunc(5*time.Second, func() { m.commander.Command("toggleconsole") })
}

func (m *QuakeManager) OnStopEzquake(msg message.Message) {
	pfmt.Println("OnStopEzquake")
	m.controller.Process.Stop(syscall.SIGTERM)

	time.AfterFunc(2*time.Second, func() {
		if m.controller.Process.IsStarted() {
			m.controller.Process.Stop(syscall.SIGKILL)
		}
	})
}

func (m *QuakeManager) OnEzquakeStopped(msg message.Message) {
	pfmt.Println("OnEzquakeStopped")
	m.serverMonitor.ClearAddress()
	m.evaluateTask.Stop()
}

func (m *QuakeManager) OnServerMatchtagChanged(msg message.Message) {
	matchtag := msg.Content.ToString()
	pfmt.Println("OnServerMatchtagChanged", matchtag)

	if strings.Contains(matchtag, "paus") {
		return
	}

	if len(matchtag) > 0 {
		m.commander.Commandf("hud_static_text_scale %f", calc.StaticTextScale(matchtag))
	}

	m.commander.Commandf("bot_set_statictext %s", matchtag)
}
