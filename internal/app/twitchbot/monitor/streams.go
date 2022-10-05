package monitor

import (
	"github.com/vikpe/qw-hub-api/types"
	"github.com/vikpe/streambot/internal/pkg/task"
	"golang.org/x/exp/slices"
)

func NewStreamsMonitor(getStreams func() []types.TwitchStream, onStreamStarted func(stream types.TwitchStream)) *task.PeriodicalTask {
	var prevStreams []types.TwitchStream

	onTick := func() {
		currentStreams := getStreams()

		if prevStreams == nil {
			prevStreams = currentStreams
			return
		}

		for _, stream := range currentStreams {
			if stream.Channel == "QuakeWorld" {
				continue
			}

			if !slices.Contains(prevStreams, stream) {
				onStreamStarted(stream)
			}
		}

		prevStreams = currentStreams
	}

	return task.NewPeriodicalTask(onTick)
}
