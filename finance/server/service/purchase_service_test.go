package service

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestCalculateAutoUserSellPrice(t *testing.T) {
	assert.Equal(t, int64(12000), calculateAutoUserSellPrice(10123, 2000))
	assert.Equal(t, int64(12500), calculateAutoUserSellPrice(10523, 2000))
	assert.Equal(t, int64(12000), calculateAutoUserSellPrice(10000, 2000))
	assert.Equal(t, int64(12500), calculateAutoUserSellPrice(10500, 2000))
	assert.Equal(t, int64(13000), calculateAutoUserSellPrice(10990, 2000))
}
