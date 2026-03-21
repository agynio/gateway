package oidcresolver

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/oidcauth"
	jose "github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestResolveFromTokenUsesExistingUser(t *testing.T) {
	provider := newOIDCTestProvider(t)
	verifier, err := oidcauth.NewVerifier(context.Background(), provider.issuer, provider.clientID)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	usersClient := &fakeUsersClient{
		getBySubjectResp: &usersv1.GetUserByOIDCSubjectResponse{User: buildUser("user-1")},
	}
	resolver, err := NewResolver(verifier, usersClient, provider.server.Client())
	if err != nil {
		t.Fatalf("failed to create resolver: %v", err)
	}

	resolved, err := resolver.ResolveFromToken(context.Background(), provider.token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertResolvedIdentity(t, resolved, "user-1")
	if usersClient.getBySubjectCalls != 1 {
		t.Fatalf("expected get-by-subject to be called once, got %d", usersClient.getBySubjectCalls)
	}
	if usersClient.lastSubject != provider.subject {
		t.Fatalf("expected subject %q, got %q", provider.subject, usersClient.lastSubject)
	}
	if usersClient.resolveCalls != 0 {
		t.Fatalf("expected resolve-or-create not to be called, got %d", usersClient.resolveCalls)
	}
	if provider.userinfoCalls != 0 {
		t.Fatalf("expected userinfo not to be called, got %d", provider.userinfoCalls)
	}
}

func TestResolveFromTokenCreatesUserOnFirstLogin(t *testing.T) {
	provider := newOIDCTestProvider(t)
	verifier, err := oidcauth.NewVerifier(context.Background(), provider.issuer, provider.clientID)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	usersClient := &fakeUsersClient{
		getBySubjectErr: status.Error(codes.NotFound, "not found"),
		resolveResp:     &usersv1.ResolveOrCreateUserResponse{User: buildUser("user-2")},
	}
	resolver, err := NewResolver(verifier, usersClient, provider.server.Client())
	if err != nil {
		t.Fatalf("failed to create resolver: %v", err)
	}

	resolved, err := resolver.ResolveFromToken(context.Background(), provider.token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertResolvedIdentity(t, resolved, "user-2")
	if usersClient.getBySubjectCalls != 1 {
		t.Fatalf("expected get-by-subject to be called once, got %d", usersClient.getBySubjectCalls)
	}
	if usersClient.resolveCalls != 1 {
		t.Fatalf("expected resolve-or-create to be called once, got %d", usersClient.resolveCalls)
	}
	if usersClient.lastResolve == nil {
		t.Fatalf("expected resolve-or-create request to be captured")
	}
	if usersClient.lastResolve.OidcSubject != provider.userInfo.Subject {
		t.Fatalf("expected resolve subject %q, got %q", provider.userInfo.Subject, usersClient.lastResolve.OidcSubject)
	}
	if usersClient.lastResolve.Name != provider.userInfo.Name {
		t.Fatalf("expected resolve name %q, got %q", provider.userInfo.Name, usersClient.lastResolve.Name)
	}
	if usersClient.lastResolve.Email != provider.userInfo.Email {
		t.Fatalf("expected resolve email %q, got %q", provider.userInfo.Email, usersClient.lastResolve.Email)
	}
	if usersClient.lastResolve.PhotoUrl != provider.userInfo.Picture {
		t.Fatalf("expected resolve photo %q, got %q", provider.userInfo.Picture, usersClient.lastResolve.PhotoUrl)
	}
	if provider.userinfoCalls != 1 {
		t.Fatalf("expected userinfo to be called once, got %d", provider.userinfoCalls)
	}
	if provider.lastAuthHeader != "Bearer "+provider.token {
		t.Fatalf("expected bearer token header, got %q", provider.lastAuthHeader)
	}
}

type oidcTestProvider struct {
	issuer           string
	clientID         string
	subject          string
	token            string
	userinfoCalls    int
	lastAuthHeader   string
	userInfo         oidc.UserInfo
	server           *httptest.Server
	userinfoEndpoint string
}

func newOIDCTestProvider(t *testing.T) *oidcTestProvider {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	keyID := "test-key"
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: key},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", keyID),
	)
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}
	publicKey := jose.JSONWebKey{Key: &key.PublicKey, KeyID: keyID, Use: oidc.KeyUseSignature, Algorithm: string(jose.RS256)}
	jwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{publicKey}}

	provider := &oidcTestProvider{
		clientID: "client-id",
		subject:  "user-123",
		userInfo: oidc.UserInfo{
			Subject: "user-123",
			UserInfoProfile: oidc.UserInfoProfile{
				Name:    "Test User",
				Picture: "https://example.com/photo.png",
			},
			UserInfoEmail: oidc.UserInfoEmail{
				Email: "test@example.com",
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(oidc.DiscoveryEndpoint, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"issuer":            provider.issuer,
			"jwks_uri":          provider.issuer + "/jwks",
			"userinfo_endpoint": provider.issuer + "/userinfo",
		})
	})
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jwks)
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		provider.userinfoCalls++
		provider.lastAuthHeader = r.Header.Get("Authorization")
		_ = json.NewEncoder(w).Encode(provider.userInfo)
	})
	provider.server = httptest.NewServer(mux)
	t.Cleanup(provider.server.Close)
	provider.issuer = provider.server.URL
	provider.userinfoEndpoint = provider.issuer + "/userinfo"

	claims := oidc.NewAccessTokenClaims(provider.issuer, provider.subject, []string{provider.clientID}, time.Now().Add(time.Hour), "jwtid", provider.clientID, time.Second)
	provider.token = signToken(t, signer, claims)

	return provider
}

type fakeUsersClient struct {
	getBySubjectResp  *usersv1.GetUserByOIDCSubjectResponse
	getBySubjectErr   error
	resolveResp       *usersv1.ResolveOrCreateUserResponse
	resolveErr        error
	getBySubjectCalls int
	resolveCalls      int
	lastSubject       string
	lastResolve       *usersv1.ResolveOrCreateUserRequest
}

func (f *fakeUsersClient) ResolveOrCreateUser(ctx context.Context, in *usersv1.ResolveOrCreateUserRequest, opts ...grpc.CallOption) (*usersv1.ResolveOrCreateUserResponse, error) {
	f.resolveCalls++
	f.lastResolve = in
	return f.resolveResp, f.resolveErr
}

func (f *fakeUsersClient) GetUser(ctx context.Context, in *usersv1.GetUserRequest, opts ...grpc.CallOption) (*usersv1.GetUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) GetUserByOIDCSubject(ctx context.Context, in *usersv1.GetUserByOIDCSubjectRequest, opts ...grpc.CallOption) (*usersv1.GetUserByOIDCSubjectResponse, error) {
	f.getBySubjectCalls++
	f.lastSubject = in.GetOidcSubject()
	return f.getBySubjectResp, f.getBySubjectErr
}

func (f *fakeUsersClient) BatchGetUsers(ctx context.Context, in *usersv1.BatchGetUsersRequest, opts ...grpc.CallOption) (*usersv1.BatchGetUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) UpdateUser(ctx context.Context, in *usersv1.UpdateUserRequest, opts ...grpc.CallOption) (*usersv1.UpdateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func buildUser(id string) *usersv1.User {
	return &usersv1.User{Meta: &usersv1.EntityMeta{Id: id}}
}

func assertResolvedIdentity(t *testing.T, resolved identity.ResolvedIdentity, expectedID string) {
	t.Helper()
	if resolved.IdentityID != expectedID {
		t.Fatalf("expected identity id %q, got %q", expectedID, resolved.IdentityID)
	}
	if resolved.IdentityType != identity.IdentityTypeUser {
		t.Fatalf("expected identity type %q, got %q", identity.IdentityTypeUser, resolved.IdentityType)
	}
	if resolved.AuthMethod != identity.AuthMethodOIDC {
		t.Fatalf("expected auth method %q, got %q", identity.AuthMethodOIDC, resolved.AuthMethod)
	}
}

func signToken(t *testing.T, signer jose.Signer, claims any) string {
	t.Helper()
	payload, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("failed to marshal claims: %v", err)
	}
	object, err := signer.Sign(payload)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	token, err := object.CompactSerialize()
	if err != nil {
		t.Fatalf("failed to serialize token: %v", err)
	}
	return token
}
