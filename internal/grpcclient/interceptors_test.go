package grpcclient

import (
	"context"
	"testing"

	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc/metadata"
)

func TestAppendIdentityMetadataWithIdentity(t *testing.T) {
	ctx := identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   "identity-123",
		IdentityType: identity.IdentityTypeUser,
		WorkloadID:   "workload-123",
	})

	ctx = appendIdentityMetadata(ctx)
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("expected outgoing metadata")
	}

	identityIDs := md.Get(identity.MetadataKeyIdentityID)
	if len(identityIDs) != 1 || identityIDs[0] != "identity-123" {
		t.Fatalf("expected identity id metadata, got %v", identityIDs)
	}

	identityTypes := md.Get(identity.MetadataKeyIdentityType)
	if len(identityTypes) != 1 || identityTypes[0] != string(identity.IdentityTypeUser) {
		t.Fatalf("expected identity type metadata, got %v", identityTypes)
	}

	workloadIDs := md.Get(identity.MetadataKeyWorkloadID)
	if len(workloadIDs) != 1 || workloadIDs[0] != "workload-123" {
		t.Fatalf("expected workload id metadata, got %v", workloadIDs)
	}
}

func TestAppendIdentityMetadataWithoutIdentity(t *testing.T) {
	ctx := appendIdentityMetadata(context.Background())
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		t.Fatalf("expected no outgoing metadata, got %v", md)
	}
}

func TestAppendIdentityMetadataPreservesExisting(t *testing.T) {
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("existing", "value"))
	ctx = identity.WithIdentity(ctx, identity.ResolvedIdentity{
		IdentityID:   "identity-456",
		IdentityType: identity.IdentityTypeAgent,
		WorkloadID:   "workload-456",
	})

	ctx = appendIdentityMetadata(ctx)
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("expected outgoing metadata")
	}

	existing := md.Get("existing")
	if len(existing) != 1 || existing[0] != "value" {
		t.Fatalf("expected existing metadata to be preserved, got %v", existing)
	}

	identityIDs := md.Get(identity.MetadataKeyIdentityID)
	if len(identityIDs) != 1 || identityIDs[0] != "identity-456" {
		t.Fatalf("expected identity id metadata, got %v", identityIDs)
	}

	identityTypes := md.Get(identity.MetadataKeyIdentityType)
	if len(identityTypes) != 1 || identityTypes[0] != string(identity.IdentityTypeAgent) {
		t.Fatalf("expected identity type metadata, got %v", identityTypes)
	}

	workloadIDs := md.Get(identity.MetadataKeyWorkloadID)
	if len(workloadIDs) != 1 || workloadIDs[0] != "workload-456" {
		t.Fatalf("expected workload id metadata, got %v", workloadIDs)
	}
}

func TestAppendIdentityMetadataReplacesExisting(t *testing.T) {
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
		identity.MetadataKeyIdentityID, "stale-identity",
		identity.MetadataKeyIdentityType, string(identity.IdentityTypeRunner),
		identity.MetadataKeyWorkloadID, "stale-workload",
	))
	ctx = identity.WithIdentity(ctx, identity.ResolvedIdentity{
		IdentityID:   "identity-789",
		IdentityType: identity.IdentityTypeUser,
	})

	ctx = appendIdentityMetadata(ctx)
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("expected outgoing metadata")
	}

	identityIDs := md.Get(identity.MetadataKeyIdentityID)
	if len(identityIDs) != 1 || identityIDs[0] != "identity-789" {
		t.Fatalf("expected one identity id metadata value, got %v", identityIDs)
	}

	identityTypes := md.Get(identity.MetadataKeyIdentityType)
	if len(identityTypes) != 1 || identityTypes[0] != string(identity.IdentityTypeUser) {
		t.Fatalf("expected one identity type metadata value, got %v", identityTypes)
	}

	workloadIDs := md.Get(identity.MetadataKeyWorkloadID)
	if len(workloadIDs) != 0 {
		t.Fatalf("expected workload metadata to be removed, got %v", workloadIDs)
	}
}
