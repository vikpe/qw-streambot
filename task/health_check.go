package task

import (
	"time"

	"github.com/vikpe/streambot/events"
)

type HealthCheckTask struct {
	isDone  bool
	onEvent EventHandler
}

func NewHealthCheckTask(onEvent EventHandler) HealthCheckTask {
	return HealthCheckTask{
		onEvent: onEvent,
	}
}

func (t *HealthCheckTask) Start(interval time.Duration) {
	t.isDone = false

	go func() {
		ticker := time.NewTicker(interval)

		for ; true; <-ticker.C {
			if t.isDone {
				return
			}

			t.onEvent(events.StreambotHealthCheck, "")
		}
	}()
}

func (t *HealthCheckTask) Stop() {
	t.isDone = true
}
