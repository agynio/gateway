package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/agynio/gateway/internal/identity"
)

type MeResponse struct {
	IdentityID   string `json:"identity_id"`
	IdentityType string `json:"identity_type"`
	TenantID     string `json:"tenant_id"`
	AuthMethod   string `json:"auth_method"`
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	resolved, ok := identity.IdentityFromContext(r.Context())
	if !ok {
		problem := NewProblem(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "identity not available")
		WriteProblem(w, problem)
		return
	}

	response := MeResponse{
		IdentityID:   resolved.IdentityID,
		IdentityType: resolved.IdentityType,
		TenantID:     resolved.TenantID,
		AuthMethod:   resolved.AuthMethod,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode /me response: %v", err)
	}
}
