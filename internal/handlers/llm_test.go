package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agynio/gateway/internal/llmclient"
	"github.com/agynio/gateway/internal/llmgen"
)

type llmCall struct {
	method string
	assert func(any)
	resp   any
	token  string
	err    error
}

type stubLLMClient struct {
	t     *testing.T
	calls []llmCall
	idx   int
}

type listProvidersArgs struct {
	pageSize  int32
	pageToken string
}

type listModelsArgs struct {
	pageSize   int32
	pageToken  string
	providerID string
}

type updateProviderArgs struct {
	id     string
	params llmclient.UpdateProviderParams
}

type updateModelArgs struct {
	id     string
	params llmclient.UpdateModelParams
}

func (s *stubLLMClient) Expect(call llmCall) {
	s.calls = append(s.calls, call)
}

func (s *stubLLMClient) AssertDone() {
	if s.idx != len(s.calls) {
		s.t.Fatalf("expected %d calls, got %d", len(s.calls), s.idx)
	}
}

func (s *stubLLMClient) nextCall(method string, req any) llmCall {
	if s.idx >= len(s.calls) {
		s.t.Fatalf("unexpected call %s", method)
	}
	call := s.calls[s.idx]
	s.idx++
	if call.method != method {
		s.t.Fatalf("unexpected method: got %s want %s", method, call.method)
	}
	if call.assert != nil {
		call.assert(req)
	}
	return call
}

func (s *stubLLMClient) CreateProvider(ctx context.Context, params llmclient.CreateProviderParams) (llmclient.LLMProvider, error) {
	call := s.nextCall("CreateProvider", params)
	if call.resp == nil {
		return llmclient.LLMProvider{}, call.err
	}
	resp, ok := call.resp.(llmclient.LLMProvider)
	if !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubLLMClient) GetProvider(ctx context.Context, id string) (llmclient.LLMProvider, error) {
	call := s.nextCall("GetProvider", id)
	if call.resp == nil {
		return llmclient.LLMProvider{}, call.err
	}
	resp, ok := call.resp.(llmclient.LLMProvider)
	if !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubLLMClient) UpdateProvider(ctx context.Context, id string, params llmclient.UpdateProviderParams) (llmclient.LLMProvider, error) {
	call := s.nextCall("UpdateProvider", updateProviderArgs{id: id, params: params})
	if call.resp == nil {
		return llmclient.LLMProvider{}, call.err
	}
	resp, ok := call.resp.(llmclient.LLMProvider)
	if !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubLLMClient) DeleteProvider(ctx context.Context, id string) error {
	call := s.nextCall("DeleteProvider", id)
	return call.err
}

func (s *stubLLMClient) ListProviders(ctx context.Context, pageSize int32, pageToken string) ([]llmclient.LLMProvider, string, error) {
	call := s.nextCall("ListProviders", listProvidersArgs{pageSize: pageSize, pageToken: pageToken})
	if call.resp == nil {
		return nil, call.token, call.err
	}
	resp, ok := call.resp.([]llmclient.LLMProvider)
	if !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.token, call.err
}

func (s *stubLLMClient) CreateModel(ctx context.Context, params llmclient.CreateModelParams) (llmclient.Model, error) {
	call := s.nextCall("CreateModel", params)
	if call.resp == nil {
		return llmclient.Model{}, call.err
	}
	resp, ok := call.resp.(llmclient.Model)
	if !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubLLMClient) GetModel(ctx context.Context, id string) (llmclient.Model, error) {
	call := s.nextCall("GetModel", id)
	if call.resp == nil {
		return llmclient.Model{}, call.err
	}
	resp, ok := call.resp.(llmclient.Model)
	if !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubLLMClient) UpdateModel(ctx context.Context, id string, params llmclient.UpdateModelParams) (llmclient.Model, error) {
	call := s.nextCall("UpdateModel", updateModelArgs{id: id, params: params})
	if call.resp == nil {
		return llmclient.Model{}, call.err
	}
	resp, ok := call.resp.(llmclient.Model)
	if !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubLLMClient) DeleteModel(ctx context.Context, id string) error {
	call := s.nextCall("DeleteModel", id)
	return call.err
}

func (s *stubLLMClient) ListModels(ctx context.Context, pageSize int32, pageToken, llmProviderID string) ([]llmclient.Model, string, error) {
	call := s.nextCall("ListModels", listModelsArgs{pageSize: pageSize, pageToken: pageToken, providerID: llmProviderID})
	if call.resp == nil {
		return nil, call.token, call.err
	}
	resp, ok := call.resp.([]llmclient.Model)
	if !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.token, call.err
}

func TestLLMGetProvidersPagination(t *testing.T) {
	stub := &stubLLMClient{t: t}
	providerID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	createdAt := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)
	updatedAt := time.Date(2024, 2, 4, 4, 5, 6, 0, time.UTC)

	stub.Expect(llmCall{
		method: "ListProviders",
		assert: func(req any) {
			input := req.(listProvidersArgs)
			if input.pageSize != 3 {
				t.Fatalf("unexpected page size: %d", input.pageSize)
			}
			if input.pageToken != "3" {
				t.Fatalf("unexpected page token: %s", input.pageToken)
			}
		},
		resp: []llmclient.LLMProvider{
			{
				ID:         providerID.String(),
				Endpoint:   "https://llm.example.com",
				AuthMethod: llmclient.AuthMethodBearer,
				CreatedAt:  createdAt,
				UpdatedAt:  &updatedAt,
			},
		},
	})

	page := 2
	perPage := 3
	h := NewLLMHandler(stub)
	resp, err := h.GetProviders(context.Background(), llmgen.GetProvidersRequestObject{Params: llmgen.GetProvidersParams{
		Page:    &page,
		PerPage: &perPage,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := resp.(llmgen.GetProviders200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if list.Page != page || list.PerPage != perPage {
		t.Fatalf("unexpected pagination: page %d perPage %d", list.Page, list.PerPage)
	}
	if list.Total != -1 {
		t.Fatalf("unexpected total: %d", list.Total)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(list.Items))
	}
	if list.Items[0].Id.String() != providerID.String() {
		t.Fatalf("unexpected provider id: %s", list.Items[0].Id.String())
	}
	if list.Items[0].AuthMethod != llmgen.Bearer {
		t.Fatalf("unexpected auth method: %s", list.Items[0].AuthMethod)
	}
	if list.Items[0].UpdatedAt == nil || !list.Items[0].UpdatedAt.Equal(updatedAt) {
		t.Fatalf("unexpected updatedAt: %v", list.Items[0].UpdatedAt)
	}

	stub.AssertDone()
}

func TestLLMPostProviders(t *testing.T) {
	stub := &stubLLMClient{t: t}
	providerID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	createdAt := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)

	stub.Expect(llmCall{
		method: "CreateProvider",
		assert: func(req any) {
			input := req.(llmclient.CreateProviderParams)
			if input.Endpoint != "https://llm.example.com" {
				t.Fatalf("unexpected endpoint: %s", input.Endpoint)
			}
			if input.AuthMethod != llmclient.AuthMethodBearer {
				t.Fatalf("unexpected auth method: %s", input.AuthMethod)
			}
			if input.Token != "secret" {
				t.Fatalf("unexpected token: %s", input.Token)
			}
		},
		resp: llmclient.LLMProvider{
			ID:         providerID.String(),
			Endpoint:   "https://llm.example.com",
			AuthMethod: llmclient.AuthMethodBearer,
			CreatedAt:  createdAt,
		},
	})

	request := llmgen.PostProvidersRequestObject{Body: &llmgen.LLMProviderCreateRequest{
		Endpoint:   "https://llm.example.com",
		AuthMethod: llmgen.Bearer,
		Token:      "secret",
	}}

	h := NewLLMHandler(stub)
	resp, err := h.PostProviders(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(llmgen.PostProviders201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Id.String() != providerID.String() {
		t.Fatalf("unexpected provider id: %s", created.Id.String())
	}
	if created.Endpoint != "https://llm.example.com" {
		t.Fatalf("unexpected endpoint: %s", created.Endpoint)
	}

	stub.AssertDone()
}

func TestLLMGetProviderByID(t *testing.T) {
	stub := &stubLLMClient{t: t}
	providerID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	stub.Expect(llmCall{
		method: "GetProvider",
		assert: func(req any) {
			if req.(string) != providerID.String() {
				t.Fatalf("unexpected provider id: %s", req.(string))
			}
		},
		resp: llmclient.LLMProvider{
			ID:         providerID.String(),
			Endpoint:   "https://llm.example.com",
			AuthMethod: llmclient.AuthMethodBearer,
			CreatedAt:  createdAt,
		},
	})

	h := NewLLMHandler(stub)
	resp, err := h.GetProvidersId(context.Background(), llmgen.GetProvidersIdRequestObject{Id: openapi_types.UUID(providerID)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	provider, ok := resp.(llmgen.GetProvidersId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if provider.Id.String() != providerID.String() {
		t.Fatalf("unexpected provider id: %s", provider.Id.String())
	}

	stub.AssertDone()
}

func TestLLMPatchProvider(t *testing.T) {
	stub := &stubLLMClient{t: t}
	providerID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	createdAt := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	endpoint := "https://llm.example.com"
	authMethod := llmgen.Bearer
	token := "secret"

	stub.Expect(llmCall{
		method: "UpdateProvider",
		assert: func(req any) {
			input := req.(updateProviderArgs)
			if input.id != providerID.String() {
				t.Fatalf("unexpected provider id: %s", input.id)
			}
			if input.params.Endpoint == nil || *input.params.Endpoint != endpoint {
				t.Fatalf("unexpected endpoint: %v", input.params.Endpoint)
			}
			if input.params.AuthMethod == nil || *input.params.AuthMethod != llmclient.AuthMethodBearer {
				t.Fatalf("unexpected auth method: %v", input.params.AuthMethod)
			}
			if input.params.Token == nil || *input.params.Token != token {
				t.Fatalf("unexpected token: %v", input.params.Token)
			}
		},
		resp: llmclient.LLMProvider{
			ID:         providerID.String(),
			Endpoint:   endpoint,
			AuthMethod: llmclient.AuthMethodBearer,
			CreatedAt:  createdAt,
		},
	})

	request := llmgen.PatchProvidersIdRequestObject{
		Id: openapi_types.UUID(providerID),
		Body: &llmgen.LLMProviderUpdateRequest{
			Endpoint:   &endpoint,
			AuthMethod: &authMethod,
			Token:      &token,
		},
	}

	h := NewLLMHandler(stub)
	resp, err := h.PatchProvidersId(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, ok := resp.(llmgen.PatchProvidersId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if updated.Id.String() != providerID.String() {
		t.Fatalf("unexpected provider id: %s", updated.Id.String())
	}

	stub.AssertDone()
}

func TestLLMDeleteProvider(t *testing.T) {
	stub := &stubLLMClient{t: t}
	providerID := uuid.MustParse("55555555-5555-5555-5555-555555555555")

	stub.Expect(llmCall{
		method: "DeleteProvider",
		assert: func(req any) {
			if req.(string) != providerID.String() {
				t.Fatalf("unexpected provider id: %s", req.(string))
			}
		},
	})

	h := NewLLMHandler(stub)
	resp, err := h.DeleteProvidersId(context.Background(), llmgen.DeleteProvidersIdRequestObject{Id: openapi_types.UUID(providerID)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := resp.(llmgen.DeleteProvidersId204Response); !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}

	stub.AssertDone()
}

func TestLLMGetModelsPagination(t *testing.T) {
	stub := &stubLLMClient{t: t}
	modelID := uuid.MustParse("66666666-6666-6666-6666-666666666666")
	providerID := uuid.MustParse("77777777-7777-7777-7777-777777777777")
	createdAt := time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)

	stub.Expect(llmCall{
		method: "ListModels",
		assert: func(req any) {
			input := req.(listModelsArgs)
			if input.pageSize != 2 {
				t.Fatalf("unexpected page size: %d", input.pageSize)
			}
			if input.pageToken != "4" {
				t.Fatalf("unexpected page token: %s", input.pageToken)
			}
			if input.providerID != providerID.String() {
				t.Fatalf("unexpected provider id: %s", input.providerID)
			}
		},
		resp: []llmclient.Model{
			{
				ID:            modelID.String(),
				Name:          "gpt-4",
				LLMProviderID: providerID.String(),
				RemoteName:    "gpt-4",
				CreatedAt:     createdAt,
			},
		},
	})

	page := 3
	perPage := 2
	providerParam := openapi_types.UUID(providerID)
	h := NewLLMHandler(stub)
	resp, err := h.GetModels(context.Background(), llmgen.GetModelsRequestObject{Params: llmgen.GetModelsParams{
		ProviderId: &providerParam,
		Page:       &page,
		PerPage:    &perPage,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := resp.(llmgen.GetModels200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if list.Total != -1 {
		t.Fatalf("unexpected total: %d", list.Total)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 model, got %d", len(list.Items))
	}
	if list.Items[0].Id.String() != modelID.String() {
		t.Fatalf("unexpected model id: %s", list.Items[0].Id.String())
	}

	stub.AssertDone()
}

func TestLLMPostModels(t *testing.T) {
	stub := &stubLLMClient{t: t}
	modelID := uuid.MustParse("88888888-8888-8888-8888-888888888888")
	providerID := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	createdAt := time.Date(2024, 4, 5, 6, 7, 8, 0, time.UTC)

	stub.Expect(llmCall{
		method: "CreateModel",
		assert: func(req any) {
			input := req.(llmclient.CreateModelParams)
			if input.Name != "gpt-4" {
				t.Fatalf("unexpected name: %s", input.Name)
			}
			if input.LLMProviderID != providerID.String() {
				t.Fatalf("unexpected provider id: %s", input.LLMProviderID)
			}
			if input.RemoteName != "gpt-4" {
				t.Fatalf("unexpected remote name: %s", input.RemoteName)
			}
		},
		resp: llmclient.Model{
			ID:            modelID.String(),
			Name:          "gpt-4",
			LLMProviderID: providerID.String(),
			RemoteName:    "gpt-4",
			CreatedAt:     createdAt,
		},
	})

	request := llmgen.PostModelsRequestObject{Body: &llmgen.ModelCreateRequest{
		Name:          "gpt-4",
		LlmProviderId: openapi_types.UUID(providerID),
		RemoteName:    "gpt-4",
	}}

	h := NewLLMHandler(stub)
	resp, err := h.PostModels(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(llmgen.PostModels201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Id.String() != modelID.String() {
		t.Fatalf("unexpected model id: %s", created.Id.String())
	}

	stub.AssertDone()
}

func TestLLMGetModelByID(t *testing.T) {
	stub := &stubLLMClient{t: t}
	modelID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	providerID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	createdAt := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)

	stub.Expect(llmCall{
		method: "GetModel",
		assert: func(req any) {
			if req.(string) != modelID.String() {
				t.Fatalf("unexpected model id: %s", req.(string))
			}
		},
		resp: llmclient.Model{
			ID:            modelID.String(),
			Name:          "gpt-4",
			LLMProviderID: providerID.String(),
			RemoteName:    "gpt-4",
			CreatedAt:     createdAt,
		},
	})

	h := NewLLMHandler(stub)
	resp, err := h.GetModelsId(context.Background(), llmgen.GetModelsIdRequestObject{Id: openapi_types.UUID(modelID)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model, ok := resp.(llmgen.GetModelsId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if model.Id.String() != modelID.String() {
		t.Fatalf("unexpected model id: %s", model.Id.String())
	}

	stub.AssertDone()
}

func TestLLMPatchModel(t *testing.T) {
	stub := &stubLLMClient{t: t}
	modelID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	providerID := uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")
	createdAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	name := "gpt-4"
	remoteName := "gpt-4"
	providerParam := openapi_types.UUID(providerID)

	stub.Expect(llmCall{
		method: "UpdateModel",
		assert: func(req any) {
			input := req.(updateModelArgs)
			if input.id != modelID.String() {
				t.Fatalf("unexpected model id: %s", input.id)
			}
			if input.params.Name == nil || *input.params.Name != name {
				t.Fatalf("unexpected name: %v", input.params.Name)
			}
			if input.params.RemoteName == nil || *input.params.RemoteName != remoteName {
				t.Fatalf("unexpected remote name: %v", input.params.RemoteName)
			}
			if input.params.LLMProviderID == nil || *input.params.LLMProviderID != providerID.String() {
				t.Fatalf("unexpected provider id: %v", input.params.LLMProviderID)
			}
		},
		resp: llmclient.Model{
			ID:            modelID.String(),
			Name:          name,
			LLMProviderID: providerID.String(),
			RemoteName:    remoteName,
			CreatedAt:     createdAt,
		},
	})

	request := llmgen.PatchModelsIdRequestObject{
		Id: openapi_types.UUID(modelID),
		Body: &llmgen.ModelUpdateRequest{
			Name:          &name,
			RemoteName:    &remoteName,
			LlmProviderId: &providerParam,
		},
	}

	h := NewLLMHandler(stub)
	resp, err := h.PatchModelsId(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, ok := resp.(llmgen.PatchModelsId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if updated.Id.String() != modelID.String() {
		t.Fatalf("unexpected model id: %s", updated.Id.String())
	}

	stub.AssertDone()
}

func TestLLMDeleteModel(t *testing.T) {
	stub := &stubLLMClient{t: t}
	modelID := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")

	stub.Expect(llmCall{
		method: "DeleteModel",
		assert: func(req any) {
			if req.(string) != modelID.String() {
				t.Fatalf("unexpected model id: %s", req.(string))
			}
		},
	})

	h := NewLLMHandler(stub)
	resp, err := h.DeleteModelsId(context.Background(), llmgen.DeleteModelsIdRequestObject{Id: openapi_types.UUID(modelID)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := resp.(llmgen.DeleteModelsId204Response); !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}

	stub.AssertDone()
}

func TestLLMGrpcErrorMapping(t *testing.T) {
	stub := &stubLLMClient{t: t}
	stub.Expect(llmCall{
		method: "ListProviders",
		err:    status.Error(codes.NotFound, "missing"),
	})

	h := NewLLMHandler(stub)
	_, err := h.GetProviders(context.Background(), llmgen.GetProvidersRequestObject{})
	if err == nil {
		t.Fatalf("expected error")
	}

	var problemErr *ProblemError
	if !errors.As(err, &problemErr) {
		t.Fatalf("expected ProblemError, got %T", err)
	}
	if problemErr.Problem.Status != 404 {
		t.Fatalf("unexpected status: %d", problemErr.Problem.Status)
	}

	stub.AssertDone()
}
