package task

import (
	"time"

	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/streambot/topics"
)

type ServerMonitor struct {
	isDone  bool
	Address string
	onEvent func(string, any)
}

func NewServerMonitor(address string, onEvent func(string, any)) ServerMonitor {
	return ServerMonitor{
		isDone:  false,
		onEvent: onEvent,
		Address: address,
	}
}

func (s *ServerMonitor) Start(interval time.Duration) {
	s.isDone = false

	go func() {
		ticker := time.NewTicker(interval)
		prevState := NewServerState(s.Address)

		for ; true; <-ticker.C {
			if s.isDone {
				return
			}

			currentState := NewServerState(s.Address)
			diff := NewServerStateDiff(currentState, prevState)

			if diff.HasChangedTitle {
				s.onEvent(topics.ServerTitleChanged, currentState.Title)
			}

			if diff.HasChangedMap {
				s.onEvent(topics.ServerMapChanged, currentState.Map)
			}

			if diff.HasChangedScore {
				s.onEvent(topics.ServerScoreChanged, currentState.Score)
			}

			prevState = currentState
		}
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
	if "" == address {
		return ServerState{
			Title: "",
			Map:   "",
			Score: 0,
		}
	}

	genericServer, err := serverstat.GetInfo(address)

	if err != nil {
		return ServerState{
			Title: "",
			Map:   "",
			Score: 0,
		}
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
	HasChangedMap   bool
	HasChangedScore bool
}

func NewServerStateDiff(current ServerState, prev ServerState) ServerStateDiff {
	diff := ServerStateDiff{}

	if current.Title != prev.Title {
		diff.HasChangedTitle = true
	}

	if current.Map != prev.Map {
		diff.HasChangedMap = true
	}

	if current.Score != prev.Score {
		diff.HasChangedScore = true
	}

	return diff
}
