package chanman

import (
	"sync"

	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
	"github.com/vikpe/streambot/internal/third_party/twitch"
	"github.com/vikpe/streambot/pkg/topic"
)

type ChannelManager struct {
	twitch     *twitch.Client
	subscriber zeromq.Subscriber
}

func NewChannelManager(clientID, accessToken, broadcasterID, subscriberAddress string) ChannelManager {
	return ChannelManager{
		twitch:     twitch.NewClient(clientID, accessToken, broadcasterID),
		subscriber: zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll),
	}
}

func (cm *ChannelManager) Start() {
	cm.subscriber.Start(func(msg message.Message) {
		switch msg.Topic {
		case topic.ServerTitleChanged:
			cm.twitch.SetTitle(msg.Content.ToString())
		}
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
