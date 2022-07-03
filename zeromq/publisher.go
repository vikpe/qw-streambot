package zeromq

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
	"github.com/vikpe/streambot/zeromq/message"
)

type Publisher struct {
	address string
}

func NewPublisher(address string) Publisher {
	return Publisher{address: address}
}

func (p Publisher) SendMessage(topic string, content ...any) {
	var msgContent any

	if len(content) > 0 {
		msgContent = content[0]
	} else {
		msgContent = ""
	}

	msg := message.NewMessage(topic, msgContent)
	PubSendMessage(p.address, msg)
}

func PubSendMessage(address string, msg message.Message) {
	pubSocket, _ := zmq.NewSocket(zmq.PUB)
	defer pubSocket.Close()
	err := pubSocket.Connect(address)
	if err != nil {
		fmt.Println("PubSendMessage: error connecting to pub socket", err)
		return
	}
	WaitForConnection()

	_, err = pubSocket.SendMessage(msg.Topic, msg.ContentType, msg.Content)
	if err != nil {
		fmt.Println("PubSendMessage: error sending message", err)
		return
	}
}
