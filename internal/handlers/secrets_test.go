package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agynio/gateway/internal/secretsclient"
)

type stubSecretsClient struct {
	t                *testing.T
	createProviderFn func(context.Context, secretsclient.CreateProviderParams) (secretsclient.SecretProvider, error)
	getProviderFn    func(context.Context, string) (secretsclient.SecretProvider, error)
	updateProviderFn func(context.Context, string, secretsclient.UpdateProviderParams) (secretsclient.SecretProvider, error)
	deleteProviderFn func(context.Context, string) error
	listProvidersFn  func(context.Context, int32, string) ([]secretsclient.SecretProvider, string, error)
	createSecretFn   func(context.Context, secretsclient.CreateSecretParams) (secretsclient.Secret, error)
	getSecretFn      func(context.Context, string) (secretsclient.Secret, error)
	updateSecretFn   func(context.Context, string, secretsclient.UpdateSecretParams) (secretsclient.Secret, error)
	deleteSecretFn   func(context.Context, string) error
	listSecretsFn    func(context.Context, int32, string, string) ([]secretsclient.Secret, string, error)
	resolveSecretFn  func(context.Context, string) (secretsclient.ResolvedSecret, error)
}

func (s *stubSecretsClient) CreateProvider(ctx context.Context, params secretsclient.CreateProviderParams) (secretsclient.SecretProvider, error) {
	if s.createProviderFn == nil {
		s.t.Fatalf("unexpected CreateProvider call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.createProviderFn(ctx, params)
}

func (s *stubSecretsClient) GetProvider(ctx context.Context, id string) (secretsclient.SecretProvider, error) {
	if s.getProviderFn == nil {
		s.t.Fatalf("unexpected GetProvider call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.getProviderFn(ctx, id)
}

func (s *stubSecretsClient) UpdateProvider(ctx context.Context, id string, params secretsclient.UpdateProviderParams) (secretsclient.SecretProvider, error) {
	if s.updateProviderFn == nil {
		s.t.Fatalf("unexpected UpdateProvider call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.updateProviderFn(ctx, id, params)
}

func (s *stubSecretsClient) DeleteProvider(ctx context.Context, id string) error {
	if s.deleteProviderFn == nil {
		s.t.Fatalf("unexpected DeleteProvider call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.deleteProviderFn(ctx, id)
}

func (s *stubSecretsClient) ListProviders(ctx context.Context, pageSize int32, pageToken string) ([]secretsclient.SecretProvider, string, error) {
	if s.listProvidersFn == nil {
		s.t.Fatalf("unexpected ListProviders call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.listProvidersFn(ctx, pageSize, pageToken)
}

func (s *stubSecretsClient) CreateSecret(ctx context.Context, params secretsclient.CreateSecretParams) (secretsclient.Secret, error) {
	if s.createSecretFn == nil {
		s.t.Fatalf("unexpected CreateSecret call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.createSecretFn(ctx, params)
}

func (s *stubSecretsClient) GetSecret(ctx context.Context, id string) (secretsclient.Secret, error) {
	if s.getSecretFn == nil {
		s.t.Fatalf("unexpected GetSecret call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.getSecretFn(ctx, id)
}

func (s *stubSecretsClient) UpdateSecret(ctx context.Context, id string, params secretsclient.UpdateSecretParams) (secretsclient.Secret, error) {
	if s.updateSecretFn == nil {
		s.t.Fatalf("unexpected UpdateSecret call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.updateSecretFn(ctx, id, params)
}

func (s *stubSecretsClient) DeleteSecret(ctx context.Context, id string) error {
	if s.deleteSecretFn == nil {
		s.t.Fatalf("unexpected DeleteSecret call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.deleteSecretFn(ctx, id)
}

func (s *stubSecretsClient) ListSecrets(ctx context.Context, pageSize int32, pageToken, secretProviderID string) ([]secretsclient.Secret, string, error) {
	if s.listSecretsFn == nil {
		s.t.Fatalf("unexpected ListSecrets call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.listSecretsFn(ctx, pageSize, pageToken, secretProviderID)
}

func (s *stubSecretsClient) ResolveSecret(ctx context.Context, id string) (secretsclient.ResolvedSecret, error) {
	if s.resolveSecretFn == nil {
		s.t.Fatalf("unexpected ResolveSecret call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.resolveSecretFn(ctx, id)
}

func TestSecretsHandlerCreateProviderSuccess(t *testing.T) {
	createdAt := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)
	expected := secretsclient.SecretProvider{
		ID:          "prov-1",
		Title:       "Vault provider",
		Description: "primary",
		Type:        secretsclient.SecretProviderTypeVault,
		Config: secretsclient.SecretProviderConfig{
			Vault: &secretsclient.VaultConfig{Address: "https://vault.example.com", Token: "token"},
		},
		CreatedAt: createdAt,
	}

	called := false
	client := &stubSecretsClient{
		t: t,
		createProviderFn: func(ctx context.Context, params secretsclient.CreateProviderParams) (secretsclient.SecretProvider, error) {
			called = true
			if params.Title != expected.Title {
				t.Fatalf("unexpected title: %s", params.Title)
			}
			if params.Description != expected.Description {
				t.Fatalf("unexpected description: %s", params.Description)
			}
			if params.Type != expected.Type {
				t.Fatalf("unexpected type: %s", params.Type)
			}
			if params.Config.Vault == nil {
				t.Fatalf("expected vault config")
			}
			if params.Config.Vault.Address != expected.Config.Vault.Address {
				t.Fatalf("unexpected address: %s", params.Config.Vault.Address)
			}
			if params.Config.Vault.Token != expected.Config.Vault.Token {
				t.Fatalf("unexpected token")
			}
			return expected, nil
		},
	}

	payload := createSecretProviderRequest{
		Title:       expected.Title,
		Description: expected.Description,
		Type:        "vault",
		Config: &secretProviderConfigRequest{
			Vault: &vaultConfigRequest{
				Address: expected.Config.Vault.Address,
				Token:   expected.Config.Vault.Token,
			},
		},
	}
	req := newJSONRequest(t, http.MethodPost, "/secrets/v1/secret-providers", payload)
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).CreateProvider(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d got %d", http.StatusCreated, result.StatusCode)
	}
	if got := result.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type application/json got %q", got)
	}
	if !called {
		t.Fatalf("expected CreateProvider to be called")
	}

	var body secretsclient.SecretProvider
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestSecretsHandlerCreateProviderValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload any
		detail  string
	}{
		{
			name:    "missing type",
			payload: createSecretProviderRequest{Title: "provider", Config: &secretProviderConfigRequest{Vault: &vaultConfigRequest{Address: "https://vault", Token: "token"}}},
			detail:  "type is required",
		},
		{
			name:    "invalid type",
			payload: createSecretProviderRequest{Title: "provider", Type: "other", Config: &secretProviderConfigRequest{Vault: &vaultConfigRequest{Address: "https://vault", Token: "token"}}},
			detail:  "type must be \"vault\"",
		},
		{
			name:    "missing config",
			payload: createSecretProviderRequest{Title: "provider", Type: "vault"},
			detail:  "config is required",
		},
		{
			name:    "missing vault",
			payload: createSecretProviderRequest{Title: "provider", Type: "vault", Config: &secretProviderConfigRequest{}},
			detail:  "config.vault is required",
		},
		{
			name:    "missing address",
			payload: createSecretProviderRequest{Title: "provider", Type: "vault", Config: &secretProviderConfigRequest{Vault: &vaultConfigRequest{Token: "token"}}},
			detail:  "vault.address is required",
		},
		{
			name:    "missing token",
			payload: createSecretProviderRequest{Title: "provider", Type: "vault", Config: &secretProviderConfigRequest{Vault: &vaultConfigRequest{Address: "https://vault"}}},
			detail:  "vault.token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubSecretsClient{t: t}
			req := newJSONRequest(t, http.MethodPost, "/secrets/v1/secret-providers", tt.payload)
			resp := httptest.NewRecorder()

			NewSecretsHandler(client).CreateProvider(resp, req)

			assertProblemResponse(t, resp, http.StatusBadRequest, tt.detail)
		})
	}
}

func TestSecretsHandlerUpdateProviderPartial(t *testing.T) {
	expected := secretsclient.SecretProvider{
		ID:          "prov-1",
		Title:       "updated",
		Description: "desc",
		Type:        secretsclient.SecretProviderTypeVault,
		Config: secretsclient.SecretProviderConfig{
			Vault: &secretsclient.VaultConfig{Address: "https://vault.example.com", Token: "token"},
		},
		CreatedAt: time.Date(2024, 4, 5, 6, 7, 8, 0, time.UTC),
	}

	called := false
	client := &stubSecretsClient{
		t: t,
		updateProviderFn: func(ctx context.Context, id string, params secretsclient.UpdateProviderParams) (secretsclient.SecretProvider, error) {
			called = true
			if id != "prov-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			if params.Title == nil || *params.Title != expected.Title {
				t.Fatalf("unexpected title")
			}
			if params.Description != nil || params.Config != nil {
				t.Fatalf("expected only title update")
			}
			return expected, nil
		},
	}

	payload := updateSecretProviderRequest{Title: ptr(expected.Title)}
	req := newJSONRequest(t, http.MethodPatch, "/secrets/v1/secret-providers/prov-1", payload)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).UpdateProvider(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}
	if !called {
		t.Fatalf("expected UpdateProvider to be called")
	}

	var body secretsclient.SecretProvider
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestSecretsHandlerUpdateProviderEmptyBody(t *testing.T) {
	client := &stubSecretsClient{t: t}
	req := httptest.NewRequest(http.MethodPatch, "/secrets/v1/secret-providers/prov-1", nil)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).UpdateProvider(resp, req)

	assertProblemResponse(t, resp, http.StatusBadRequest, "request body is required")
}

func TestSecretsHandlerUpdateProviderValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload any
		detail  string
	}{
		{
			name:    "empty title",
			payload: updateSecretProviderRequest{Title: ptr(" ")},
			detail:  "title must not be empty",
		},
		{
			name:    "empty description",
			payload: updateSecretProviderRequest{Description: ptr(" ")},
			detail:  "description must not be empty",
		},
		{
			name:    "missing config vault",
			payload: updateSecretProviderRequest{Config: &secretProviderConfigRequest{}},
			detail:  "config.vault is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubSecretsClient{t: t}
			req := newJSONRequest(t, http.MethodPatch, "/secrets/v1/secret-providers/prov-1", tt.payload)
			req = withURLParam(req, "providerId", "prov-1")
			resp := httptest.NewRecorder()

			NewSecretsHandler(client).UpdateProvider(resp, req)

			assertProblemResponse(t, resp, http.StatusBadRequest, tt.detail)
		})
	}
}

func TestSecretsHandlerGetProviderSuccess(t *testing.T) {
	expected := secretsclient.SecretProvider{
		ID:          "prov-1",
		Title:       "Vault",
		Description: "primary",
		Type:        secretsclient.SecretProviderTypeVault,
		Config: secretsclient.SecretProviderConfig{
			Vault: &secretsclient.VaultConfig{Address: "https://vault.example.com", Token: "token"},
		},
		CreatedAt: time.Date(2024, 6, 7, 8, 9, 10, 0, time.UTC),
	}

	client := &stubSecretsClient{
		t: t,
		getProviderFn: func(ctx context.Context, id string) (secretsclient.SecretProvider, error) {
			if id != "prov-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return expected, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/secrets/v1/secret-providers/prov-1", nil)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).GetProvider(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}

	var body secretsclient.SecretProvider
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestSecretsHandlerGetProviderNotFound(t *testing.T) {
	client := &stubSecretsClient{
		t: t,
		getProviderFn: func(ctx context.Context, id string) (secretsclient.SecretProvider, error) {
			return secretsclient.SecretProvider{}, status.Error(codes.NotFound, "missing")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/secrets/v1/secret-providers/prov-1", nil)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).GetProvider(resp, req)

	assertProblemResponse(t, resp, http.StatusNotFound, "missing")
}

func TestSecretsHandlerDeleteProviderSuccess(t *testing.T) {
	called := false
	client := &stubSecretsClient{
		t: t,
		deleteProviderFn: func(ctx context.Context, id string) error {
			called = true
			if id != "prov-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/secrets/v1/secret-providers/prov-1", nil)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).DeleteProvider(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d got %d", http.StatusNoContent, result.StatusCode)
	}
	if !called {
		t.Fatalf("expected DeleteProvider to be called")
	}
}

func TestSecretsHandlerDeleteProviderNotFound(t *testing.T) {
	client := &stubSecretsClient{
		t: t,
		deleteProviderFn: func(ctx context.Context, id string) error {
			return status.Error(codes.NotFound, "missing")
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/secrets/v1/secret-providers/prov-1", nil)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).DeleteProvider(resp, req)

	assertProblemResponse(t, resp, http.StatusNotFound, "missing")
}

func TestSecretsHandlerListProvidersDefaults(t *testing.T) {
	items := []secretsclient.SecretProvider{{
		ID:          "prov-1",
		Title:       "Vault",
		Description: "primary",
		Type:        secretsclient.SecretProviderTypeVault,
		Config: secretsclient.SecretProviderConfig{
			Vault: &secretsclient.VaultConfig{Address: "https://vault.example.com", Token: "token"},
		},
		CreatedAt: time.Date(2024, 8, 9, 10, 11, 12, 0, time.UTC),
	}}

	client := &stubSecretsClient{
		t: t,
		listProvidersFn: func(ctx context.Context, pageSize int32, pageToken string) ([]secretsclient.SecretProvider, string, error) {
			if pageSize != defaultPageSize {
				t.Fatalf("unexpected page size: %d", pageSize)
			}
			if pageToken != "" {
				t.Fatalf("unexpected page token: %s", pageToken)
			}
			return items, "next", nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/secrets/v1/secret-providers", nil)
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).ListProviders(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}

	var body secretProviderListResponse
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body.Items, items) {
		t.Fatalf("unexpected items: %#v", body.Items)
	}
	if body.NextPageToken != "next" {
		t.Fatalf("unexpected nextPageToken: %s", body.NextPageToken)
	}
}

func TestSecretsHandlerListProvidersInvalidParams(t *testing.T) {
	client := &stubSecretsClient{t: t}
	req := httptest.NewRequest(http.MethodGet, "/secrets/v1/secret-providers?pageSize=0", nil)
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).ListProviders(resp, req)

	assertProblemResponse(t, resp, http.StatusBadRequest, "pageSize must be a positive integer")
}

func TestSecretsHandlerCreateSecretSuccess(t *testing.T) {
	createdAt := time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)
	expected := secretsclient.Secret{
		ID:               "secret-1",
		Title:            "DB password",
		Description:      "db",
		SecretProviderID: "prov-1",
		RemoteName:       "db/password",
		CreatedAt:        createdAt,
	}

	called := false
	client := &stubSecretsClient{
		t: t,
		createSecretFn: func(ctx context.Context, params secretsclient.CreateSecretParams) (secretsclient.Secret, error) {
			called = true
			if params.Title != expected.Title {
				t.Fatalf("unexpected title: %s", params.Title)
			}
			if params.Description != expected.Description {
				t.Fatalf("unexpected description: %s", params.Description)
			}
			if params.SecretProviderID != expected.SecretProviderID {
				t.Fatalf("unexpected provider id: %s", params.SecretProviderID)
			}
			if params.RemoteName != expected.RemoteName {
				t.Fatalf("unexpected remote name: %s", params.RemoteName)
			}
			return expected, nil
		},
	}

	payload := createSecretRequest{
		Title:            expected.Title,
		Description:      expected.Description,
		SecretProviderID: expected.SecretProviderID,
		RemoteName:       expected.RemoteName,
	}
	req := newJSONRequest(t, http.MethodPost, "/secrets/v1/secrets", payload)
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).CreateSecret(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d got %d", http.StatusCreated, result.StatusCode)
	}
	if got := result.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type application/json got %q", got)
	}
	if !called {
		t.Fatalf("expected CreateSecret to be called")
	}

	var body secretsclient.Secret
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestSecretsHandlerCreateSecretValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload any
		detail  string
	}{
		{
			name:    "missing provider",
			payload: createSecretRequest{RemoteName: "db/password"},
			detail:  "secretProviderId is required",
		},
		{
			name:    "missing remote name",
			payload: createSecretRequest{SecretProviderID: "prov-1"},
			detail:  "remoteName is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubSecretsClient{t: t}
			req := newJSONRequest(t, http.MethodPost, "/secrets/v1/secrets", tt.payload)
			resp := httptest.NewRecorder()

			NewSecretsHandler(client).CreateSecret(resp, req)

			assertProblemResponse(t, resp, http.StatusBadRequest, tt.detail)
		})
	}
}

func TestSecretsHandlerUpdateSecretPartial(t *testing.T) {
	expected := secretsclient.Secret{
		ID:               "secret-1",
		Title:            "DB password",
		Description:      "db",
		SecretProviderID: "prov-1",
		RemoteName:       "db/password/updated",
		CreatedAt:        time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC),
	}

	called := false
	client := &stubSecretsClient{
		t: t,
		updateSecretFn: func(ctx context.Context, id string, params secretsclient.UpdateSecretParams) (secretsclient.Secret, error) {
			called = true
			if id != "secret-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			if params.RemoteName == nil || *params.RemoteName != expected.RemoteName {
				t.Fatalf("unexpected remote name")
			}
			if params.Title != nil || params.Description != nil || params.SecretProviderID != nil {
				t.Fatalf("expected only remote name update")
			}
			return expected, nil
		},
	}

	payload := updateSecretRequest{RemoteName: ptr(expected.RemoteName)}
	req := newJSONRequest(t, http.MethodPatch, "/secrets/v1/secrets/secret-1", payload)
	req = withURLParam(req, "secretId", "secret-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).UpdateSecret(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}
	if !called {
		t.Fatalf("expected UpdateSecret to be called")
	}

	var body secretsclient.Secret
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestSecretsHandlerUpdateSecretEmptyBody(t *testing.T) {
	client := &stubSecretsClient{t: t}
	req := httptest.NewRequest(http.MethodPatch, "/secrets/v1/secrets/secret-1", nil)
	req = withURLParam(req, "secretId", "secret-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).UpdateSecret(resp, req)

	assertProblemResponse(t, resp, http.StatusBadRequest, "request body is required")
}

func TestSecretsHandlerUpdateSecretValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload any
		detail  string
	}{
		{
			name:    "empty secret provider",
			payload: updateSecretRequest{SecretProviderID: ptr(" ")},
			detail:  "secretProviderId must not be empty",
		},
		{
			name:    "empty remote name",
			payload: updateSecretRequest{RemoteName: ptr(" ")},
			detail:  "remoteName must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubSecretsClient{t: t}
			req := newJSONRequest(t, http.MethodPatch, "/secrets/v1/secrets/secret-1", tt.payload)
			req = withURLParam(req, "secretId", "secret-1")
			resp := httptest.NewRecorder()

			NewSecretsHandler(client).UpdateSecret(resp, req)

			assertProblemResponse(t, resp, http.StatusBadRequest, tt.detail)
		})
	}
}

func TestSecretsHandlerGetSecretSuccess(t *testing.T) {
	expected := secretsclient.Secret{
		ID:               "secret-1",
		Title:            "DB password",
		Description:      "db",
		SecretProviderID: "prov-1",
		RemoteName:       "db/password",
		CreatedAt:        time.Date(2024, 6, 7, 8, 9, 10, 0, time.UTC),
	}

	client := &stubSecretsClient{
		t: t,
		getSecretFn: func(ctx context.Context, id string) (secretsclient.Secret, error) {
			if id != "secret-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return expected, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/secrets/v1/secrets/secret-1", nil)
	req = withURLParam(req, "secretId", "secret-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).GetSecret(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}

	var body secretsclient.Secret
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestSecretsHandlerGetSecretNotFound(t *testing.T) {
	client := &stubSecretsClient{
		t: t,
		getSecretFn: func(ctx context.Context, id string) (secretsclient.Secret, error) {
			return secretsclient.Secret{}, status.Error(codes.NotFound, "missing")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/secrets/v1/secrets/secret-1", nil)
	req = withURLParam(req, "secretId", "secret-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).GetSecret(resp, req)

	assertProblemResponse(t, resp, http.StatusNotFound, "missing")
}

func TestSecretsHandlerDeleteSecretSuccess(t *testing.T) {
	called := false
	client := &stubSecretsClient{
		t: t,
		deleteSecretFn: func(ctx context.Context, id string) error {
			called = true
			if id != "secret-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/secrets/v1/secrets/secret-1", nil)
	req = withURLParam(req, "secretId", "secret-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).DeleteSecret(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d got %d", http.StatusNoContent, result.StatusCode)
	}
	if !called {
		t.Fatalf("expected DeleteSecret to be called")
	}
}

func TestSecretsHandlerDeleteSecretNotFound(t *testing.T) {
	client := &stubSecretsClient{
		t: t,
		deleteSecretFn: func(ctx context.Context, id string) error {
			return status.Error(codes.NotFound, "missing")
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/secrets/v1/secrets/secret-1", nil)
	req = withURLParam(req, "secretId", "secret-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).DeleteSecret(resp, req)

	assertProblemResponse(t, resp, http.StatusNotFound, "missing")
}

func TestSecretsHandlerListSecretsDefaults(t *testing.T) {
	items := []secretsclient.Secret{{
		ID:               "secret-1",
		Title:            "DB password",
		Description:      "db",
		SecretProviderID: "prov-1",
		RemoteName:       "db/password",
		CreatedAt:        time.Date(2024, 7, 8, 9, 10, 11, 0, time.UTC),
	}}

	client := &stubSecretsClient{
		t: t,
		listSecretsFn: func(ctx context.Context, pageSize int32, pageToken, secretProviderID string) ([]secretsclient.Secret, string, error) {
			if pageSize != defaultPageSize {
				t.Fatalf("unexpected page size: %d", pageSize)
			}
			if pageToken != "" {
				t.Fatalf("unexpected page token: %s", pageToken)
			}
			if secretProviderID != "prov-1" {
				t.Fatalf("unexpected provider id: %s", secretProviderID)
			}
			return items, "more", nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/secrets/v1/secrets?secretProviderId=prov-1", nil)
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).ListSecrets(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}

	var body secretListResponse
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body.Items, items) {
		t.Fatalf("unexpected items: %#v", body.Items)
	}
	if body.NextPageToken != "more" {
		t.Fatalf("unexpected nextPageToken: %s", body.NextPageToken)
	}
}

func TestSecretsHandlerListSecretsInvalidParams(t *testing.T) {
	client := &stubSecretsClient{t: t}
	req := httptest.NewRequest(http.MethodGet, "/secrets/v1/secrets?pageSize=bad", nil)
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).ListSecrets(resp, req)

	assertProblemResponse(t, resp, http.StatusBadRequest, "pageSize must be a positive integer")
}

func TestSecretsHandlerResolveSecretSuccess(t *testing.T) {
	client := &stubSecretsClient{
		t: t,
		resolveSecretFn: func(ctx context.Context, id string) (secretsclient.ResolvedSecret, error) {
			if id != "secret-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return secretsclient.ResolvedSecret{Value: "resolved"}, nil
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/secrets/v1/secrets/secret-1/resolve", nil)
	req = withURLParam(req, "secretId", "secret-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).ResolveSecret(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}

	var body secretsclient.ResolvedSecret
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Value != "resolved" {
		t.Fatalf("unexpected value: %s", body.Value)
	}
}

func TestSecretsHandlerResolveSecretNotFound(t *testing.T) {
	client := &stubSecretsClient{
		t: t,
		resolveSecretFn: func(ctx context.Context, id string) (secretsclient.ResolvedSecret, error) {
			return secretsclient.ResolvedSecret{}, status.Error(codes.NotFound, "missing")
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/secrets/v1/secrets/secret-1/resolve", nil)
	req = withURLParam(req, "secretId", "secret-1")
	resp := httptest.NewRecorder()

	NewSecretsHandler(client).ResolveSecret(resp, req)

	assertProblemResponse(t, resp, http.StatusNotFound, "missing")
}

func TestWriteSecretsGRPCErrorMappings(t *testing.T) {
	tests := []struct {
		code       codes.Code
		statusCode int
	}{
		{code: codes.NotFound, statusCode: http.StatusNotFound},
		{code: codes.InvalidArgument, statusCode: http.StatusBadRequest},
		{code: codes.AlreadyExists, statusCode: http.StatusConflict},
		{code: codes.FailedPrecondition, statusCode: http.StatusPreconditionFailed},
		{code: codes.Unavailable, statusCode: http.StatusServiceUnavailable},
		{code: codes.PermissionDenied, statusCode: http.StatusForbidden},
		{code: codes.Unauthenticated, statusCode: http.StatusUnauthorized},
		{code: codes.Internal, statusCode: http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			resp := httptest.NewRecorder()
			writeSecretsGRPCError(resp, status.Error(tt.code, "detail"))
			assertProblemResponse(t, resp, tt.statusCode, "detail")
		})
	}
}

func TestParseProviderConfig(t *testing.T) {
	config, err := parseProviderConfig(&secretProviderConfigRequest{
		Vault: &vaultConfigRequest{Address: "https://vault.example.com", Token: "token"},
	}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config == nil || config.Vault == nil {
		t.Fatalf("expected config")
	}
	if config.Vault.Address != "https://vault.example.com" {
		t.Fatalf("unexpected address: %s", config.Vault.Address)
	}
	if config.Vault.Token != "token" {
		t.Fatalf("unexpected token")
	}

	if _, err := parseProviderConfig(nil, true); err == nil {
		t.Fatalf("expected error")
	}
}
