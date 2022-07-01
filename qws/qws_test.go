package qws_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/qws"
)

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
