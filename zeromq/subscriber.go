package zeromq

import (
	"fmt"
	"os"
	"strings"

	zmq "github.com/pebbe/zmq4"
)

type Subscriber struct {
	address   string
	topics    []string
	onMessage MessageHandler
}

func NewSubscriber(address string, topics []string, onMessage MessageHandler) Subscriber {
	return Subscriber{
		address:   address,
		topics:    topics,
		onMessage: onMessage,
	}
}

func (s Subscriber) Start() {
	subSocket, _ := zmq.NewSocket(zmq.SUB)
	defer subSocket.Close()
	subSocket.Connect(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"))
	WaitForConnection()
	subSocket.SetSubscribe(strings.Join(s.topics, " "))

	for {
		zmqMsg, err := subSocket.RecvMessage(0)

		if err != nil {
			fmt.Println("Error recieving message", err)
		} else {
			msg := NewMessage(zmqMsg)
			s.onMessage(msg)
		}
	}
}
