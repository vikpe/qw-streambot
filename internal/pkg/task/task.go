package task

import (
	"time"
)

type PeriodicalTask struct {
	isActive bool
	onTick   func()
}

func NewPeriodicalTask(callback func()) *PeriodicalTask {
	return &PeriodicalTask{
		isActive: false,
		onTick:   callback,
	}
}

func (t *PeriodicalTask) Start(interval time.Duration) {
	if t.isActive {
		return
	}

	t.isActive = true
	ticker := time.NewTicker(interval)

	for ; true; <-ticker.C {
		if !t.isActive {
			return
		}

		t.onTick()
	}
}

func (t *PeriodicalTask) Stop() {
	t.isActive = false
}
