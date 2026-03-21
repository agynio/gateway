package identity

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	MetadataKeyTenantID     = "x-agyn-tenant-id"
	MetadataKeyIdentityID   = "x-agyn-identity-id"
	MetadataKeyIdentityType = "x-agyn-identity-type"
	MetadataKeyAuthMethod   = "x-agyn-auth-method"
)

func AppendToOutgoingContext(ctx context.Context) context.Context {
	resolved, ok := IdentityFromContext(ctx)
	if !ok {
		return ctx
	}

	return metadata.AppendToOutgoingContext(
		ctx,
		MetadataKeyTenantID, resolved.TenantID,
		MetadataKeyIdentityID, resolved.IdentityID,
		MetadataKeyIdentityType, string(resolved.IdentityType),
		MetadataKeyAuthMethod, string(resolved.AuthMethod),
	)
}

func IdentityFromIncomingContext(ctx context.Context) (ResolvedIdentity, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ResolvedIdentity{}, fmt.Errorf("missing identity metadata")
	}

	tenantID, err := requiredMetadataValue(md, MetadataKeyTenantID)
	if err != nil {
		return ResolvedIdentity{}, err
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

	authMethodValue, err := requiredMetadataValue(md, MetadataKeyAuthMethod)
	if err != nil {
		return ResolvedIdentity{}, err
	}

	return ResolvedIdentity{
		TenantID:     tenantID,
		IdentityID:   identityID,
		IdentityType: identityType,
		AuthMethod:   AuthMethod(authMethodValue),
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
