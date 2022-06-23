package tasks

import (
	"time"
)

type HealthCheckTask struct {
	isDone bool
}

func NewHealthCheckTask() HealthCheckTask {
	return HealthCheckTask{}
}

func (t *HealthCheckTask) Start(interval time.Duration) {
	t.isDone = false
	ticker := time.NewTicker(interval)

	go func() {
		for ; true; <-ticker.C {
			if t.isDone {
				return
			}

			// do stuff
		}
	}()
}

func (t *HealthCheckTask) Stop() {
	t.isDone = true
}
