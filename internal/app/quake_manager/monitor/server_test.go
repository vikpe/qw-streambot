package monitor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/streambot/internal/app/quake_manager/monitor"
	"github.com/vikpe/streambot/internal/pkg/zeromq/mock"
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
		{
			Settings: qsettings.Settings{"map": "dm6"},
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
	publisherMock := mock.NewPublisherMock()
	serverMonitor = monitor.NewServerMonitor(getInfo, publisherMock.SendMessage)
	go serverMonitor.Start(time.Microsecond)
	time.Sleep(time.Millisecond * 20)

	expectCalls := [][]any{
		{"server.title_changed", "1on1: x vs y [dm2]"},
		{"server.matchtag_changed", "kombat"},
		{"server.title_changed", "kombat / 1on1: x vs y [dm2]"},
		{"server.matchtag_changed", ""},
		{"server.map_changed", "dm6"},
		{"server.title_changed", ""},
	}
	assert.Equal(t, expectCalls, publisherMock.SendMessageCalls)
}

func TestServerMonitor_Address(t *testing.T) {
	getInfo := func(address string) mvdsv.Mvdsv { return mvdsv.Mvdsv{} }
	onEvent := func(topic string, data ...any) {}
	serverMonitor := monitor.NewServerMonitor(getInfo, onEvent)

	serverMonitor.SetAddress("qw.foppa.dk:27501")
	assert.Equal(t, "qw.foppa.dk:27501", serverMonitor.GetAddress())
}

func TestServerMonitor_GetTimeConnected(t *testing.T) {
	getInfo := func(address string) mvdsv.Mvdsv { return mvdsv.Mvdsv{} }
	onEvent := func(topic string, data ...any) {}
	serverMonitor := monitor.NewServerMonitor(getInfo, onEvent)

	assert.Equal(t, int64(0), serverMonitor.GetTimeConnected().Milliseconds())
	serverMonitor.SetAddress("qw.foppa.dk:27501")
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, int64(10), serverMonitor.GetTimeConnected().Milliseconds())
}
