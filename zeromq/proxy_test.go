package zeromq_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/zeromq"
)

func TestEndToEnd(t *testing.T) {
	wg := sync.WaitGroup{}

	topicsToSend := []string{"alpha", "beta", "gamma"}

	// proxy
	go func() {
		zeromq.NewProxy("tcp://*:5555", "tcp://*:5556").Start()
	}()
	zeromq.WaitForConnection()

	// subscriber
	topicsReceived := make([]string, 0)

	go func() {
		subscriber := zeromq.NewSubscriber("tcp://localhost:5556", zeromq.TopicsAll, func(message zeromq.Message) {
			topicsReceived = append(topicsReceived, message.Topic)

			if len(topicsReceived) == len(topicsToSend) {
				wg.Done()
			}
		})
		subscriber.Start()
	}()
	zeromq.WaitForConnection()

	// publisher
	go func() {
		publisher := zeromq.NewPublisher("tcp://localhost:5555")

		for _, topic := range topicsToSend {
			publisher.SendMessage(topic, "")
		}
	}()

	wg.Add(1)
	wg.Wait()

	// assertions
	assert.Equal(t, topicsToSend, topicsReceived)
}
