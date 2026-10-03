package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProduct_isInHHMM(t *testing.T) {
	assert.False(t, isInHHMM(time.Now(), "00:00", "00:00"))

	hhmmss := func(hhmm string) time.Time {
		t, _ := time.Parse("15:04:05", hhmm+":00")
		return t
	}

	assert.False(t, isInHHMM(hhmmss("23:58"), "00:00", "00:00"))
	assert.False(t, isInHHMM(hhmmss("00:00"), "00:00", "00:00"))
	assert.False(t, isInHHMM(hhmmss("00:01"), "00:00", "00:00"))

	assert.True(t, isInHHMM(hhmmss("23:59"), "00:01", "00:00"))
	assert.True(t, isInHHMM(hhmmss("23:58"), "23:50", "00:30"))
	assert.True(t, isInHHMM(hhmmss("00:01"), "23:50", "00:30"))
	assert.True(t, isInHHMM(hhmmss("23:00"), "21:00", "05:00"))
	assert.True(t, isInHHMM(hhmmss("00:00"), "21:00", "05:00"))
	assert.True(t, isInHHMM(hhmmss("01:00"), "21:00", "05:00"))
	assert.True(t, isInHHMM(hhmmss("02:00"), "21:00", "05:00"))
	assert.True(t, isInHHMM(hhmmss("03:00"), "21:00", "05:00"))

	assert.True(t, isInHHMM(hhmmss("4:0"), "21:0", "5:0"))
	assert.True(t, isInHHMM(hhmmss("5:0"), "21:0", "5:0"))
	assert.True(t, isInHHMM(hhmmss("6:0"), "21:0", "5:0"))

}
