package zeromq

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

type Subscriber struct {
	address   string
	topics    string
	onMessage message.Handler
	stopChan  chan os.Signal
	OnStarted func()
	OnError   func(error)
	OnStopped func(os.Signal)
}

func NewSubscriber(address string, topics string) *Subscriber {
	return &Subscriber{
		address:   address,
		topics:    topics,
		stopChan:  make(chan os.Signal, 1),
		OnStarted: func() {},
		OnError:   func(err error) {},
		OnStopped: func(sig os.Signal) {},
	}
}

func (s *Subscriber) Start(onMessage message.Handler) {
	// catch SIGETRM and SIGINTERRUPT
	s.stopChan = make(chan os.Signal, 1)
	signal.Notify(s.stopChan, syscall.SIGTERM, syscall.SIGINT)

	var err error

	go func() {
		subSocket, _ := zmq.NewSocket(zmq.SUB)
		defer subSocket.Close()
		subSocket.Connect(s.address)
		WaitForConnection()
		subSocket.SetSubscribe(s.topics)

		s.OnStarted()

		for {
			zmqMsg, err := subSocket.RecvMessage(0)

			if err != nil {
				fmt.Println("Error recieving message", err)
			} else {
				msg, err := message.NewMessageFromFrames(zmqMsg)

				if err != nil {
					fmt.Println(err)
					return
				}

				onMessage(msg)
			}
		}
	}()
	sig := <-s.stopChan

	if err != nil {
		s.OnError(err)
	}

	s.OnStopped(sig)
}

func (s *Subscriber) Stop() {
	if s.stopChan == nil {
		return
	}
	s.stopChan <- syscall.SIGINT
	time.Sleep(10 * time.Millisecond)
}
