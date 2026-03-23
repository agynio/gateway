package oidctestutil

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type Provider struct {
	Issuer           string
	ClientID         string
	Subject          string
	Token            string
	UserInfo         oidc.UserInfo
	UserinfoEndpoint string
	Server           *httptest.Server
	UserinfoCalls    int
	LastAuthHeader   string
	signer           jose.Signer
}

type Option func(*providerConfig)

type providerConfig struct {
	clientID       string
	subject        string
	userInfo       oidc.UserInfo
	userInfoSet    bool
	userinfoStatus int
	userinfoBody   string
}

func WithClientID(clientID string) Option {
	return func(cfg *providerConfig) {
		cfg.clientID = clientID
	}
}

func WithSubject(subject string) Option {
	return func(cfg *providerConfig) {
		cfg.subject = subject
	}
}

func WithUserInfo(userInfo oidc.UserInfo) Option {
	return func(cfg *providerConfig) {
		cfg.userInfo = userInfo
		cfg.userInfoSet = true
	}
}

func WithUserInfoResponse(status int, body string) Option {
	return func(cfg *providerConfig) {
		cfg.userinfoStatus = status
		cfg.userinfoBody = body
	}
}

func NewProvider(t *testing.T, opts ...Option) *Provider {
	t.Helper()

	cfg := providerConfig{
		clientID:       "client-id",
		subject:        "user-123",
		userinfoStatus: http.StatusOK,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if !cfg.userInfoSet {
		cfg.userInfo = defaultUserInfo(cfg.subject)
	}

	key := NewRSAKey(t)
	const keyID = "test-key"
	signer := NewSigner(t, key, keyID)
	jwks := NewJWKS(key, keyID)

	provider := &Provider{
		ClientID: cfg.clientID,
		Subject:  cfg.subject,
		UserInfo: cfg.userInfo,
		signer:   signer,
	}

	var issuer string
	mux := http.NewServeMux()
	mux.HandleFunc(oidc.DiscoveryEndpoint, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"issuer":            issuer,
			"jwks_uri":          issuer + "/jwks",
			"userinfo_endpoint": issuer + "/userinfo",
		})
	})
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jwks)
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		provider.UserinfoCalls++
		provider.LastAuthHeader = r.Header.Get("Authorization")
		if cfg.userinfoStatus != http.StatusOK {
			w.WriteHeader(cfg.userinfoStatus)
			if cfg.userinfoBody != "" {
				_, _ = io.WriteString(w, cfg.userinfoBody)
			}
			return
		}
		if cfg.userinfoBody != "" {
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, cfg.userinfoBody)
			return
		}
		_ = json.NewEncoder(w).Encode(provider.UserInfo)
	})
	provider.Server = httptest.NewServer(mux)
	t.Cleanup(provider.Server.Close)
	issuer = provider.Server.URL

	provider.Issuer = issuer
	provider.UserinfoEndpoint = issuer + "/userinfo"
	claims := oidc.NewAccessTokenClaims(
		issuer,
		provider.Subject,
		[]string{issuer},
		time.Now().Add(time.Hour),
		"jwtid",
		provider.ClientID,
		time.Second,
	)
	provider.Token = SignToken(t, signer, claims)

	return provider
}

func (p *Provider) SignAccessToken(t *testing.T, claims *oidc.AccessTokenClaims) string {
	t.Helper()

	return SignToken(t, p.signer, claims)
}

func NewRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	return key
}

func NewSigner(t *testing.T, key *rsa.PrivateKey, keyID string) jose.Signer {
	t.Helper()

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: key},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", keyID),
	)
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}
	return signer
}

func NewJWKS(key *rsa.PrivateKey, keyID string) jose.JSONWebKeySet {
	publicKey := jose.JSONWebKey{Key: &key.PublicKey, KeyID: keyID, Use: oidc.KeyUseSignature, Algorithm: string(jose.RS256)}
	return jose.JSONWebKeySet{Keys: []jose.JSONWebKey{publicKey}}
}

func SignToken(t *testing.T, signer jose.Signer, claims any) string {
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

func defaultUserInfo(subject string) oidc.UserInfo {
	return oidc.UserInfo{
		Subject: subject,
		UserInfoProfile: oidc.UserInfoProfile{
			Name:    "Test User",
			Picture: "https://example.com/photo.png",
		},
		UserInfoEmail: oidc.UserInfoEmail{
			Email: "test@example.com",
		},
	}
}
