package monitor

import (
	"time"

	"github.com/vikpe/streambot/pkg/zeromq"
	"github.com/vikpe/streambot/topic"
)

type ProcessMonitor struct {
	isDone       bool
	getIsStarted func() bool
	onEvent      zeromq.EventHandler
	prevState    bool
}

func NewProcessMonitor(getIsStarted func() bool, onEvent zeromq.EventHandler) ProcessMonitor {
	return ProcessMonitor{
		isDone:       false,
		getIsStarted: getIsStarted,
		onEvent:      onEvent,
		prevState:    false,
	}
}

func (p *ProcessMonitor) Start(interval time.Duration) {
	p.isDone = false

	go func() {
		ticker := time.NewTicker(interval)
		p.prevState = p.getIsStarted()

		for ; true; <-ticker.C {
			if p.isDone {
				return
			}

			p.CompareStates()
		}
	}()
}

func (p *ProcessMonitor) CompareStates() {
	isStarted := p.getIsStarted()

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
