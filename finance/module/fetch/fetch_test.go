package fetch

import (
	"testing"
)

func TestFetch(t *testing.T) {
	resp, err := Get("https://icanhazip.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Status: %d", resp.Status)

	for key, val := range resp.Headers {
		t.Logf("Header %s: %s", key, val)
	}

	t.Logf("Body: %s", string(resp.Body))
}
