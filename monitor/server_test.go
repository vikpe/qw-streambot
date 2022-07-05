package monitor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/streambot/monitor"
	"github.com/vikpe/streambot/zeromq/test_helpers"
)

func TestServerMonitor_Address(t *testing.T) {
	getInfo := func(address string) mvdsv.Mvdsv { return mvdsv.Mvdsv{} }
	onEvent := func(topic string, data ...any) {}
	serverMonitor := monitor.NewServerMonitor(getInfo, onEvent)

	serverMonitor.SetAddress("qw.foppa.dk:27501")
	assert.Equal(t, "qw.foppa.dk:27501", serverMonitor.GetAddress())

	timeFormat := "2006:01:02 15:04:05"
	assert.Equal(t, time.Now().Format(timeFormat), serverMonitor.GetAddressTimestamp().Format(timeFormat))
}

func TestServerMonitor_CompareStates(t *testing.T) {
	// mock server info responses
	mInfo := []mvdsv.Mvdsv{
		{Title: "1on1: x vs y [dm2]"},
		{Title: "1on1: x vs y [dm2]"},
		{
			Title:    "kombat / 1on1: x vs y [dm2]",
			Settings: qsettings.Settings{"matchtag": "kombat"},
		},
	}

	infoCallCount := 0
	getInfo := func(address string) mvdsv.Mvdsv {
		server := mInfo[infoCallCount]
		infoCallCount++
		return server
	}
	publisherMock := test_helpers.NewPublisherMock()
	serverMonitor := monitor.NewServerMonitor(getInfo, publisherMock.SendMessage)

	for i := 0; i < len(mInfo); i++ {
		serverMonitor.CompareStates()
	}

	expectCalls := [][]any{
		{"server.title_changed", "1on1: x vs y [dm2]"},
		{"server.matchtag_changed", "kombat"},
		{"server.title_changed", "kombat / 1on1: x vs y [dm2]"},
	}
	assert.Equal(t, expectCalls, publisherMock.SendMessageCalls)
}
