//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	sdk "github.com/openziti/sdk-golang"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/stretchr/testify/require"
)

const defaultZitiServiceName = "gateway"

type mePayload struct {
	IdentityID   string `json:"identity_id"`
	IdentityType string `json:"identity_type"`
}

func TestZitiMeEndpointAuthenticated(t *testing.T) {
	identityFile := zitiIdentityFile()
	if identityFile == "" {
		t.Skip("ziti identity file not configured")
	}
	if _, err := os.Stat(identityFile); err != nil {
		t.Skipf("ziti identity file unavailable: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	zitiContext, err := ziti.NewContextFromFile(identityFile)
	if err != nil {
		t.Skipf("ziti context unavailable: %v", err)
	}
	defer zitiContext.Close()

	client := sdk.NewHttpClient(zitiContext, nil)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, zitiGatewayURL()+"/me", nil)
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}

	var payload mePayload
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.IdentityID == "" {
		t.Fatalf("expected identity_id field")
	}
	if payload.IdentityType == "" {
		t.Fatalf("expected identity_type field")
	}
}

func TestZitiMeEndpointUnauthenticated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := newClient()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, gatewayURL+"/me", nil)
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.StatusCode)
	}
}

func zitiIdentityFile() string {
	if value := strings.TrimSpace(os.Getenv("ZITI_E2E_IDENTITY_FILE")); value != "" {
		return value
	}
	return strings.TrimSpace(os.Getenv("ZITI_IDENTITY_FILE"))
}

func zitiGatewayURL() string {
	if value := strings.TrimSpace(os.Getenv("ZITI_GATEWAY_URL")); value != "" {
		return value
	}

	serviceName := strings.TrimSpace(os.Getenv("ZITI_SERVICE_NAME"))
	if serviceName == "" {
		serviceName = defaultZitiServiceName
	}
	return "http://" + serviceName
}
