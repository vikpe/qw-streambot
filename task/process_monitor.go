package task

import (
	"time"

	"github.com/vikpe/streambot/ezquake"
	"github.com/vikpe/streambot/topics"
)

type ProcessMonitor struct {
	isDone  bool
	process *ezquake.Process
	onEvent EventHandler
}

func NewProcessMonitor(process *ezquake.Process, onEvent EventHandler) ProcessMonitor {
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
		prevState := NewProcessState(*p.process)

		for ; true; <-ticker.C {
			if p.isDone {
				return
			}

			currentState := NewProcessState(*p.process)
			diff := NewProcessDiff(currentState, prevState)

			if diff.HasStarted {
				p.onEvent(topics.ClientStart, "")

			} else if diff.HasStopped {
				p.onEvent(topics.ClientStop, "")
			}

			prevState = currentState
		}
	}()
}

func (p *ProcessMonitor) Stop() {
	p.isDone = true
}

type ProcessState struct {
	IsStarted bool
}

func NewProcessState(process ezquake.Process) ProcessState {
	return ProcessState{
		IsStarted: process.IsStarted(),
	}
}

type ProcessDiff struct {
	HasStarted bool
	HasStopped bool
}

func NewProcessDiff(current ProcessState, prev ProcessState) ProcessDiff {
	diff := ProcessDiff{}

	if current.IsStarted && !prev.IsStarted {
		diff.HasStarted = true
	} else if !current.IsStarted && prev.IsStarted {
		diff.HasStopped = true
	}

	return diff
}
