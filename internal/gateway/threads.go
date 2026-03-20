package gateway

import (
	"context"

	"connectrpc.com/connect"
	threadsv1 "github.com/agynio/gateway/gen/agynio/api/threads/v1"
)

func (g *Gateway) CreateThread(ctx context.Context, req *connect.Request[threadsv1.CreateThreadRequest]) (*connect.Response[threadsv1.CreateThreadResponse], error) {
	resp, err := g.threads.CreateThread(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ArchiveThread(ctx context.Context, req *connect.Request[threadsv1.ArchiveThreadRequest]) (*connect.Response[threadsv1.ArchiveThreadResponse], error) {
	resp, err := g.threads.ArchiveThread(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) AddParticipant(ctx context.Context, req *connect.Request[threadsv1.AddParticipantRequest]) (*connect.Response[threadsv1.AddParticipantResponse], error) {
	resp, err := g.threads.AddParticipant(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) SendMessage(ctx context.Context, req *connect.Request[threadsv1.SendMessageRequest]) (*connect.Response[threadsv1.SendMessageResponse], error) {
	resp, err := g.threads.SendMessage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetThreads(ctx context.Context, req *connect.Request[threadsv1.GetThreadsRequest]) (*connect.Response[threadsv1.GetThreadsResponse], error) {
	resp, err := g.threads.GetThreads(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetMessages(ctx context.Context, req *connect.Request[threadsv1.GetMessagesRequest]) (*connect.Response[threadsv1.GetMessagesResponse], error) {
	resp, err := g.threads.GetMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetUnackedMessages(ctx context.Context, req *connect.Request[threadsv1.GetUnackedMessagesRequest]) (*connect.Response[threadsv1.GetUnackedMessagesResponse], error) {
	resp, err := g.threads.GetUnackedMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) AckMessages(ctx context.Context, req *connect.Request[threadsv1.AckMessagesRequest]) (*connect.Response[threadsv1.AckMessagesResponse], error) {
	resp, err := g.threads.AckMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
