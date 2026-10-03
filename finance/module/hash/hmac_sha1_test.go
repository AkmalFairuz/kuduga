package hash

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestHmacSha1(t *testing.T) {
	assert.Equal(t, HmacSha1("secret_data", "this_is_secret"), "36ecdd9daed2f3ba681a84e1a71fd3facabfc08c")
}
