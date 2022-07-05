package monitor

import (
	"time"

	"github.com/vikpe/streambot/ezquake"
	"github.com/vikpe/streambot/zeromq"
	"github.com/vikpe/streambot/zeromq/topic"
)

type ProcessMonitor struct {
	isDone  bool
	process *ezquake.Process
	onEvent zeromq.EventHandler
}

func NewProcessMonitor(process *ezquake.Process, onEvent zeromq.EventHandler) ProcessMonitor {
	return ProcessMonitor{
		isDone:  false,
		process: process,
		onEvent: onEvent,
	}
}

func (p *ProcessMonitor) Start(interval time.Duration) {
	p.isDone = false

	go func() {
		ticker := time.NewTicker(interval)
		prevIsStarted := p.process.IsStarted()

		for ; true; <-ticker.C {
			if p.isDone {
				return
			}

			isStarted := p.process.IsStarted()

			if isStarted && !prevIsStarted {
				p.onEvent(topic.EzquakeStarted)

			} else if !isStarted && prevIsStarted {
				p.onEvent(topic.EzquakeStopped)
			}

			prevIsStarted = isStarted
		}
	}()
}

func (p *ProcessMonitor) Stop() {
	p.isDone = true
}
