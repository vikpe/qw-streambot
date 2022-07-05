package commander_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/zeromq/commander"
	"github.com/vikpe/streambot/zeromq/test_helpers"
)

func TestCommander_Autotrack(t *testing.T) {
	publisher := test_helpers.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Autotrack()

	expectedCalls := [][]any{{"ezquake.command", "bot_track"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_Command(t *testing.T) {
	publisher := test_helpers.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Command("console")

	expectedCalls := [][]any{{"ezquake.command", "console"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_DisableAuto(t *testing.T) {
	publisher := test_helpers.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.DisableAuto()

	expectedCalls := [][]any{{"streambot.disable_auto"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_EnableAuto(t *testing.T) {
	publisher := test_helpers.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.EnableAuto()

	expectedCalls := [][]any{{"streambot.enable_auto"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_Lastscores(t *testing.T) {
	publisher := test_helpers.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Lastscores()

	expectedCalls := [][]any{{"ezquake.script", "lastscores"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_SuggestServer(t *testing.T) {
	publisher := test_helpers.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	server := mvdsv.Mvdsv{Address: "qw.fopp.dk:27501"}
	cmder.SuggestServer(server)

	expectedCalls := [][]any{{"streambot.suggest_server", server}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_Showscores(t *testing.T) {
	publisher := test_helpers.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Showscores()

	expectedCalls := [][]any{{"ezquake.script", "showscores"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_StopEzquake(t *testing.T) {
	publisher := test_helpers.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.StopEzquake()

	expectedCalls := [][]any{{"ezquake.stop"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_Track(t *testing.T) {
	publisher := test_helpers.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Track("xantom")

	expectedCalls := [][]any{{"ezquake.command", "bot_track xantom"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}
