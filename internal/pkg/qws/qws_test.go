package qws_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/streambot/internal/pkg/qws"
)

func TestIsRelevantServer(t *testing.T) {
	t.Run("no - excluded mode", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			SpectatorSlots: slots.New(4, 0),
			Mode:           qmode.Mode("fortress"),
		}
		assert.False(t, qws.IsRelevantServer(server))
	})

	t.Run("no - far away region without qtv", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			SpectatorSlots: slots.New(4, 0),
			Geo:            geo.Location{Region: "Oceania"},
		}
		assert.False(t, qws.IsRelevantServer(server))
	})

	t.Run("yes", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			SpectatorSlots: slots.New(4, 0),
			Geo:            geo.Location{Region: "Europe"},
		}
		assert.True(t, qws.IsRelevantServer(server))
	})
}

func TestServerScoreBonus(t *testing.T) {
	t.Run("not XonX", func(t *testing.T) {
		server := mvdsv.Mvdsv{Mode: qmode.Mode("ffa")}
		assert.Equal(t, 0, qws.ServerScoreBonus(server))
	})

	t.Run("not official", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			Mode:     qmode.Mode("1on1"),
			Settings: qsettings.Settings{},
		}
		assert.Equal(t, 0, qws.ServerScoreBonus(server))
	})

	t.Run("1on1 with 1 free slot", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			Mode:        qmode.Mode("1on1"),
			PlayerSlots: slots.New(2, 1),
			Settings:    qsettings.Settings{"matchtag": "getquad"},
		}
		assert.Equal(t, 0, qws.ServerScoreBonus(server))
	})

	t.Run("2on2 with 2 free slots", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			Mode:        qmode.Mode("2on2"),
			PlayerSlots: slots.New(4, 2),
			Settings:    qsettings.Settings{"matchtag": "getquad"},
		}
		assert.Equal(t, 0, qws.ServerScoreBonus(server))
	})

	t.Run("4on4 with 3 free slots", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			Mode:        qmode.Mode("4on4"),
			PlayerSlots: slots.New(8, 5),
			Settings:    qsettings.Settings{"matchtag": "getquad"},
		}
		assert.Equal(t, 0, qws.ServerScoreBonus(server))
	})

	t.Run("4on4 with 1 free slot", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			Mode:        qmode.Mode("4on4"),
			PlayerSlots: slots.New(8, 7),
			Settings:    qsettings.Settings{"matchtag": "getquad"},
		}
		assert.Equal(t, 30, qws.ServerScoreBonus(server))
	})
}
