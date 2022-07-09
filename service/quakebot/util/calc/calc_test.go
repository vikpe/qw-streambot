package calc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/service/quakebot/util/calc"
)

func TestClampFloat64(t *testing.T) {
	t.Run("below min", func(t *testing.T) {
		assert.Equal(t, 5.0, calc.ClampFloat64(2.5, 5.0, 10.0))
	})

	t.Run("above max", func(t *testing.T) {
		assert.Equal(t, 10.0, calc.ClampFloat64(15, 5.0, 10.0))
	})
}

func TestRoundFloat64(t *testing.T) {
	assert.Equal(t, 3.33, calc.RoundFloat64(10.0/3, 2))
	assert.Equal(t, 3.0, calc.RoundFloat64(10.0/3, 0))
}

func TestStaticTextScale(t *testing.T) {
	assert.Equal(t, 1.5, calc.StaticTextScale(""))
	assert.Equal(t, 1.5, calc.StaticTextScale("getquad semi"))
	assert.Equal(t, 1.25, calc.StaticTextScale("getquad semigetquad semi"))
	assert.Equal(t, 1.0, calc.StaticTextScale("getquad semigetquad semigetquad semi"))
}
