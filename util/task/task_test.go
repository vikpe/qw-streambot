package task_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/util/task"
)

func TestPeriodicalTask(t *testing.T) {
	callCount := 0

	countTask := task.NewPeriodicalTask(func() {
		callCount++
	})
	interval := time.Millisecond
	countTask.Start(interval)
	time.Sleep(4 * interval)
	countTask.Stop()

	assert.Equal(t, 5, callCount)
}
