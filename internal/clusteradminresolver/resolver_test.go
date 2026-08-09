package clusteradminresolver

import (
	"context"
	"testing"

	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewResolverValidation(t *testing.T) {
	if _, err := NewResolver("", "identity-1"); err == nil {
		t.Fatalf("expected error for missing token")
	}
	if _, err := NewResolver("token-1", ""); err == nil {
		t.Fatalf("expected error for missing identity id")
	}
	resolver, err := NewResolver("token-1", "identity-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolver == nil {
		t.Fatalf("expected resolver to be created")
	}
}

func TestResolverMatches(t *testing.T) {
	resolver, err := NewResolver("token-1", "identity-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resolver.Matches("token-1") {
		t.Fatalf("expected token to match")
	}
	if resolver.Matches("token-2") {
		t.Fatalf("expected token to not match")
	}
}

func TestResolverResolveFromToken(t *testing.T) {
	resolver, err := NewResolver("token-1", "identity-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = resolver.ResolveFromToken(context.Background(), "token-2")
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected unauthenticated error, got %v", err)
	}

	resolved, err := resolver.ResolveFromToken(context.Background(), "token-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.IdentityID != "identity-1" {
		t.Fatalf("expected identity id %q, got %q", "identity-1", resolved.IdentityID)
	}
	// Not a user: it has no account in Users, and resolving it as one would put
	// an unnameable member in every organization it touches.
	if resolved.IdentityType != identity.IdentityTypePlatform {
		t.Fatalf("expected identity type %q, got %q", identity.IdentityTypePlatform, resolved.IdentityType)
	}
}
