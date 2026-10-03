package service

import (
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestStringRemoveNonAlphanumeric(t *testing.T) {
	assert.Equal(t, helper.StringRemoveNonAlphanumeric("abc.123"), "abc123")
	assert.Equal(t, helper.StringRemoveNonAlphanumeric("abc.123!@#$%^&*()"), "abc123")
}
