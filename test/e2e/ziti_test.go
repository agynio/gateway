//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestZitiMeEndpointAuthenticated(t *testing.T) {
	require.NotNil(t, zitiHTTPClient)
	require.NotEmpty(t, zitiIdentityID)
	require.NotEmpty(t, zitiServiceID)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, zitiGatewayURL+"/me", nil)
	require.NoError(t, err)

	response, err := zitiHTTPClient.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}

	var payload mePayload
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))
	require.NotEmpty(t, payload.IdentityID)
	require.NotEmpty(t, payload.IdentityType)
}
