package handlers

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agynio/gateway/internal/llmclient"
	"github.com/agynio/gateway/internal/llmgen"
)

const (
	llmBasePath    = "/llm/v1"
	llmTotalAbsent = -1
)

func LLMBasePath() string {
	return llmBasePath
}

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

// LLMHandler implements llmgen.StrictServerInterface for the LLM API.
type LLMHandler struct {
	client LLMClient
}

func NewLLMHandler(client LLMClient) *LLMHandler {
	if client == nil {
		panic("llm client is required")
	}
	return &LLMHandler{client: client}
}

func (h *LLMHandler) GetModels(ctx context.Context, request llmgen.GetModelsRequestObject) (llmgen.GetModelsResponseObject, error) {
	page, perPage := normalizePagination(request.Params.Page, request.Params.PerPage)
	pageToken := llmPageToken(page, perPage)
	providerID := ""
	if request.Params.ProviderId != nil {
		providerID = request.Params.ProviderId.String()
	}

	models, _, err := h.client.ListModels(ctx, int32(perPage), pageToken, providerID)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]llmgen.Model, 0, len(models))
	for _, model := range models {
		converted, err := modelFromClient(model)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := llmgen.PaginatedModels{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   llmTotalAbsent,
	}

	return llmgen.GetModels200JSONResponse(payload), nil
}

func (h *LLMHandler) PostModels(ctx context.Context, request llmgen.PostModelsRequestObject) (llmgen.PostModelsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	params, err := modelCreateToParams(*request.Body)
	if err != nil {
		return nil, requestProblem(err)
	}

	model, err := h.client.CreateModel(ctx, params)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	converted, err := modelFromClient(model)
	if err != nil {
		return nil, responseProblem(err)
	}

	return llmgen.PostModels201JSONResponse(converted), nil
}

func (h *LLMHandler) DeleteModelsId(ctx context.Context, request llmgen.DeleteModelsIdRequestObject) (llmgen.DeleteModelsIdResponseObject, error) {
	if err := h.client.DeleteModel(ctx, request.Id.String()); err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return llmgen.DeleteModelsId204Response{}, nil
}

func (h *LLMHandler) GetModelsId(ctx context.Context, request llmgen.GetModelsIdRequestObject) (llmgen.GetModelsIdResponseObject, error) {
	model, err := h.client.GetModel(ctx, request.Id.String())
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	converted, err := modelFromClient(model)
	if err != nil {
		return nil, responseProblem(err)
	}

	return llmgen.GetModelsId200JSONResponse(converted), nil
}

func (h *LLMHandler) PatchModelsId(ctx context.Context, request llmgen.PatchModelsIdRequestObject) (llmgen.PatchModelsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	params, err := modelUpdateToParams(*request.Body)
	if err != nil {
		return nil, requestProblem(err)
	}

	model, err := h.client.UpdateModel(ctx, request.Id.String(), params)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	converted, err := modelFromClient(model)
	if err != nil {
		return nil, responseProblem(err)
	}

	return llmgen.PatchModelsId200JSONResponse(converted), nil
}

func (h *LLMHandler) GetProviders(ctx context.Context, request llmgen.GetProvidersRequestObject) (llmgen.GetProvidersResponseObject, error) {
	page, perPage := normalizePagination(request.Params.Page, request.Params.PerPage)
	pageToken := llmPageToken(page, perPage)

	providers, _, err := h.client.ListProviders(ctx, int32(perPage), pageToken)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]llmgen.LLMProvider, 0, len(providers))
	for _, provider := range providers {
		converted, err := providerFromClient(provider)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := llmgen.PaginatedLLMProviders{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   llmTotalAbsent,
	}

	return llmgen.GetProviders200JSONResponse(payload), nil
}

func (h *LLMHandler) PostProviders(ctx context.Context, request llmgen.PostProvidersRequestObject) (llmgen.PostProvidersResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	params, err := providerCreateToParams(*request.Body)
	if err != nil {
		return nil, requestProblem(err)
	}

	provider, err := h.client.CreateProvider(ctx, params)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	converted, err := providerFromClient(provider)
	if err != nil {
		return nil, responseProblem(err)
	}

	return llmgen.PostProviders201JSONResponse(converted), nil
}

func (h *LLMHandler) DeleteProvidersId(ctx context.Context, request llmgen.DeleteProvidersIdRequestObject) (llmgen.DeleteProvidersIdResponseObject, error) {
	if err := h.client.DeleteProvider(ctx, request.Id.String()); err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return llmgen.DeleteProvidersId204Response{}, nil
}

func (h *LLMHandler) GetProvidersId(ctx context.Context, request llmgen.GetProvidersIdRequestObject) (llmgen.GetProvidersIdResponseObject, error) {
	provider, err := h.client.GetProvider(ctx, request.Id.String())
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	converted, err := providerFromClient(provider)
	if err != nil {
		return nil, responseProblem(err)
	}

	return llmgen.GetProvidersId200JSONResponse(converted), nil
}

func (h *LLMHandler) PatchProvidersId(ctx context.Context, request llmgen.PatchProvidersIdRequestObject) (llmgen.PatchProvidersIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	params, err := providerUpdateToParams(*request.Body)
	if err != nil {
		return nil, requestProblem(err)
	}

	provider, err := h.client.UpdateProvider(ctx, request.Id.String(), params)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	converted, err := providerFromClient(provider)
	if err != nil {
		return nil, responseProblem(err)
	}

	return llmgen.PatchProvidersId200JSONResponse(converted), nil
}

func (h *LLMHandler) PostResponses(ctx context.Context, request llmgen.PostResponsesRequestObject) (llmgen.PostResponsesResponseObject, error) {
	return nil, grpcErrorToProblem(status.Error(codes.Unimplemented, "responses endpoint must be proxied"))
}

func llmPageToken(page, perPage int) string {
	if page <= 1 {
		return ""
	}
	return strconv.Itoa((page - 1) * perPage)
}
