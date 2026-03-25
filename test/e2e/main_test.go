//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	sdk "github.com/openziti/sdk-golang"
	"github.com/openziti/sdk-golang/ziti"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gatewayv1connect "github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
	zitimgmtv1 "github.com/agynio/gateway/gen/agynio/api/ziti_management/v1"
)

const (
	mockAuthTokenURL     = "https://mockauth.dev/r/301ebb13-15a8-48f4-baac-e3fa25be29fc/oidc/token"
	mockAuthClientID     = "client_MU95KU3gHQf5Ir7p"
	mockAuthClientSecret = "XPKka2i9uzISrKZ95zxli8sY51BK4eTJ"
	llmGRPCTarget        = "llm:50051"
	zitiManagementTarget = "ziti-management:50051"
	zitiGatewayURL       = "http://gateway"
)

var (
	gatewayURL = envOrDefault("GATEWAY_URL", "http://gateway-gateway:8080")

	apiTokenCreds  apiTokenCredentials
	agentModelID   string
	zitiHTTPClient *http.Client
	zitiIdentityID string
	zitiServiceID  string
)

type apiTokenCredentials struct {
	token      string
	identityID string
}

type mePayload struct {
	IdentityID   string `json:"identity_id"`
	IdentityType string `json:"identity_type"`
}

func TestMain(m *testing.M) {
	cleanup := &cleanupStack{}
	ctx := context.Background()

	if err := setupCredentials(ctx, cleanup); err != nil {
		exitWithSetupError(cleanup, fmt.Errorf("setup credentials: %w", err))
	}
	if err := setupModelID(ctx, cleanup); err != nil {
		exitWithSetupError(cleanup, fmt.Errorf("setup model id: %w", err))
	}
	if err := setupZitiIdentity(ctx, cleanup); err != nil {
		exitWithSetupError(cleanup, fmt.Errorf("setup ziti identity: %w", err))
	}

	exitCode := m.Run()
	cleanup.Run()
	os.Exit(exitCode)
}

type cleanupStack struct {
	fns []func()
}

func (c *cleanupStack) Add(fn func()) {
	c.fns = append(c.fns, fn)
}

func (c *cleanupStack) Run() {
	for i := len(c.fns) - 1; i >= 0; i-- {
		c.fns[i]()
	}
}

func exitWithSetupError(cleanup *cleanupStack, err error) {
	cleanup.Run()
	fmt.Fprintf(os.Stderr, "e2e setup failed: %v\n", err)
	os.Exit(1)
}

func setupCredentials(ctx context.Context, cleanup *cleanupStack) error {
	requestCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	accessToken, err := requestOIDCAccessToken(requestCtx)
	if err != nil {
		return err
	}

	meCtx, meCancel := context.WithTimeout(ctx, 15*time.Second)
	defer meCancel()
	identityID, err := fetchIdentityID(meCtx, accessToken)
	if err != nil {
		return err
	}

	apiTokenCreds.identityID = identityID

	createCtx, createCancel := context.WithTimeout(ctx, 15*time.Second)
	defer createCancel()

	client := gatewayv1connect.NewUsersGatewayClient(newAuthenticatedClient(accessToken), gatewayURL)
	createResp, err := client.CreateAPIToken(createCtx, connect.NewRequest(&usersv1.CreateAPITokenRequest{
		Name: fmt.Sprintf("e2e-api-token-%d", time.Now().UnixNano()),
	}))
	if err != nil {
		if connect.CodeOf(err) == connect.CodeUnimplemented {
			apiTokenCreds.token = accessToken
			return nil
		}
		return fmt.Errorf("create api token: %w", err)
	}

	plaintextToken := strings.TrimSpace(createResp.Msg.GetPlaintextToken())
	if plaintextToken == "" {
		return fmt.Errorf("api token plaintext missing")
	}

	tokenID := strings.TrimSpace(createResp.Msg.GetToken().GetId())
	if tokenID == "" {
		return fmt.Errorf("api token id missing")
	}

	apiTokenCreds.token = plaintextToken

	cleanup.Add(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cleanupCancel()
		_, _ = client.RevokeAPIToken(cleanupCtx, connect.NewRequest(&usersv1.RevokeAPITokenRequest{TokenId: tokenID}))
	})

	return nil
}

func requestOIDCAccessToken(ctx context.Context) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("username", "e2e-test-user")
	form.Set("scope", "openid profile email")
	form.Set("client_id", mockAuthClientID)
	form.Set("client_secret", mockAuthClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mockAuthTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := newClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("mockauth token request failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}

	accessToken := strings.TrimSpace(payload.AccessToken)
	if accessToken == "" {
		return "", fmt.Errorf("mockauth access_token missing")
	}

	return accessToken, nil
}

func fetchIdentityID(ctx context.Context, accessToken string) (string, error) {
	client := newAuthenticatedClient(accessToken)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, gatewayURL+"/me", nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("me endpoint failed: status %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload mePayload
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", err
	}

	identityID := strings.TrimSpace(payload.IdentityID)
	if identityID == "" {
		return "", fmt.Errorf("identity_id missing")
	}

	return identityID, nil
}

func setupModelID(ctx context.Context, cleanup *cleanupStack) error {
	conn, err := grpc.NewClient(llmGRPCTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client := llmv1.NewLLMServiceClient(conn)

	var providerID string
	var modelID string

	cleanup.Add(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		if modelID != "" {
			_, _ = client.DeleteModel(cleanupCtx, &llmv1.DeleteModelRequest{Id: modelID})
		}
		if providerID != "" {
			_, _ = client.DeleteLLMProvider(cleanupCtx, &llmv1.DeleteLLMProviderRequest{Id: providerID})
		}
		_ = conn.Close()
	})

	providerCtx, providerCancel := context.WithTimeout(ctx, 15*time.Second)
	defer providerCancel()

	providerResp, err := client.CreateLLMProvider(providerCtx, &llmv1.CreateLLMProviderRequest{
		Endpoint:   "https://testllm.dev/v1/org/agynio/suite/agn/responses",
		Token:      "unused",
		AuthMethod: llmv1.AuthMethod_AUTH_METHOD_BEARER,
	})
	if err != nil {
		return fmt.Errorf("create llm provider: %w", err)
	}

	providerID = strings.TrimSpace(providerResp.GetProvider().GetMeta().GetId())
	if providerID == "" {
		return fmt.Errorf("llm provider id missing")
	}

	modelCtx, modelCancel := context.WithTimeout(ctx, 15*time.Second)
	defer modelCancel()

	modelResp, err := client.CreateModel(modelCtx, &llmv1.CreateModelRequest{
		Name:          "e2e-test-model",
		LlmProviderId: providerID,
		RemoteName:    "simple-hello",
	})
	if err != nil {
		return fmt.Errorf("create llm model: %w", err)
	}

	modelID = strings.TrimSpace(modelResp.GetModel().GetMeta().GetId())
	if modelID == "" {
		return fmt.Errorf("llm model id missing")
	}

	agentModelID = modelID
	return nil
}

func setupZitiIdentity(ctx context.Context, cleanup *cleanupStack) error {
	conn, err := grpc.NewClient(zitiManagementTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client := zitimgmtv1.NewZitiManagementServiceClient(conn)

	var zitiContext ziti.Context
	cleanup.Add(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		if zitiIdentityID != "" && zitiServiceID != "" {
			_, _ = client.DeleteAppIdentity(cleanupCtx, &zitimgmtv1.DeleteAppIdentityRequest{
				ZitiIdentityId: zitiIdentityID,
				ZitiServiceId:  zitiServiceID,
			})
		}
		if zitiContext != nil {
			zitiContext.Close()
		}
		_ = conn.Close()
	})

	createCtx, createCancel := context.WithTimeout(ctx, 30*time.Second)
	defer createCancel()

	resp, err := client.CreateAppIdentity(createCtx, &zitimgmtv1.CreateAppIdentityRequest{
		IdentityId: uuid.NewString(),
		Slug:       "e2e-test",
	})
	if err != nil {
		return fmt.Errorf("create app identity: %w", err)
	}

	zitiIdentityID = strings.TrimSpace(resp.GetZitiIdentityId())
	zitiServiceID = strings.TrimSpace(resp.GetZitiServiceId())
	identityJSON := resp.GetIdentityJson()
	if zitiIdentityID == "" {
		return fmt.Errorf("ziti identity id missing")
	}
	if zitiServiceID == "" {
		return fmt.Errorf("ziti service id missing")
	}
	if len(identityJSON) == 0 {
		return fmt.Errorf("ziti identity json missing")
	}

	zitiConfig := &ziti.Config{}
	if err := json.Unmarshal(identityJSON, zitiConfig); err != nil {
		return fmt.Errorf("parse ziti identity: %w", err)
	}

	zitiContext, err = ziti.NewContext(zitiConfig)
	if err != nil {
		return fmt.Errorf("create ziti context: %w", err)
	}

	zitiHTTPClient = sdk.NewHttpClient(zitiContext, nil)
	return nil
}

func newClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
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

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
