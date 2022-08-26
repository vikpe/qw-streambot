package zeromq

import (
	"errors"
	"fmt"

	zmq "github.com/pebbe/zmq4"
	"github.com/vikpe/streambot/internal/pkg/service"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

func StartSubscriber(address, topics string, onMessage message.Handler) error {
	subSocket, _ := zmq.NewSocket(zmq.SUB)
	defer subSocket.Close()
	subSocket.Connect(address)
	WaitForConnection()
	subSocket.SetSubscribe(topics)

	for {
		zmqMsg, err := subSocket.RecvMessage(0)

		if err != nil {
			return errors.New(fmt.Sprintf("Error receiving message: %s", err))
		}

		msg, err := message.NewMessageFromFrames(zmqMsg)

		if err != nil {
			return err
		}

		onMessage(msg)
	}
}

type Subscriber struct {
	*service.Service
	OnMessage message.Handler
}

func NewSubscriber(address, topics string) *Subscriber {
	sub := Subscriber{
		Service:   service.New(),
		OnMessage: func(msg message.Message) {},
	}
	sub.Service.Work = func() error {
		return StartSubscriber(address, topics, sub.OnMessage)
	}

	return &sub
}
