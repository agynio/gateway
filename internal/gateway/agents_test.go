package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	"github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	"google.golang.org/grpc"
)

type roleAgentsClient struct {
	fakeAgentsClient
	setAgentRoleReq     *agentsv1.SetAgentRoleRequest
	removeAgentRoleReq  *agentsv1.RemoveAgentRoleRequest
	listAgentRolesReq   *agentsv1.ListAgentRolesRequest
	listMyAgentRolesReq *agentsv1.ListMyAgentRolesRequest
}

func (c *roleAgentsClient) SetAgentRole(ctx context.Context, req *agentsv1.SetAgentRoleRequest, opts ...grpc.CallOption) (*agentsv1.SetAgentRoleResponse, error) {
	c.setAgentRoleReq = req
	return &agentsv1.SetAgentRoleResponse{
		Assignment: &agentsv1.AgentRoleAssignment{
			AgentId:    req.GetAgentId(),
			IdentityId: req.GetIdentityId(),
			Role:       req.GetRole(),
		},
	}, nil
}

func (c *roleAgentsClient) RemoveAgentRole(ctx context.Context, req *agentsv1.RemoveAgentRoleRequest, opts ...grpc.CallOption) (*agentsv1.RemoveAgentRoleResponse, error) {
	c.removeAgentRoleReq = req
	return &agentsv1.RemoveAgentRoleResponse{}, nil
}

func (c *roleAgentsClient) ListAgentRoles(ctx context.Context, req *agentsv1.ListAgentRolesRequest, opts ...grpc.CallOption) (*agentsv1.ListAgentRolesResponse, error) {
	c.listAgentRolesReq = req
	return &agentsv1.ListAgentRolesResponse{
		Assignments: []*agentsv1.AgentRoleAssignment{{
			AgentId:    req.GetAgentId(),
			IdentityId: "identity-1",
			Role:       agentsv1.AgentRole_AGENT_ROLE_OWNER,
		}},
	}, nil
}

func (c *roleAgentsClient) ListMyAgentRoles(ctx context.Context, req *agentsv1.ListMyAgentRolesRequest, opts ...grpc.CallOption) (*agentsv1.ListMyAgentRolesResponse, error) {
	c.listMyAgentRolesReq = req
	return &agentsv1.ListMyAgentRolesResponse{
		Assignments: []*agentsv1.AgentRoleAssignment{{
			AgentId:    "agent-1",
			IdentityId: "identity-1",
			Role:       agentsv1.AgentRole_AGENT_ROLE_PARTICIPANT,
		}},
	}, nil
}

func TestAgentsGatewayRoleRPCsAreRegisteredAndForwarded(t *testing.T) {
	agentsClient := &roleAgentsClient{}
	gateway := &Gateway{agents: agentsClient}
	server := newAgentsGatewayTestServer(gateway)
	t.Cleanup(server.Close)

	client := gatewayv1connect.NewAgentsGatewayClient(server.Client(), server.URL)

	setResp, err := client.SetAgentRole(context.Background(), connect.NewRequest(&agentsv1.SetAgentRoleRequest{
		AgentId:    "agent-1",
		IdentityId: "identity-1",
		Role:       agentsv1.AgentRole_AGENT_ROLE_MAINTAINER,
	}))
	if err != nil {
		t.Fatalf("SetAgentRole: %v", err)
	}
	if agentsClient.setAgentRoleReq == nil || agentsClient.setAgentRoleReq.GetAgentId() != "agent-1" {
		t.Fatalf("expected SetAgentRole request to be forwarded")
	}
	if setResp.Msg.GetAssignment().GetRole() != agentsv1.AgentRole_AGENT_ROLE_MAINTAINER {
		t.Fatalf("expected SetAgentRole response role to be forwarded")
	}

	_, err = client.RemoveAgentRole(context.Background(), connect.NewRequest(&agentsv1.RemoveAgentRoleRequest{
		AgentId:    "agent-2",
		IdentityId: "identity-2",
	}))
	if err != nil {
		t.Fatalf("RemoveAgentRole: %v", err)
	}
	if agentsClient.removeAgentRoleReq == nil || agentsClient.removeAgentRoleReq.GetIdentityId() != "identity-2" {
		t.Fatalf("expected RemoveAgentRole request to be forwarded")
	}

	listResp, err := client.ListAgentRoles(context.Background(), connect.NewRequest(&agentsv1.ListAgentRolesRequest{AgentId: "agent-3"}))
	if err != nil {
		t.Fatalf("ListAgentRoles: %v", err)
	}
	if agentsClient.listAgentRolesReq == nil || agentsClient.listAgentRolesReq.GetAgentId() != "agent-3" {
		t.Fatalf("expected ListAgentRoles request to be forwarded")
	}
	if got := listResp.Msg.GetAssignments()[0].GetAgentId(); got != "agent-3" {
		t.Fatalf("expected ListAgentRoles response agent id agent-3, got %q", got)
	}

	myRolesResp, err := client.ListMyAgentRoles(context.Background(), connect.NewRequest(&agentsv1.ListMyAgentRolesRequest{OrganizationId: "org-1"}))
	if err != nil {
		t.Fatalf("ListMyAgentRoles: %v", err)
	}
	if agentsClient.listMyAgentRolesReq == nil || agentsClient.listMyAgentRolesReq.GetOrganizationId() != "org-1" {
		t.Fatalf("expected ListMyAgentRoles request to be forwarded")
	}
	if got := myRolesResp.Msg.GetAssignments()[0].GetRole(); got != agentsv1.AgentRole_AGENT_ROLE_PARTICIPANT {
		t.Fatalf("expected ListMyAgentRoles response role participant, got %v", got)
	}
}

func newAgentsGatewayTestServer(gateway *Gateway) *httptest.Server {
	handlerPath, handler := gatewayv1connect.NewAgentsGatewayHandler(gateway)
	mux := http.NewServeMux()
	mux.Handle(handlerPath, handler)
	return httptest.NewServer(mux)
}
