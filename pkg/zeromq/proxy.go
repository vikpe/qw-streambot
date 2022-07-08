package zeromq

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	zmq "github.com/pebbe/zmq4"
)

type Proxy struct {
	frontendAddress string
	backendAddress  string
	stopChan        chan os.Signal
	OnStart         func()
	OnError         func(error)
	OnStop          func(os.Signal)
}

func NewProxy(frontend string, backend string) Proxy {
	return Proxy{
		frontendAddress: frontend,
		backendAddress:  backend,
		OnStart:         func() {},
		OnError:         func(err error) {},
		OnStop:          func(sig os.Signal) {},
	}
}

func (p *Proxy) Start() {
	// catch SIGETRM and SIGINTERRUPT
	p.stopChan = make(chan os.Signal, 1)
	signal.Notify(p.stopChan, syscall.SIGTERM, syscall.SIGINT)

	var err error

	go func() {
		// frontend - endpoint for publishers
		frontend, _ := zmq.NewSocket(zmq.XSUB)
		defer frontend.Close()
		err = frontend.Bind(p.frontendAddress)

		if err != nil {
			err = errors.New(fmt.Sprintf("unable to bind to frontend (%s)", err.Error()))
			return
		}

		// backend - endpoint for subscribers
		backend, _ := zmq.NewSocket(zmq.XPUB)
		defer backend.Close()
		err = backend.Bind(p.backendAddress)

		if err != nil {
			err = errors.New(fmt.Sprintf("unable to bind to backend (%s)", err.Error()))
			return
		}

		// run until interrupt
		p.OnStart()
		err = zmq.Proxy(frontend, backend, nil)

		if err != nil {
			err = errors.New(fmt.Sprintf("proxy interrupted: (%s)", err.Error()))
			return
		}
	}()
	sig := <-p.stopChan

	if err != nil {
		p.OnError(err)
	}

	p.OnStop(sig)
}

func (p *Proxy) Stop() {
	if p.stopChan == nil {
		return
	}
	p.stopChan <- syscall.SIGINT
	time.Sleep(10 * time.Millisecond)
}
