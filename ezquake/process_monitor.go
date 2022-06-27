package ezquake

import (
	"time"
)

type eventHandler func(string, any)

type ProcessMonitor struct {
	isDone       bool
	process      *Process
	eventHandler eventHandler
	interval     time.Duration
}

func NewProcessMonitor(process *Process, eventHandler eventHandler, interval time.Duration) ProcessMonitor {
	return ProcessMonitor{
		isDone:       false,
		process:      process,
		eventHandler: eventHandler,
		interval:     interval,
	}
}

func (t *ProcessMonitor) Start() {
	t.isDone = false
	ticker := time.NewTicker(t.interval)
	prevState := NewProcessState(*t.process)

	go func() {
		for ; true; <-ticker.C {
			if t.isDone {
				return
			}

			currentState := NewProcessState(*t.process)
			diff := NewProcessDiff(currentState, prevState)

			if diff.HasStarted {
				t.eventHandler(EventProcessStart, "")

			} else if diff.HasStopped {
				t.eventHandler(EventProcessStop, "")
			}

			prevState = currentState
		}
	}()
}

func (t *ProcessMonitor) Stop() {
	t.isDone = true
}

type ProcessState struct {
	IsStarted bool
}

func NewProcessState(process Process) ProcessState {
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
