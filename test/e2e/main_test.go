package e2e

import (
	"os"
	"testing"
)

var gatewayURL = envOrDefault("GATEWAY_URL", "http://gateway-gateway:8080")

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
