package service_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/internal/pkg/service"
)

func TestService(t *testing.T) {
	t.Run("error in Work()", func(t *testing.T) {
		myService := service.New()
		myService.Work = func() error { return errors.New("fail") }

		var workError error
		myService.OnError = func(err error) { workError = err }

		go myService.Start()
		time.Sleep(time.Millisecond * 10)
		myService.Stop()

		assert.ErrorContains(t, workError, "fail")
	})

	t.Run("no error in Work()", func(t *testing.T) {
		myService := service.New()
		myService.Work = func() error { return nil }

		var workError error
		myService.OnError = func(err error) { workError = err }

		go myService.Start()
		time.Sleep(time.Millisecond * 10)
		myService.Stop()

		assert.Nil(t, workError)
	})

	t.Run("callbacks (OnStarted/OnStopped)", func(t *testing.T) {
		hasStarted := false
		hasStopped := false

		myService := service.New()
		myService.OnStarted = func() { hasStarted = true }
		myService.OnStopped = func(signal os.Signal) { hasStopped = true }

		go myService.Start()
		time.Sleep(time.Millisecond * 10)
		myService.Stop()

		assert.True(t, hasStarted)
		assert.True(t, hasStopped)
	})
}
