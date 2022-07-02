package zeromq

import (
	"fmt"

	"github.com/goccy/go-json"
	zmq "github.com/pebbe/zmq4"
)

type Publisher struct {
	address string
}

func NewPublisher(address string) Publisher {
	return Publisher{address: address}
}

func (p Publisher) SendMessage(topic string, data any) {
	PubSendMessage(p.address, topic, data)
}

func PubSendMessage(address string, topic string, data any) {
	pubSocket, _ := zmq.NewSocket(zmq.PUB)
	defer pubSocket.Close()
	err := pubSocket.Connect(address)
	if err != nil {
		fmt.Println("error connecting to pub socket", err)
		return
	}
	WaitForConnection()

	dataAsJson, _ := json.Marshal(data)
	_, err = pubSocket.SendMessage(topic, dataAsJson, fmt.Sprintf("%T", data))
	if err != nil {
		fmt.Println("error sending message", err)
		return
	}
}
