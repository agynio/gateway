package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agynio/gateway/internal/gen"
	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/ziticonn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type stubIdentityResolver struct {
	identity identity.ResolvedIdentity
	err      error
	seen     string
}

func (s *stubIdentityResolver) ResolveIdentity(ctx context.Context, sourceIdentity string) (identity.ResolvedIdentity, error) {
	s.seen = sourceIdentity
	return s.identity, s.err
}

func TestAuthMiddlewareMissingSourceIdentity(t *testing.T) {
	resolver := &stubIdentityResolver{}
	middleware := NewAuthMiddleware(resolver)

	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	response := httptest.NewRecorder()
	called := false

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	handler.ServeHTTP(response, request)

	if called {
		t.Fatalf("expected handler not to be called")
	}
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
	if response.Header().Get("Content-Type") != problemContentType {
		t.Fatalf("unexpected content type: %s", response.Header().Get("Content-Type"))
	}

	var problemResponse gen.Problem
	if err := json.NewDecoder(response.Body).Decode(&problemResponse); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if problemResponse.Status != http.StatusUnauthorized {
		t.Fatalf("unexpected problem status: %d", problemResponse.Status)
	}
}

func TestAuthMiddlewareResolvesIdentity(t *testing.T) {
	resolved := identity.ResolvedIdentity{
		IdentityID:   "id-1",
		IdentityType: "user",
		TenantID:     "tenant-1",
		AuthMethod:   "ziti",
	}
	resolver := &stubIdentityResolver{identity: resolved}
	middleware := NewAuthMiddleware(resolver)

	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	request = request.WithContext(ziticonn.WithSourceIdentity(request.Context(), "source-1"))
	response := httptest.NewRecorder()

	var got identity.ResolvedIdentity
	var ok bool
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok = identity.IdentityFromContext(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", response.Code)
	}
	if resolver.seen != "source-1" {
		t.Fatalf("unexpected source identity: %s", resolver.seen)
	}
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if got != resolved {
		t.Fatalf("unexpected resolved identity: %+v", got)
	}
}

func TestAuthMiddlewareResolverError(t *testing.T) {
	resolver := &stubIdentityResolver{err: status.Error(codes.Unauthenticated, "unauthenticated")}
	middleware := NewAuthMiddleware(resolver)

	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	request = request.WithContext(ziticonn.WithSourceIdentity(request.Context(), "source-1"))
	response := httptest.NewRecorder()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}
