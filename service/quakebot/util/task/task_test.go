package task_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/service/quakebot/util/task"
)

func TestPeriodicalTask(t *testing.T) {
	count := 0
	countTask := task.NewPeriodicalTask(func() { count++ })
	interval := time.Millisecond * 20
	countTask.Start(interval)
	time.Sleep(4 * interval)
	countTask.Stop()

	assert.Equal(t, 5, count)
}
