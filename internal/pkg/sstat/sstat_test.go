package sstat_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/streambot/internal/pkg/sstat"
)

func TestGetMvdsvServer(t *testing.T) {
	t.Run("invalid server address", func(t *testing.T) {
		assert.Equal(t, mvdsv.Mvdsv{}, sstat.GetMvdsvServer("foo"))
	})
}
