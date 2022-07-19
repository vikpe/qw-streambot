package twitch_manager

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/vikpe/streambot/internal/comms/topic"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

type TwitchManager struct {
	apiClient     *helix.Client
	broadcasterID string
	subscriber    *zeromq.SubscriberService
	stopChan      chan os.Signal
	OnStarted     func()
	OnStopped     func(os.Signal)
	OnError       func(error)
}

func New(clientID, accessToken, broadcasterID, subscriberAddress string) (*TwitchManager, error) {
	apiClient, err := helix.NewClient(&helix.Options{ClientID: clientID, AppAccessToken: accessToken})

	if err != nil {
		fmt.Println("twitch api client error", err)
		return &TwitchManager{}, err
	}

	subService := zeromq.NewSubscriberService(subscriberAddress, zeromq.TopicsAll)
	manager := TwitchManager{
		apiClient:     apiClient,
		broadcasterID: broadcasterID,
		subscriber:    subService,
		OnStarted:     func() {},
		OnStopped:     func(os.Signal) {},
		OnError:       func(error) {},
	}
	subService.OnMessage = manager.OnMessage

	return &manager, nil
}

func (m *TwitchManager) Start() {
	m.OnStarted()

	m.stopChan = make(chan os.Signal, 1)
	signal.Notify(m.stopChan, syscall.SIGTERM, syscall.SIGINT)
	go m.subscriber.Service.Start()
	sig := <-m.stopChan

	m.OnStopped(sig)
}

func (m *TwitchManager) Stop() {
	if m.stopChan == nil {
		return
	}
	m.stopChan <- syscall.SIGINT
	time.Sleep(30 * time.Millisecond)
}

func (m *TwitchManager) OnMessage(msg message.Message) {
	var err error

	switch msg.Topic {
	case topic.ServerTitleChanged:
		err = m.SetTitle(msg.Content.ToString())
	}

	if err != nil {
		m.OnError(err)
	}
}

func (m *TwitchManager) SetTitle(title string) error {
	const quakeGameId = "7348"

	_, err := m.apiClient.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: m.broadcasterID,
		Title:         title,
		GameID:        quakeGameId,
	})

	return err
}
