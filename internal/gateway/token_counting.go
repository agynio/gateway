package gateway

import (
	"context"

	"connectrpc.com/connect"
	tokencountingv1 "github.com/agynio/gateway/gen/agynio/api/token_counting/v1"
)

func (g *Gateway) CountTokens(ctx context.Context, req *connect.Request[tokencountingv1.CountTokensRequest]) (*connect.Response[tokencountingv1.CountTokensResponse], error) {
	resp, err := g.tokenCounting.CountTokens(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
