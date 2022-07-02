package zeromq

import (
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
	OnStop          func(os.Signal)
}

func NewProxy(frontend string, backend string) Proxy {
	return Proxy{
		frontendAddress: frontend,
		backendAddress:  backend,
		OnStart:         func() {},
		OnStop:          func(sig os.Signal) {},
	}
}

func (p *Proxy) Start() error {
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
			fmt.Println("unable to connect to frontend", err.Error())
			return
		}

		// backend - endpoint for subscribers
		backend, _ := zmq.NewSocket(zmq.XPUB)
		defer backend.Close()
		err = backend.Bind(p.backendAddress)

		if err != nil {
			fmt.Println("unable to bind to backend", err.Error())
			return
		}

		// run until interrupt
		p.OnStart()
		err = zmq.Proxy(frontend, backend, nil)

		if err != nil {
			fmt.Println("proxy interrupted:", err.Error())
		}
	}()
	sig := <-p.stopChan

	p.OnStop(sig)

	return err
}

func (p *Proxy) Stop() {
	if p.stopChan == nil {
		return
	}
	p.stopChan <- syscall.SIGINT
	time.Sleep(10 * time.Millisecond)
}
