package identity

import (
	"context"
	"testing"
)

func TestIdentityContextRoundTrip(t *testing.T) {
	ctx := context.Background()
	input := ResolvedIdentity{
		IdentityID:   "id-123",
		IdentityType: IdentityTypeUser,
		TenantID:     "tenant-1",
		AuthMethod:   AuthMethodZiti,
	}

	ctx = WithIdentity(ctx, input)

	got, ok := IdentityFromContext(ctx)
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if got != input {
		t.Fatalf("unexpected identity: %+v", got)
	}
}

func TestIdentityContextMissing(t *testing.T) {
	_, ok := IdentityFromContext(context.Background())
	if ok {
		t.Fatalf("expected no identity")
	}
}

func TestParseAuthMethod(t *testing.T) {
	method, err := ParseAuthMethod("  ziti ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != AuthMethodZiti {
		t.Fatalf("expected %q, got %q", AuthMethodZiti, method)
	}
}

func TestParseAuthMethodInvalid(t *testing.T) {
	_, err := ParseAuthMethod("banana")
	if err == nil {
		t.Fatalf("expected error")
	}
}
