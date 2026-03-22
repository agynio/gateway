//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	gatewayv1connect "github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const invalidAPIToken = "agyn_invalid"

type apiTokenCredentials struct {
	token      string
	identityID string
}

func TestAPIToken_MeEndpoint(t *testing.T) {
	config := apiTokenConfig()
	if config.token == "" || config.identityID == "" {
		t.Skip("api token identity config not provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := newAuthenticatedClient(config.token)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, gatewayURL+"/me", nil)
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}

	var payload mePayload
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))
	assert.Equal(t, config.identityID, payload.IdentityID)
	assert.Equal(t, "user", payload.IdentityType)
}

func TestAPIToken_MeEndpointInvalidToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := newAuthenticatedClient(invalidAPIToken)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, gatewayURL+"/me", nil)
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.StatusCode)
	}
}

func TestAPIToken_ConnectRPCEndpointAuthenticated(t *testing.T) {
	config := apiTokenConfig()
	if config.token == "" {
		t.Skip("api token not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := gatewayv1connect.NewAgentsGatewayClient(newAuthenticatedClient(config.token), gatewayURL)
	resp, err := client.ListAgents(ctx, connect.NewRequest(&agentsv1.ListAgentsRequest{}))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Msg.Agents)
}

func TestAPIToken_ConnectRPCEndpointInvalidToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := gatewayv1connect.NewAgentsGatewayClient(newAuthenticatedClient(invalidAPIToken), gatewayURL)
	_, err := client.ListAgents(ctx, connect.NewRequest(&agentsv1.ListAgentsRequest{}))
	require.Error(t, err)
	assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
}

func TestUsersGateway_CreateAndRevokeAPIToken(t *testing.T) {
	config := apiTokenConfig()
	if config.token == "" {
		t.Skip("api token not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := gatewayv1connect.NewUsersGatewayClient(newAuthenticatedClient(config.token), gatewayURL)
	createResp, err := client.CreateAPIToken(ctx, connect.NewRequest(&usersv1.CreateAPITokenRequest{
		Name: fmt.Sprintf("e2e-api-token-%d", time.Now().UnixNano()),
	}))
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.NotNil(t, createResp.Msg)
	require.NotNil(t, createResp.Msg.Token)
	require.NotEmpty(t, createResp.Msg.Token.Id)

	tokenID := createResp.Msg.Token.Id
	t.Cleanup(func() {
		if tokenID == "" {
			return
		}
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cleanupCancel()
		_, _ = client.RevokeAPIToken(cleanupCtx, connect.NewRequest(&usersv1.RevokeAPITokenRequest{TokenId: tokenID}))
	})

	listResp, err := client.ListAPITokens(ctx, connect.NewRequest(&usersv1.ListAPITokensRequest{}))
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.NotNil(t, listResp.Msg)
	assert.True(t, hasTokenID(listResp.Msg.Tokens, tokenID), "expected token to be listed")

	_, err = client.RevokeAPIToken(ctx, connect.NewRequest(&usersv1.RevokeAPITokenRequest{TokenId: tokenID}))
	require.NoError(t, err)

	listResp, err = client.ListAPITokens(ctx, connect.NewRequest(&usersv1.ListAPITokensRequest{}))
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.NotNil(t, listResp.Msg)
	assert.False(t, hasTokenID(listResp.Msg.Tokens, tokenID), "expected token to be revoked")
}

func TestUsersGateway_ListAPITokens(t *testing.T) {
	config := apiTokenConfig()
	if config.token == "" {
		t.Skip("api token not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := gatewayv1connect.NewUsersGatewayClient(newAuthenticatedClient(config.token), gatewayURL)
	resp, err := client.ListAPITokens(ctx, connect.NewRequest(&usersv1.ListAPITokensRequest{}))
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Msg)
	assert.NotNil(t, resp.Msg.Tokens)
}

func TestUsersGateway_RevokeAPITokenNotFound(t *testing.T) {
	config := apiTokenConfig()
	if config.token == "" {
		t.Skip("api token not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := gatewayv1connect.NewUsersGatewayClient(newAuthenticatedClient(config.token), gatewayURL)
	_, err := client.RevokeAPIToken(ctx, connect.NewRequest(&usersv1.RevokeAPITokenRequest{TokenId: uuid.NewString()}))
	require.Error(t, err)
	assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
}

func TestUsersGateway_CreateAPITokenUnauthenticated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := gatewayv1connect.NewUsersGatewayClient(newClient(), gatewayURL)
	_, err := client.CreateAPIToken(ctx, connect.NewRequest(&usersv1.CreateAPITokenRequest{Name: "e2e-unauth"}))
	require.Error(t, err)
	assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
}

func TestAPIToken_CreatedTokenAuthenticates(t *testing.T) {
	config := apiTokenConfig()
	if config.token == "" || config.identityID == "" {
		t.Skip("api token identity config not provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := gatewayv1connect.NewUsersGatewayClient(newAuthenticatedClient(config.token), gatewayURL)
	createResp, err := client.CreateAPIToken(ctx, connect.NewRequest(&usersv1.CreateAPITokenRequest{
		Name: fmt.Sprintf("e2e-roundtrip-token-%d", time.Now().UnixNano()),
	}))
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.NotNil(t, createResp.Msg)
	require.NotNil(t, createResp.Msg.Token)
	require.NotEmpty(t, createResp.Msg.Token.Id)
	require.NotEmpty(t, createResp.Msg.PlaintextToken)
	assert.True(t, strings.HasPrefix(createResp.Msg.PlaintextToken, "agyn_"))

	tokenID := createResp.Msg.Token.Id
	t.Cleanup(func() {
		if tokenID == "" {
			return
		}
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cleanupCancel()
		_, _ = client.RevokeAPIToken(cleanupCtx, connect.NewRequest(&usersv1.RevokeAPITokenRequest{TokenId: tokenID}))
	})

	meCtx, meCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer meCancel()

	authClient := newAuthenticatedClient(createResp.Msg.PlaintextToken)
	request, err := http.NewRequestWithContext(meCtx, http.MethodGet, gatewayURL+"/me", nil)
	require.NoError(t, err)

	response, err := authClient.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}

	var payload mePayload
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))
	assert.Equal(t, config.identityID, payload.IdentityID)
	assert.Equal(t, "user", payload.IdentityType)
}

func apiTokenConfig() apiTokenCredentials {
	token := strings.TrimSpace(os.Getenv("E2E_API_TOKEN"))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("API_TOKEN"))
	}

	identityID := strings.TrimSpace(os.Getenv("E2E_API_TOKEN_IDENTITY_ID"))
	if identityID == "" {
		identityID = strings.TrimSpace(os.Getenv("API_TOKEN_IDENTITY_ID"))
	}

	return apiTokenCredentials{token: token, identityID: identityID}
}

func hasTokenID(tokens []*usersv1.APIToken, tokenID string) bool {
	for _, token := range tokens {
		if token.GetId() == tokenID {
			return true
		}
	}
	return false
}

func newAuthenticatedClient(token string) *http.Client {
	client := newClient()
	client.Transport = bearerTransport{token: token, base: client.Transport}
	return client
}

type bearerTransport struct {
	token string
	base  http.RoundTripper
}

func (t bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	clone := req.Clone(req.Context())
	clone.Header.Set("Authorization", "Bearer "+t.token)
	return base.RoundTrip(clone)
}
