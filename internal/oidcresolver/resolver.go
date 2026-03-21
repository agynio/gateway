package oidcresolver

import (
	"context"
	"fmt"
	"strings"

	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/oidcauth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Resolver struct {
	verifier    *oidcauth.Verifier
	usersClient usersv1.UsersServiceClient
}

func NewResolver(verifier *oidcauth.Verifier, usersClient usersv1.UsersServiceClient) (*Resolver, error) {
	if verifier == nil {
		return nil, fmt.Errorf("verifier is required")
	}
	if usersClient == nil {
		return nil, fmt.Errorf("users client is required")
	}

	return &Resolver{
		verifier:    verifier,
		usersClient: usersClient,
	}, nil
}

func (r *Resolver) ResolveFromToken(ctx context.Context, accessToken string) (identity.ResolvedIdentity, error) {
	claims, err := r.verifier.Verify(ctx, accessToken)
	if err != nil {
		return identity.ResolvedIdentity{}, status.Errorf(codes.Unauthenticated, "invalid bearer token: %v", err)
	}

	response, err := r.usersClient.ResolveOrCreateUser(ctx, &usersv1.ResolveOrCreateUserRequest{
		OidcSubject: claims.Subject,
		Name:        claims.Name,
		Email:       claims.Email,
		PhotoUrl:    claims.Picture,
	})
	if err != nil {
		return identity.ResolvedIdentity{}, err
	}

	user := response.GetUser()
	if user == nil {
		return identity.ResolvedIdentity{}, fmt.Errorf("user missing from response")
	}
	meta := user.GetMeta()
	if meta == nil {
		return identity.ResolvedIdentity{}, fmt.Errorf("user metadata missing from response")
	}
	identityID := strings.TrimSpace(meta.GetId())
	if identityID == "" {
		return identity.ResolvedIdentity{}, fmt.Errorf("identity id missing")
	}

	return identity.ResolvedIdentity{
		IdentityID:   identityID,
		IdentityType: identity.IdentityTypeUser,
		AuthMethod:   identity.AuthMethodOIDC,
	}, nil
}
