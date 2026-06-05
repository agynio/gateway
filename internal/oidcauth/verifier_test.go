package oidcauth

import (
	"context"
	"testing"
	"time"

	"github.com/agynio/gateway/internal/oidctestutil"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func newVerifier(t *testing.T, provider *oidctestutil.Provider) *Verifier {
	t.Helper()

	verifier, err := NewVerifier(context.Background(), provider.Issuer, provider.ClientID)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	return verifier
}

func TestVerifierVerifySubject(t *testing.T) {
	provider := oidctestutil.NewProvider(t)
	verifier := newVerifier(t, provider)

	claims, err := verifier.Verify(context.Background(), provider.Token)
	if err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}
	if claims.Subject != provider.Subject {
		t.Fatalf("expected subject %q, got %q", provider.Subject, claims.Subject)
	}
	if verifier.UserinfoEndpoint() != provider.UserinfoEndpoint {
		t.Fatalf("expected userinfo endpoint %q, got %q", provider.UserinfoEndpoint, verifier.UserinfoEndpoint())
	}
}

func TestVerifierVerifyAcceptsResourceServerAudience(t *testing.T) {
	provider := oidctestutil.NewProvider(t)
	verifier := newVerifier(t, provider)

	claims := oidc.NewAccessTokenClaims(
		provider.Issuer,
		provider.Subject,
		[]string{"https://api.example.com"},
		time.Now().Add(time.Hour),
		"jwtid",
		provider.ClientID,
		time.Second,
	)
	token := provider.SignAccessToken(t, claims)

	if _, err := verifier.Verify(context.Background(), token); err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}
}

func TestVerifierVerifyRequiresConfiguredAudience(t *testing.T) {
	provider := oidctestutil.NewProvider(t)

	verifier, err := NewVerifierWithAudience(context.Background(), provider.Issuer, provider.ClientID, "https://api.example.com")
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	claims := oidc.NewAccessTokenClaims(
		provider.Issuer,
		provider.Subject,
		[]string{"https://other.example.com"},
		time.Now().Add(time.Hour),
		"jwtid",
		provider.ClientID,
		time.Second,
	)
	token := provider.SignAccessToken(t, claims)

	if _, err := verifier.Verify(context.Background(), token); err == nil {
		t.Fatal("expected audience validation error")
	}
}

func TestVerifierVerifyAcceptsConfiguredAudience(t *testing.T) {
	provider := oidctestutil.NewProvider(t)

	verifier, err := NewVerifierWithAudience(context.Background(), provider.Issuer, provider.ClientID, "https://api.example.com")
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	claims := oidc.NewAccessTokenClaims(
		provider.Issuer,
		provider.Subject,
		[]string{"https://api.example.com"},
		time.Now().Add(time.Hour),
		"jwtid",
		provider.ClientID,
		time.Second,
	)
	token := provider.SignAccessToken(t, claims)

	if _, err := verifier.Verify(context.Background(), token); err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}
}

func TestVerifierVerifyAcceptsEmptyAudience(t *testing.T) {
	provider := oidctestutil.NewProvider(t)
	verifier := newVerifier(t, provider)

	now := time.Now().UTC().Add(-time.Second)
	claims := &oidc.AccessTokenClaims{
		TokenClaims: oidc.TokenClaims{
			Issuer:     provider.Issuer,
			Subject:    provider.Subject,
			Audience:   oidc.Audience{},
			Expiration: oidc.FromTime(now.Add(time.Hour)),
			IssuedAt:   oidc.FromTime(now),
			NotBefore:  oidc.FromTime(now),
			ClientID:   provider.ClientID,
			JWTID:      "jwtid",
		},
	}
	token := provider.SignAccessToken(t, claims)

	if _, err := verifier.Verify(context.Background(), token); err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}
}
