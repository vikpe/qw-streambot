package monitor

import (
	"time"

	"github.com/vikpe/streambot/ezquake"
	"github.com/vikpe/streambot/zeromq/topic"
)

type ProcessMonitor struct {
	isDone  bool
	process *ezquake.Process
	onEvent func(string, any)
}

func NewProcessMonitor(process *ezquake.Process, onEvent func(string, any)) ProcessMonitor {
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
		prevState := newProcessState(*p.process)

		for ; true; <-ticker.C {
			if p.isDone {
				return
			}

			currentState := newProcessState(*p.process)
			diff := newProcessDiff(currentState, prevState)

			if diff.HasStarted {
				p.onEvent(topic.EzquakeStarted, "")

			} else if diff.HasStopped {
				p.onEvent(topic.EzquakeStopped, "")
			}

			prevState = currentState
		}
	}()
}

func (p *ProcessMonitor) Stop() {
	p.isDone = true
}

type processState struct {
	IsStarted bool
}

func newProcessState(process ezquake.Process) processState {
	return processState{
		IsStarted: process.IsStarted(),
	}
}

type processDiff struct {
	HasStarted bool
	HasStopped bool
}

func newProcessDiff(current processState, prev processState) processDiff {
	diff := processDiff{}

	if current.IsStarted && !prev.IsStarted {
		diff.HasStarted = true
	} else if !current.IsStarted && prev.IsStarted {
		diff.HasStopped = true
	}

	return diff
}
