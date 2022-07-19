package zeromq

import (
	"errors"
	"fmt"

	zmq "github.com/pebbe/zmq4"
	"github.com/vikpe/streambot/internal/pkg/service"
)

type Proxy struct {
	frontendAddress string
	backendAddress  string
}

func NewProxy(frontendAddress, backendAddress string) *Proxy {
	return &Proxy{
		frontendAddress: frontendAddress,
		backendAddress:  backendAddress,
	}
}

func (p *Proxy) Start() error {
	// frontend - endpoint for publishers
	frontend, _ := zmq.NewSocket(zmq.XSUB)
	defer frontend.Close()
	err := frontend.Bind(p.frontendAddress)

	if err != nil {
		return errors.New(fmt.Sprintf("unable to bind to frontend (%s)", err.Error()))
	}

	// backend - endpoint for subscribers
	backend, _ := zmq.NewSocket(zmq.XPUB)
	defer backend.Close()
	err = backend.Bind(p.backendAddress)

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

type ProxyService struct {
	*Proxy
	*service.Service
}

func NewProxyService(frontendAddress, backendAddress string) *ProxyService {
	proxy := NewProxy(frontendAddress, backendAddress)

	proxyService := service.New()
	proxyService.Work = func() error {
		return proxy.Start()
	}

	return &ProxyService{
		Proxy:   proxy,
		Service: proxyService,
	}
}
