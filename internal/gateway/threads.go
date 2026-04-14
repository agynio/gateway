package gateway

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	threadsv1 "github.com/agynio/gateway/gen/agynio/api/threads/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ThreadsGateway forwards thread RPCs using the shared Gateway clients.
// It is a distinct type to avoid method name collisions with ChatGateway RPCs
// while keeping all clients on the central Gateway struct.
type ThreadsGateway Gateway

func NewThreadsGateway(gateway *Gateway) *ThreadsGateway {
	return (*ThreadsGateway)(gateway)
}

func (g *ThreadsGateway) CreateThread(ctx context.Context, req *connect.Request[threadsv1.CreateThreadRequest]) (*connect.Response[threadsv1.CreateThreadResponse], error) {
	ctx, err := g.withOrganizationContext(ctx)
	if err != nil {
		return nil, toConnectError(err)
	}
	resp, err := g.threads.CreateThread(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) ArchiveThread(ctx context.Context, req *connect.Request[threadsv1.ArchiveThreadRequest]) (*connect.Response[threadsv1.ArchiveThreadResponse], error) {
	ctx, err := g.withOrganizationContext(ctx)
	if err != nil {
		return nil, toConnectError(err)
	}
	resp, err := g.threads.ArchiveThread(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) AddParticipant(ctx context.Context, req *connect.Request[threadsv1.AddParticipantRequest]) (*connect.Response[threadsv1.AddParticipantResponse], error) {
	ctx, err := g.withOrganizationContext(ctx)
	if err != nil {
		return nil, toConnectError(err)
	}
	resp, err := g.threads.AddParticipant(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) SendMessage(ctx context.Context, req *connect.Request[threadsv1.SendMessageRequest]) (*connect.Response[threadsv1.SendMessageResponse], error) {
	ctx, err := g.withOrganizationContext(ctx)
	if err != nil {
		return nil, toConnectError(err)
	}
	resp, err := g.threads.SendMessage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) GetThreads(ctx context.Context, req *connect.Request[threadsv1.GetThreadsRequest]) (*connect.Response[threadsv1.GetThreadsResponse], error) {
	ctx, err := g.withOrganizationContext(ctx)
	if err != nil {
		return nil, toConnectError(err)
	}
	resp, err := g.threads.GetThreads(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) GetMessages(ctx context.Context, req *connect.Request[threadsv1.GetMessagesRequest]) (*connect.Response[threadsv1.GetMessagesResponse], error) {
	ctx, err := g.withOrganizationContext(ctx)
	if err != nil {
		return nil, toConnectError(err)
	}
	resp, err := g.threads.GetMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) GetUnackedMessages(ctx context.Context, req *connect.Request[threadsv1.GetUnackedMessagesRequest]) (*connect.Response[threadsv1.GetUnackedMessagesResponse], error) {
	ctx, err := g.withOrganizationContext(ctx)
	if err != nil {
		return nil, toConnectError(err)
	}
	resp, err := g.threads.GetUnackedMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) AckMessages(ctx context.Context, req *connect.Request[threadsv1.AckMessagesRequest]) (*connect.Response[threadsv1.AckMessagesResponse], error) {
	ctx, err := g.withOrganizationContext(ctx)
	if err != nil {
		return nil, toConnectError(err)
	}
	resp, err := g.threads.AckMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) withOrganizationContext(ctx context.Context) (context.Context, error) {
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok {
		return ctx, nil
	}
	if resolved.IdentityType != identity.IdentityTypeAgent {
		return ctx, nil
	}
	identityID := strings.TrimSpace(resolved.IdentityID)
	if identityID == "" {
		return ctx, status.Error(codes.Unauthenticated, "identity id missing")
	}
	resp, err := g.agents.GetAgent(ctx, &agentsv1.GetAgentRequest{Id: identityID})
	if err != nil {
		return ctx, err
	}
	agent := resp.GetAgent()
	if agent == nil {
		return ctx, status.Error(codes.Internal, "agent response missing")
	}
	organizationID := strings.TrimSpace(agent.GetOrganizationId())
	if organizationID == "" {
		return ctx, status.Error(codes.Internal, "agent organization_id missing")
	}
	return metadata.AppendToOutgoingContext(ctx, identity.MetadataKeyOrganizationID, organizationID), nil
}
