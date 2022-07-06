package zeromq_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	zeromq2 "github.com/vikpe/streambot/pkg/zeromq"
	"github.com/vikpe/streambot/pkg/zeromq/message"
)

func TestEndToEnd(t *testing.T) {
	wg := sync.WaitGroup{}

	topicsToSend := []string{"domain.topic1", "domain.topic2", "domain.topic3"}

	// proxy
	go func() {
		proxy := zeromq2.NewProxy("tcp://*:5555", "tcp://*:5556")
		proxy.Start()
	}()
	zeromq2.WaitForConnection()

	// subscriber
	topicsReceived := make([]string, 0)

	go func() {
		subscriber := zeromq2.NewSubscriber("tcp://localhost:5556", zeromq2.TopicsAll)
		subscriber.Start(func(message message.Message) {
			topicsReceived = append(topicsReceived, message.Topic)

			if len(topicsReceived) == len(topicsToSend) {
				wg.Done()
			}
		})
	}()
	zeromq2.WaitForConnection()

	// publisher
	go func() {
		publisher := zeromq2.NewPublisher("tcp://localhost:5555")

		for _, topic := range topicsToSend {
			publisher.SendMessage(topic)
		}
	}()

	wg.Add(1)
	wg.Wait()

	// assertions
	assert.Equal(t, topicsToSend, topicsReceived)
}
