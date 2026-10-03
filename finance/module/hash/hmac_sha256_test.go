package hash

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestHmacSha256(t *testing.T) {
	assert.Equal(t, HmacSha256("secret_data", "this_is_secret"), "834e0ed5949a4f064d6500cbe1e08ae4777496dae0c88f257f72691fc655cbb4")
}
