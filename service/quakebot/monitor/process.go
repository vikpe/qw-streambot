package monitor

import (
	"time"

	"github.com/vikpe/streambot/pkg/topic"
	"github.com/vikpe/streambot/pkg/zeromq"
)

type ProcessMonitor struct {
	isDone           bool
	processIsStarted func() bool
	onEvent          zeromq.EventHandler
	prevState        bool
}

func NewProcessMonitor(processIsStarted func() bool, onEvent zeromq.EventHandler) ProcessMonitor {
	return ProcessMonitor{
		isDone:           false,
		processIsStarted: processIsStarted,
		onEvent:          onEvent,
		prevState:        false,
	}
}

func (p *ProcessMonitor) Start(interval time.Duration) {
	p.isDone = false

	go func() {
		ticker := time.NewTicker(interval)
		p.prevState = p.processIsStarted()

		for ; true; <-ticker.C {
			if p.isDone {
				return
			}

			p.CompareStates()
		}
	}()
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
