package clusteradminresolver

import (
	"context"
	"crypto/subtle"
	"fmt"

	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Resolver struct {
	token      string
	identityID string
}

func NewResolver(token, identityID string) (*Resolver, error) {
	if token == "" {
		return nil, fmt.Errorf("cluster admin token is required")
	}
	if identityID == "" {
		return nil, fmt.Errorf("cluster admin identity id is required")
	}
	return &Resolver{token: token, identityID: identityID}, nil
}

func (r *Resolver) Matches(accessToken string) bool {
	return subtle.ConstantTimeCompare([]byte(accessToken), []byte(r.token)) == 1
}

func (r *Resolver) ResolveFromToken(_ context.Context, accessToken string) (identity.ResolvedIdentity, error) {
	if !r.Matches(accessToken) {
		return identity.ResolvedIdentity{}, status.Error(codes.Unauthenticated, "invalid cluster admin token")
	}
	return identity.ResolvedIdentity{
		IdentityID:   r.identityID,
		IdentityType: identity.IdentityTypePlatform,
	}, nil
}
