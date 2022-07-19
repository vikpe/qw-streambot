package zeromq

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

type Subscriber struct {
	address   string
	topics    string
	OnMessage message.Handler
}

func NewSubscriber(address string, topics string) *Subscriber {
	return &Subscriber{
		address: address,
		topics:  topics,
	}
}

func (s *Subscriber) Start() {
	subSocket, _ := zmq.NewSocket(zmq.SUB)
	defer subSocket.Close()
	subSocket.Connect(s.address)
	WaitForConnection()
	subSocket.SetSubscribe(s.topics)

	for {
		zmqMsg, err := subSocket.RecvMessage(0)

		if err != nil {
			fmt.Println("Error recieving message", err)
		} else {
			msg, err := message.NewMessageFromFrames(zmqMsg)

			if err != nil {
				fmt.Println(err)
				return
			}

			s.OnMessage(msg)
		}
	}
}
