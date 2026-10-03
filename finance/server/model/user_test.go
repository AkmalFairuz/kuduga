package model

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestIsValidUsername(t *testing.T) {
	assert.Equal(t, IsValidUsername("abc"), true)
	assert.Equal(t, IsValidUsername("abc123"), true)
	assert.Equal(t, IsValidUsername("abc.123"), false)
	assert.Equal(t, IsValidUsername("dfhsiufy893748728479248294729472"), false)
	assert.Equal(t, IsValidUsername("abc!@#$%^&*()"), false)
	assert.Equal(t, IsValidUsername("a"), false)
}
