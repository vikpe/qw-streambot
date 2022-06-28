package qws_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/streambot/qws"
)

func TestRequiresPassword(t *testing.T) {
	testCases := map[int]bool{
		0: false,
		4: false,
		5: false,
		2: true,
		3: true,
		6: true,
		7: true,
	}

	for needpass, expect := range testCases {
		t.Run(fmt.Sprintf("needpass=%d", needpass), func(t *testing.T) {
			assert.Equal(t, expect, qws.RequiresPassword(needpass))
		})
	}
}

func TestIsSpeccable(t *testing.T) {
	t.Run("yes - has qtv stream", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			QtvStream: qtvstream.QtvStream{Url: "2@troopers.fi:28000"},
		}
		assert.True(t, qws.IsRelevantServer(server))
	})

	t.Run("yes - has free spectator slots and no password", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			SpectatorSlots: slots.New(4, 3),
		}
		assert.True(t, qws.IsRelevantServer(server))
	})

	t.Run("no - has free spectator slots and password", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			SpectatorSlots: slots.New(4, 3),
			Settings:       qsettings.Settings{"needpass": "2"},
		}
		assert.False(t, qws.IsRelevantServer(server))
	})
}

func TestIsRelevantServer(t *testing.T) {
	t.Run("no - excluded region", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			Geo: geo.Info{Region: "South America"},
		}
		assert.False(t, qws.IsRelevantServer(server))
	})

	t.Run("yes", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			Geo: geo.Info{Region: "Europe"},
		}
		assert.False(t, qws.IsRelevantServer(server))
	})
}
