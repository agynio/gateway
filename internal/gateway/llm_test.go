package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
)

func TestParseTestModelInputSuccess(t *testing.T) {
	input, err := parseTestModelInput(&llmv1.ResolveModelResponse{
		Protocol: llmv1.Protocol_PROTOCOL_RESPONSES,
	}, "model-id", "https://proxy.example.com", "agyn_proxy_token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.endpoint != "https://proxy.example.com/v1/responses" {
		t.Fatalf("expected endpoint to be set")
	}
	if input.modelID != "model-id" {
		t.Fatalf("expected model id to be set")
	}
	if input.protocol != llmv1.Protocol_PROTOCOL_RESPONSES {
		t.Fatalf("expected protocol to be set")
	}
	if input.authToken != "agyn_proxy_token" {
		t.Fatalf("expected auth token to be set")
	}
}

func TestParseTestModelInputErrors(t *testing.T) {
	base := func() (*llmv1.ResolveModelResponse, string, string, string) {
		return &llmv1.ResolveModelResponse{
			Protocol: llmv1.Protocol_PROTOCOL_RESPONSES,
		}, "model-id", "https://proxy.example.com", "agyn_proxy_token"
	}
	for _, tt := range []struct {
		name    string
		mutate  func(*llmv1.ResolveModelResponse, *string, *string, *string)
		code    connect.Code
		message string
	}{
		{
			name: "missing model id",
			mutate: func(_ *llmv1.ResolveModelResponse, modelID *string, _ *string, _ *string) {
				*modelID = ""
			},
			code:    connect.CodeInvalidArgument,
			message: "model_id is required",
		},
		{
			name: "unsupported protocol",
			mutate: func(resp *llmv1.ResolveModelResponse, _ *string, _ *string, _ *string) {
				resp.Protocol = llmv1.Protocol_PROTOCOL_UNSPECIFIED
			},
			code:    connect.CodeFailedPrecondition,
			message: "unsupported model protocol",
		},
		{
			name: "missing proxy url",
			mutate: func(_ *llmv1.ResolveModelResponse, _ *string, proxyURL *string, _ *string) {
				*proxyURL = ""
			},
			code:    connect.CodeFailedPrecondition,
			message: "llm proxy url missing",
		},
		{
			name: "invalid proxy url",
			mutate: func(_ *llmv1.ResolveModelResponse, _ *string, proxyURL *string, _ *string) {
				*proxyURL = "http://"
			},
			code:    connect.CodeFailedPrecondition,
			message: "invalid llm proxy url",
		},
		{
			name: "missing token",
			mutate: func(_ *llmv1.ResolveModelResponse, _ *string, _ *string, token *string) {
				*token = ""
			},
			code:    connect.CodeFailedPrecondition,
			message: "llm proxy token missing",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			resolved, modelID, proxyURL, token := base()
			tt.mutate(resolved, &modelID, &proxyURL, &token)
			_, err := parseTestModelInput(resolved, modelID, proxyURL, token)
			if err == nil {
				t.Fatalf("expected error")
			}
			var testErr *testModelError
			if !errors.As(err, &testErr) {
				t.Fatalf("expected testModelError, got %T", err)
			}
			if testErr.code != tt.code {
				t.Fatalf("expected code %v, got %v", tt.code, testErr.code)
			}
			if testErr.message != tt.message {
				t.Fatalf("expected message %q, got %q", tt.message, testErr.message)
			}
		})
	}
}

func TestBuildTestModelRequestResponses(t *testing.T) {
	body, headers, err := buildTestModelRequest(testModelInput{
		endpoint:  "https://example.com",
		modelID:   "model-id",
		authToken: "agyn_proxy_token",
		protocol:  llmv1.Protocol_PROTOCOL_RESPONSES,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := headers.Get("Authorization"); got != "Bearer agyn_proxy_token" {
		t.Fatalf("expected bearer auth header, got %q", got)
	}
	if got := headers.Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type header, got %q", got)
	}
	if got := headers.Get("x-api-key"); got != "" {
		t.Fatalf("unexpected x-api-key header: %q", got)
	}
	if got := headers.Get("anthropic-version"); got != "" {
		t.Fatalf("unexpected anthropic header: %q", got)
	}
	var payload responsesRequest
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to parse payload: %v", err)
	}
	if payload.Model != "model-id" {
		t.Fatalf("expected model name to match")
	}
	if payload.Input != testModelPrompt {
		t.Fatalf("expected prompt to match")
	}
}

func TestBuildTestModelRequestAnthropic(t *testing.T) {
	body, headers, err := buildTestModelRequest(testModelInput{
		endpoint:  "https://example.com",
		modelID:   "model-id",
		authToken: "agyn_proxy_token",
		protocol:  llmv1.Protocol_PROTOCOL_ANTHROPIC_MESSAGES,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := headers.Get("Authorization"); got != "Bearer agyn_proxy_token" {
		t.Fatalf("expected bearer auth header, got %q", got)
	}
	if got := headers.Get("x-api-key"); got != "" {
		t.Fatalf("unexpected x-api-key header, got %q", got)
	}
	if got := headers.Get("anthropic-version"); got != testModelAnthropicHeader {
		t.Fatalf("expected anthropic version header, got %q", got)
	}
	var payload anthropicRequest
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to parse payload: %v", err)
	}
	if payload.Model != "model-id" {
		t.Fatalf("expected model name to match")
	}
	if payload.MaxTokens != testModelAnthropicMaxTokens {
		t.Fatalf("expected max tokens to match")
	}
	if len(payload.Messages) != 1 {
		t.Fatalf("expected one message")
	}
	if payload.Messages[0].Role != "user" || payload.Messages[0].Content != testModelPrompt {
		t.Fatalf("expected user prompt message")
	}
}

func TestParseResponsesOutput(t *testing.T) {
	body, err := json.Marshal(responsesResponse{
		Output: []responsesOutput{
			{Content: []responsesContent{{Text: ""}}},
			{Content: []responsesContent{{Text: "  response  "}}},
		},
	})
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}
	text, err := parseResponsesOutput(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "response" {
		t.Fatalf("expected trimmed response text, got %q", text)
	}
}

func TestParseAnthropicOutput(t *testing.T) {
	body, err := json.Marshal(anthropicResponse{
		Content: []anthropicContent{
			{Text: ""},
			{Text: "  response  "},
		},
	})
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}
	text, err := parseAnthropicOutput(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "response" {
		t.Fatalf("expected trimmed response text, got %q", text)
	}
}

func TestFormatUpstreamErrorTruncates(t *testing.T) {
	body := strings.Repeat("a", testModelMaxErrorBodyBytes+10)
	message := formatUpstreamError(http.StatusBadGateway, []byte(body))
	prefix := fmt.Sprintf("request failed with status %d: ", http.StatusBadGateway)
	if !strings.HasPrefix(message, prefix) {
		t.Fatalf("expected prefix %q", prefix)
	}
	trimmed := strings.TrimPrefix(message, prefix)
	if !strings.HasSuffix(trimmed, "...") {
		t.Fatalf("expected truncated message suffix")
	}
	if len(trimmed) != testModelMaxErrorBodyBytes+3 {
		t.Fatalf("expected trimmed length %d, got %d", testModelMaxErrorBodyBytes+3, len(trimmed))
	}
}

func TestTestModelSuccessResponses(t *testing.T) {
	errCh := make(chan error, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			select {
			case errCh <- fmt.Errorf("expected POST, got %s", r.Method):
			default:
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/v1/responses" {
			select {
			case errCh <- fmt.Errorf("expected path /v1/responses, got %s", r.URL.Path):
			default:
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer agyn_proxy_token" {
			select {
			case errCh <- fmt.Errorf("expected bearer auth, got %q", got):
			default:
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			select {
			case errCh <- fmt.Errorf("expected content type, got %q", got):
			default:
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			select {
			case errCh <- fmt.Errorf("failed to read body: %v", err):
			default:
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var payload responsesRequest
		if err := json.Unmarshal(body, &payload); err != nil {
			select {
			case errCh <- fmt.Errorf("invalid payload: %v", err):
			default:
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if payload.Model != "model-id" || payload.Input != testModelPrompt {
			select {
			case errCh <- fmt.Errorf("unexpected payload: %#v", payload):
			default:
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(responsesResponse{
			Output: []responsesOutput{{Content: []responsesContent{{Text: "ok"}}}},
		}); err != nil {
			select {
			case errCh <- fmt.Errorf("failed to write response: %v", err):
			default:
			}
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	output, err := testModel(ctx, testModelInput{
		endpoint:  server.URL + "/v1/responses",
		modelID:   "model-id",
		authToken: "agyn_proxy_token",
		protocol:  llmv1.Protocol_PROTOCOL_RESPONSES,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "ok" {
		t.Fatalf("expected output to match, got %q", output)
	}
	select {
	case err := <-errCh:
		t.Fatalf("server validation failed: %v", err)
	default:
	}
}
