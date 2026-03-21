package oidcauth

import (
	"context"
	"testing"

	"github.com/agynio/gateway/internal/oidctestutil"
)

func TestVerifierVerifySubject(t *testing.T) {
	provider := oidctestutil.NewProvider(t)

	verifier, err := NewVerifier(context.Background(), provider.Issuer, provider.ClientID)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

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
