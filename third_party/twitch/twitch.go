package twitch

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
)

const quakeGameId = "7348"

type Client struct {
	client        *helix.Client
	broadcasterId string
}

func NewClient(clientID string, accessToken string, broadcasterID string) Client {
	client, _ := helix.NewClient(&helix.Options{ClientID: clientID, AppAccessToken: accessToken})

	return Client{
		client:        client,
		broadcasterId: broadcasterID,
	}
}

func (a Client) SetTitle(title string) {
	_, err := a.client.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: a.broadcasterId,
		Title:         title,
		GameID:        quakeGameId,
	})

	if err != nil {
		fmt.Println(err)
		return
	}
}
