package gateway

import (
	"context"

	"connectrpc.com/connect"
	chatv1 "github.com/agynio/gateway/gen/agynio/api/chat/v1"
)

func (g *Gateway) CreateChat(ctx context.Context, req *connect.Request[chatv1.CreateChatRequest]) (*connect.Response[chatv1.CreateChatResponse], error) {
	resp, err := g.chat.CreateChat(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetChats(ctx context.Context, req *connect.Request[chatv1.GetChatsRequest]) (*connect.Response[chatv1.GetChatsResponse], error) {
	resp, err := g.chat.GetChats(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateChat(ctx context.Context, req *connect.Request[chatv1.UpdateChatRequest]) (*connect.Response[chatv1.UpdateChatResponse], error) {
	resp, err := g.chat.UpdateChat(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteChat(ctx context.Context, req *connect.Request[chatv1.DeleteChatRequest]) (*connect.Response[chatv1.DeleteChatResponse], error) {
	resp, err := g.chat.DeleteChat(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetMessages(ctx context.Context, req *connect.Request[chatv1.GetMessagesRequest]) (*connect.Response[chatv1.GetMessagesResponse], error) {
	resp, err := g.chat.GetMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) SendMessage(ctx context.Context, req *connect.Request[chatv1.SendMessageRequest]) (*connect.Response[chatv1.SendMessageResponse], error) {
	resp, err := g.chat.SendMessage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) MarkAsRead(ctx context.Context, req *connect.Request[chatv1.MarkAsReadRequest]) (*connect.Response[chatv1.MarkAsReadResponse], error) {
	resp, err := g.chat.MarkAsRead(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
