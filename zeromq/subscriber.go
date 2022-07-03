package zeromq

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
	"github.com/vikpe/streambot/zeromq/message"
)

type Subscriber struct {
	address   string
	topics    string
	onMessage message.Handler
}

func NewSubscriber(address string, topics string) Subscriber {
	return Subscriber{
		address: address,
		topics:  topics,
	}
}

func (s *Subscriber) Start(onMessage message.Handler) {
	go func() {
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
				msg, _ := message.NewMessageFromParts(zmqMsg)
				onMessage(msg)
			}
		}
	}()
}
