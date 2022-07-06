package ezquake_test

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/pkg/ezquake"
	"github.com/vikpe/streambot/pkg/ezquake/shell/test_helpers"
)

func TestProcess_GetProcessID(t *testing.T) {
	t.Run("no process found", func(t *testing.T) {
		exec := test_helpers.NewExecMock()
		exec.Output["pgrep"] = ""
		process := ezquake.NewProcess("/home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64")
		process.ExecCommand = exec.Command

		assert.Equal(t, 0, process.ID())
		assert.False(t, process.IsStarted())
		assert.True(t, process.IsStopped())
	})

	t.Run("process found", func(t *testing.T) {
		exec := test_helpers.NewExecMock()
		exec.Output["pgrep"] = "1818481"
		process := ezquake.NewProcess("/home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64")
		process.ExecCommand = exec.Command

		assert.Equal(t, 1818481, process.ID())
		assert.True(t, process.IsStarted())
		assert.False(t, process.IsStopped())
	})
}

func TestProcess_Stop(t *testing.T) {
	t.Run("not started", func(t *testing.T) {
		exec := test_helpers.NewExecMock()
		exec.Output["pgrep"] = ""
		process := ezquake.NewProcess("/home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64")
		process.ExecCommand = exec.Command

		process.Stop(syscall.SIGTERM)
		assert.Equal(t, []string{
			"pgrep -fo /home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64",
		}, exec.Calls)
	})

	t.Run("started", func(t *testing.T) {
		exec := test_helpers.NewExecMock()
		exec.Output["pgrep"] = "1818481"
		process := ezquake.NewProcess("/home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64")
		process.ExecCommand = exec.Command

		process.Stop(syscall.SIGTERM)
		assert.Equal(t, []string{
			"pgrep -fo /home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64",
			"kill -s 15 1818481",
		}, exec.Calls)
	})
}
