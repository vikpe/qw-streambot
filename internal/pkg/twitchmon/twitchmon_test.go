package twitchmon_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/types"
	"github.com/vikpe/streambot/internal/pkg/twitchmon"
)

func TestNewTwitchMonitor(t *testing.T) {
	chanQuakeWorld := types.TwitchStream{
		Channel: "QuakeWorld",
	}

	// mock getStreams
	callCount := 0
	mockGetStreams := func() []types.TwitchStream {
		callCount++

		if 1 == callCount {
			return make([]types.TwitchStream, 0)
		}

		return []types.TwitchStream{
			chanQuakeWorld,
		}
	}

	// mock event handler
	calls := make([]types.TwitchStream, 0)

	onStreamStarted := func(stream types.TwitchStream) {
		calls = append(calls, stream)
	}

	// run tests
	monitor := twitchmon.New(mockGetStreams, onStreamStarted)
	interval := 10 * time.Millisecond
	go monitor.Start(interval)
	defer monitor.Stop()

	expectedCalls := []types.TwitchStream{}
	assert.Equal(t, expectedCalls, calls)

	time.Sleep(interval)
	expectedCalls = []types.TwitchStream{chanQuakeWorld}
	assert.Equal(t, expectedCalls, calls)

	time.Sleep(interval)
	expectedCalls = []types.TwitchStream{chanQuakeWorld}
	assert.Equal(t, expectedCalls, calls)
}
