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
	pubSocket.Connect(address)
	WaitForConnection()

	dataAsJson, _ := json.Marshal(data)
	pubSocket.SendMessage(topic, dataAsJson, fmt.Sprintf("%T", data))
}
