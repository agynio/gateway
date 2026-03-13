package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agynio/gateway/internal/llmclient"
)

type stubLLMResponsesClient struct {
	t                         *testing.T
	createResponse            func(ctx context.Context, modelID string, body []byte) (llmclient.CreateResponseResult, error)
	createResponseStream      func(ctx context.Context, modelID string, body []byte) (llmclient.ResponseStream, error)
	createResponseCalls       int
	createResponseStreamCalls int
}

func (s *stubLLMResponsesClient) CreateResponse(ctx context.Context, modelID string, body []byte) (llmclient.CreateResponseResult, error) {
	s.createResponseCalls++
	if s.createResponse == nil {
		s.t.Fatalf("unexpected CreateResponse call")
	}
	return s.createResponse(ctx, modelID, body)
}

func (s *stubLLMResponsesClient) CreateResponseStream(ctx context.Context, modelID string, body []byte) (llmclient.ResponseStream, error) {
	s.createResponseStreamCalls++
	if s.createResponseStream == nil {
		s.t.Fatalf("unexpected CreateResponseStream call")
	}
	return s.createResponseStream(ctx, modelID, body)
}

type stubResponseStream struct {
	events []llmclient.StreamEvent
	idx    int
	closed bool
}

func (s *stubResponseStream) Recv() (llmclient.StreamEvent, error) {
	if s.idx >= len(s.events) {
		return llmclient.StreamEvent{}, io.EOF
	}
	event := s.events[s.idx]
	s.idx++
	return event, nil
}

func (s *stubResponseStream) Close() {
	s.closed = true
}

func TestParseResponsesBody(t *testing.T) {
	tests := []struct {
		name        string
		body        []byte
		wantModelID string
		wantStream  bool
		wantErr     string
		errContains bool
	}{
		{
			name:    "empty body",
			body:    []byte(""),
			wantErr: "request body is empty",
		},
		{
			name:        "invalid json",
			body:        []byte("{"),
			wantErr:     "invalid JSON",
			errContains: true,
		},
		{
			name:    "missing model id",
			body:    []byte(`{"stream":true}`),
			wantErr: "model_id is required",
		},
		{
			name:        "non streaming",
			body:        []byte(`{"model_id":"model-1"}`),
			wantModelID: "model-1",
			wantStream:  false,
		},
		{
			name:        "streaming true",
			body:        []byte(`{"model_id":"model-1","stream":true}`),
			wantModelID: "model-1",
			wantStream:  true,
		},
		{
			name:        "streaming false",
			body:        []byte(`{"model_id":"model-1","stream":false}`),
			wantModelID: "model-1",
			wantStream:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelID, stream, err := parseResponsesBody(tt.body)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error")
				}
				if tt.errContains {
					if !strings.Contains(err.Error(), tt.wantErr) {
						t.Fatalf("unexpected error: %v", err)
					}
					return
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("unexpected error: got %q want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if modelID != tt.wantModelID {
				t.Fatalf("unexpected modelID: got %q want %q", modelID, tt.wantModelID)
			}
			if stream != tt.wantStream {
				t.Fatalf("unexpected stream: got %v want %v", stream, tt.wantStream)
			}
		})
	}
}

func TestLLMResponsesHandlerUnarySuccess(t *testing.T) {
	body := []byte(`{"model_id":"model-1","input":"hi"}`)
	responseBody := []byte(`{"id":"resp"}`)
	stub := &stubLLMResponsesClient{t: t}
	stub.createResponse = func(ctx context.Context, modelID string, payload []byte) (llmclient.CreateResponseResult, error) {
		if modelID != "model-1" {
			t.Fatalf("unexpected modelID: %s", modelID)
		}
		if !bytes.Equal(payload, body) {
			t.Fatalf("unexpected body: %s", payload)
		}
		return llmclient.CreateResponseResult{Body: responseBody}, nil
	}

	handler := NewLLMResponsesHandler(stub)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/llm/v1/responses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(recorder, req)

	if stub.createResponseCalls != 1 {
		t.Fatalf("expected CreateResponse to be called once")
	}
	if stub.createResponseStreamCalls != 0 {
		t.Fatalf("expected CreateResponseStream not to be called")
	}
	resp := recorder.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected content type: %q", got)
	}
	if got := recorder.Body.Bytes(); !bytes.Equal(got, responseBody) {
		t.Fatalf("unexpected body: %s", got)
	}
}

func TestLLMResponsesHandlerUnaryError(t *testing.T) {
	body := []byte(`{"model_id":"model-1","input":"hi"}`)
	stub := &stubLLMResponsesClient{t: t}
	stub.createResponse = func(ctx context.Context, modelID string, payload []byte) (llmclient.CreateResponseResult, error) {
		if modelID != "model-1" {
			t.Fatalf("unexpected modelID: %s", modelID)
		}
		if !bytes.Equal(payload, body) {
			t.Fatalf("unexpected body: %s", payload)
		}
		return llmclient.CreateResponseResult{}, status.Error(codes.NotFound, "missing")
	}

	handler := NewLLMResponsesHandler(stub)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/llm/v1/responses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(recorder, req)

	if stub.createResponseCalls != 1 {
		t.Fatalf("expected CreateResponse to be called once")
	}
	assertProblemResponse(t, recorder, http.StatusNotFound, "missing")
}

func TestLLMResponsesHandlerStreamSuccess(t *testing.T) {
	body := []byte(`{"model_id":"model-1","stream":true}`)
	stream := &stubResponseStream{events: []llmclient.StreamEvent{
		{EventType: "message", Data: []byte(`{"chunk":"hi"}`)},
		{EventType: "", Data: []byte("[DONE]")},
	}}
	stub := &stubLLMResponsesClient{t: t}
	stub.createResponseStream = func(ctx context.Context, modelID string, payload []byte) (llmclient.ResponseStream, error) {
		if modelID != "model-1" {
			t.Fatalf("unexpected modelID: %s", modelID)
		}
		if !bytes.Equal(payload, body) {
			t.Fatalf("unexpected body: %s", payload)
		}
		return stream, nil
	}

	handler := NewLLMResponsesHandler(stub)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/llm/v1/responses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(recorder, req)

	if stub.createResponseStreamCalls != 1 {
		t.Fatalf("expected CreateResponseStream to be called once")
	}
	if stub.createResponseCalls != 0 {
		t.Fatalf("expected CreateResponse not to be called")
	}
	if !stream.closed {
		t.Fatalf("expected stream to be closed")
	}
	resp := recorder.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "text/event-stream" {
		t.Fatalf("unexpected content type: %q", got)
	}
	expected := "event: message\ndata: {\"chunk\":\"hi\"}\n\ndata: [DONE]\n\n"
	if got := recorder.Body.String(); got != expected {
		t.Fatalf("unexpected body: %q", got)
	}
}

func TestLLMResponsesHandlerStreamError(t *testing.T) {
	body := []byte(`{"model_id":"model-1","stream":true}`)
	stub := &stubLLMResponsesClient{t: t}
	stub.createResponseStream = func(ctx context.Context, modelID string, payload []byte) (llmclient.ResponseStream, error) {
		if modelID != "model-1" {
			t.Fatalf("unexpected modelID: %s", modelID)
		}
		if !bytes.Equal(payload, body) {
			t.Fatalf("unexpected body: %s", payload)
		}
		return nil, status.Error(codes.Unavailable, "down")
	}

	handler := NewLLMResponsesHandler(stub)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/llm/v1/responses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(recorder, req)

	if stub.createResponseStreamCalls != 1 {
		t.Fatalf("expected CreateResponseStream to be called once")
	}
	assertProblemResponse(t, recorder, http.StatusServiceUnavailable, "down")
}
