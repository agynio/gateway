package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/agynio/gateway/internal/llmclient"
)

const maxResponsesBodySize = 1 << 20 // 1 MiB

// LLMResponsesClient defines the interface needed by the responses handler.
type LLMResponsesClient interface {
	CreateResponse(ctx context.Context, modelID string, body []byte) (llmclient.CreateResponseResult, error)
	CreateResponseStream(ctx context.Context, modelID string, body []byte) (llmclient.ResponseStream, error)
}

// LLMResponsesHandler handles POST /llm/v1/responses by calling the LLM gRPC
// service. Streaming requests receive SSE; non-streaming requests receive JSON.
// This handler is mounted as a raw http.Handler because SSE streaming is
// incompatible with the oapi-codegen strict server return-value model.
type LLMResponsesHandler struct {
	client LLMResponsesClient
}

func NewLLMResponsesHandler(client LLMResponsesClient) *LLMResponsesHandler {
	return &LLMResponsesHandler{client: client}
}

func (h *LLMResponsesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxResponsesBodySize)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		problem := NewProblem(http.StatusBadRequest, "Bad Request", "failed to read request body")
		WriteProblem(w, problem)
		return
	}

	modelID, stream, err := parseResponsesBody(body)
	if err != nil {
		problem := NewProblem(http.StatusBadRequest, "Bad Request", err.Error())
		WriteProblem(w, problem)
		return
	}

	if stream {
		h.handleStream(w, r, modelID, body)
	} else {
		h.handleUnary(w, r, modelID, body)
	}
}

func (h *LLMResponsesHandler) handleUnary(w http.ResponseWriter, r *http.Request, modelID string, body []byte) {
	result, err := h.client.CreateResponse(r.Context(), modelID, body)
	if err != nil {
		handleGRPCError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result.Body)
}

func (h *LLMResponsesHandler) handleStream(w http.ResponseWriter, r *http.Request, modelID string, body []byte) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		problem := NewProblem(http.StatusInternalServerError, "Internal Server Error", "streaming not supported")
		WriteProblem(w, problem)
		return
	}

	stream, err := h.client.CreateResponseStream(r.Context(), modelID, body)
	if err != nil {
		handleGRPCError(w, err)
		return
	}
	defer stream.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	for {
		event, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			// Stream already started — can't change status code.
			// Log and close.
			log.Printf("llm response stream error: %v", err)
			return
		}
		if err := writeSSEEvent(w, event); err != nil {
			log.Printf("llm response stream write error: %v", err)
			return
		}
		flusher.Flush()
	}
}

func writeSSEEvent(w io.Writer, event llmclient.StreamEvent) error {
	if event.EventType != "" {
		if _, err := fmt.Fprintf(w, "event: %s\n", event.EventType); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, "data: %s\n\n", event.Data); err != nil {
		return err
	}
	return nil
}

// parseResponsesBody extracts model_id and stream flag from the request body.
func parseResponsesBody(body []byte) (modelID string, stream bool, err error) {
	if len(body) == 0 {
		return "", false, fmt.Errorf("request body is empty")
	}
	var payload struct {
		ModelID string `json:"model_id"`
		Stream  *bool  `json:"stream,omitempty"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", false, fmt.Errorf("invalid JSON: %w", err)
	}
	if payload.ModelID == "" {
		return "", false, fmt.Errorf("model_id is required")
	}
	if payload.Stream != nil {
		stream = *payload.Stream
	}
	return payload.ModelID, stream, nil
}

// handleGRPCError converts a gRPC error to a Problem response.
func handleGRPCError(w http.ResponseWriter, err error) {
	problemErr := grpcErrorToProblem(err)
	WriteProblem(w, problemErr.Problem)
}
