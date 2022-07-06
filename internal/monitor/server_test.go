package monitor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/streambot/internal/monitor"
	"github.com/vikpe/streambot/pkg/zeromq/test_helpers"
)

func TestServerMonitor(t *testing.T) {
	var serverMonitor *monitor.ServerMonitor

	// return mocked values and stop after x calls
	mockedResponses := []mvdsv.Mvdsv{
		{Title: "1on1: x vs y [dm2]"},
		{Title: "1on1: x vs y [dm2]"},
		{
			Title:    "kombat / 1on1: x vs y [dm2]",
			Settings: qsettings.Settings{"matchtag": "kombat"},
		},
	}

	callCount := 0
	getInfo := func(address string) mvdsv.Mvdsv {
		server := mockedResponses[callCount]
		callCount++

		if callCount >= len(mockedResponses) {
			serverMonitor.Stop()
		}

		return server
	}
	publisherMock := test_helpers.NewPublisherMock()
	serverMonitor = monitor.NewServerMonitor(getInfo, publisherMock.SendMessage)
	serverMonitor.Start(time.Microsecond)
	time.Sleep(time.Millisecond * 20)

	expectCalls := [][]any{
		{"server.title_changed", "1on1: x vs y [dm2]"},
		{"server.matchtag_changed", "kombat"},
		{"server.title_changed", "kombat / 1on1: x vs y [dm2]"},
	}
	assert.Equal(t, expectCalls, publisherMock.SendMessageCalls)
}

func TestServerMonitor_Address(t *testing.T) {
	getInfo := func(address string) mvdsv.Mvdsv { return mvdsv.Mvdsv{} }
	onEvent := func(topic string, data ...any) {}
	serverMonitor := monitor.NewServerMonitor(getInfo, onEvent)

	serverMonitor.SetAddress("qw.foppa.dk:27501")
	assert.Equal(t, "qw.foppa.dk:27501", serverMonitor.GetAddress())

	timeFormat := "2006:01:02 15:04:05"
	assert.Equal(t, time.Now().Format(timeFormat), serverMonitor.GetAddressTimestamp().Format(timeFormat))
}
