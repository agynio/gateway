package identity

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	MetadataKeyIdentityID        = "x-identity-id"
	MetadataKeyIdentityType      = "x-identity-type"
	MetadataKeyWorkloadID        = "x-workload-id"
	MetadataKeyOrganizationID    = "x-organization-id"
	MetadataKeyAppInstallationID = "x-app-installation-id"
)

func AppendToOutgoingContext(ctx context.Context) context.Context {
	resolved, ok := IdentityFromContext(ctx)
	if !ok {
		return ctx
	}

	entries := []string{
		MetadataKeyIdentityID, resolved.IdentityID,
		MetadataKeyIdentityType, string(resolved.IdentityType),
	}
	workloadID := strings.TrimSpace(resolved.WorkloadID)
	if workloadID != "" {
		entries = append(entries, MetadataKeyWorkloadID, workloadID)
	}

	return metadata.AppendToOutgoingContext(ctx, entries...)
}

func IdentityFromIncomingContext(ctx context.Context) (ResolvedIdentity, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ResolvedIdentity{}, fmt.Errorf("missing identity metadata")
	}

	identityID, err := requiredMetadataValue(md, MetadataKeyIdentityID)
	if err != nil {
		return ResolvedIdentity{}, err
	}

	identityTypeValue, err := requiredMetadataValue(md, MetadataKeyIdentityType)
	if err != nil {
		return ResolvedIdentity{}, err
	}
	identityType, err := ParseIdentityType(identityTypeValue)
	if err != nil {
		return ResolvedIdentity{}, err
	}
	workloadID := optionalMetadataValue(md, MetadataKeyWorkloadID)

	return ResolvedIdentity{
		IdentityID:   identityID,
		IdentityType: identityType,
		WorkloadID:   workloadID,
	}, nil
}

func requiredMetadataValue(md metadata.MD, key string) (string, error) {
	values := md.Get(key)
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed, nil
		}
	}
	return "", fmt.Errorf("missing %s metadata", key)
}

func optionalMetadataValue(md metadata.MD, key string) string {
	values := md.Get(key)
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
