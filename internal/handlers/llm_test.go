package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agynio/gateway/internal/llmclient"
)

type stubLLMClient struct {
	t                *testing.T
	createProviderFn func(context.Context, llmclient.CreateProviderParams) (llmclient.LLMProvider, error)
	getProviderFn    func(context.Context, string) (llmclient.LLMProvider, error)
	updateProviderFn func(context.Context, string, llmclient.UpdateProviderParams) (llmclient.LLMProvider, error)
	deleteProviderFn func(context.Context, string) error
	listProvidersFn  func(context.Context, int32, string) ([]llmclient.LLMProvider, string, error)
	createModelFn    func(context.Context, llmclient.CreateModelParams) (llmclient.Model, error)
	getModelFn       func(context.Context, string) (llmclient.Model, error)
	updateModelFn    func(context.Context, string, llmclient.UpdateModelParams) (llmclient.Model, error)
	deleteModelFn    func(context.Context, string) error
	listModelsFn     func(context.Context, int32, string, string) ([]llmclient.Model, string, error)
}

func (s *stubLLMClient) CreateProvider(ctx context.Context, params llmclient.CreateProviderParams) (llmclient.LLMProvider, error) {
	if s.createProviderFn == nil {
		s.t.Fatalf("unexpected CreateProvider call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.createProviderFn(ctx, params)
}

func (s *stubLLMClient) GetProvider(ctx context.Context, id string) (llmclient.LLMProvider, error) {
	if s.getProviderFn == nil {
		s.t.Fatalf("unexpected GetProvider call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.getProviderFn(ctx, id)
}

func (s *stubLLMClient) UpdateProvider(ctx context.Context, id string, params llmclient.UpdateProviderParams) (llmclient.LLMProvider, error) {
	if s.updateProviderFn == nil {
		s.t.Fatalf("unexpected UpdateProvider call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.updateProviderFn(ctx, id, params)
}

func (s *stubLLMClient) DeleteProvider(ctx context.Context, id string) error {
	if s.deleteProviderFn == nil {
		s.t.Fatalf("unexpected DeleteProvider call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.deleteProviderFn(ctx, id)
}

func (s *stubLLMClient) ListProviders(ctx context.Context, pageSize int32, pageToken string) ([]llmclient.LLMProvider, string, error) {
	if s.listProvidersFn == nil {
		s.t.Fatalf("unexpected ListProviders call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.listProvidersFn(ctx, pageSize, pageToken)
}

func (s *stubLLMClient) CreateModel(ctx context.Context, params llmclient.CreateModelParams) (llmclient.Model, error) {
	if s.createModelFn == nil {
		s.t.Fatalf("unexpected CreateModel call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.createModelFn(ctx, params)
}

func (s *stubLLMClient) GetModel(ctx context.Context, id string) (llmclient.Model, error) {
	if s.getModelFn == nil {
		s.t.Fatalf("unexpected GetModel call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.getModelFn(ctx, id)
}

func (s *stubLLMClient) UpdateModel(ctx context.Context, id string, params llmclient.UpdateModelParams) (llmclient.Model, error) {
	if s.updateModelFn == nil {
		s.t.Fatalf("unexpected UpdateModel call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.updateModelFn(ctx, id, params)
}

func (s *stubLLMClient) DeleteModel(ctx context.Context, id string) error {
	if s.deleteModelFn == nil {
		s.t.Fatalf("unexpected DeleteModel call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.deleteModelFn(ctx, id)
}

func (s *stubLLMClient) ListModels(ctx context.Context, pageSize int32, pageToken, llmProviderID string) ([]llmclient.Model, string, error) {
	if s.listModelsFn == nil {
		s.t.Fatalf("unexpected ListModels call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	return s.listModelsFn(ctx, pageSize, pageToken, llmProviderID)
}

func TestLLMHandlerCreateProviderSuccess(t *testing.T) {
	createdAt := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)
	expected := llmclient.LLMProvider{
		ID:         "prov-1",
		Endpoint:   "https://llm.example.com",
		AuthMethod: llmclient.AuthMethodBearer,
		CreatedAt:  createdAt,
	}

	called := false
	client := &stubLLMClient{
		t: t,
		createProviderFn: func(ctx context.Context, params llmclient.CreateProviderParams) (llmclient.LLMProvider, error) {
			called = true
			if params.Endpoint != expected.Endpoint {
				t.Fatalf("unexpected endpoint: %s", params.Endpoint)
			}
			if params.AuthMethod != expected.AuthMethod {
				t.Fatalf("unexpected auth method: %s", params.AuthMethod)
			}
			if params.Token != "secret" {
				t.Fatalf("unexpected token")
			}
			return expected, nil
		},
	}

	payload := createProviderRequest{
		Endpoint:   expected.Endpoint,
		AuthMethod: "bearer",
		Token:      "secret",
	}
	req := newJSONRequest(t, http.MethodPost, "/llm/v1/providers", payload)
	resp := httptest.NewRecorder()

	NewLLMHandler(client).CreateProvider(resp, req)

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

	var body llmclient.LLMProvider
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestLLMHandlerCreateProviderValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload any
		detail  string
	}{
		{
			name:    "missing endpoint",
			payload: createProviderRequest{AuthMethod: "bearer", Token: "secret"},
			detail:  "endpoint is required",
		},
		{
			name:    "missing auth method",
			payload: createProviderRequest{Endpoint: "https://llm.example.com", Token: "secret"},
			detail:  "authMethod is required",
		},
		{
			name:    "missing token",
			payload: createProviderRequest{Endpoint: "https://llm.example.com", AuthMethod: "bearer"},
			detail:  "token is required for bearer auth",
		},
		{
			name:    "invalid auth method",
			payload: createProviderRequest{Endpoint: "https://llm.example.com", AuthMethod: "basic", Token: "secret"},
			detail:  "authMethod must be \"bearer\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubLLMClient{t: t}
			req := newJSONRequest(t, http.MethodPost, "/llm/v1/providers", tt.payload)
			resp := httptest.NewRecorder()

			NewLLMHandler(client).CreateProvider(resp, req)

			assertProblemResponse(t, resp, http.StatusBadRequest, tt.detail)
		})
	}
}

func TestLLMHandlerCreateModelSuccess(t *testing.T) {
	createdAt := time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)
	expected := llmclient.Model{
		ID:            "model-1",
		Name:          "gpt-4",
		LLMProviderID: "prov-1",
		RemoteName:    "gpt-4",
		CreatedAt:     createdAt,
	}

	called := false
	client := &stubLLMClient{
		t: t,
		createModelFn: func(ctx context.Context, params llmclient.CreateModelParams) (llmclient.Model, error) {
			called = true
			if params.Name != expected.Name {
				t.Fatalf("unexpected name: %s", params.Name)
			}
			if params.LLMProviderID != expected.LLMProviderID {
				t.Fatalf("unexpected provider id: %s", params.LLMProviderID)
			}
			if params.RemoteName != expected.RemoteName {
				t.Fatalf("unexpected remote name: %s", params.RemoteName)
			}
			return expected, nil
		},
	}

	payload := createModelRequest{
		Name:          expected.Name,
		LLMProviderID: expected.LLMProviderID,
		RemoteName:    expected.RemoteName,
	}
	req := newJSONRequest(t, http.MethodPost, "/llm/v1/models", payload)
	resp := httptest.NewRecorder()

	NewLLMHandler(client).CreateModel(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d got %d", http.StatusCreated, result.StatusCode)
	}
	if got := result.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type application/json got %q", got)
	}
	if !called {
		t.Fatalf("expected CreateModel to be called")
	}

	var body llmclient.Model
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestLLMHandlerCreateModelValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload any
		detail  string
	}{
		{
			name:    "missing name",
			payload: createModelRequest{LLMProviderID: "prov-1", RemoteName: "gpt-4"},
			detail:  "name is required",
		},
		{
			name:    "missing provider",
			payload: createModelRequest{Name: "gpt-4", RemoteName: "gpt-4"},
			detail:  "llmProviderId is required",
		},
		{
			name:    "missing remote name",
			payload: createModelRequest{Name: "gpt-4", LLMProviderID: "prov-1"},
			detail:  "remoteName is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubLLMClient{t: t}
			req := newJSONRequest(t, http.MethodPost, "/llm/v1/models", tt.payload)
			resp := httptest.NewRecorder()

			NewLLMHandler(client).CreateModel(resp, req)

			assertProblemResponse(t, resp, http.StatusBadRequest, tt.detail)
		})
	}
}

func TestLLMHandlerUpdateProviderPartial(t *testing.T) {
	expected := llmclient.LLMProvider{
		ID:         "prov-1",
		Endpoint:   "https://updated.example.com",
		AuthMethod: llmclient.AuthMethodBearer,
		CreatedAt:  time.Date(2024, 4, 5, 6, 7, 8, 0, time.UTC),
	}

	called := false
	client := &stubLLMClient{
		t: t,
		updateProviderFn: func(ctx context.Context, id string, params llmclient.UpdateProviderParams) (llmclient.LLMProvider, error) {
			called = true
			if id != "prov-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			if params.Endpoint == nil || *params.Endpoint != expected.Endpoint {
				t.Fatalf("unexpected endpoint")
			}
			if params.AuthMethod != nil || params.Token != nil {
				t.Fatalf("expected only endpoint update")
			}
			return expected, nil
		},
	}

	payload := updateProviderRequest{Endpoint: ptr(expected.Endpoint)}
	req := newJSONRequest(t, http.MethodPatch, "/llm/v1/providers/prov-1", payload)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).UpdateProvider(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}
	if !called {
		t.Fatalf("expected UpdateProvider to be called")
	}

	var body llmclient.LLMProvider
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestLLMHandlerUpdateProviderEmptyBody(t *testing.T) {
	client := &stubLLMClient{t: t}
	req := httptest.NewRequest(http.MethodPatch, "/llm/v1/providers/prov-1", nil)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).UpdateProvider(resp, req)

	assertProblemResponse(t, resp, http.StatusBadRequest, "request body is required")
}

func TestLLMHandlerUpdateModelPartial(t *testing.T) {
	expected := llmclient.Model{
		ID:            "model-1",
		Name:          "gpt-4",
		LLMProviderID: "prov-1",
		RemoteName:    "gpt-4-turbo",
		CreatedAt:     time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC),
	}

	called := false
	client := &stubLLMClient{
		t: t,
		updateModelFn: func(ctx context.Context, id string, params llmclient.UpdateModelParams) (llmclient.Model, error) {
			called = true
			if id != "model-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			if params.RemoteName == nil || *params.RemoteName != expected.RemoteName {
				t.Fatalf("unexpected remote name")
			}
			if params.Name != nil || params.LLMProviderID != nil {
				t.Fatalf("expected only remote name update")
			}
			return expected, nil
		},
	}

	payload := updateModelRequest{RemoteName: ptr(expected.RemoteName)}
	req := newJSONRequest(t, http.MethodPatch, "/llm/v1/models/model-1", payload)
	req = withURLParam(req, "modelId", "model-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).UpdateModel(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}
	if !called {
		t.Fatalf("expected UpdateModel to be called")
	}

	var body llmclient.Model
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestLLMHandlerUpdateModelEmptyBody(t *testing.T) {
	client := &stubLLMClient{t: t}
	req := httptest.NewRequest(http.MethodPatch, "/llm/v1/models/model-1", nil)
	req = withURLParam(req, "modelId", "model-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).UpdateModel(resp, req)

	assertProblemResponse(t, resp, http.StatusBadRequest, "request body is required")
}

func TestLLMHandlerGetProviderSuccess(t *testing.T) {
	expected := llmclient.LLMProvider{
		ID:         "prov-1",
		Endpoint:   "https://llm.example.com",
		AuthMethod: llmclient.AuthMethodBearer,
		CreatedAt:  time.Date(2024, 6, 7, 8, 9, 10, 0, time.UTC),
	}

	client := &stubLLMClient{
		t: t,
		getProviderFn: func(ctx context.Context, id string) (llmclient.LLMProvider, error) {
			if id != "prov-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return expected, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers/prov-1", nil)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).GetProvider(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}

	var body llmclient.LLMProvider
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestLLMHandlerGetProviderNotFound(t *testing.T) {
	client := &stubLLMClient{
		t: t,
		getProviderFn: func(ctx context.Context, id string) (llmclient.LLMProvider, error) {
			return llmclient.LLMProvider{}, status.Error(codes.NotFound, "missing")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers/prov-1", nil)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).GetProvider(resp, req)

	assertProblemResponse(t, resp, http.StatusNotFound, "missing")
}

func TestLLMHandlerGetModelSuccess(t *testing.T) {
	expected := llmclient.Model{
		ID:            "model-1",
		Name:          "gpt-4",
		LLMProviderID: "prov-1",
		RemoteName:    "gpt-4",
		CreatedAt:     time.Date(2024, 7, 8, 9, 10, 11, 0, time.UTC),
	}

	client := &stubLLMClient{
		t: t,
		getModelFn: func(ctx context.Context, id string) (llmclient.Model, error) {
			if id != "model-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return expected, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/llm/v1/models/model-1", nil)
	req = withURLParam(req, "modelId", "model-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).GetModel(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}

	var body llmclient.Model
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestLLMHandlerGetModelNotFound(t *testing.T) {
	client := &stubLLMClient{
		t: t,
		getModelFn: func(ctx context.Context, id string) (llmclient.Model, error) {
			return llmclient.Model{}, status.Error(codes.NotFound, "missing")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/llm/v1/models/model-1", nil)
	req = withURLParam(req, "modelId", "model-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).GetModel(resp, req)

	assertProblemResponse(t, resp, http.StatusNotFound, "missing")
}

func TestLLMHandlerDeleteProviderSuccess(t *testing.T) {
	called := false
	client := &stubLLMClient{
		t: t,
		deleteProviderFn: func(ctx context.Context, id string) error {
			called = true
			if id != "prov-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/llm/v1/providers/prov-1", nil)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).DeleteProvider(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d got %d", http.StatusNoContent, result.StatusCode)
	}
	if !called {
		t.Fatalf("expected DeleteProvider to be called")
	}
}

func TestLLMHandlerDeleteProviderNotFound(t *testing.T) {
	client := &stubLLMClient{
		t: t,
		deleteProviderFn: func(ctx context.Context, id string) error {
			return status.Error(codes.NotFound, "missing")
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/llm/v1/providers/prov-1", nil)
	req = withURLParam(req, "providerId", "prov-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).DeleteProvider(resp, req)

	assertProblemResponse(t, resp, http.StatusNotFound, "missing")
}

func TestLLMHandlerDeleteModelSuccess(t *testing.T) {
	called := false
	client := &stubLLMClient{
		t: t,
		deleteModelFn: func(ctx context.Context, id string) error {
			called = true
			if id != "model-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/llm/v1/models/model-1", nil)
	req = withURLParam(req, "modelId", "model-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).DeleteModel(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d got %d", http.StatusNoContent, result.StatusCode)
	}
	if !called {
		t.Fatalf("expected DeleteModel to be called")
	}
}

func TestLLMHandlerDeleteModelNotFound(t *testing.T) {
	client := &stubLLMClient{
		t: t,
		deleteModelFn: func(ctx context.Context, id string) error {
			return status.Error(codes.NotFound, "missing")
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/llm/v1/models/model-1", nil)
	req = withURLParam(req, "modelId", "model-1")
	resp := httptest.NewRecorder()

	NewLLMHandler(client).DeleteModel(resp, req)

	assertProblemResponse(t, resp, http.StatusNotFound, "missing")
}

func TestLLMHandlerListProvidersDefaults(t *testing.T) {
	items := []llmclient.LLMProvider{{
		ID:         "prov-1",
		Endpoint:   "https://llm.example.com",
		AuthMethod: llmclient.AuthMethodBearer,
		CreatedAt:  time.Date(2024, 8, 9, 10, 11, 12, 0, time.UTC),
	}}

	client := &stubLLMClient{
		t: t,
		listProvidersFn: func(ctx context.Context, pageSize int32, pageToken string) ([]llmclient.LLMProvider, string, error) {
			if pageSize != defaultPageSize {
				t.Fatalf("unexpected page size: %d", pageSize)
			}
			if pageToken != "" {
				t.Fatalf("unexpected page token: %s", pageToken)
			}
			return items, "next", nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers", nil)
	resp := httptest.NewRecorder()

	NewLLMHandler(client).ListProviders(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}

	var body providerListResponse
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

func TestLLMHandlerListProvidersCustomPage(t *testing.T) {
	client := &stubLLMClient{
		t: t,
		listProvidersFn: func(ctx context.Context, pageSize int32, pageToken string) ([]llmclient.LLMProvider, string, error) {
			if pageSize != 50 {
				t.Fatalf("unexpected page size: %d", pageSize)
			}
			if pageToken != "cursor" {
				t.Fatalf("unexpected page token: %s", pageToken)
			}
			return []llmclient.LLMProvider{}, "", nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers?pageSize=50&pageToken=cursor", nil)
	resp := httptest.NewRecorder()

	NewLLMHandler(client).ListProviders(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}
}

func TestLLMHandlerListProvidersInvalidParams(t *testing.T) {
	client := &stubLLMClient{t: t}
	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers?pageSize=0", nil)
	resp := httptest.NewRecorder()

	NewLLMHandler(client).ListProviders(resp, req)

	assertProblemResponse(t, resp, http.StatusBadRequest, "pageSize must be a positive integer")
}

func TestLLMHandlerListModelsDefaults(t *testing.T) {
	items := []llmclient.Model{{
		ID:            "model-1",
		Name:          "gpt-4",
		LLMProviderID: "prov-1",
		RemoteName:    "gpt-4",
		CreatedAt:     time.Date(2024, 9, 10, 11, 12, 13, 0, time.UTC),
	}}

	client := &stubLLMClient{
		t: t,
		listModelsFn: func(ctx context.Context, pageSize int32, pageToken, llmProviderID string) ([]llmclient.Model, string, error) {
			if pageSize != defaultPageSize {
				t.Fatalf("unexpected page size: %d", pageSize)
			}
			if pageToken != "" {
				t.Fatalf("unexpected page token: %s", pageToken)
			}
			if llmProviderID != "prov-1" {
				t.Fatalf("unexpected provider id: %s", llmProviderID)
			}
			return items, "more", nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/llm/v1/models?llmProviderId=prov-1", nil)
	resp := httptest.NewRecorder()

	NewLLMHandler(client).ListModels(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}

	var body modelListResponse
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

func TestLLMHandlerListModelsCustomPage(t *testing.T) {
	client := &stubLLMClient{
		t: t,
		listModelsFn: func(ctx context.Context, pageSize int32, pageToken, llmProviderID string) ([]llmclient.Model, string, error) {
			if pageSize != 10 {
				t.Fatalf("unexpected page size: %d", pageSize)
			}
			if pageToken != "cursor" {
				t.Fatalf("unexpected page token: %s", pageToken)
			}
			if llmProviderID != "" {
				t.Fatalf("unexpected provider id: %s", llmProviderID)
			}
			return []llmclient.Model{}, "", nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/llm/v1/models?pageSize=10&pageToken=cursor", nil)
	resp := httptest.NewRecorder()

	NewLLMHandler(client).ListModels(resp, req)

	result := resp.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, result.StatusCode)
	}
}

func TestLLMHandlerListModelsInvalidParams(t *testing.T) {
	client := &stubLLMClient{t: t}
	req := httptest.NewRequest(http.MethodGet, "/llm/v1/models?pageSize=bad", nil)
	resp := httptest.NewRecorder()

	NewLLMHandler(client).ListModels(resp, req)

	assertProblemResponse(t, resp, http.StatusBadRequest, "pageSize must be a positive integer")
}

func TestWriteLLMGRPCErrorMappings(t *testing.T) {
	tests := []struct {
		code       codes.Code
		statusCode int
	}{
		{code: codes.NotFound, statusCode: http.StatusNotFound},
		{code: codes.InvalidArgument, statusCode: http.StatusBadRequest},
		{code: codes.AlreadyExists, statusCode: http.StatusConflict},
		{code: codes.FailedPrecondition, statusCode: http.StatusConflict},
		{code: codes.Unavailable, statusCode: http.StatusServiceUnavailable},
		{code: codes.PermissionDenied, statusCode: http.StatusForbidden},
		{code: codes.Unauthenticated, statusCode: http.StatusUnauthorized},
		{code: codes.Internal, statusCode: http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			resp := httptest.NewRecorder()
			writeLLMGRPCError(resp, status.Error(tt.code, "detail"))
			assertProblemResponse(t, resp, tt.statusCode, "detail")
		})
	}
}

func TestParsePaginationDefaults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers", nil)
	pageSize, pageToken, err := parsePagination(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pageSize != defaultPageSize {
		t.Fatalf("unexpected page size: %d", pageSize)
	}
	if pageToken != "" {
		t.Fatalf("unexpected page token: %s", pageToken)
	}
}

func TestParsePaginationPageSizeOverMax(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers?pageSize=500", nil)
	pageSize, _, err := parsePagination(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pageSize != maxPageSize {
		t.Fatalf("expected page size %d got %d", maxPageSize, pageSize)
	}
}

func TestParsePaginationPageSizeZero(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers?pageSize=0", nil)
	if _, _, err := parsePagination(req); err == nil {
		t.Fatalf("expected error")
	}
}

func TestParsePaginationPageSizeNegative(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers?pageSize=-5", nil)
	if _, _, err := parsePagination(req); err == nil {
		t.Fatalf("expected error")
	}
}

func newJSONRequest(t *testing.T, method, path string, payload any) *http.Request {
	t.Helper()

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func withURLParam(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}
