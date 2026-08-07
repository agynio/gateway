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

// Subscriptions. ResolveSubscription and CountSubscriptionsReferencingSecret
// are deliberately absent: they are internal calls over Istio, and a
// subscription credential is not something an external caller should reach
// through the Gateway.

func (g *Gateway) CreateSubscription(ctx context.Context, req *connect.Request[llmv1.CreateSubscriptionRequest]) (*connect.Response[llmv1.CreateSubscriptionResponse], error) {
	resp, err := g.llm.CreateSubscription(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetSubscription(ctx context.Context, req *connect.Request[llmv1.GetSubscriptionRequest]) (*connect.Response[llmv1.GetSubscriptionResponse], error) {
	resp, err := g.llm.GetSubscription(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateSubscription(ctx context.Context, req *connect.Request[llmv1.UpdateSubscriptionRequest]) (*connect.Response[llmv1.UpdateSubscriptionResponse], error) {
	resp, err := g.llm.UpdateSubscription(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteSubscription(ctx context.Context, req *connect.Request[llmv1.DeleteSubscriptionRequest]) (*connect.Response[llmv1.DeleteSubscriptionResponse], error) {
	resp, err := g.llm.DeleteSubscription(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListSubscriptions(ctx context.Context, req *connect.Request[llmv1.ListSubscriptionsRequest]) (*connect.Response[llmv1.ListSubscriptionsResponse], error) {
	resp, err := g.llm.ListSubscriptions(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateSubscriptionAttachment(ctx context.Context, req *connect.Request[llmv1.CreateSubscriptionAttachmentRequest]) (*connect.Response[llmv1.CreateSubscriptionAttachmentResponse], error) {
	resp, err := g.llm.CreateSubscriptionAttachment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteSubscriptionAttachment(ctx context.Context, req *connect.Request[llmv1.DeleteSubscriptionAttachmentRequest]) (*connect.Response[llmv1.DeleteSubscriptionAttachmentResponse], error) {
	resp, err := g.llm.DeleteSubscriptionAttachment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListSubscriptionAttachments(ctx context.Context, req *connect.Request[llmv1.ListSubscriptionAttachmentsRequest]) (*connect.Response[llmv1.ListSubscriptionAttachmentsResponse], error) {
	resp, err := g.llm.ListSubscriptionAttachments(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
