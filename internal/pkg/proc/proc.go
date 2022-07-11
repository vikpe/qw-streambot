package proc

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/vikpe/streambot/internal/pkg/exc"
)

type ProcessController struct {
	Path        string
	ExecCommand func(command string) string
}

func NewProcessController(path string) *ProcessController {
	return &ProcessController{
		Path:        path,
		ExecCommand: exc.GetOutput,
	}
}

func (p ProcessController) GetID() int {
	pregCmd := fmt.Sprintf("pgrep -fo %s", p.Path)
	prepOutput := p.ExecCommand(pregCmd)
	id, err := strconv.Atoi(prepOutput)
	if err != nil {
		return 0
	}

	return id
}

func (p ProcessController) Stop(signal syscall.Signal) {
	pid := p.GetID()

	if 0 == pid {
		return
	}

	killCmd := fmt.Sprintf("kill -s %d %d", int(signal), pid)
	p.ExecCommand(killCmd)
}

func (p ProcessController) IsStarted() bool {
	return 0 != p.GetID()
}

func (p ProcessController) IsStopped() bool {
	return !p.IsStarted()
}
