package zeromq_test

import (
	"os"
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
	var proxyStarted bool
	var proxyStopped bool
	proxy.OnStarted = func() { proxyStarted = true }
	proxy.OnStopped = func(sig os.Signal) { proxyStopped = true }

	go proxy.Start()
	zeromq.WaitForConnection()

	// subscriber
	wg := sync.WaitGroup{}
	messagesRecieved := make([]message.Message, 0)
	subscriber := zeromq.NewSubscriber("tcp://localhost:5556", zeromq.TopicsAll)
	var subscriberStarted bool
	var subscriberStopped bool
	subscriber.OnStarted = func() { subscriberStarted = true }
	subscriber.OnStopped = func(sig os.Signal) { subscriberStopped = true }

	go subscriber.Start(func(msg message.Message) {
		messagesRecieved = append(messagesRecieved, msg)

		if len(messagesRecieved) == len(messagesToSend) {
			proxy.Stop()
			subscriber.Stop()
			wg.Done()
		}
	})
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
	assert.True(t, proxyStarted)
	assert.True(t, proxyStopped)
	assert.True(t, subscriberStarted)
	assert.True(t, subscriberStopped)

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
