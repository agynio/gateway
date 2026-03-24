package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/agynio/gateway/internal/clusteradminresolver"
	"github.com/agynio/gateway/internal/httpauth"
	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/ziticonn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fakeZitiResolver struct {
	resolved   identity.ResolvedIdentity
	err        error
	calls      int
	lastSource string
}

func (f *fakeZitiResolver) ResolveIdentity(ctx context.Context, sourceIdentity string) (identity.ResolvedIdentity, error) {
	f.calls++
	f.lastSource = sourceIdentity
	if f.err != nil {
		return identity.ResolvedIdentity{}, f.err
	}
	return f.resolved, nil
}

type fakeOIDCResolver struct {
	resolved  identity.ResolvedIdentity
	err       error
	calls     int
	lastToken string
}

func (f *fakeOIDCResolver) ResolveFromToken(ctx context.Context, accessToken string) (identity.ResolvedIdentity, error) {
	f.calls++
	f.lastToken = accessToken
	if f.err != nil {
		return identity.ResolvedIdentity{}, f.err
	}
	return f.resolved, nil
}

type fakeAPITokenResolver struct {
	resolved  identity.ResolvedIdentity
	err       error
	calls     int
	lastToken string
}

func (f *fakeAPITokenResolver) ResolveFromToken(ctx context.Context, accessToken string) (identity.ResolvedIdentity, error) {
	f.calls++
	f.lastToken = accessToken
	if f.err != nil {
		return identity.ResolvedIdentity{}, f.err
	}
	return f.resolved, nil
}

func TestAuthInterceptorWrapUnarySuccess(t *testing.T) {
	resolver := &fakeZitiResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-1",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	interceptor := authInterceptor{zitiResolver: resolver}

	called := false
	wrapped := interceptor.WrapUnary(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		called = true
		resolved, ok := identity.IdentityFromContext(ctx)
		if !ok {
			t.Fatalf("expected identity in context")
		}
		if resolved.IdentityID != resolver.resolved.IdentityID {
			t.Fatalf("expected identity %q, got %q", resolver.resolved.IdentityID, resolved.IdentityID)
		}
		return connect.NewResponse(&emptypb.Empty{}), nil
	})

	ctx := ziticonn.WithSourceIdentity(context.Background(), "source-identity")
	resp, err := wrapped(ctx, connect.NewRequest(&emptypb.Empty{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if !called {
		t.Fatalf("expected next handler to be called")
	}
	if resolver.calls != 1 {
		t.Fatalf("expected resolver to be called once, got %d", resolver.calls)
	}
}

func TestAuthInterceptorWrapUnaryError(t *testing.T) {
	resolver := &fakeZitiResolver{err: status.Error(codes.Unauthenticated, "missing")}
	interceptor := authInterceptor{zitiResolver: resolver}

	called := false
	wrapped := interceptor.WrapUnary(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		called = true
		return connect.NewResponse(&emptypb.Empty{}), nil
	})

	ctx := ziticonn.WithSourceIdentity(context.Background(), "source-identity")
	_, err := wrapped(ctx, connect.NewRequest(&emptypb.Empty{}))
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Fatalf("expected CodeUnauthenticated, got %v", connect.CodeOf(err))
	}
	if called {
		t.Fatalf("expected next handler not to be called")
	}
}

func TestAuthInterceptorWrapUnaryBearer(t *testing.T) {
	resolver := &fakeOIDCResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-oidc",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	interceptor := authInterceptor{oidcResolver: resolver}

	called := false
	wrapped := interceptor.WrapUnary(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		called = true
		resolved, ok := identity.IdentityFromContext(ctx)
		if !ok {
			t.Fatalf("expected identity in context")
		}
		if resolved.IdentityID != resolver.resolved.IdentityID {
			t.Fatalf("expected identity %q, got %q", resolver.resolved.IdentityID, resolved.IdentityID)
		}
		return connect.NewResponse(&emptypb.Empty{}), nil
	})

	req := connect.NewRequest(&emptypb.Empty{})
	req.Header().Set("Authorization", "Bearer token-abc")
	resp, err := wrapped(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if !called {
		t.Fatalf("expected next handler to be called")
	}
	if resolver.calls != 1 {
		t.Fatalf("expected resolver to be called once, got %d", resolver.calls)
	}
	if resolver.lastToken != "token-abc" {
		t.Fatalf("expected token to be propagated, got %q", resolver.lastToken)
	}
}

func TestRecoveryInterceptorWrapUnary(t *testing.T) {
	interceptor := NewRecoveryInterceptor()
	wrapped := interceptor.WrapUnary(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		panic("boom")
	})

	_, err := wrapped(context.Background(), connect.NewRequest(&emptypb.Empty{}))
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeInternal {
		t.Fatalf("expected CodeInternal, got %v", connect.CodeOf(err))
	}
}

func TestNewAuthMiddlewareOptionsPassthrough(t *testing.T) {
	resolver := &fakeZitiResolver{}
	middleware := NewAuthMiddleware(resolver, nil, nil, nil)

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodOptions, "/me", nil)
	resp := httptest.NewRecorder()
	middleware(handler).ServeHTTP(resp, req)

	if !called {
		t.Fatalf("expected handler to be called")
	}
	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, resp.Code)
	}
	if resolver.calls != 0 {
		t.Fatalf("expected resolver to be skipped for OPTIONS, got %d", resolver.calls)
	}
}

func TestNewAuthMiddlewareMissingIdentity(t *testing.T) {
	resolver := &fakeZitiResolver{}
	middleware := NewAuthMiddleware(resolver, nil, nil, nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := identity.IdentityFromContext(r.Context()); ok {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	resp := httptest.NewRecorder()
	middleware(handler).ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.Code)
	}
	if resolver.calls != 0 {
		t.Fatalf("expected resolver to not be called, got %d", resolver.calls)
	}
}

func TestNewAuthMiddlewareSuccess(t *testing.T) {
	resolver := &fakeZitiResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-2",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	middleware := NewAuthMiddleware(resolver, nil, nil, nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := identity.IdentityFromContext(r.Context()); !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(ziticonn.WithSourceIdentity(req.Context(), "source-identity"))
	resp := httptest.NewRecorder()
	middleware(handler).ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if resolver.calls != 1 {
		t.Fatalf("expected resolver to be called once, got %d", resolver.calls)
	}
}

func TestNewAuthMiddlewareBearer(t *testing.T) {
	resolver := &fakeOIDCResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-oidc",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	middleware := NewAuthMiddleware(nil, resolver, nil, nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := identity.IdentityFromContext(r.Context()); !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer token-xyz")
	resp := httptest.NewRecorder()
	middleware(handler).ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if resolver.calls != 1 {
		t.Fatalf("expected resolver to be called once, got %d", resolver.calls)
	}
}

func TestResolveIdentityWithoutSource(t *testing.T) {
	zitiResolver := &fakeZitiResolver{}
	oidcResolver := &fakeOIDCResolver{}
	ctx, err := resolveIdentity(context.Background(), zitiResolver, oidcResolver, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := identity.IdentityFromContext(ctx); ok {
		t.Fatalf("expected no identity in context")
	}
	if zitiResolver.calls != 0 {
		t.Fatalf("expected ziti resolver not to be called, got %d", zitiResolver.calls)
	}
	if oidcResolver.calls != 0 {
		t.Fatalf("expected oidc resolver not to be called, got %d", oidcResolver.calls)
	}
}

func TestResolveIdentityWithSource(t *testing.T) {
	resolver := &fakeZitiResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-3",
			IdentityType: identity.IdentityTypeAgent,
		},
	}
	ctx := ziticonn.WithSourceIdentity(context.Background(), "source-identity")
	ctx, err := resolveIdentity(ctx, resolver, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if resolved.IdentityID != resolver.resolved.IdentityID {
		t.Fatalf("expected identity %q, got %q", resolver.resolved.IdentityID, resolved.IdentityID)
	}
	if resolver.calls != 1 {
		t.Fatalf("expected resolver to be called once, got %d", resolver.calls)
	}
}

func TestResolveIdentityWithBearer(t *testing.T) {
	resolver := &fakeOIDCResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-bearer",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	ctx := httpauth.WithBearerToken(context.Background(), "token-bearer")
	ctx, err := resolveIdentity(ctx, nil, resolver, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if resolved.IdentityID != resolver.resolved.IdentityID {
		t.Fatalf("expected identity %q, got %q", resolver.resolved.IdentityID, resolved.IdentityID)
	}
	if resolver.calls != 1 {
		t.Fatalf("expected resolver to be called once, got %d", resolver.calls)
	}
	if resolver.lastToken != "token-bearer" {
		t.Fatalf("expected token to be propagated, got %q", resolver.lastToken)
	}
}

func TestResolveIdentityWithClusterAdminToken(t *testing.T) {
	clusterResolver, err := clusteradminresolver.NewResolver("agyn_cluster", "identity-cluster")
	if err != nil {
		t.Fatalf("failed to create cluster admin resolver: %v", err)
	}
	apiResolver := &fakeAPITokenResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-api",
			IdentityType: identity.IdentityTypeUser,
		},
	}

	ctx := httpauth.WithBearerToken(context.Background(), "agyn_cluster")
	ctx, err = resolveIdentity(ctx, nil, nil, apiResolver, clusterResolver)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if resolved.IdentityID != "identity-cluster" {
		t.Fatalf("expected identity %q, got %q", "identity-cluster", resolved.IdentityID)
	}
	if apiResolver.calls != 0 {
		t.Fatalf("expected api token resolver not to be called, got %d", apiResolver.calls)
	}
}

func TestResolveIdentityClusterAdminFallbackToAPIToken(t *testing.T) {
	clusterResolver, err := clusteradminresolver.NewResolver("cluster-token", "identity-cluster")
	if err != nil {
		t.Fatalf("failed to create cluster admin resolver: %v", err)
	}
	apiResolver := &fakeAPITokenResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-api",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	oidcResolver := &fakeOIDCResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-oidc",
			IdentityType: identity.IdentityTypeUser,
		},
	}

	ctx := httpauth.WithBearerToken(context.Background(), "agyn_token")
	ctx, err = resolveIdentity(ctx, nil, oidcResolver, apiResolver, clusterResolver)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if resolved.IdentityID != apiResolver.resolved.IdentityID {
		t.Fatalf("expected identity %q, got %q", apiResolver.resolved.IdentityID, resolved.IdentityID)
	}
	if apiResolver.calls != 1 {
		t.Fatalf("expected api token resolver to be called once, got %d", apiResolver.calls)
	}
	if oidcResolver.calls != 0 {
		t.Fatalf("expected oidc resolver not to be called, got %d", oidcResolver.calls)
	}
}

func TestResolveIdentityClusterAdminResolverNil(t *testing.T) {
	oidcResolver := &fakeOIDCResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-oidc",
			IdentityType: identity.IdentityTypeUser,
		},
	}

	ctx := httpauth.WithBearerToken(context.Background(), "token-oidc")
	ctx, err := resolveIdentity(ctx, nil, oidcResolver, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if resolved.IdentityID != oidcResolver.resolved.IdentityID {
		t.Fatalf("expected identity %q, got %q", oidcResolver.resolved.IdentityID, resolved.IdentityID)
	}
	if oidcResolver.calls != 1 {
		t.Fatalf("expected oidc resolver to be called once, got %d", oidcResolver.calls)
	}
}

func TestResolveIdentityWithAPIToken(t *testing.T) {
	apiResolver := &fakeAPITokenResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-api",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	oidcResolver := &fakeOIDCResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-oidc",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	ctx := httpauth.WithBearerToken(context.Background(), "agyn_token")
	ctx, err := resolveIdentity(ctx, nil, oidcResolver, apiResolver, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if resolved.IdentityID != apiResolver.resolved.IdentityID {
		t.Fatalf("expected identity %q, got %q", apiResolver.resolved.IdentityID, resolved.IdentityID)
	}
	if apiResolver.calls != 1 {
		t.Fatalf("expected api token resolver to be called once, got %d", apiResolver.calls)
	}
	if apiResolver.lastToken != "agyn_token" {
		t.Fatalf("expected token to be propagated, got %q", apiResolver.lastToken)
	}
	if oidcResolver.calls != 0 {
		t.Fatalf("expected oidc resolver not to be called, got %d", oidcResolver.calls)
	}
}

func TestResolveIdentityOIDCTokenNotPrefixed(t *testing.T) {
	apiResolver := &fakeAPITokenResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-api",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	oidcResolver := &fakeOIDCResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-oidc",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	ctx := httpauth.WithBearerToken(context.Background(), "token-oidc")
	ctx, err := resolveIdentity(ctx, nil, oidcResolver, apiResolver, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if resolved.IdentityID != oidcResolver.resolved.IdentityID {
		t.Fatalf("expected identity %q, got %q", oidcResolver.resolved.IdentityID, resolved.IdentityID)
	}
	if oidcResolver.calls != 1 {
		t.Fatalf("expected oidc resolver to be called once, got %d", oidcResolver.calls)
	}
	if oidcResolver.lastToken != "token-oidc" {
		t.Fatalf("expected token to be propagated, got %q", oidcResolver.lastToken)
	}
	if apiResolver.calls != 0 {
		t.Fatalf("expected api token resolver not to be called, got %d", apiResolver.calls)
	}
}

func TestResolveIdentityAPITokenResolverNil(t *testing.T) {
	oidcResolver := &fakeOIDCResolver{}
	ctx := httpauth.WithBearerToken(context.Background(), "agyn_missing")
	_, err := resolveIdentity(ctx, nil, oidcResolver, nil, nil)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected unauthenticated error, got %v", err)
	}
	if oidcResolver.calls != 0 {
		t.Fatalf("expected oidc resolver not to be called, got %d", oidcResolver.calls)
	}
}

func TestResolveIdentityAPITokenError(t *testing.T) {
	apiResolver := &fakeAPITokenResolver{err: status.Error(codes.Internal, "boom")}
	oidcResolver := &fakeOIDCResolver{}
	ctx := httpauth.WithBearerToken(context.Background(), "agyn_error")
	_, err := resolveIdentity(ctx, nil, oidcResolver, apiResolver, nil)
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected internal error, got %v", err)
	}
	if apiResolver.calls != 1 {
		t.Fatalf("expected api token resolver to be called once, got %d", apiResolver.calls)
	}
	if oidcResolver.calls != 0 {
		t.Fatalf("expected oidc resolver not to be called, got %d", oidcResolver.calls)
	}
}

func TestResolveIdentityPrefersZiti(t *testing.T) {
	zitiResolver := &fakeZitiResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-ziti",
			IdentityType: identity.IdentityTypeUser,
		},
	}
	oidcResolver := &fakeOIDCResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-oidc",
			IdentityType: identity.IdentityTypeUser,
		},
	}

	ctx := ziticonn.WithSourceIdentity(context.Background(), "source-identity")
	ctx = httpauth.WithBearerToken(ctx, "token-preferred")
	ctx, err := resolveIdentity(ctx, zitiResolver, oidcResolver, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if resolved.IdentityID != zitiResolver.resolved.IdentityID {
		t.Fatalf("expected identity %q, got %q", zitiResolver.resolved.IdentityID, resolved.IdentityID)
	}
	if zitiResolver.calls != 1 {
		t.Fatalf("expected ziti resolver to be called once, got %d", zitiResolver.calls)
	}
	if oidcResolver.calls != 0 {
		t.Fatalf("expected oidc resolver not to be called, got %d", oidcResolver.calls)
	}
}
