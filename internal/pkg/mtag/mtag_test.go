package mtag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/internal/pkg/mtag"
)

func TestIsOfficial(t *testing.T) {
	testCases := map[string]bool{
		"getquad":      true,
		"getquad ":     true,
		"getquad semi": true,
		"kombat":       true,
		"kombat 2on2":  true,
		"getquad 6":    true,
		"gq6":          true,
		"qwdl":         true,
		"qwduel":       true,

		"":            false,
		"lolz":        false,
		"some kombat": false,
	}

	for matchtag, expect := range testCases {
		t.Run(matchtag, func(t *testing.T) {
			assert.Equal(t, expect, mtag.IsOfficial(matchtag))
		})
	}
}
