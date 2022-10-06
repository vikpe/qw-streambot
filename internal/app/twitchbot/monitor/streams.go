package monitor

import (
	"github.com/ssoroka/slice"
	"github.com/vikpe/qw-hub-api/types"
	"github.com/vikpe/streambot/internal/pkg/task"
	"golang.org/x/exp/slices"
)

func NewStreamsMonitor(getStreams func() []types.TwitchStream, onStreamStarted func(stream types.TwitchStream)) *task.PeriodicalTask {
	var prevChannels []string

	onTick := func() {
		streams := getStreams()
		currentChannels := slice.Map[types.TwitchStream, string](streams, func(stream types.TwitchStream) string {
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
