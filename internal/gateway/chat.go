package gateway

import (
	"context"

	"connectrpc.com/connect"
	chatv1 "github.com/agynio/gateway/gen/agynio/api/chat/v1"
)

type ChatGateway struct {
	client chatv1.ChatServiceClient
}

func NewChatGateway(client chatv1.ChatServiceClient) *ChatGateway {
	if client == nil {
		panic("chat client is required")
	}
	return &ChatGateway{client: client}
}

func (g *ChatGateway) CreateChat(ctx context.Context, req *connect.Request[chatv1.CreateChatRequest]) (*connect.Response[chatv1.CreateChatResponse], error) {
	resp, err := g.client.CreateChat(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ChatGateway) GetChats(ctx context.Context, req *connect.Request[chatv1.GetChatsRequest]) (*connect.Response[chatv1.GetChatsResponse], error) {
	resp, err := g.client.GetChats(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ChatGateway) GetMessages(ctx context.Context, req *connect.Request[chatv1.GetMessagesRequest]) (*connect.Response[chatv1.GetMessagesResponse], error) {
	resp, err := g.client.GetMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ChatGateway) SendMessage(ctx context.Context, req *connect.Request[chatv1.SendMessageRequest]) (*connect.Response[chatv1.SendMessageResponse], error) {
	resp, err := g.client.SendMessage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ChatGateway) MarkAsRead(ctx context.Context, req *connect.Request[chatv1.MarkAsReadRequest]) (*connect.Response[chatv1.MarkAsReadResponse], error) {
	resp, err := g.client.MarkAsRead(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
