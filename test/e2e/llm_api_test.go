//go:build e2e

package e2e

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	gatewayv1connect "github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
	"github.com/stretchr/testify/require"
)

func TestLLMGateway_CreateResponse(t *testing.T) {
	modelID, body := llmRequestConfig()
	if modelID == "" || len(body) == 0 {
		t.Skip("llm request config not provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := gatewayv1connect.NewLLMGatewayClient(newClient(), gatewayURL)
	resp, err := client.CreateResponse(ctx, connect.NewRequest(&llmv1.CreateResponseRequest{ModelId: modelID, Body: body}))
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestLLMGateway_CreateResponseStream(t *testing.T) {
	modelID, body := llmRequestConfig()
	if modelID == "" || len(body) == 0 {
		t.Skip("llm request config not provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := gatewayv1connect.NewLLMGatewayClient(newClient(), gatewayURL)
	stream, err := client.CreateResponseStream(ctx, connect.NewRequest(&llmv1.CreateResponseStreamRequest{ModelId: modelID, Body: body}))
	require.NoError(t, err)

	for stream.Receive() {
		msg := stream.Msg()
		if msg != nil {
			return
		}
	}

	require.NoError(t, stream.Err())
}

func llmRequestConfig() (string, []byte) {
	modelID := strings.TrimSpace(os.Getenv("E2E_LLM_MODEL_ID"))
	if modelID == "" {
		modelID = strings.TrimSpace(os.Getenv("LLM_MODEL_ID"))
	}

	requestBody := strings.TrimSpace(os.Getenv("E2E_LLM_REQUEST_BODY"))
	if requestBody == "" {
		requestBody = strings.TrimSpace(os.Getenv("LLM_REQUEST_BODY"))
	}

	if modelID == "" || requestBody == "" {
		return "", nil
	}

	return modelID, []byte(requestBody)
}
