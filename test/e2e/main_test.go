//go:build e2e

package e2e

import (
	"net/http"
	"os"
	"testing"
	"time"
)

var gatewayURL = envOrDefault("GATEWAY_URL", "http://gateway-gateway:8080")

func newClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
