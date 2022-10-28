package monitor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
	"github.com/vikpe/streambot/internal/app/twitchbot/monitor"
)

func TestStreamsMonitor(t *testing.T) {
	// mock getStreams
	chan1 := twitch.Stream{Channel: "1", ViewerCount: 1}
	chan2 := twitch.Stream{Channel: "2", ViewerCount: 1}
	chan2Updated := twitch.Stream{Channel: "2", ViewerCount: 2}
	chan3 := twitch.Stream{Channel: "3", ViewerCount: 1}
	chanQW := twitch.Stream{Channel: "QuakeWorld", ViewerCount: 1}
	mockedResult := [][]twitch.Stream{
		{chan1},                              // call 1
		{chan1, chan2, chanQW},               // call 2
		{chan1, chan2Updated, chan3, chanQW}, // call 3
	}

	callCount := 0
	mockGetStreams := func() []twitch.Stream {
		callCount++

		if callCount <= len(mockedResult) {
			return mockedResult[callCount-1]
		}

		return []twitch.Stream{}
	}

	// mock event handler
	calls := make([]twitch.Stream, 0)

	onStreamStarted := func(stream twitch.Stream) {
		calls = append(calls, stream)
	}

	// run tests
	streamsMonitor := monitor.NewStreamsMonitor(mockGetStreams, onStreamStarted)
	interval := 10 * time.Millisecond
	go streamsMonitor.Start(interval)
	defer streamsMonitor.Stop()
	time.Sleep(time.Millisecond)

	expectedCalls := []twitch.Stream{}
	assert.Equal(t, expectedCalls, calls)

	time.Sleep(interval)
	expectedCalls = []twitch.Stream{chan2}
	assert.Equal(t, expectedCalls, calls)

	time.Sleep(interval + time.Millisecond)
	expectedCalls = []twitch.Stream{chan2, chan3}
	assert.Equal(t, expectedCalls, calls)
}
