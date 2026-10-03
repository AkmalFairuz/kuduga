package helper

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestFormatRupiah(t *testing.T) {
	assert.Equal(t, FormatRupiah(123999), "Rp123.999")
	assert.Equal(t, FormatRupiah(123000000), "Rp123.000.000")
	assert.Equal(t, FormatRupiah(10), "Rp10")
	assert.Equal(t, FormatRupiah(10333), "Rp10.333")
	assert.Equal(t, FormatRupiah(10333444), "Rp10.333.444")
}
