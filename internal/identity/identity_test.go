package identity

import (
	"context"
	"testing"
)

func TestIdentityContextRoundTrip(t *testing.T) {
	ctx := context.Background()
	input := ResolvedIdentity{
		IdentityID:   "id-123",
		IdentityType: "user",
		TenantID:     "tenant-1",
		AuthMethod:   "ziti",
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
