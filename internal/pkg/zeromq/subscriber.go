package zeromq

import (
	"errors"
	"fmt"

	zmq "github.com/pebbe/zmq4"
	"github.com/vikpe/streambot/internal/pkg/service"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

type Subscriber struct {
	address   string
	topics    string
	OnMessage message.Handler
}

func NewSubscriber(address string, topics string) *Subscriber {
	return &Subscriber{
		address:   address,
		topics:    topics,
		OnMessage: func(msg message.Message) {},
	}
}

func (s *Subscriber) Start() error {
	subSocket, _ := zmq.NewSocket(zmq.SUB)
	defer subSocket.Close()
	subSocket.Connect(s.address)
	WaitForConnection()
	subSocket.SetSubscribe(s.topics)

	for {
		zmqMsg, err := subSocket.RecvMessage(0)

		if err != nil {
			return errors.New(fmt.Sprintf("Error recieving message: %s", err))
		}

		msg, err := message.NewMessageFromFrames(zmqMsg)

		if err != nil {
			return err
		}

		s.OnMessage(msg)
	}
}

type SubscriberService struct {
	*Subscriber
	*service.Service
}

func NewSubscriberService(address, topics string) *SubscriberService {
	sub := NewSubscriber(address, topics)
	subService := service.New()
	subService.Work = func() error {
		return sub.Start()
	}

	return &SubscriberService{
		Subscriber: sub,
		Service:    subService,
	}
}
