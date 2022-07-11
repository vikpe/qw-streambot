package sstat_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/internal/pkg/sstat"
	"github.com/vikpe/udphelper"
)

func TestGetMvdsvServer(t *testing.T) {
	t.Run("invalid server address", func(t *testing.T) {
		assert.Equal(t, mvdsv.Mvdsv{}, sstat.GetMvdsvServer("foo"))
	})

	t.Run("valid server", func(t *testing.T) {
		response := append([]byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}, []byte(`*version\MVDSV 0.35-dev`)...)
		go func() { udphelper.New(":28501").Respond(response) }()
		time.Sleep(time.Millisecond * 10)

		server := sstat.GetMvdsvServer("localhost:28501")
		assert.Equal(t, "localhost:28501", server.Address)
		assert.Equal(t, "MVDSV 0.35-dev", server.Settings.Get("*version", ""))
	})
}
