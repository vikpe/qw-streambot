package commander_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/zeromq/commander"
)

type call struct {
	topic string
	args  []any
}

type publisherMock struct {
	calls []call
}

func (s *publisherMock) SendMessage(topic string, args ...any) {
	s.calls = append(s.calls, call{topic: topic, args: args})
}

func TestCommander_Autotrack(t *testing.T) {
	publisher := publisherMock{}
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Autotrack()

	expectedCalls := []call{{
		topic: "ezquake.command", args: []any{"bot_track"},
	}}
	assert.Equal(t, expectedCalls, publisher.calls)
}

func TestCommander_DisableAuto(t *testing.T) {
	publisher := publisherMock{}
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.DisableAuto()

	expectedCalls := []call{{
		topic: "streambot.disable_auto", args: nil,
	}}
	assert.Equal(t, expectedCalls, publisher.calls)
}

func TestCommander_EnableAuto(t *testing.T) {
	publisher := publisherMock{}
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.EnableAuto()

	expectedCalls := []call{{
		topic: "streambot.enable_auto", args: nil,
	}}
	assert.Equal(t, expectedCalls, publisher.calls)
}

func TestCommander_SuggestServer(t *testing.T) {
	publisher := publisherMock{}
	cmder := commander.NewCommander(publisher.SendMessage)
	server := mvdsv.Mvdsv{Address: "qw.fopp.dk:27501"}
	cmder.SuggestServer(server)

	expectedCalls := []call{{
		topic: "streambot.suggest_server", args: []any{server},
	}}
	assert.Equal(t, expectedCalls, publisher.calls)
}

func TestCommander_StopEzquake(t *testing.T) {
	publisher := publisherMock{}
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.StopEzquake()

	expectedCalls := []call{{
		topic: "ezquake.stop", args: nil,
	}}
	assert.Equal(t, expectedCalls, publisher.calls)
}

func TestCommander_Track(t *testing.T) {
	publisher := publisherMock{}
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Track("xantom")

	expectedCalls := []call{{
		topic: "ezquake.command", args: []any{"bot_track xantom"},
	}}
	assert.Equal(t, expectedCalls, publisher.calls)
}
