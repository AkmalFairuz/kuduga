package hash

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestSha256(t *testing.T) {
	assert.Equal(t, "6dd2d64d09f76ba044154313caf82d6219b3cbb6a32672ba564875a8c1bd634b", Sha256("this is secret!"))
}
