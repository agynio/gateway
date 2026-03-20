package gateway

import (
	"context"

	"connectrpc.com/connect"
	agentstatev1 "github.com/agynio/gateway/gen/agynio/api/agent_state/v1"
)

func (g *Gateway) AppendConversationMessages(ctx context.Context, req *connect.Request[agentstatev1.AppendConversationMessagesRequest]) (*connect.Response[agentstatev1.AppendConversationMessagesResponse], error) {
	resp, err := g.agentState.AppendConversationMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListConversationMessages(ctx context.Context, req *connect.Request[agentstatev1.ListConversationMessagesRequest]) (*connect.Response[agentstatev1.ListConversationMessagesResponse], error) {
	resp, err := g.agentState.ListConversationMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ReplaceConversationMessagesRange(ctx context.Context, req *connect.Request[agentstatev1.ReplaceConversationMessagesRangeRequest]) (*connect.Response[agentstatev1.ReplaceConversationMessagesRangeResponse], error) {
	resp, err := g.agentState.ReplaceConversationMessagesRange(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteConversationMessagesRange(ctx context.Context, req *connect.Request[agentstatev1.DeleteConversationMessagesRangeRequest]) (*connect.Response[agentstatev1.DeleteConversationMessagesRangeResponse], error) {
	resp, err := g.agentState.DeleteConversationMessagesRange(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetConversation(ctx context.Context, req *connect.Request[agentstatev1.GetConversationRequest]) (*connect.Response[agentstatev1.GetConversationResponse], error) {
	resp, err := g.agentState.GetConversation(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListConversationMessageIds(ctx context.Context, req *connect.Request[agentstatev1.ListConversationMessageIdsRequest]) (*connect.Response[agentstatev1.ListConversationMessageIdsResponse], error) {
	resp, err := g.agentState.ListConversationMessageIds(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListSnapshotMessages(ctx context.Context, req *connect.Request[agentstatev1.ListSnapshotMessagesRequest]) (*connect.Response[agentstatev1.ListSnapshotMessagesResponse], error) {
	resp, err := g.agentState.ListSnapshotMessages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListSnapshotMessageIds(ctx context.Context, req *connect.Request[agentstatev1.ListSnapshotMessageIdsRequest]) (*connect.Response[agentstatev1.ListSnapshotMessageIdsResponse], error) {
	resp, err := g.agentState.ListSnapshotMessageIds(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateSnapshot(ctx context.Context, req *connect.Request[agentstatev1.CreateSnapshotRequest]) (*connect.Response[agentstatev1.CreateSnapshotResponse], error) {
	resp, err := g.agentState.CreateSnapshot(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
