package twitch_manager

import (
	"errors"
	"fmt"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/vikpe/streambot/internal/comms/topic"
	"github.com/vikpe/streambot/internal/pkg/zeromq"
	"github.com/vikpe/streambot/internal/pkg/zeromq/message"
)

func New(clientID, accessToken, broadcasterID, subscriberAddress string) (*zeromq.Subscriber, error) {
	apiClient, err := helix.NewClient(&helix.Options{
		ClientID:       clientID,
		AppAccessToken: accessToken,
		RateLimitFunc:  rateLimitCallback,
	})

	subscriber := zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll)

	if err != nil {
		err := errors.New(fmt.Sprintf("twitch api client error: %s", err))
		return subscriber, err
	}

	subscriber.OnMessage = func(msg message.Message) {
		var err error

		switch msg.Topic {
		case topic.ServerTitleChanged:
			err = SetTitle(apiClient, broadcasterID, msg.Content.ToString())
		}

		if err != nil {
			subscriber.OnError(err)
		}
	}

	return subscriber, nil
}

func rateLimitCallback(lastResponse *helix.Response) error {
	if lastResponse.GetRateLimitRemaining() > 0 {
		return nil
	}

	nextReset := int64(lastResponse.GetRateLimitReset())
	currentTime := time.Now().Unix()

	if nextReset > currentTime {
		timeToNextReset := time.Duration(nextReset - currentTime)
		time.Sleep(timeToNextReset * time.Second)
	}

	return nil
}

func SetTitle(apiClient *helix.Client, broadcasterID, title string) error {
	const quakeGameId = "7348"

	_, err := apiClient.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: broadcasterID,
		Title:         title,
		GameID:        quakeGameId,
	})

	return err
}
