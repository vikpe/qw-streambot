package ezquake

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/vikpe/streambot/pkg/ezquake/shell"
)

type Process struct {
	Path        string
	ExecCommand func(command string) string
}

func NewProcess(path string) Process {
	return Process{
		Path:        path,
		ExecCommand: shell.ExecCommand,
	}
}

func (p Process) ID() int {
	pregCmd := fmt.Sprintf("pgrep -fo %s", p.Path)
	prepOutput := p.ExecCommand(pregCmd)
	id, err := strconv.Atoi(prepOutput)
	if err != nil {
		return 0
	}

	return id
}

func (p Process) Stop(signal syscall.Signal) {
	pid := p.ID()

	if 0 == pid {
		return
	}

	killCmd := fmt.Sprintf("kill -s %d %d", int(signal), pid)
	p.ExecCommand(killCmd)
}

func (p Process) IsStarted() bool {
	return 0 != p.ID()
}

func (p Process) IsStopped() bool {
	return !p.IsStarted()
}
