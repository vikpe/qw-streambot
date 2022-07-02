package zeromq

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	zmq "github.com/pebbe/zmq4"
)

type Proxy struct {
	frontendAddress string
	backendAddress  string
	OnStart         func()
	OnStop          func(os.Signal)
	OnError         func(error)
}

func NewProxy(frontend string, backend string) Proxy {
	return Proxy{
		frontendAddress: frontend,
		backendAddress:  backend,
		OnStart:         func() {},
		OnStop:          func(sig os.Signal) {},
		OnError:         func(err error) {},
	}
}

func (p Proxy) Start() error {
	// catch SIGETRM and SIGINTERRUPT
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

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
			return
		}
	}()
	sig := <-cancelChan

	if err != nil {
		p.OnError(err)
	}

	p.OnStop(sig)

	return err
}
