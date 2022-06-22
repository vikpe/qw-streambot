package main

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

type ZmqProxy struct {
	frontendAddress string
	backendAddress  string
}

func (p ZmqProxy) Start() {
	frontend, _ := zmq.NewSocket(zmq.XSUB)
	defer frontend.Close()
	frontend.Connect(p.frontendAddress)

	backend, _ := zmq.NewSocket(zmq.XPUB)
	defer backend.Close()
	backend.Bind(p.backendAddress)

	// run until interrupt
	fmt.Println("proxy started")
	err := zmq.Proxy(frontend, backend, nil)
	fmt.Println("proxy interrupted:", err)
}
