package zeromq

import (
	"errors"
	"fmt"

	zmq "github.com/pebbe/zmq4"
	"github.com/vikpe/streambot/internal/pkg/service"
)

func StartProxy(frontendAddress, backendAddress string) error {
	// frontend - endpoint for publishers
	frontend, _ := zmq.NewSocket(zmq.XSUB)
	defer frontend.Close()
	err := frontend.Bind(frontendAddress)

	if err != nil {
		return errors.New(fmt.Sprintf("unable to bind to frontend (%s)", err.Error()))
	}

	// backend - endpoint for subscribers
	backend, _ := zmq.NewSocket(zmq.XPUB)
	defer backend.Close()
	err = backend.Bind(backendAddress)

	if err != nil {
		return errors.New(fmt.Sprintf("unable to bind to backend (%s)", err.Error()))
	}

	// run until interrupt
	err = zmq.Proxy(frontend, backend, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("proxyService interrupted (%s)", err.Error()))
	}

	return nil
}

func NewProxyService(frontendAddress, backendAddress string) *service.Service {
	proxyService := service.New()
	proxyService.Work = func() error {
		return StartProxy(frontendAddress, backendAddress)
	}
	return proxyService
}
