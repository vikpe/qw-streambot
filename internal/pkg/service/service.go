package service

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Service struct {
	stopChan  chan os.Signal
	Work      func() error
	OnStarted func()
	OnError   func(error)
	OnStopped func(signal os.Signal)
}

func New() *Service {
	return &Service{
		stopChan:  make(chan os.Signal, 1),
		Work:      func() error { return nil },
		OnStarted: func() {},
		OnError:   func(error) {},
		OnStopped: func(os.Signal) {},
	}
}

func (s *Service) Start() {
	// catch SIGETRM and SIGINTERRUPT
	s.stopChan = make(chan os.Signal, 1)
	signal.Notify(s.stopChan, syscall.SIGTERM, syscall.SIGINT)

	var err error

	go func() {
		err = s.Work()
	}()
	s.OnStarted()
	stopSignal := <-s.stopChan

	if err != nil {
		s.OnError(err)
	}

	s.OnStopped(stopSignal)
}

func (s *Service) Stop() {
	s.stopChan <- syscall.SIGINT
	time.Sleep(10 * time.Millisecond) //
}
