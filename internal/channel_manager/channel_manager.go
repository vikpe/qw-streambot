package channel_manager

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
	"github.com/vikpe/streambot/pkg/topic"
)

const quakeGameId = "7348"

type ChannelManager struct {
	apiClient     *helix.Client
	broadcasterID string
	subscriber    zeromq.Subscriber
	stopChan      chan os.Signal
	OnStarted     func()
	OnStopped     func(os.Signal)
	OnError       func(error)
}

func NewChannelManager(clientID, accessToken, broadcasterID, subscriberAddress string) (ChannelManager, error) {
	apiClient, err := helix.NewClient(&helix.Options{ClientID: clientID, AppAccessToken: accessToken})

	if err != nil {
		fmt.Println("twitch api client error", err)
		return ChannelManager{}, err
	}

	return ChannelManager{
		apiClient:     apiClient,
		broadcasterID: broadcasterID,
		subscriber:    zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll),
		OnStarted:     func() {},
		OnStopped:     func(os.Signal) {},
		OnError:       func(error) {},
	}, nil
}

func (m *ChannelManager) Start() {
	m.OnStarted()

	go func() {
		m.subscriber.Start(func(msg message.Message) {
			var err error
			switch msg.Topic {
			case topic.ServerTitleChanged:
				err = m.SetTitle(msg.Content.ToString())
			}

			if err != nil {
				m.OnError(err)
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
