package qws_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/streambot/pkg/qws"
)

func TestIsRelevantServer(t *testing.T) {
	t.Run("no - excluded mode", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			SpectatorSlots: slots.New(4, 0),
			Mode:           qmode.Mode("fortress"),
		}
		assert.False(t, qws.IsRelevantServer(server))
	})

	t.Run("no - excluded region", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			SpectatorSlots: slots.New(4, 0),
			Geo:            geo.Info{Region: "South America"},
		}
		assert.False(t, qws.IsRelevantServer(server))
	})

	t.Run("yes", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			SpectatorSlots: slots.New(4, 0),
			Geo:            geo.Info{Region: "Europe"},
		}
		assert.True(t, qws.IsRelevantServer(server))
	})
}
