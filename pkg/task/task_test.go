package task_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/pkg/task"
)

func TestPeriodicalTask(t *testing.T) {
	count := 0
	countTask := task.NewPeriodicalTask(func() { count++ })
	interval := time.Millisecond
	countTask.Start(interval)
	time.Sleep(4 * interval)
	countTask.Stop()

	assert.Equal(t, 5, count)
}
