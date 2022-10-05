package task_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/internal/pkg/task"
)

func TestPeriodicalTask(t *testing.T) {
	count := 0
	countTask := task.NewPeriodicalTask(func() { count++ })
	interval := 20 * time.Millisecond
	go countTask.Start(interval)
	go countTask.Start(interval) // calls to start while started should have no effect
	time.Sleep(time.Millisecond)

	time.Sleep(2 * interval)
	assert.Equal(t, 3, count)

	countTask.Stop()
	assert.Equal(t, 3, count)

	time.Sleep(interval)
	assert.Equal(t, 3, count)
}
