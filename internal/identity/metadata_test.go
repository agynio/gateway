package identity

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestAppendToOutgoingContext(t *testing.T) {
	input := ResolvedIdentity{
		IdentityID:   "identity-123",
		IdentityType: IdentityTypeAgent,
		TenantID:     "tenant-9",
		AuthMethod:   AuthMethodOIDC,
	}

	ctx := WithIdentity(context.Background(), input)
	ctx = AppendToOutgoingContext(ctx)

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("expected outgoing metadata")
	}

	assertMetadataValue(t, md, MetadataKeyTenantID, input.TenantID)
	assertMetadataValue(t, md, MetadataKeyIdentityID, input.IdentityID)
	assertMetadataValue(t, md, MetadataKeyIdentityType, string(input.IdentityType))
	assertMetadataValue(t, md, MetadataKeyAuthMethod, string(input.AuthMethod))
}

func TestAppendToOutgoingContextMissingIdentity(t *testing.T) {
	ctx := AppendToOutgoingContext(context.Background())

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return
	}

	if len(md) != 0 {
		t.Fatalf("expected no outgoing metadata, got %v", md)
	}
}

func TestIdentityFromIncomingContext(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		MetadataKeyTenantID, "tenant-7",
		MetadataKeyIdentityID, "identity-777",
		MetadataKeyIdentityType, string(IdentityTypeUser),
		MetadataKeyAuthMethod, string(AuthMethodZiti),
	))

	got, err := IdentityFromIncomingContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := ResolvedIdentity{
		IdentityID:   "identity-777",
		IdentityType: IdentityTypeUser,
		TenantID:     "tenant-7",
		AuthMethod:   AuthMethodZiti,
	}
	if got != expected {
		t.Fatalf("unexpected identity: %+v", got)
	}
}

func TestIdentityFromIncomingContextMissingField(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		MetadataKeyTenantID, "tenant-7",
		MetadataKeyIdentityID, "identity-777",
		MetadataKeyAuthMethod, string(AuthMethodZiti),
	))

	_, err := IdentityFromIncomingContext(ctx)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestIdentityFromIncomingContextInvalidType(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		MetadataKeyTenantID, "tenant-7",
		MetadataKeyIdentityID, "identity-777",
		MetadataKeyIdentityType, "unknown",
		MetadataKeyAuthMethod, string(AuthMethodZiti),
	))

	_, err := IdentityFromIncomingContext(ctx)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func assertMetadataValue(t *testing.T, md metadata.MD, key, expected string) {
	t.Helper()
	values := md.Get(key)
	if len(values) != 1 {
		t.Fatalf("expected 1 value for %s, got %v", key, values)
	}
	if values[0] != expected {
		t.Fatalf("expected %s=%q, got %q", key, expected, values[0])
	}
}
