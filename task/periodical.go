package task

import (
	"time"
)

type PeriodicalTask struct {
	isDone bool
	onTick func()
}

func NewPeriodicalTask(onTick func()) PeriodicalTask {
	return PeriodicalTask{
		onTick: onTick,
	}
}

func (t *PeriodicalTask) Start(interval time.Duration) {
	t.isDone = false

	go func() {
		ticker := time.NewTicker(interval)

		for ; true; <-ticker.C {
			if t.isDone {
				return
			}

			t.onTick()
		}
	}()
}

func (t *PeriodicalTask) Stop() {
	t.isDone = true
}
