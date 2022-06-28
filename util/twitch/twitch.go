package twitch

import (
	"github.com/nicklaw5/helix/v2"
)

type Client struct {
	client        *helix.Client
	broadcasterId string
}

func NewClient(clientID string, accessToken string, broadcasterID string) Client {
	client, _ := helix.NewClient(&helix.Options{ClientID: clientID, UserAccessToken: accessToken})

	return Client{
		client:        client,
		broadcasterId: broadcasterID,
	}
}

func (a Client) SetTitle(title string) {
	a.client.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: a.broadcasterId,
		Title:         title,
	})
}
