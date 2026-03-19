package identity

import "context"

type ResolvedIdentity struct {
	IdentityID   string
	IdentityType string
	TenantID     string
	AuthMethod   string
}

type contextKey struct{}

func WithIdentity(ctx context.Context, identity ResolvedIdentity) context.Context {
	return context.WithValue(ctx, contextKey{}, identity)
}

func IdentityFromContext(ctx context.Context) (ResolvedIdentity, bool) {
	identity, ok := ctx.Value(contextKey{}).(ResolvedIdentity)
	return identity, ok
}
