package proc_test

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/pkg/proc"
	"github.com/vikpe/streambot/pkg/proc/shell/mock"
)

func TestProcess_GetProcessID(t *testing.T) {
	t.Run("no process found", func(t *testing.T) {
		exec := mock.NewExecMock()
		exec.Output["pgrep"] = ""
		process := proc.NewProcessController("/home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64")
		process.ExecCommand = exec.Command

		assert.Equal(t, 0, process.GetID())
		assert.False(t, process.IsStarted())
		assert.True(t, process.IsStopped())
	})

	t.Run("process found", func(t *testing.T) {
		exec := mock.NewExecMock()
		exec.Output["pgrep"] = "1818481"
		process := proc.NewProcessController("/home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64")
		process.ExecCommand = exec.Command

		assert.Equal(t, 1818481, process.GetID())
		assert.True(t, process.IsStarted())
		assert.False(t, process.IsStopped())
	})
}

func TestProcess_Stop(t *testing.T) {
	t.Run("not started", func(t *testing.T) {
		exec := mock.NewExecMock()
		exec.Output["pgrep"] = ""
		process := proc.NewProcessController("/home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64")
		process.ExecCommand = exec.Command

		process.Stop(syscall.SIGTERM)
		assert.Equal(t, []string{
			"pgrep -fo /home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64",
		}, exec.Calls)
	})

	t.Run("started", func(t *testing.T) {
		exec := mock.NewExecMock()
		exec.Output["pgrep"] = "1818481"
		process := proc.NewProcessController("/home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64")
		process.ExecCommand = exec.Command

		process.Stop(syscall.SIGTERM)
		assert.Equal(t, []string{
			"pgrep -fo /home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64",
			"kill -s 15 1818481",
		}, exec.Calls)
	})
}
