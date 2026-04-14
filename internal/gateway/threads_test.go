package gateway

import (
	"context"
	"testing"

	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type fakeAgentsClient struct {
	agentsv1.AgentsServiceClient
	getAgentReq   *agentsv1.GetAgentRequest
	getAgentResp  *agentsv1.GetAgentResponse
	getAgentErr   error
	getAgentCalls int
}

func (f *fakeAgentsClient) GetAgent(ctx context.Context, in *agentsv1.GetAgentRequest, opts ...grpc.CallOption) (*agentsv1.GetAgentResponse, error) {
	f.getAgentCalls++
	f.getAgentReq = in
	if f.getAgentErr != nil {
		return nil, f.getAgentErr
	}
	if f.getAgentResp == nil {
		f.getAgentResp = &agentsv1.GetAgentResponse{}
	}
	return f.getAgentResp, nil
}

func newThreadsGateway(agents agentsv1.AgentsServiceClient) *ThreadsGateway {
	gateway := New(agents, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	return NewThreadsGateway(gateway)
}

func TestWithOrganizationContextNoIdentity(t *testing.T) {
	client := &fakeAgentsClient{}
	gateway := newThreadsGateway(client)

	ctx := context.Background()
	gotCtx, err := gateway.withOrganizationContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCtx != ctx {
		t.Fatalf("expected context to be passed through")
	}
	if md, ok := metadata.FromOutgoingContext(gotCtx); ok {
		if values := md.Get(identity.MetadataKeyOrganizationID); len(values) > 0 {
			t.Fatalf("expected no organization metadata")
		}
	}
	if client.getAgentCalls != 0 {
		t.Fatalf("expected agent not to be called, got %d", client.getAgentCalls)
	}
}

func TestWithOrganizationContextNonAgentIdentity(t *testing.T) {
	client := &fakeAgentsClient{}
	gateway := newThreadsGateway(client)

	ctx := identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   "user-1",
		IdentityType: identity.IdentityTypeUser,
	})
	gotCtx, err := gateway.withOrganizationContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCtx != ctx {
		t.Fatalf("expected context to be passed through")
	}
	if md, ok := metadata.FromOutgoingContext(gotCtx); ok {
		if values := md.Get(identity.MetadataKeyOrganizationID); len(values) > 0 {
			t.Fatalf("expected no organization metadata")
		}
	}
	if client.getAgentCalls != 0 {
		t.Fatalf("expected agent not to be called, got %d", client.getAgentCalls)
	}
}

func TestWithOrganizationContextAgentIdentity(t *testing.T) {
	client := &fakeAgentsClient{
		getAgentResp: &agentsv1.GetAgentResponse{
			Agent: &agentsv1.Agent{
				OrganizationId: "org-1",
			},
		},
	}
	gateway := newThreadsGateway(client)
	ctx := identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   "agent-1",
		IdentityType: identity.IdentityTypeAgent,
	})

	gotCtx, err := gateway.withOrganizationContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.getAgentCalls != 1 {
		t.Fatalf("expected agent lookup once, got %d", client.getAgentCalls)
	}
	if client.getAgentReq == nil || client.getAgentReq.GetId() != "agent-1" {
		t.Fatalf("expected agent id to be forwarded")
	}
	md, ok := metadata.FromOutgoingContext(gotCtx)
	if !ok {
		t.Fatalf("expected outgoing metadata")
	}
	values := md.Get(identity.MetadataKeyOrganizationID)
	if len(values) != 1 || values[0] != "org-1" {
		t.Fatalf("expected organization metadata to be injected, got %#v", values)
	}
}

func TestWithOrganizationContextAgentMissingOrg(t *testing.T) {
	client := &fakeAgentsClient{
		getAgentResp: &agentsv1.GetAgentResponse{
			Agent: &agentsv1.Agent{},
		},
	}
	gateway := newThreadsGateway(client)
	ctx := identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   "agent-1",
		IdentityType: identity.IdentityTypeAgent,
	})

	_, err := gateway.withOrganizationContext(ctx)
	if err == nil {
		t.Fatalf("expected error")
	}
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected internal error, got %v", status.Code(err))
	}
	if client.getAgentCalls != 1 {
		t.Fatalf("expected agent lookup once, got %d", client.getAgentCalls)
	}
}
