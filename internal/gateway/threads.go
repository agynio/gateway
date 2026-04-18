package gateway

import (
	"context"

	"connectrpc.com/connect"
	threadsv1 "github.com/agynio/gateway/gen/agynio/api/threads/v1"
)

// ThreadsGateway forwards thread RPCs using the shared Gateway clients.
// It is a distinct type to avoid method name collisions with ChatGateway RPCs
// while keeping all clients on the central Gateway struct.
type ThreadsGateway Gateway

func NewThreadsGateway(gateway *Gateway) *ThreadsGateway {
	return (*ThreadsGateway)(gateway)
}

func (g *ThreadsGateway) CreateThread(ctx context.Context, req *connect.Request[threadsv1.CreateThreadRequest]) (*connect.Response[threadsv1.CreateThreadResponse], error) {
	resp, err := g.threads.CreateThread(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) ArchiveThread(ctx context.Context, req *connect.Request[threadsv1.ArchiveThreadRequest]) (*connect.Response[threadsv1.ArchiveThreadResponse], error) {
	resp, err := g.threads.ArchiveThread(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) AddParticipant(ctx context.Context, req *connect.Request[threadsv1.AddParticipantRequest]) (*connect.Response[threadsv1.AddParticipantResponse], error) {
	resp, err := g.threads.AddParticipant(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) SendMessage(ctx context.Context, req *connect.Request[threadsv1.SendMessageRequest]) (*connect.Response[threadsv1.SendMessageResponse], error) {
	resp, err := g.threads.SendMessage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) GetThreads(ctx context.Context, req *connect.Request[threadsv1.GetThreadsRequest]) (*connect.Response[threadsv1.GetThreadsResponse], error) {
	resp, err := g.threads.GetThreads(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) ListOrganizationThreads(ctx context.Context, req *connect.Request[threadsv1.ListOrganizationThreadsRequest]) (*connect.Response[threadsv1.ListOrganizationThreadsResponse], error) {
	resp, err := g.threads.ListOrganizationThreads(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) GetThread(ctx context.Context, req *connect.Request[threadsv1.GetThreadRequest]) (*connect.Response[threadsv1.GetThreadResponse], error) {
	resp, err := g.threads.GetThread(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) GetMessages(ctx context.Context, req *connect.Request[threadsv1.GetMessagesRequest]) (*connect.Response[threadsv1.GetMessagesResponse], error) {
	resp, err := g.threads.GetMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) GetUnackedMessages(ctx context.Context, req *connect.Request[threadsv1.GetUnackedMessagesRequest]) (*connect.Response[threadsv1.GetUnackedMessagesResponse], error) {
	resp, err := g.threads.GetUnackedMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ThreadsGateway) AckMessages(ctx context.Context, req *connect.Request[threadsv1.AckMessagesRequest]) (*connect.Response[threadsv1.AckMessagesResponse], error) {
	resp, err := g.threads.AckMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
