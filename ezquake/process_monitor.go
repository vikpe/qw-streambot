package ezquake

import (
	"fmt"
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

func (p *ProcessMonitor) Start() {
	p.isDone = false
	ticker := time.NewTicker(p.interval)
	prevState := NewProcessState(*p.process)

	go func() {
		for ; true; <-ticker.C {
			if p.isDone {
				return
			}

			currentState := NewProcessState(*p.process)
			diff := NewProcessDiff(currentState, prevState)

			if diff.HasStarted {
				p.eventHandler(EventProcessStart, "")

			} else if diff.HasStopped {
				p.eventHandler(EventProcessStop, "")
			}

			prevState = currentState
		}
	}()
}

func (p *ProcessMonitor) Stop() {
	p.isDone = true
	fmt.Println("yes we stopped..")
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
