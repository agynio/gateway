package gateway

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
	"google.golang.org/grpc"
)

type fakeLLMClient struct {
	testModelReq  *llmv1.TestModelRequest
	testModelResp *llmv1.TestModelResponse
}

func (f *fakeLLMClient) CreateLLMProvider(context.Context, *llmv1.CreateLLMProviderRequest, ...grpc.CallOption) (*llmv1.CreateLLMProviderResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) GetLLMProvider(context.Context, *llmv1.GetLLMProviderRequest, ...grpc.CallOption) (*llmv1.GetLLMProviderResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) UpdateLLMProvider(context.Context, *llmv1.UpdateLLMProviderRequest, ...grpc.CallOption) (*llmv1.UpdateLLMProviderResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) DeleteLLMProvider(context.Context, *llmv1.DeleteLLMProviderRequest, ...grpc.CallOption) (*llmv1.DeleteLLMProviderResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) ListLLMProviders(context.Context, *llmv1.ListLLMProvidersRequest, ...grpc.CallOption) (*llmv1.ListLLMProvidersResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) CreateModel(context.Context, *llmv1.CreateModelRequest, ...grpc.CallOption) (*llmv1.CreateModelResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) GetModel(context.Context, *llmv1.GetModelRequest, ...grpc.CallOption) (*llmv1.GetModelResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) UpdateModel(context.Context, *llmv1.UpdateModelRequest, ...grpc.CallOption) (*llmv1.UpdateModelResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) DeleteModel(context.Context, *llmv1.DeleteModelRequest, ...grpc.CallOption) (*llmv1.DeleteModelResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) ListModels(context.Context, *llmv1.ListModelsRequest, ...grpc.CallOption) (*llmv1.ListModelsResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) ResolveModel(context.Context, *llmv1.ResolveModelRequest, ...grpc.CallOption) (*llmv1.ResolveModelResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeLLMClient) TestModel(ctx context.Context, req *llmv1.TestModelRequest, _ ...grpc.CallOption) (*llmv1.TestModelResponse, error) {
	f.testModelReq = req
	return f.testModelResp, nil
}

func TestGatewayTestModelForwards(t *testing.T) {
	fake := &fakeLLMClient{testModelResp: &llmv1.TestModelResponse{OutputText: "ok"}}
	gateway := &Gateway{llm: fake}

	resp, err := gateway.TestModel(context.Background(), connect.NewRequest(&llmv1.TestModelRequest{ModelId: "model-1"}))
	if err != nil {
		t.Fatalf("TestModel: %v", err)
	}
	if fake.testModelReq == nil || fake.testModelReq.GetModelId() != "model-1" {
		t.Fatalf("expected forwarded model id")
	}
	if resp.Msg.GetOutputText() != "ok" {
		t.Fatalf("expected output ok, got %q", resp.Msg.GetOutputText())
	}
}
