package oidcauth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func TestVerifierVerifySubject(t *testing.T) {
	provider := newTestOIDCProvider(t)

	verifier, err := NewVerifier(context.Background(), provider.issuer, provider.clientID)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	claims, err := verifier.Verify(context.Background(), provider.token)
	if err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}
	if claims.Subject != provider.subject {
		t.Fatalf("expected subject %q, got %q", provider.subject, claims.Subject)
	}
	if verifier.UserinfoEndpoint() != provider.userinfoEndpoint {
		t.Fatalf("expected userinfo endpoint %q, got %q", provider.userinfoEndpoint, verifier.UserinfoEndpoint())
	}
}

type testOIDCProvider struct {
	issuer           string
	clientID         string
	subject          string
	token            string
	userinfoEndpoint string
}

func newTestOIDCProvider(t *testing.T) testOIDCProvider {
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

	clientID := "client-id"
	subject := "user-123"

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
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	issuer = server.URL

	claims := oidc.NewAccessTokenClaims(issuer, subject, []string{clientID}, time.Now().Add(time.Hour), "jwtid", clientID, time.Second)
	token := signToken(t, signer, claims)

	return testOIDCProvider{
		issuer:           issuer,
		clientID:         clientID,
		subject:          subject,
		token:            token,
		userinfoEndpoint: issuer + "/userinfo",
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
