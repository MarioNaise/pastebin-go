package pastebin

import (
	"testing"
	"time"
)

func TestGetHttpClient(t *testing.T) {
	httpClient = nil
	client := getHttpClient()
	if client == nil {
		t.Error("Expected non-nil HTTP client")
	}
	if client.Timeout != 10*time.Second {
		t.Errorf("Expected timeout of 10s, got %v", client.Timeout)
	}

	newClient := getHttpClient()
	if newClient != client {
		t.Error("Expected the same HTTP client instance")
	}
}
