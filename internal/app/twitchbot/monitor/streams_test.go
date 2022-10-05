package monitor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/types"
	"github.com/vikpe/streambot/internal/app/twitchbot/monitor"
)

func TestStreamsMonitor(t *testing.T) {
	// mock getStreams
	chan1 := types.TwitchStream{Channel: "1"}
	chan2 := types.TwitchStream{Channel: "2"}
	chan3 := types.TwitchStream{Channel: "3"}
	chanQW := types.TwitchStream{Channel: "QuakeWorld"}
	mockedResult := [][]types.TwitchStream{
		{chan1},                       // call 1
		{chan1, chan2, chanQW},        // call 2
		{chan1, chan2, chan3, chanQW}, // call 3
	}

	callCount := 0
	mockGetStreams := func() []types.TwitchStream {
		callCount++

		if callCount <= len(mockedResult) {
			return mockedResult[callCount-1]
		}

		return []types.TwitchStream{}
	}

	// mock event handler
	calls := make([]types.TwitchStream, 0)

	onStreamStarted := func(stream types.TwitchStream) {
		calls = append(calls, stream)
	}

	// run tests
	streamsMonitor := monitor.NewStreamsMonitor(mockGetStreams, onStreamStarted)
	interval := 10 * time.Millisecond
	go streamsMonitor.Start(interval)
	defer streamsMonitor.Stop()
	time.Sleep(time.Millisecond)

	expectedCalls := []types.TwitchStream{}
	assert.Equal(t, expectedCalls, calls)

	time.Sleep(interval)
	expectedCalls = []types.TwitchStream{chan2}
	assert.Equal(t, expectedCalls, calls)

	time.Sleep(interval + time.Millisecond)
	expectedCalls = []types.TwitchStream{chan2, chan3}
	assert.Equal(t, expectedCalls, calls)
}
