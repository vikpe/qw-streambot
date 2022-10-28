package monitor

import (
	"github.com/ssoroka/slice"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
	"github.com/vikpe/streambot/internal/pkg/task"
	"golang.org/x/exp/slices"
)

func NewStreamsMonitor(getStreams func() []twitch.Stream, onStreamStarted func(stream twitch.Stream)) *task.PeriodicalTask {
	var prevChannels []string

	onTick := func() {
		streams := getStreams()
		currentChannels := slice.Map[twitch.Stream, string](streams, func(stream twitch.Stream) string {
			return stream.Channel
		})

		if prevChannels == nil {
			prevChannels = currentChannels
			return
		}

		for _, stream := range streams {
			if stream.Channel == "QuakeWorld" {
				continue
			}

			if !slices.Contains(prevChannels, stream.Channel) {
				onStreamStarted(stream)
			}
		}

		prevChannels = currentChannels
	}

	return task.NewPeriodicalTask(onTick)
}
