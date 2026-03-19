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

func TestZitiMeEndpoint(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := newClient()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, gatewayURL+"/me", nil)
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unexpected status %d", response.StatusCode)
	}

	if response.StatusCode == http.StatusOK {
		var payload map[string]any
		if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if _, ok := payload["identity_id"]; !ok {
			t.Fatalf("expected identity_id field")
		}
		if _, ok := payload["identity_type"]; !ok {
			t.Fatalf("expected identity_type field")
		}
		if _, ok := payload["tenant_id"]; !ok {
			t.Fatalf("expected tenant_id field")
		}
		if _, ok := payload["auth_method"]; !ok {
			t.Fatalf("expected auth_method field")
		}
	}
}
