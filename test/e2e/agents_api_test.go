//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"connectrpc.com/connect"
	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	gatewayv1connect "github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentsGateway_ListAgents(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := gatewayv1connect.NewAgentsGatewayClient(newClient(), gatewayURL)
	resp, err := client.ListAgents(ctx, connect.NewRequest(&agentsv1.ListAgentsRequest{}))
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Msg.Agents)
}

func TestAgentsGateway_CreateAndDeleteAgent(t *testing.T) {
	modelID := agentModelID
	require.NotEmpty(t, modelID)

	orgID := organizationID
	require.NotEmpty(t, orgID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := gatewayv1connect.NewAgentsGatewayClient(newClient(), gatewayURL)
	createReq := &agentsv1.CreateAgentRequest{
		Name:           fmt.Sprintf("e2e-smoke-agent-%d", time.Now().UnixNano()),
		Role:           "assistant",
		Model:          modelID,
		Configuration:  "{}",
		Image:          "agent-image",
		InitImage:      "agent-init-image",
		OrganizationId: orgID,
	}
	createResp, err := client.CreateAgent(ctx, connect.NewRequest(createReq))
	require.NoError(t, err)
	require.NotNil(t, createResp.Msg)
	require.NotNil(t, createResp.Msg.Agent)
	require.NotNil(t, createResp.Msg.Agent.Meta)

	agentID := createResp.Msg.Agent.Meta.Id
	require.NotEmpty(t, agentID)

	_, err = client.DeleteAgent(ctx, connect.NewRequest(&agentsv1.DeleteAgentRequest{Id: agentID}))
	require.NoError(t, err)
}

func TestAgentsGateway_ListMcps(t *testing.T) {
	modelID := agentModelID
	require.NotEmpty(t, modelID)

	orgID := organizationID
	require.NotEmpty(t, orgID)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := gatewayv1connect.NewAgentsGatewayClient(newClient(), gatewayURL)
	createResp, err := client.CreateAgent(ctx, connect.NewRequest(&agentsv1.CreateAgentRequest{
		Name:           fmt.Sprintf("e2e-mcp-agent-%d", time.Now().UnixNano()),
		Role:           "assistant",
		Model:          modelID,
		Configuration:  "{}",
		Image:          "agent-image",
		InitImage:      "agent-init-image",
		OrganizationId: orgID,
	}))
	require.NoError(t, err)
	require.NotNil(t, createResp.Msg)
	require.NotNil(t, createResp.Msg.Agent)
	require.NotNil(t, createResp.Msg.Agent.Meta)

	agentID := createResp.Msg.Agent.Meta.Id
	require.NotEmpty(t, agentID)

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cleanupCancel()
		_, _ = client.DeleteAgent(cleanupCtx, connect.NewRequest(&agentsv1.DeleteAgentRequest{Id: agentID}))
	})

	resp, err := client.ListMcps(ctx, connect.NewRequest(&agentsv1.ListMcpsRequest{AgentId: agentID}))
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Msg)
}

func TestAgentsGateway_InvalidPayloadReturnsClientError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := gatewayv1connect.NewAgentsGatewayClient(newClient(), gatewayURL)
	_, err := client.GetAgent(ctx, connect.NewRequest(&agentsv1.GetAgentRequest{}))
	require.Error(t, err)
}
