package channel_manager

import (
	"os"
	"syscall"
	"time"

	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
	"github.com/vikpe/streambot/internal/third_party/twitch"
	"github.com/vikpe/streambot/pkg/topic"
)

type ChannelManager struct {
	twitchApi  *twitch.Client
	subscriber zeromq.Subscriber
	stopChan   chan os.Signal
	OnStarted  func()
	OnStopped  func(os.Signal)
}

func NewChannelManager(clientID, accessToken, broadcasterID, subscriberAddress string) ChannelManager {
	return ChannelManager{
		twitchApi:  twitch.NewClient(clientID, accessToken, broadcasterID),
		subscriber: zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll),
		OnStarted:  func() {},
		OnStopped:  func(os.Signal) {},
	}
}

func (m *ChannelManager) Start() {
	m.OnStarted()

	go func() {
		m.subscriber.Start(func(msg message.Message) {
			switch msg.Topic {
			case topic.ServerTitleChanged:
				m.twitchApi.SetTitle(msg.Content.ToString())
			}
		})
	}()

	sig := <-m.stopChan
	m.OnStopped(sig)
}

func (m *ChannelManager) SetTitle(title string) error {
	_, err := m.apiClient.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: m.broadcasterID,
		Title:         title,
		GameID:        quakeGameId,
	})

	return err
}

func (m *ChannelManager) Stop() {
	if m.stopChan == nil {
		return
	}
	m.stopChan <- syscall.SIGINT
	time.Sleep(50 * time.Millisecond)
}
