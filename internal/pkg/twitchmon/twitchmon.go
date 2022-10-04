package twitchmon

import (
	"time"

	"github.com/vikpe/qw-hub-api/types"
	"golang.org/x/exp/slices"
)

type TwitchMonitor struct {
	isDone          bool
	OnStreamStarted func(stream types.TwitchStream)
	getStreams      func() []types.TwitchStream
	prevStreams     []types.TwitchStream
}

func New(getStreams func() []types.TwitchStream, callbackFunc func(stream types.TwitchStream)) *TwitchMonitor {
	return &TwitchMonitor{
		getStreams:      getStreams,
		OnStreamStarted: callbackFunc,
		prevStreams:     nil,
	}
}

func (t *TwitchMonitor) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	t.prevStreams = t.getStreams()
	t.isDone = false

	for ; true; <-ticker.C {
		if t.isDone {
			return
		}

		t.CompareStates()
	}
}

func (t *TwitchMonitor) CompareStates() {
	currentStreams := t.getStreams()

	for _, stream := range currentStreams {
		if !slices.Contains(t.prevStreams, stream) {
			t.OnStreamStarted(stream)
		}
	}

	t.prevStreams = currentStreams
}

func (t *TwitchMonitor) Stop() {
	t.isDone = true
}
