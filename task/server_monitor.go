package task

import (
	"time"

	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/streambot/topics"
)

type ServerMonitor struct {
	isDone           bool
	onEvent          func(string, any)
	address          string
	addressTimestamp time.Time
	prevState        ServerState
}

func NewServerMonitor(onEvent func(string, any)) ServerMonitor {
	return ServerMonitor{
		isDone:    false,
		onEvent:   onEvent,
		address:   "",
		prevState: NewServerState(""),
	}
}

func (s *ServerMonitor) SetAddress(address string) {
	s.address = address
	s.prevState = NewServerState("")
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

		for ; true; <-ticker.C {
			if s.isDone {
				return
			}

			currentState := NewServerState(s.address)
			diff := NewServerStateDiff(currentState, s.prevState)

			if diff.HasChangedTitle {
				s.onEvent(topics.ServerTitleChanged, currentState.Title)
			}

			s.prevState = currentState
		}

		defer ticker.Stop()
	}()
}

func (s *ServerMonitor) Stop() {
	s.isDone = true
}

type ServerState struct {
	Title string
	Map   string
	Score int
}

func NewServerState(address string) ServerState {
	nullState := ServerState{
		Title: "",
		Map:   "",
		Score: 0,
	}
	if "" == address {
		return nullState
	}

	genericServer, err := serverstat.GetInfo(address)

	if err != nil {
		return nullState
	}

	server := convert.ToMvdsv(genericServer)

	return ServerState{
		Title: server.Title,
		Map:   server.Settings.Get("map", ""),
		Score: server.Score,
	}
}

type ServerStateDiff struct {
	HasChangedTitle bool
}

func NewServerStateDiff(current ServerState, prev ServerState) ServerStateDiff {
	diff := ServerStateDiff{}

	if current.Title != prev.Title {
		diff.HasChangedTitle = true
	}

	return diff
}
