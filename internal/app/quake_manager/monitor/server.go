package monitor

import (
	"time"

	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/internal/comms/topic"
)

type MvdsvProvider func(address string) mvdsv.Mvdsv

type ServerMonitor struct {
	isDone          bool
	onEvent         func(string, ...any)
	getInfo         MvdsvProvider
	address         string
	serverTimestamp time.Time
	prevState       serverState
}

func NewServerMonitor(getInfo MvdsvProvider, onEvent func(topic string, data ...any)) *ServerMonitor {
	return &ServerMonitor{
		isDone:          false,
		getInfo:         getInfo,
		onEvent:         onEvent,
		address:         "",
		prevState:       serverState{},
		serverTimestamp: time.Time{},
	}
}

func (s *ServerMonitor) SetAddress(address string) {
	s.address = address

	if "" == address {
		s.serverTimestamp = time.Time{}
	} else {
		s.touchServerTimestamp()
	}
}

func (s *ServerMonitor) touchServerTimestamp() {
	s.serverTimestamp = time.Now()
}

func (s *ServerMonitor) GetAddress() string {
	return s.address
}

func (s *ServerMonitor) ClearAddress() {
	s.SetAddress("")
}

func (s *ServerMonitor) IsConnected() bool {
	return s.address != ""
}

func (s *ServerMonitor) GetTimeConnected() time.Duration {
	if s.serverTimestamp.IsZero() {
		return 0
	}

	return time.Now().Sub(s.serverTimestamp)
}

func (s *ServerMonitor) Start(interval time.Duration) {
	s.isDone = false

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		if s.isDone {
			return
		}

		s.CompareStates()
	}
}

func (s *ServerMonitor) CompareStates() {
	currentState := newServerState(s.getInfo, s.address)

	if currentState.Matchtag != s.prevState.Matchtag {
		s.onEvent(topic.ServerMatchtagChanged, currentState.Matchtag)
	}

	if currentState.Map != s.prevState.Map {
		s.touchServerTimestamp()
	}

	if currentState.Title != s.prevState.Title {
		s.onEvent(topic.ServerTitleChanged, currentState.Title)
	}

	s.prevState = currentState
}

func (s *ServerMonitor) Stop() {
	s.isDone = true
}

type serverState struct {
	Matchtag string
	Map      string
	Title    string
}

func newServerState(getInfo MvdsvProvider, address string) serverState {
	server := getInfo(address)

	return serverState{
		Matchtag: server.Settings.Get("matchtag", ""),
		Map:      server.Settings.Get("map", ""),
		Title:    server.Title,
	}
}
