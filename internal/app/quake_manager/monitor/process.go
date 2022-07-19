package monitor

import (
	"time"

	"github.com/vikpe/streambot/internal/comms/topic"
)

type ProcessMonitor struct {
	isDone           bool
	processIsStarted func() bool
	onEvent          func(string, ...any)
	prevState        bool
}

func NewProcessMonitor(processIsStarted func() bool, onEvent func(topic string, data ...any)) ProcessMonitor {
	return ProcessMonitor{
		isDone:           false,
		processIsStarted: processIsStarted,
		onEvent:          onEvent,
		prevState:        false,
	}
}

func (p *ProcessMonitor) Start(interval time.Duration) {
	p.isDone = false

	ticker := time.NewTicker(interval)
	p.prevState = p.processIsStarted()

	for ; true; <-ticker.C {
		if p.isDone {
			return
		}

		p.CompareStates()
	}
}

func (p *ProcessMonitor) CompareStates() {
	isStarted := p.processIsStarted()

	if isStarted && !p.prevState {
		p.onEvent(topic.EzquakeStarted)

	} else if !isStarted && p.prevState {
		p.onEvent(topic.EzquakeStopped)
	}

	p.prevState = isStarted
}

func (p *ProcessMonitor) Stop() {
	p.isDone = true
}
