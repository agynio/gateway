//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

func TestTeamAPI_ListAgents(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := newClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gatewayURL+"/team/v1/agents", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	_, ok := body["items"]
	assert.True(t, ok, "response should contain 'items' key")
}

func TestTeamAPI_CreateAndDeleteAgent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := newClient()

	payload := `{"title":"e2e-smoke-agent","config":{}}`
	createReq, err := http.NewRequestWithContext(ctx, http.MethodPost, gatewayURL+"/team/v1/agents", strings.NewReader(payload))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := client.Do(createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	var created map[string]interface{}
	err = json.NewDecoder(createResp.Body).Decode(&created)
	require.NoError(t, err)

	agentID, ok := created["id"].(string)
	require.True(t, ok, "created agent must have a string 'id'")
	require.NotEmpty(t, agentID)

	deleteReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, gatewayURL+"/team/v1/agents/"+agentID, nil)
	require.NoError(t, err)

	deleteResp, err := client.Do(deleteReq)
	require.NoError(t, err)
	defer deleteResp.Body.Close()

	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)
}

func TestTeamAPI_ListMcpServers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := newClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gatewayURL+"/team/v1/mcp-servers", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestTeamAPI_InvalidPayloadReturnsClientError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := newClient()
	payload := `{"unknown_field":"value"}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, gatewayURL+"/team/v1/agents", strings.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	assert.GreaterOrEqual(t, resp.StatusCode, 400)
	assert.Less(t, resp.StatusCode, 500)
}
