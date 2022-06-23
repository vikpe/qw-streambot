package zeromq

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/goccy/go-json"
	zmq "github.com/pebbe/zmq4"
	"github.com/vikpe/streambot/util/term"
)

var pp = term.NewPrettyPrinter("zmq", color.FgHiMagenta)

const ConnectionTimeout = time.Millisecond * 10

func WaitForConnection() {
	time.Sleep(ConnectionTimeout)
}

type Proxy struct {
	frontendAddress string
	backendAddress  string
}

func NewProxy(frontend string, backend string) Proxy {
	pp.Print("NEWPROXY", frontend, backend)

	return Proxy{
		frontendAddress: frontend,
		backendAddress:  backend,
	}
}

func (p Proxy) Start() error {
	pp.Print("PROXY START")

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

func (p Publisher) SendMessage(topic string, data any) {
	PubSendMessage(p.address, topic, data)
}

func PubSendMessage(address string, topic string, data any) {
	pubSocket, _ := zmq.NewSocket(zmq.PUB)
	defer pubSocket.Close()
	pubSocket.Connect(address)
	WaitForConnection()

	dataAsJson, _ := json.Marshal(data)
	pubSocket.SendMessage(topic, dataAsJson, fmt.Sprintf("%T", data))
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
		if rawMsg, err := subSocket.RecvMessage(0); err != nil {
			pp.Print("Error recieving message", err)
		} else {
			var topic string
			var dataType string
			var data any

			topic = rawMsg[0]

			if 3 == len(rawMsg) {
				dataType = rawMsg[2]
				data = rawMsg[1]
			} else {
				if 2 == len(rawMsg) {
					data = rawMsg[1]
				} else {
					data = ""
				}
				dataType = "string"
			}

			pp.Print(topic, fmt.Sprintf("(%s)", dataType), data)
		}
	}
}

type Message struct {
	Topic string
	Data  string
}

func (m Message) String() string {
	var target string
	m.To(target)
	return target
}

func (m Message) Int() int {
	var target int
	m.To(target)
	return target
}

func (m Message) To(target interface{}) {
	json.Unmarshal([]byte(m.Data), &target)
}
