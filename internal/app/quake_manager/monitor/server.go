package monitor

import (
	"time"

	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/internal/comms/topic"
)

type MvdsvProvider func(address string) mvdsv.Mvdsv

type ServerMonitor struct {
	isDone               bool
	onEvent              func(string, ...any)
	getInfo              MvdsvProvider
	address              string
	connectionTimestamp  time.Time
	lastStartedTimestamp time.Time
	prevState            serverState
}

func NewServerMonitor(getInfo MvdsvProvider, onEvent func(topic string, data ...any)) *ServerMonitor {
	return &ServerMonitor{
		isDone:               false,
		getInfo:              getInfo,
		onEvent:              onEvent,
		address:              "",
		prevState:            serverState{},
		connectionTimestamp:  time.Time{},
		lastStartedTimestamp: time.Time{},
	}
}

func (s *ServerMonitor) SetAddress(address string) {
	s.address = address

	if "" == address {
		s.connectionTimestamp = time.Time{}
	} else {
		s.connectionTimestamp = time.Now()
	}

	s.lastStartedTimestamp = s.connectionTimestamp
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

func (s *ServerMonitor) GetConnectionDuration() time.Duration {
	if s.connectionTimestamp.IsZero() {
		return 0
	}

	return time.Since(s.connectionTimestamp)
}

func (s *ServerMonitor) GetIdleDuration() time.Duration {
	if s.lastStartedTimestamp.IsZero() {
		return 0
	}

	return time.Since(s.lastStartedTimestamp)
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

	if currentState.Title != s.prevState.Title {
		s.onEvent(topic.ServerTitleChanged, currentState.Title)
	}

	if currentState.IsStarted {
		s.lastStartedTimestamp = time.Now()
	}

	s.prevState = currentState
}

func (s *ServerMonitor) Stop() {
	s.isDone = true
}

type serverState struct {
	Matchtag  string
	IsStarted bool
	Title     string
}

func newServerState(getInfo MvdsvProvider, address string) serverState {
	server := getInfo(address)

	return serverState{
		Matchtag:  server.Settings.Get("matchtag", ""),
		IsStarted: server.Status.IsStarted(),
		Title:     server.Title,
	}
}
