package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agynio/gateway/internal/identity"
)

func TestMeHandlerSuccess(t *testing.T) {
	resolved := identity.ResolvedIdentity{
		IdentityID:   "identity-1",
		IdentityType: identity.IdentityTypeUser,
		TenantID:     "tenant-1",
		AuthMethod:   identity.AuthMethodOIDC,
	}
	ctx := identity.WithIdentity(context.Background(), resolved)

	req := httptest.NewRequest(http.MethodGet, "/me", nil).WithContext(ctx)
	resp := httptest.NewRecorder()
	MeHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if resp.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected application/json content type")
	}

	var payload MeResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.IdentityID != resolved.IdentityID {
		t.Fatalf("expected identity_id %q, got %q", resolved.IdentityID, payload.IdentityID)
	}
	if payload.IdentityType != resolved.IdentityType {
		t.Fatalf("expected identity_type %q, got %q", resolved.IdentityType, payload.IdentityType)
	}
	if payload.TenantID != resolved.TenantID {
		t.Fatalf("expected tenant_id %q, got %q", resolved.TenantID, payload.TenantID)
	}
	if payload.AuthMethod != resolved.AuthMethod {
		t.Fatalf("expected auth_method %q, got %q", resolved.AuthMethod, payload.AuthMethod)
	}
}

func TestMeHandlerMissingIdentity(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	resp := httptest.NewRecorder()
	MeHandler(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.Code)
	}
	if resp.Header().Get("Content-Type") != problemContentType {
		t.Fatalf("expected problem content type")
	}

	var payload problemResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Status != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, payload.Status)
	}
	if payload.Title == "" {
		t.Fatalf("expected title")
	}
	if payload.Detail == "" {
		t.Fatalf("expected detail")
	}
}

func TestWriteProblemDefaultsDetail(t *testing.T) {
	resp := httptest.NewRecorder()
	writeProblem(resp, http.StatusBadRequest, " ")

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}
	if resp.Header().Get("Content-Type") != problemContentType {
		t.Fatalf("expected problem content type")
	}

	var payload problemResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Detail != http.StatusText(http.StatusBadRequest) {
		t.Fatalf("expected detail %q, got %q", http.StatusText(http.StatusBadRequest), payload.Detail)
	}
}
