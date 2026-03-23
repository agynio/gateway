package gateway

import (
	"context"

	"connectrpc.com/connect"
	tracingv1 "github.com/agynio/gateway/gen/agynio/api/tracing/v1"
)

func (g *Gateway) ListSpans(ctx context.Context, req *connect.Request[tracingv1.ListSpansRequest]) (*connect.Response[tracingv1.ListSpansResponse], error) {
	resp, err := g.tracing.ListSpans(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetSpan(ctx context.Context, req *connect.Request[tracingv1.GetSpanRequest]) (*connect.Response[tracingv1.GetSpanResponse], error) {
	resp, err := g.tracing.GetSpan(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetTrace(ctx context.Context, req *connect.Request[tracingv1.GetTraceRequest]) (*connect.Response[tracingv1.GetTraceResponse], error) {
	resp, err := g.tracing.GetTrace(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
