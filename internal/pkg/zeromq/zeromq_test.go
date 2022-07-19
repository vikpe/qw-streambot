package zeromq_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

func TestEndToEnd(t *testing.T) {
	type outMessage struct {
		Topic   string
		Content any
	}

	messagesToSend := []outMessage{
		{"domain.topic1", "hello world"},
		{"domain.topic2", []string{"hello", "world"}},
		{"domain.topic3", 666},
	}

	// proxy
	proxy := zeromq.NewProxy("tcp://*:5555", "tcp://*:5556")
	go proxy.Start()
	zeromq.WaitForConnection()

	// subscriber
	wg := sync.WaitGroup{}
	messagesRecieved := make([]message.Message, 0)
	subscriber := zeromq.NewSubscriberService("tcp://localhost:5556", zeromq.TopicsAll)
	subscriber.OnMessage = func(msg message.Message) {
		messagesRecieved = append(messagesRecieved, msg)

		if len(messagesRecieved) == len(messagesToSend) {
			proxy.Stop()
			subscriber.Service.Stop()
			wg.Done()
		}
	}

	go subscriber.Service.Start()
	zeromq.WaitForConnection()

	// publisher
	go func() {
		publisher := zeromq.NewPublisher("tcp://localhost:5555")

		for _, msg := range messagesToSend {
			publisher.SendMessage(msg.Topic, msg.Content)
		}
	}()

	wg.Add(1)
	wg.Wait()

	// assertions
	// message 1
	assert.Equal(t, messagesToSend[0].Topic, messagesRecieved[0].Topic)
	assert.Equal(t, messagesToSend[0].Content, messagesRecieved[0].Content.ToString())

	// message 2
	assert.Equal(t, messagesToSend[1].Topic, messagesRecieved[1].Topic)
	var message2Content []string
	messagesRecieved[1].Content.To(&message2Content)
	assert.Equal(t, messagesToSend[1].Content, message2Content)

	// message 3
	assert.Equal(t, messagesToSend[2].Topic, messagesRecieved[2].Topic)
	var message3Content int
	messagesRecieved[2].Content.To(&message3Content)
	assert.Equal(t, messagesToSend[2].Content, message3Content)
}
