package monitor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/monitor"
	"github.com/vikpe/streambot/zeromq/test_helpers"
)

func TestProcessMonitor(t *testing.T) {
	var processMonitor monitor.ProcessMonitor

	// return mocked values and stop after x calls
	mockedResponses := []bool{false, false, true, false, true}
	callCount := 0
	getIsStarted := func() bool {
		value := mockedResponses[callCount]
		callCount++

		if callCount >= len(mockedResponses) {
			processMonitor.Stop()
		}

		return value
	}

	// run monitor
	publisherMock := test_helpers.NewPublisherMock()
	processMonitor = monitor.NewProcessMonitor(getIsStarted, publisherMock.SendMessage)
	processMonitor.Start(time.Microsecond)
	time.Sleep(time.Millisecond)

	expectCalls := [][]any{
		{"ezquake.started"},
		{"ezquake.stopped"},
		{"ezquake.started"},
	}
	assert.Equal(t, expectCalls, publisherMock.SendMessageCalls)
}
