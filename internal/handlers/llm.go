package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agynio/gateway/internal/llmclient"
)

type LLMClient interface {
	CreateProvider(ctx context.Context, params llmclient.CreateProviderParams) (llmclient.LLMProvider, error)
	GetProvider(ctx context.Context, id string) (llmclient.LLMProvider, error)
	UpdateProvider(ctx context.Context, id string, params llmclient.UpdateProviderParams) (llmclient.LLMProvider, error)
	DeleteProvider(ctx context.Context, id string) error
	ListProviders(ctx context.Context, pageSize int32, pageToken string) ([]llmclient.LLMProvider, string, error)
	CreateModel(ctx context.Context, params llmclient.CreateModelParams) (llmclient.Model, error)
	GetModel(ctx context.Context, id string) (llmclient.Model, error)
	UpdateModel(ctx context.Context, id string, params llmclient.UpdateModelParams) (llmclient.Model, error)
	DeleteModel(ctx context.Context, id string) error
	ListModels(ctx context.Context, pageSize int32, pageToken, llmProviderID string) ([]llmclient.Model, string, error)
}

type LLMHandler struct {
	client LLMClient
}

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

func NewLLMHandler(client LLMClient) *LLMHandler {
	if client == nil {
		panic("llm client is required")
	}
	return &LLMHandler{client: client}
}

type createProviderRequest struct {
	Endpoint   string `json:"endpoint"`
	AuthMethod string `json:"authMethod"`
	Token      string `json:"token"`
}

type updateProviderRequest struct {
	Endpoint   *string `json:"endpoint"`
	AuthMethod *string `json:"authMethod"`
	Token      *string `json:"token"`
}

type createModelRequest struct {
	Name          string `json:"name"`
	LLMProviderID string `json:"llmProviderId"`
	RemoteName    string `json:"remoteName"`
}

type updateModelRequest struct {
	Name          *string `json:"name"`
	LLMProviderID *string `json:"llmProviderId"`
	RemoteName    *string `json:"remoteName"`
}

type providerListResponse struct {
	Items         []llmclient.LLMProvider `json:"items"`
	NextPageToken string                  `json:"nextPageToken"`
}

type modelListResponse struct {
	Items         []llmclient.Model `json:"items"`
	NextPageToken string            `json:"nextPageToken"`
}

func (h *LLMHandler) CreateProvider(w http.ResponseWriter, r *http.Request) {
	var payload createProviderRequest
	if err := decodeLLMJSON(r, &payload); err != nil {
		writeLLMBadRequest(w, err)
		return
	}

	endpoint := strings.TrimSpace(payload.Endpoint)
	if endpoint == "" {
		writeLLMValidationMessage(w, "endpoint is required")
		return
	}

	authMethod, err := parseAuthMethod(payload.AuthMethod)
	if err != nil {
		writeLLMBadRequest(w, err)
		return
	}

	token := strings.TrimSpace(payload.Token)
	if authMethod == llmclient.AuthMethodBearer && token == "" {
		writeLLMValidationMessage(w, "token is required for bearer auth")
		return
	}

	provider, err := h.client.CreateProvider(r.Context(), llmclient.CreateProviderParams{
		Endpoint:   endpoint,
		AuthMethod: authMethod,
		Token:      token,
	})
	if err != nil {
		writeLLMGRPCError(w, err)
		return
	}

	writeLLMJSON(w, http.StatusCreated, provider)
}

func (h *LLMHandler) GetProvider(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "providerId"))
	if providerID == "" {
		writeLLMValidationMessage(w, "providerId is required")
		return
	}

	provider, err := h.client.GetProvider(r.Context(), providerID)
	if err != nil {
		writeLLMGRPCError(w, err)
		return
	}

	writeLLMJSON(w, http.StatusOK, provider)
}

func (h *LLMHandler) UpdateProvider(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "providerId"))
	if providerID == "" {
		writeLLMValidationMessage(w, "providerId is required")
		return
	}

	var payload updateProviderRequest
	if err := decodeLLMJSON(r, &payload); err != nil {
		writeLLMBadRequest(w, err)
		return
	}

	params := llmclient.UpdateProviderParams{}
	if payload.Endpoint != nil {
		endpoint := strings.TrimSpace(*payload.Endpoint)
		if endpoint == "" {
			writeLLMValidationMessage(w, "endpoint must not be empty")
			return
		}
		params.Endpoint = &endpoint
	}
	if payload.AuthMethod != nil {
		method, err := parseAuthMethod(*payload.AuthMethod)
		if err != nil {
			writeLLMBadRequest(w, err)
			return
		}
		params.AuthMethod = &method
	}
	if payload.Token != nil {
		token := strings.TrimSpace(*payload.Token)
		if token == "" {
			writeLLMValidationMessage(w, "token must not be empty")
			return
		}
		params.Token = &token
	}

	if params.Endpoint == nil && params.AuthMethod == nil && params.Token == nil {
		writeLLMValidationMessage(w, "at least one field must be provided")
		return
	}

	provider, err := h.client.UpdateProvider(r.Context(), providerID, params)
	if err != nil {
		writeLLMGRPCError(w, err)
		return
	}

	writeLLMJSON(w, http.StatusOK, provider)
}

func (h *LLMHandler) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "providerId"))
	if providerID == "" {
		writeLLMValidationMessage(w, "providerId is required")
		return
	}

	if err := h.client.DeleteProvider(r.Context(), providerID); err != nil {
		writeLLMGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *LLMHandler) ListProviders(w http.ResponseWriter, r *http.Request) {
	pageSize, pageToken, err := parsePagination(r)
	if err != nil {
		writeLLMBadRequest(w, err)
		return
	}

	providers, nextToken, err := h.client.ListProviders(r.Context(), int32(pageSize), pageToken)
	if err != nil {
		writeLLMGRPCError(w, err)
		return
	}

	writeLLMJSON(w, http.StatusOK, providerListResponse{
		Items:         providers,
		NextPageToken: nextToken,
	})
}

func (h *LLMHandler) CreateModel(w http.ResponseWriter, r *http.Request) {
	var payload createModelRequest
	if err := decodeLLMJSON(r, &payload); err != nil {
		writeLLMBadRequest(w, err)
		return
	}

	name := strings.TrimSpace(payload.Name)
	if name == "" {
		writeLLMValidationMessage(w, "name is required")
		return
	}

	providerID := strings.TrimSpace(payload.LLMProviderID)
	if providerID == "" {
		writeLLMValidationMessage(w, "llmProviderId is required")
		return
	}

	remoteName := strings.TrimSpace(payload.RemoteName)
	if remoteName == "" {
		writeLLMValidationMessage(w, "remoteName is required")
		return
	}

	model, err := h.client.CreateModel(r.Context(), llmclient.CreateModelParams{
		Name:          name,
		LLMProviderID: providerID,
		RemoteName:    remoteName,
	})
	if err != nil {
		writeLLMGRPCError(w, err)
		return
	}

	writeLLMJSON(w, http.StatusCreated, model)
}

func (h *LLMHandler) GetModel(w http.ResponseWriter, r *http.Request) {
	modelID := strings.TrimSpace(chi.URLParam(r, "modelId"))
	if modelID == "" {
		writeLLMValidationMessage(w, "modelId is required")
		return
	}

	model, err := h.client.GetModel(r.Context(), modelID)
	if err != nil {
		writeLLMGRPCError(w, err)
		return
	}

	writeLLMJSON(w, http.StatusOK, model)
}

func (h *LLMHandler) UpdateModel(w http.ResponseWriter, r *http.Request) {
	modelID := strings.TrimSpace(chi.URLParam(r, "modelId"))
	if modelID == "" {
		writeLLMValidationMessage(w, "modelId is required")
		return
	}

	var payload updateModelRequest
	if err := decodeLLMJSON(r, &payload); err != nil {
		writeLLMBadRequest(w, err)
		return
	}

	params := llmclient.UpdateModelParams{}
	if payload.Name != nil {
		name := strings.TrimSpace(*payload.Name)
		if name == "" {
			writeLLMValidationMessage(w, "name must not be empty")
			return
		}
		params.Name = &name
	}
	if payload.LLMProviderID != nil {
		providerID := strings.TrimSpace(*payload.LLMProviderID)
		if providerID == "" {
			writeLLMValidationMessage(w, "llmProviderId must not be empty")
			return
		}
		params.LLMProviderID = &providerID
	}
	if payload.RemoteName != nil {
		remoteName := strings.TrimSpace(*payload.RemoteName)
		if remoteName == "" {
			writeLLMValidationMessage(w, "remoteName must not be empty")
			return
		}
		params.RemoteName = &remoteName
	}

	if params.Name == nil && params.LLMProviderID == nil && params.RemoteName == nil {
		writeLLMValidationMessage(w, "at least one field must be provided")
		return
	}

	model, err := h.client.UpdateModel(r.Context(), modelID, params)
	if err != nil {
		writeLLMGRPCError(w, err)
		return
	}

	writeLLMJSON(w, http.StatusOK, model)
}

func (h *LLMHandler) DeleteModel(w http.ResponseWriter, r *http.Request) {
	modelID := strings.TrimSpace(chi.URLParam(r, "modelId"))
	if modelID == "" {
		writeLLMValidationMessage(w, "modelId is required")
		return
	}

	if err := h.client.DeleteModel(r.Context(), modelID); err != nil {
		writeLLMGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *LLMHandler) ListModels(w http.ResponseWriter, r *http.Request) {
	pageSize, pageToken, err := parsePagination(r)
	if err != nil {
		writeLLMBadRequest(w, err)
		return
	}

	providerID := strings.TrimSpace(r.URL.Query().Get("llmProviderId"))
	models, nextToken, err := h.client.ListModels(r.Context(), int32(pageSize), pageToken, providerID)
	if err != nil {
		writeLLMGRPCError(w, err)
		return
	}

	writeLLMJSON(w, http.StatusOK, modelListResponse{
		Items:         models,
		NextPageToken: nextToken,
	})
}

func parsePagination(r *http.Request) (int, string, error) {
	pageSize := defaultPageSize
	pageToken := strings.TrimSpace(r.URL.Query().Get("pageToken"))

	if raw := strings.TrimSpace(r.URL.Query().Get("pageSize")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			return 0, "", fmt.Errorf("pageSize must be a positive integer")
		}
		pageSize = parsed
	}

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return pageSize, pageToken, nil
}

func writeLLMGRPCError(w http.ResponseWriter, err error) {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		log.Printf("llm grpc error: %v", err)
		problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), err.Error())
		WriteProblem(w, problem)
		return
	}

	statusCode := http.StatusBadGateway
	switch grpcStatus.Code() {
	case codes.NotFound:
		statusCode = http.StatusNotFound
	case codes.InvalidArgument:
		statusCode = http.StatusBadRequest
	case codes.AlreadyExists, codes.FailedPrecondition:
		statusCode = http.StatusConflict
	case codes.Unavailable:
		statusCode = http.StatusServiceUnavailable
	case codes.PermissionDenied:
		statusCode = http.StatusForbidden
	case codes.Unauthenticated:
		statusCode = http.StatusUnauthorized
	}

	if statusCode >= http.StatusInternalServerError {
		log.Printf("llm grpc error: %v", err)
	}

	message := strings.TrimSpace(grpcStatus.Message())
	if message == "" {
		message = http.StatusText(statusCode)
	}
	problem := NewProblem(statusCode, http.StatusText(statusCode), message)
	WriteProblem(w, problem)
}

func decodeLLMJSON(r *http.Request, out any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("request body is required")
		}
		return err
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return fmt.Errorf("request body must contain a single JSON object")
	}
	return nil
}

func parseAuthMethod(raw string) (llmclient.AuthMethod, error) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return "", fmt.Errorf("authMethod is required")
	}
	normalized = strings.ToLower(normalized)
	switch normalized {
	case string(llmclient.AuthMethodBearer):
		return llmclient.AuthMethodBearer, nil
	default:
		return "", fmt.Errorf("authMethod must be \"bearer\"")
	}
}

func writeLLMBadRequest(w http.ResponseWriter, err error) {
	message := strings.TrimSpace(err.Error())
	if message == "" {
		message = "invalid request"
	}
	writeLLMValidationMessage(w, message)
}

func writeLLMValidationMessage(w http.ResponseWriter, message string) {
	problem := NewProblem(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), message)
	WriteProblem(w, problem)
}

func writeLLMJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode llm response: %v", err)
	}
}
