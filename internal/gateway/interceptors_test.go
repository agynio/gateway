package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/ziticonn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fakeResolver struct {
	resolved   identity.ResolvedIdentity
	err        error
	calls      int
	lastSource string
}

func (f *fakeResolver) ResolveIdentity(ctx context.Context, sourceIdentity string) (identity.ResolvedIdentity, error) {
	f.calls++
	f.lastSource = sourceIdentity
	if f.err != nil {
		return identity.ResolvedIdentity{}, f.err
	}
	return f.resolved, nil
}

func TestAuthInterceptorWrapUnarySuccess(t *testing.T) {
	resolver := &fakeResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-1",
			IdentityType: identity.IdentityTypeUser,
			TenantID:     "tenant-1",
			AuthMethod:   identity.AuthMethodOIDC,
		},
	}
	interceptor := authInterceptor{resolver: resolver}

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
	resolver := &fakeResolver{err: status.Error(codes.Unauthenticated, "missing")}
	interceptor := authInterceptor{resolver: resolver}

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
	resolver := &fakeResolver{}
	middleware := NewAuthMiddleware(resolver)

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
	resolver := &fakeResolver{}
	middleware := NewAuthMiddleware(resolver)

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
	resolver := &fakeResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-2",
			IdentityType: identity.IdentityTypeUser,
			TenantID:     "tenant-2",
			AuthMethod:   identity.AuthMethodZiti,
		},
	}
	middleware := NewAuthMiddleware(resolver)

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

func TestResolveIdentityWithoutSource(t *testing.T) {
	resolver := &fakeResolver{}
	ctx, err := resolveIdentity(context.Background(), resolver)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := identity.IdentityFromContext(ctx); ok {
		t.Fatalf("expected no identity in context")
	}
	if resolver.calls != 0 {
		t.Fatalf("expected resolver not to be called, got %d", resolver.calls)
	}
}

func TestResolveIdentityWithSource(t *testing.T) {
	resolver := &fakeResolver{
		resolved: identity.ResolvedIdentity{
			IdentityID:   "identity-3",
			IdentityType: identity.IdentityTypeAgent,
			TenantID:     "tenant-3",
			AuthMethod:   identity.AuthMethodZiti,
		},
	}
	ctx := ziticonn.WithSourceIdentity(context.Background(), "source-identity")
	ctx, err := resolveIdentity(ctx, resolver)
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
