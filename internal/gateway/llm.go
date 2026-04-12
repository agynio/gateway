package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
)

const (
	testModelPrompt          = "Hello, world"
	testModelTimeout         = 30 * time.Second
	testModelAnthropicHeader = "2023-06-01"
)

type testModelInput struct {
	endpoint   string
	token      string
	remoteName string
	protocol   llmv1.Protocol
	authMethod llmv1.AuthMethod
}

type responsesRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responsesResponse struct {
	Output []responsesOutput `json:"output"`
}

type responsesOutput struct {
	Content []responsesContent `json:"content"`
}

type responsesContent struct {
	Text string `json:"text"`
}

type anthropicResponse struct {
	Content []anthropicContent `json:"content"`
}

type anthropicContent struct {
	Text string `json:"text"`
}

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
	modelID := strings.TrimSpace(req.Msg.GetModelId())
	if modelID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("model_id is required"))
	}

	resolved, err := g.llm.ResolveModel(ctx, &llmv1.ResolveModelRequest{ModelId: modelID})
	if err != nil {
		return connect.NewResponse(&llmv1.TestModelResponse{ErrorMessage: toConnectError(err).Message()}), nil
	}

	input, err := parseTestModelInput(resolved)
	if err != nil {
		return connect.NewResponse(&llmv1.TestModelResponse{ErrorMessage: err.Error()}), nil
	}

	outputText, errorMessage := testModel(ctx, input)
	return connect.NewResponse(&llmv1.TestModelResponse{OutputText: outputText, ErrorMessage: errorMessage}), nil
}

func parseTestModelInput(resolved *llmv1.ResolveModelResponse) (testModelInput, error) {
	if resolved == nil {
		return testModelInput{}, errors.New("model resolution missing")
	}

	endpoint := strings.TrimSpace(resolved.GetEndpoint())
	if endpoint == "" {
		return testModelInput{}, errors.New("model endpoint missing")
	}

	remoteName := strings.TrimSpace(resolved.GetRemoteName())
	if remoteName == "" {
		return testModelInput{}, errors.New("model remote name missing")
	}

	protocol := resolved.GetProtocol()
	switch protocol {
	case llmv1.Protocol_PROTOCOL_RESPONSES, llmv1.Protocol_PROTOCOL_ANTHROPIC_MESSAGES:
	default:
		return testModelInput{}, errors.New("unsupported model protocol")
	}

	authMethod := resolved.GetAuthMethod()
	switch authMethod {
	case llmv1.AuthMethod_AUTH_METHOD_BEARER, llmv1.AuthMethod_AUTH_METHOD_X_API_KEY:
	default:
		return testModelInput{}, errors.New("unsupported auth method")
	}

	token := strings.TrimSpace(resolved.GetToken())
	if token == "" {
		return testModelInput{}, errors.New("model token missing")
	}

	return testModelInput{
		endpoint:   endpoint,
		remoteName: remoteName,
		protocol:   protocol,
		authMethod: authMethod,
		token:      token,
	}, nil
}

func testModel(ctx context.Context, input testModelInput) (string, string) {
	body, headers, err := buildTestModelRequest(input)
	if err != nil {
		return "", err.Error()
	}

	requestCtx, cancel := context.WithTimeout(ctx, testModelTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(requestCtx, http.MethodPost, input.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Sprintf("invalid request: %v", err)
	}
	request.Header = headers

	client := &http.Client{Timeout: testModelTimeout}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Sprintf("request failed: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Sprintf("failed to read response: %v", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", formatUpstreamError(response.StatusCode, responseBody)
	}

	outputText, err := parseTestModelOutput(input.protocol, responseBody)
	if err != nil {
		return "", err.Error()
	}

	return outputText, ""
}

func buildTestModelRequest(input testModelInput) ([]byte, http.Header, error) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	if input.protocol == llmv1.Protocol_PROTOCOL_ANTHROPIC_MESSAGES {
		headers.Set("anthropic-version", testModelAnthropicHeader)
	}

	switch input.authMethod {
	case llmv1.AuthMethod_AUTH_METHOD_BEARER:
		headers.Set("Authorization", fmt.Sprintf("Bearer %s", input.token))
	case llmv1.AuthMethod_AUTH_METHOD_X_API_KEY:
		headers.Set("x-api-key", input.token)
	default:
		return nil, nil, errors.New("unsupported auth method")
	}

	var payload any
	switch input.protocol {
	case llmv1.Protocol_PROTOCOL_RESPONSES:
		payload = responsesRequest{Model: input.remoteName, Input: testModelPrompt}
	case llmv1.Protocol_PROTOCOL_ANTHROPIC_MESSAGES:
		payload = anthropicRequest{
			Model:     input.remoteName,
			MaxTokens: 256,
			Messages: []anthropicMessage{
				{Role: "user", Content: testModelPrompt},
			},
		}
	default:
		return nil, nil, errors.New("unsupported model protocol")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode request: %w", err)
	}

	return body, headers, nil
}

func parseTestModelOutput(protocol llmv1.Protocol, responseBody []byte) (string, error) {
	switch protocol {
	case llmv1.Protocol_PROTOCOL_RESPONSES:
		return parseResponsesOutput(responseBody)
	case llmv1.Protocol_PROTOCOL_ANTHROPIC_MESSAGES:
		return parseAnthropicOutput(responseBody)
	default:
		return "", errors.New("unsupported model protocol")
	}
}

func parseResponsesOutput(responseBody []byte) (string, error) {
	var payload responsesResponse
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return "", fmt.Errorf("failed to parse responses output: %w", err)
	}

	for _, output := range payload.Output {
		for _, content := range output.Content {
			text := strings.TrimSpace(content.Text)
			if text != "" {
				return text, nil
			}
		}
	}

	return "", errors.New("response output text missing")
}

func parseAnthropicOutput(responseBody []byte) (string, error) {
	var payload anthropicResponse
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return "", fmt.Errorf("failed to parse anthropic output: %w", err)
	}

	for _, content := range payload.Content {
		text := strings.TrimSpace(content.Text)
		if text != "" {
			return text, nil
		}
	}

	return "", errors.New("response output text missing")
}

func formatUpstreamError(statusCode int, responseBody []byte) string {
	trimmed := strings.TrimSpace(string(responseBody))
	if trimmed == "" {
		return fmt.Sprintf("request failed with status %d", statusCode)
	}
	return fmt.Sprintf("request failed with status %d: %s", statusCode, trimmed)
}
