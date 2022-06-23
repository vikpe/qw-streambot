package zeromq

import (
	"fmt"
	"os"
	"time"

	zmq "github.com/pebbe/zmq4"
)

const ConnectionTimeout = time.Millisecond * 10

func WaitForConnection() {
	time.Sleep(ConnectionTimeout)
}

type Proxy struct {
	frontendAddress string
	backendAddress  string
}

func NewProxy(frontend string, backend string) Proxy {
	fmt.Println("NEWPROXY", frontend, backend)

	return Proxy{
		frontendAddress: frontend,
		backendAddress:  backend,
	}
}

func (p Proxy) Start() error {
	fmt.Println("PROXY START")

	// frontend - endpoint for publishers
	frontend, _ := zmq.NewSocket(zmq.XSUB)
	defer frontend.Close()
	err := frontend.Bind(p.frontendAddress)

	if err != nil {
		fmt.Println("unable to connect to frontend")
		return err
	}

	// backend - endpoint for subscribers
	backend, _ := zmq.NewSocket(zmq.XPUB)
	defer backend.Close()
	err = backend.Bind(p.backendAddress)

	if err != nil {
		fmt.Println("unable to bind to backend")
		return err
	}

	// run until interrupt
	fmt.Println("proxy started")
	err = zmq.Proxy(frontend, backend, nil)

	if err != nil {
		fmt.Println("proxy interrupted:", err)
		return err
	}

	return nil
}

type Publisher struct {
	address string
}

func NewPublisher(address string) Publisher {
	return Publisher{address: address}
}

func (p Publisher) SendMessage(message ...string) {
	PubSendMessage(p.address, message...)
}

func PubSendMessage(address string, message ...string) {
	pubSocket, _ := zmq.NewSocket(zmq.PUB)
	defer pubSocket.Close()
	pubSocket.Connect(address)
	WaitForConnection()

	pubSocket.SendMessage(message)
}

type Subscriber struct {
	address string
	topics  string
}

func NewSubscriber(address string, topics string) Subscriber {
	return Subscriber{
		address: address,
		topics:  topics,
	}
}

func (s Subscriber) Start() {
	subSocket, _ := zmq.NewSocket(zmq.SUB)
	defer subSocket.Close()
	subSocket.Connect(os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"))
	WaitForConnection()

	subSocket.SetSubscribe(s.topics)

	for {
		if msg, err := subSocket.RecvMessage(0); err != nil {
			fmt.Println("Error recieving message", err)
		} else {
			fmt.Println("Received message:", msg)
		}
	}
}

/*func EncodeMessage(message string) []byte {
	result, err := json.Marshal(message)
	if err != nil {
		return nil
	}
	return result
}

func DecodeMessage(encodedMessage string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(encodedMessage, &result)
	if err != nil {
		return make(map[string]interface{}, 0)
	}
	return result
}*/
