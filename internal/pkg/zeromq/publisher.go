package zeromq

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

type Publisher struct {
	address string
}

func NewPublisher(address string) *Publisher {
	return &Publisher{address: address}
}

func (p *Publisher) SendMessage(topic string, content ...any) {
	var msgContent any

	if len(content) > 0 {
		msgContent = content[0]
	} else {
		msgContent = ""
	}

	msg := message.NewMessage(topic, msgContent)
	pubSendMessage(p.address, msg.Topic, msg.ContentType, string(msg.Content))
}

func pubSendMessage(address string, parts ...any) {
	pubSocket, _ := zmq.NewSocket(zmq.PUB)
	defer pubSocket.Close()
	err := pubSocket.Connect(address)
	if err != nil {
		fmt.Println("pubSendMessage: error connecting to pub socket", err)
		return
	}
	WaitForConnection()

	_, err = pubSocket.SendMessage(parts...)
	if err != nil {
		fmt.Println("pubSendMessage: error sending message", err)
		return
	}
}
