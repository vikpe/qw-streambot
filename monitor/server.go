package monitor

import (
	"time"

	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/streambot/zeromq/topic"
)

type ServerMonitor struct {
	isDone           bool
	onEvent          func(string, ...any)
	address          string
	addressTimestamp time.Time
	prevState        serverState
}

func NewServerMonitor(onEvent func(string, ...any)) ServerMonitor {
	return ServerMonitor{
		isDone:    false,
		onEvent:   onEvent,
		address:   "",
		prevState: newServerState(""),
	}
}

func (s *ServerMonitor) SetAddress(address string) {
	s.address = address
	s.prevState = newServerState("")
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

			currentState := newServerState(s.address)
			diff := newServerStateDiff(currentState, s.prevState)

			if diff.HasChangedTitle {
				s.onEvent(topic.ServerTitleChanged, currentState.Title)
			}

			if diff.HasChangedMatchtag {
				s.onEvent(topic.ServerMatchtagChanged, currentState.Matchtag)
			}

			s.prevState = currentState
		}

		defer ticker.Stop()
	}()
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

func newServerState(address string) serverState {
	nullState := serverState{
		Map:      "",
		Matchtag: "",
		Score:    0,
		Title:    "",
	}
	if "" == address {
		return nullState
	}

	genericServer, err := serverstat.GetInfo(address)

	if err != nil {
		return nullState
	}

	server := convert.ToMvdsv(genericServer)

	return serverState{
		Map:      server.Settings.Get("map", ""),
		Matchtag: server.Settings.Get("matchtag", ""),
		Score:    server.Score,
		Title:    server.Title,
	}
}

type serverStateDiff struct {
	HasChangedTitle    bool
	HasChangedMatchtag bool
}

func newServerStateDiff(current serverState, prev serverState) serverStateDiff {
	diff := serverStateDiff{}

	if current.Title != prev.Title {
		diff.HasChangedTitle = true
	}

	if current.Matchtag != prev.Matchtag {
		diff.HasChangedMatchtag = true
	}

	return diff
}
