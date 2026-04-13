package gateway

import (
	"context"

	"connectrpc.com/connect"
	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
)

func (g *Gateway) CreateLLMProvider(ctx context.Context, req *connect.Request[llmv1.CreateLLMProviderRequest]) (*connect.Response[llmv1.CreateLLMProviderResponse], error) {
	resp, err := g.llm.CreateLLMProvider(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetLLMProvider(ctx context.Context, req *connect.Request[llmv1.GetLLMProviderRequest]) (*connect.Response[llmv1.GetLLMProviderResponse], error) {
	resp, err := g.llm.GetLLMProvider(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateLLMProvider(ctx context.Context, req *connect.Request[llmv1.UpdateLLMProviderRequest]) (*connect.Response[llmv1.UpdateLLMProviderResponse], error) {
	resp, err := g.llm.UpdateLLMProvider(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteLLMProvider(ctx context.Context, req *connect.Request[llmv1.DeleteLLMProviderRequest]) (*connect.Response[llmv1.DeleteLLMProviderResponse], error) {
	resp, err := g.llm.DeleteLLMProvider(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListLLMProviders(ctx context.Context, req *connect.Request[llmv1.ListLLMProvidersRequest]) (*connect.Response[llmv1.ListLLMProvidersResponse], error) {
	resp, err := g.llm.ListLLMProviders(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateModel(ctx context.Context, req *connect.Request[llmv1.CreateModelRequest]) (*connect.Response[llmv1.CreateModelResponse], error) {
	resp, err := g.llm.CreateModel(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetModel(ctx context.Context, req *connect.Request[llmv1.GetModelRequest]) (*connect.Response[llmv1.GetModelResponse], error) {
	resp, err := g.llm.GetModel(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateModel(ctx context.Context, req *connect.Request[llmv1.UpdateModelRequest]) (*connect.Response[llmv1.UpdateModelResponse], error) {
	resp, err := g.llm.UpdateModel(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteModel(ctx context.Context, req *connect.Request[llmv1.DeleteModelRequest]) (*connect.Response[llmv1.DeleteModelResponse], error) {
	resp, err := g.llm.DeleteModel(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListModels(ctx context.Context, req *connect.Request[llmv1.ListModelsRequest]) (*connect.Response[llmv1.ListModelsResponse], error) {
	resp, err := g.llm.ListModels(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) TestModel(ctx context.Context, req *connect.Request[llmv1.TestModelRequest]) (*connect.Response[llmv1.TestModelResponse], error) {
	resp, err := g.llm.TestModel(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
