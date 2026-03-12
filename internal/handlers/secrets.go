package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/agynio/gateway/internal/secretsclient"
)

type SecretsClient interface {
	CreateProvider(ctx context.Context, params secretsclient.CreateProviderParams) (secretsclient.SecretProvider, error)
	GetProvider(ctx context.Context, id string) (secretsclient.SecretProvider, error)
	UpdateProvider(ctx context.Context, id string, params secretsclient.UpdateProviderParams) (secretsclient.SecretProvider, error)
	DeleteProvider(ctx context.Context, id string) error
	ListProviders(ctx context.Context, pageSize int32, pageToken string) ([]secretsclient.SecretProvider, string, error)
	CreateSecret(ctx context.Context, params secretsclient.CreateSecretParams) (secretsclient.Secret, error)
	GetSecret(ctx context.Context, id string) (secretsclient.Secret, error)
	UpdateSecret(ctx context.Context, id string, params secretsclient.UpdateSecretParams) (secretsclient.Secret, error)
	DeleteSecret(ctx context.Context, id string) error
	ListSecrets(ctx context.Context, pageSize int32, pageToken, secretProviderID string) ([]secretsclient.Secret, string, error)
	ResolveSecret(ctx context.Context, id string) (secretsclient.ResolvedSecret, error)
}

type SecretsHandler struct {
	client SecretsClient
}

func NewSecretsHandler(client SecretsClient) *SecretsHandler {
	if client == nil {
		panic("secrets client is required")
	}
	return &SecretsHandler{client: client}
}

type vaultConfigRequest struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

type secretProviderConfigRequest struct {
	Vault *vaultConfigRequest `json:"vault"`
}

type createSecretProviderRequest struct {
	Title       string                       `json:"title"`
	Description string                       `json:"description"`
	Type        string                       `json:"type"`
	Config      *secretProviderConfigRequest `json:"config"`
}

type updateSecretProviderRequest struct {
	Title       *string                      `json:"title"`
	Description *string                      `json:"description"`
	Config      *secretProviderConfigRequest `json:"config"`
}

type createSecretRequest struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	SecretProviderID string `json:"secretProviderId"`
	RemoteName       string `json:"remoteName"`
}

type updateSecretRequest struct {
	Title            *string `json:"title"`
	Description      *string `json:"description"`
	SecretProviderID *string `json:"secretProviderId"`
	RemoteName       *string `json:"remoteName"`
}

type secretProviderListResponse struct {
	Items         []secretsclient.SecretProvider `json:"items"`
	NextPageToken string                         `json:"nextPageToken"`
}

type secretListResponse struct {
	Items         []secretsclient.Secret `json:"items"`
	NextPageToken string                 `json:"nextPageToken"`
}

func (h *SecretsHandler) CreateProvider(w http.ResponseWriter, r *http.Request) {
	var payload createSecretProviderRequest
	if err := decodeSecretsJSON(r, &payload); err != nil {
		writeSecretsBadRequest(w, err)
		return
	}

	title := strings.TrimSpace(payload.Title)
	description := strings.TrimSpace(payload.Description)

	providerType, err := parseSecretProviderType(payload.Type)
	if err != nil {
		writeSecretsBadRequest(w, err)
		return
	}

	config, err := parseProviderConfig(payload.Config, true)
	if err != nil {
		writeSecretsBadRequest(w, err)
		return
	}

	provider, err := h.client.CreateProvider(r.Context(), secretsclient.CreateProviderParams{
		Title:       title,
		Description: description,
		Type:        providerType,
		Config:      *config,
	})
	if err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	writeSecretsJSON(w, http.StatusCreated, provider)
}

func (h *SecretsHandler) GetProvider(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "providerId"))
	if providerID == "" {
		writeSecretsValidationMessage(w, "providerId is required")
		return
	}

	provider, err := h.client.GetProvider(r.Context(), providerID)
	if err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	writeSecretsJSON(w, http.StatusOK, provider)
}

func (h *SecretsHandler) UpdateProvider(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "providerId"))
	if providerID == "" {
		writeSecretsValidationMessage(w, "providerId is required")
		return
	}

	var payload updateSecretProviderRequest
	if err := decodeSecretsJSON(r, &payload); err != nil {
		writeSecretsBadRequest(w, err)
		return
	}

	params := secretsclient.UpdateProviderParams{}
	if payload.Title != nil {
		title := strings.TrimSpace(*payload.Title)
		if title == "" {
			writeSecretsValidationMessage(w, "title must not be empty")
			return
		}
		params.Title = &title
	}
	if payload.Description != nil {
		description := strings.TrimSpace(*payload.Description)
		if description == "" {
			writeSecretsValidationMessage(w, "description must not be empty")
			return
		}
		params.Description = &description
	}
	if payload.Config != nil {
		config, err := parseProviderConfig(payload.Config, true)
		if err != nil {
			writeSecretsBadRequest(w, err)
			return
		}
		params.Config = config
	}

	if params.Title == nil && params.Description == nil && params.Config == nil {
		writeSecretsValidationMessage(w, "at least one field must be provided")
		return
	}

	provider, err := h.client.UpdateProvider(r.Context(), providerID, params)
	if err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	writeSecretsJSON(w, http.StatusOK, provider)
}

func (h *SecretsHandler) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "providerId"))
	if providerID == "" {
		writeSecretsValidationMessage(w, "providerId is required")
		return
	}

	if err := h.client.DeleteProvider(r.Context(), providerID); err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SecretsHandler) ListProviders(w http.ResponseWriter, r *http.Request) {
	pageSize, pageToken, err := parsePagination(r)
	if err != nil {
		writeSecretsBadRequest(w, err)
		return
	}

	providers, nextToken, err := h.client.ListProviders(r.Context(), int32(pageSize), pageToken)
	if err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	writeSecretsJSON(w, http.StatusOK, secretProviderListResponse{
		Items:         providers,
		NextPageToken: nextToken,
	})
}

func (h *SecretsHandler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	var payload createSecretRequest
	if err := decodeSecretsJSON(r, &payload); err != nil {
		writeSecretsBadRequest(w, err)
		return
	}

	title := strings.TrimSpace(payload.Title)
	description := strings.TrimSpace(payload.Description)
	providerID := strings.TrimSpace(payload.SecretProviderID)
	if providerID == "" {
		writeSecretsValidationMessage(w, "secretProviderId is required")
		return
	}
	remoteName := strings.TrimSpace(payload.RemoteName)
	if remoteName == "" {
		writeSecretsValidationMessage(w, "remoteName is required")
		return
	}

	secret, err := h.client.CreateSecret(r.Context(), secretsclient.CreateSecretParams{
		Title:            title,
		Description:      description,
		SecretProviderID: providerID,
		RemoteName:       remoteName,
	})
	if err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	writeSecretsJSON(w, http.StatusCreated, secret)
}

func (h *SecretsHandler) GetSecret(w http.ResponseWriter, r *http.Request) {
	secretID := strings.TrimSpace(chi.URLParam(r, "secretId"))
	if secretID == "" {
		writeSecretsValidationMessage(w, "secretId is required")
		return
	}

	secret, err := h.client.GetSecret(r.Context(), secretID)
	if err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	writeSecretsJSON(w, http.StatusOK, secret)
}

func (h *SecretsHandler) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	secretID := strings.TrimSpace(chi.URLParam(r, "secretId"))
	if secretID == "" {
		writeSecretsValidationMessage(w, "secretId is required")
		return
	}

	var payload updateSecretRequest
	if err := decodeSecretsJSON(r, &payload); err != nil {
		writeSecretsBadRequest(w, err)
		return
	}

	params := secretsclient.UpdateSecretParams{}
	if payload.Title != nil {
		title := strings.TrimSpace(*payload.Title)
		if title == "" {
			writeSecretsValidationMessage(w, "title must not be empty")
			return
		}
		params.Title = &title
	}
	if payload.Description != nil {
		description := strings.TrimSpace(*payload.Description)
		if description == "" {
			writeSecretsValidationMessage(w, "description must not be empty")
			return
		}
		params.Description = &description
	}
	if payload.SecretProviderID != nil {
		providerID := strings.TrimSpace(*payload.SecretProviderID)
		if providerID == "" {
			writeSecretsValidationMessage(w, "secretProviderId must not be empty")
			return
		}
		params.SecretProviderID = &providerID
	}
	if payload.RemoteName != nil {
		remoteName := strings.TrimSpace(*payload.RemoteName)
		if remoteName == "" {
			writeSecretsValidationMessage(w, "remoteName must not be empty")
			return
		}
		params.RemoteName = &remoteName
	}

	if params.Title == nil && params.Description == nil && params.SecretProviderID == nil && params.RemoteName == nil {
		writeSecretsValidationMessage(w, "at least one field must be provided")
		return
	}

	secret, err := h.client.UpdateSecret(r.Context(), secretID, params)
	if err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	writeSecretsJSON(w, http.StatusOK, secret)
}

func (h *SecretsHandler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	secretID := strings.TrimSpace(chi.URLParam(r, "secretId"))
	if secretID == "" {
		writeSecretsValidationMessage(w, "secretId is required")
		return
	}

	if err := h.client.DeleteSecret(r.Context(), secretID); err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SecretsHandler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	pageSize, pageToken, err := parsePagination(r)
	if err != nil {
		writeSecretsBadRequest(w, err)
		return
	}

	providerID := strings.TrimSpace(r.URL.Query().Get("secretProviderId"))
	secrets, nextToken, err := h.client.ListSecrets(r.Context(), int32(pageSize), pageToken, providerID)
	if err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	writeSecretsJSON(w, http.StatusOK, secretListResponse{
		Items:         secrets,
		NextPageToken: nextToken,
	})
}

func (h *SecretsHandler) ResolveSecret(w http.ResponseWriter, r *http.Request) {
	secretID := strings.TrimSpace(chi.URLParam(r, "secretId"))
	if secretID == "" {
		writeSecretsValidationMessage(w, "secretId is required")
		return
	}

	value, err := h.client.ResolveSecret(r.Context(), secretID)
	if err != nil {
		writeSecretsGRPCError(w, err)
		return
	}

	writeSecretsJSON(w, http.StatusOK, value)
}

func parseSecretProviderType(raw string) (secretsclient.SecretProviderType, error) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return "", fmt.Errorf("type is required")
	}
	normalized = strings.ToLower(normalized)
	if normalized == string(secretsclient.SecretProviderTypeVault) {
		return secretsclient.SecretProviderTypeVault, nil
	}
	return "", fmt.Errorf("type must be \"vault\"")
}

func parseProviderConfig(payload *secretProviderConfigRequest, required bool) (*secretsclient.SecretProviderConfig, error) {
	if payload == nil {
		if required {
			return nil, fmt.Errorf("config is required")
		}
		return nil, nil
	}
	if payload.Vault == nil {
		return nil, fmt.Errorf("config.vault is required")
	}

	address := strings.TrimSpace(payload.Vault.Address)
	if address == "" {
		return nil, fmt.Errorf("vault.address is required")
	}
	token := strings.TrimSpace(payload.Vault.Token)
	if token == "" {
		return nil, fmt.Errorf("vault.token is required")
	}

	config := secretsclient.SecretProviderConfig{
		Vault: &secretsclient.VaultConfig{
			Address: address,
			Token:   token,
		},
	}
	return &config, nil
}

func writeSecretsGRPCError(w http.ResponseWriter, err error) {
	if err == nil {
		problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), "upstream request failed")
		WriteProblem(w, problem)
		return
	}

	problemErr := grpcErrorToProblem(err)
	if problemErr == nil {
		problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), "upstream request failed")
		WriteProblem(w, problem)
		return
	}

	if problemErr.Problem.Status >= http.StatusInternalServerError {
		log.Printf("secrets grpc error: %v", err)
	}

	WriteProblem(w, problemErr.Problem)
}

func decodeSecretsJSON(r *http.Request, out any) error {
	return decodeJSONBody(r, out)
}

func writeSecretsBadRequest(w http.ResponseWriter, err error) {
	writeBadRequest(w, err)
}

func writeSecretsValidationMessage(w http.ResponseWriter, message string) {
	writeValidationMessage(w, message)
}

func writeSecretsJSON(w http.ResponseWriter, status int, payload any) {
	writeJSONResponse(w, status, payload, "secrets")
}
