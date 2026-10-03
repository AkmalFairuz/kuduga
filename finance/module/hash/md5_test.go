package hash

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestMd5(t *testing.T) {
	assert.Equal(t, "938c2cc0dcc05f2b68c4287040cfcf71", Md5("frog"))
}
