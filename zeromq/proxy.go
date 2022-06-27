package zeromq

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

type Proxy struct {
	frontendAddress string
	backendAddress  string
	OnStart         func()
}

func NewProxy(frontend string, backend string) Proxy {
	return Proxy{
		frontendAddress: frontend,
		backendAddress:  backend,
	}
}

func (p Proxy) Start() error {
	pp.Print("PROXY START")

	// frontend - endpoint for publishers
	frontend, _ := zmq.NewSocket(zmq.XSUB)
	defer frontend.Close()
	err := frontend.Bind(p.frontendAddress)

	if err != nil {
		fmt.Println("unable to connect to frontend")
		return err
	}

	// backend - endpoint for subscribers
	backend, _ := zmq.NewSocket(zmq.XPUB)
	defer backend.Close()
	err = backend.Bind(p.backendAddress)

	if err != nil {
		fmt.Println("unable to bind to backend")
		return err
	}

	// run until interrupt
	p.OnStart()
	err = zmq.Proxy(frontend, backend, nil)

	if err != nil {
		fmt.Println("proxy interrupted:", err)
		return err
	}

	return nil
}
