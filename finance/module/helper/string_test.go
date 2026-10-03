package helper

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestStringToInt(t *testing.T) {
	assert.Equal(t, StringToInt("123"), 123)
	assert.Equal(t, StringToInt("-123"), -123)
}
