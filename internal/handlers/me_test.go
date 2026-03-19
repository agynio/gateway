package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agynio/gateway/internal/gen"
	"github.com/agynio/gateway/internal/identity"
)

func TestMeHandlerReturnsIdentity(t *testing.T) {
	resolved := identity.ResolvedIdentity{
		IdentityID:   "id-1",
		IdentityType: "user",
		TenantID:     "tenant-1",
		AuthMethod:   "ziti",
	}

	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	request = request.WithContext(identity.WithIdentity(request.Context(), resolved))
	response := httptest.NewRecorder()

	MeHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type: %s", response.Header().Get("Content-Type"))
	}

	var body MeResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.IdentityID != resolved.IdentityID {
		t.Fatalf("unexpected identity_id: %s", body.IdentityID)
	}
	if body.IdentityType != resolved.IdentityType {
		t.Fatalf("unexpected identity_type: %s", body.IdentityType)
	}
	if body.TenantID != resolved.TenantID {
		t.Fatalf("unexpected tenant_id: %s", body.TenantID)
	}
	if body.AuthMethod != resolved.AuthMethod {
		t.Fatalf("unexpected auth_method: %s", body.AuthMethod)
	}
}

func TestMeHandlerMissingIdentity(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	response := httptest.NewRecorder()

	MeHandler(response, request)

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
