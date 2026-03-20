package gateway

import (
	"context"

	"connectrpc.com/connect"
	secretsv1 "github.com/agynio/gateway/gen/agynio/api/secrets/v1"
)

func (g *Gateway) ResolveSecret(ctx context.Context, req *connect.Request[secretsv1.ResolveSecretRequest]) (*connect.Response[secretsv1.ResolveSecretResponse], error) {
	resp, err := g.secrets.ResolveSecret(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
