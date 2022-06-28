package zeromq

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
	zmq "github.com/pebbe/zmq4"
)

const TopicsAll = ""
const ConnectionTimeout = time.Millisecond * 10

func WaitForConnection() {
	time.Sleep(ConnectionTimeout)
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
