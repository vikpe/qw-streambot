package ezquake

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"github.com/vikpe/streambot/util/shell"
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

func (p Process) TcpAddress() string {
	pid := p.ID()

	if 0 == pid {
		return ""
	}

	pidNeedle := fmt.Sprintf("pid=%d", pid)
	ssCmd := "ss -nptH -o state established dport eq 28000"
	ssOutput := p.ExecCommand(ssCmd)
	// 0          0              192.168.2.194:41706           46.227.68.148:28000      users:(("ezquake-linux-x",pid=1818481,fd=53))

	if !strings.Contains(ssOutput, pidNeedle) {
		return ""
	}

	ssFields := strings.Fields(ssOutput)
	const destAddressIndex = 3
	return ssFields[destAddressIndex]
}
