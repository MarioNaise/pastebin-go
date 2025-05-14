package pastebin

import (
	"net/http"
	"time"
)

var httpClient *http.Client

func getHttpClient() *http.Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	return httpClient
}
