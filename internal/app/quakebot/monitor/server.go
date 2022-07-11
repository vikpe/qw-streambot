package monitor

import (
	"time"

	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/internal/comms/topic"
)

type MvdsvProvider func(address string) mvdsv.Mvdsv

type ServerMonitor struct {
	isDone           bool
	onEvent          func(string, ...any)
	getInfo          MvdsvProvider
	address          string
	addressTimestamp time.Time
	prevState        serverState
}

func NewServerMonitor(getInfo MvdsvProvider, onEvent func(topic string, data ...any)) *ServerMonitor {
	return &ServerMonitor{
		isDone:    false,
		getInfo:   getInfo,
		onEvent:   onEvent,
		address:   "",
		prevState: serverState{},
	}
}

func (s *ServerMonitor) SetAddress(address string) {
	s.address = address
	s.prevState = serverState{}
	s.addressTimestamp = time.Now()
}

func (s *ServerMonitor) GetAddress() string {
	return s.address
}

func (s *ServerMonitor) GetAddressTimestamp() time.Time {
	return s.addressTimestamp
}

func (s *ServerMonitor) Start(interval time.Duration) {
	s.isDone = false

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for ; true; <-ticker.C {
			if s.isDone {
				return
			}

			s.CompareStates()
		}
	}()
}

func (s *ServerMonitor) CompareStates() {
	currentState := newServerState(s.getInfo, s.address)

	if currentState.Matchtag != s.prevState.Matchtag {
		s.onEvent(topic.ServerMatchtagChanged, currentState.Matchtag)
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
	Map      string
	Matchtag string
	Score    int
	Title    string
}

func newServerState(getInfo MvdsvProvider, address string) serverState {
	server := getInfo(address)

	return serverState{
		Map:      server.Settings.Get("map", ""),
		Matchtag: server.Settings.Get("matchtag", ""),
		Score:    server.Score,
		Title:    server.Title,
	}
}
