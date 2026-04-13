//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"errors"
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
	identityv1 "github.com/agynio/gateway/gen/agynio/api/identity/v1"
	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
	organizationsv1 "github.com/agynio/gateway/gen/agynio/api/organizations/v1"
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
	organizationID string
	zitiHTTPClient *http.Client
	zitiAppID      string
	zitiIdentityID string
	zitiServiceID  string
)

var errZitiServiceIDNotReady = errors.New("ziti service id not ready")

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
	bootstrapCtx, bootstrapCancel := context.WithTimeout(ctx, 15*time.Second)
	defer bootstrapCancel()

	accessToken, identityID, err := resolveBootstrapCredentials(bootstrapCtx)
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

func resolveBootstrapCredentials(ctx context.Context) (string, string, error) {
	clusterToken := strings.TrimSpace(os.Getenv("CLUSTER_ADMIN_TOKEN"))
	clusterIdentityID := strings.TrimSpace(os.Getenv("CLUSTER_ADMIN_IDENTITY_ID"))
	if clusterToken != "" || clusterIdentityID != "" {
		if clusterToken == "" || clusterIdentityID == "" {
			return "", "", fmt.Errorf("CLUSTER_ADMIN_TOKEN and CLUSTER_ADMIN_IDENTITY_ID must both be set for e2e")
		}
		return clusterToken, clusterIdentityID, nil
	}

	accessToken, err := requestOIDCAccessToken(ctx)
	if err != nil {
		return "", "", err
	}

	identityID, err := fetchIdentityID(ctx, accessToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, identityID, nil
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

	orgCtx, orgCancel := context.WithTimeout(ctx, 15*time.Second)
	defer orgCancel()

	resolvedOrganizationID, err := resolveOrganizationID(orgCtx)
	if err != nil {
		return err
	}
	resolvedOrganizationID = strings.TrimSpace(resolvedOrganizationID)
	if resolvedOrganizationID == "" {
		return fmt.Errorf("organization id missing")
	}
	organizationID = resolvedOrganizationID

	providerResp, err := client.CreateLLMProvider(providerCtx, &llmv1.CreateLLMProviderRequest{
		Endpoint:       "https://testllm.dev/v1/org/agynio/suite/agn/responses",
		Token:          "unused",
		AuthMethod:     llmv1.AuthMethod_AUTH_METHOD_BEARER,
		OrganizationId: resolvedOrganizationID,
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
		Name:           "e2e-test-model",
		LlmProviderId:  providerID,
		RemoteName:     "simple-hello",
		OrganizationId: resolvedOrganizationID,
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

func resolveOrganizationID(ctx context.Context) (string, error) {
	accessToken := strings.TrimSpace(apiTokenCreds.token)
	if accessToken == "" {
		return "", fmt.Errorf("api token missing")
	}

	client := gatewayv1connect.NewOrganizationsGatewayClient(newAuthenticatedClient(accessToken), gatewayURL)
	listResp, err := client.ListOrganizations(ctx, connect.NewRequest(&organizationsv1.ListOrganizationsRequest{PageSize: 50}))
	if err == nil {
		orgID := firstOrganizationID(listResp.Msg.GetOrganizations())
		if orgID != "" {
			return orgID, nil
		}
	}

	identityID := strings.TrimSpace(apiTokenCreds.identityID)
	if identityID == "" {
		return "", fmt.Errorf("identity id missing")
	}

	accessibleResp, accessibleErr := client.ListAccessibleOrganizations(ctx, connect.NewRequest(&organizationsv1.ListAccessibleOrganizationsRequest{
		IdentityId: identityID,
	}))
	if accessibleErr != nil {
		if err != nil {
			return "", fmt.Errorf("list organizations: %v; list accessible organizations: %w", err, accessibleErr)
		}
		return "", fmt.Errorf("list accessible organizations: %w", accessibleErr)
	}

	orgID := firstOrganizationID(accessibleResp.Msg.GetOrganizations())
	if orgID == "" {
		return "", fmt.Errorf("organization id missing")
	}

	return orgID, nil
}

func firstOrganizationID(organizations []*organizationsv1.Organization) string {
	for _, organization := range organizations {
		id := strings.TrimSpace(organization.GetId())
		if id != "" {
			return id
		}
	}
	return ""
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
		if zitiAppID != "" && zitiServiceID != "" {
			_, _ = client.DeleteAppIdentity(cleanupCtx, &zitimgmtv1.DeleteAppIdentityRequest{
				IdentityId:    zitiAppID,
				ZitiServiceId: zitiServiceID,
			})
		}
		if zitiContext != nil {
			zitiContext.Close()
		}
		_ = conn.Close()
	})

	createCtx, createCancel := context.WithTimeout(ctx, 30*time.Second)
	defer createCancel()

	appIdentityID := uuid.NewString()
	resp, err := client.CreateAppIdentity(createCtx, &zitimgmtv1.CreateAppIdentityRequest{
		IdentityId: appIdentityID,
		Slug:       "e2e-test",
	})
	if err != nil {
		return fmt.Errorf("create app identity: %w", err)
	}

	zitiIdentityID = strings.TrimSpace(resp.GetZitiIdentityId())
	zitiAppID = appIdentityID
	identityJSON := resp.GetIdentityJson()
	if zitiIdentityID == "" {
		return fmt.Errorf("ziti identity id missing")
	}
	if len(identityJSON) == 0 {
		return fmt.Errorf("ziti identity json missing")
	}

	serviceCtx, serviceCancel := context.WithTimeout(ctx, 15*time.Second)
	defer serviceCancel()
	serviceID, err := resolveZitiServiceID(serviceCtx, client, zitiIdentityID)
	if err != nil {
		if errors.Is(err, errZitiServiceIDNotReady) || errors.Is(err, context.DeadlineExceeded) {
			zitiServiceID = ""
		} else {
			return err
		}
	} else {
		zitiServiceID = serviceID
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

func resolveZitiServiceID(ctx context.Context, client zitimgmtv1.ZitiManagementServiceClient, zitiIdentityID string) (string, error) {
	trimmedIdentity := strings.TrimSpace(zitiIdentityID)
	if trimmedIdentity == "" {
		return "", fmt.Errorf("ziti identity id missing")
	}

	deadline := time.Now().Add(15 * time.Second)
	for {
		serviceID, err := lookupZitiServiceID(ctx, client, trimmedIdentity)
		if err == nil {
			return serviceID, nil
		}
		if !errors.Is(err, errZitiServiceIDNotReady) {
			return "", err
		}
		if time.Now().After(deadline) {
			return "", errZitiServiceIDNotReady
		}
		timer := time.NewTimer(time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return "", fmt.Errorf("resolve ziti service id: %w", ctx.Err())
		case <-timer.C:
		}
	}
}

func lookupZitiServiceID(ctx context.Context, client zitimgmtv1.ZitiManagementServiceClient, zitiIdentityID string) (string, error) {
	pageToken := ""
	for {
		resp, err := client.ListManagedIdentities(ctx, &zitimgmtv1.ListManagedIdentitiesRequest{
			IdentityType: identityv1.IdentityType_IDENTITY_TYPE_APP,
			PageSize:     200,
			PageToken:    pageToken,
		})
		if err != nil {
			return "", fmt.Errorf("list managed identities: %w", err)
		}

		for _, identity := range resp.GetIdentities() {
			if strings.TrimSpace(identity.GetZitiIdentityId()) == zitiIdentityID {
				serviceID := strings.TrimSpace(identity.GetZitiServiceId())
				if serviceID == "" {
					return "", errZitiServiceIDNotReady
				}
				return serviceID, nil
			}
		}

		pageToken = strings.TrimSpace(resp.GetNextPageToken())
		if pageToken == "" {
			break
		}
	}

	return "", errZitiServiceIDNotReady
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
