package zeromq

import (
	"os"

	zmq "github.com/pebbe/zmq4"
)

type Subscriber struct {
	address   string
	topics    string
	onMessage EventHandler
}

func NewSubscriber(address string, topics string, onEvent EventHandler) Subscriber {
	return Subscriber{
		address:   address,
		topics:    topics,
		onMessage: onEvent,
	}
}

func (s Subscriber) Start() {
	subSocket, _ := zmq.NewSocket(zmq.SUB)
	defer subSocket.Close()
	subSocket.Connect(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"))
	WaitForConnection()
	subSocket.SetSubscribe(s.topics)

	for {
		if rawMsg, err := subSocket.RecvMessage(0); err != nil {
			pp.Print("Error recieving message", err)
		} else {
			msg := ParseEvent(rawMsg)
			s.onMessage(msg)
		}
	}
}
