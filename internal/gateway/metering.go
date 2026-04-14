package gateway

import (
	"context"

	"connectrpc.com/connect"
	meteringv1 "github.com/agynio/gateway/gen/agynio/api/metering/v1"
)

type MeteringGateway struct {
	metering meteringv1.MeteringServiceClient
}

func NewMeteringGateway(metering meteringv1.MeteringServiceClient) *MeteringGateway {
	return &MeteringGateway{metering: metering}
}

func (g *MeteringGateway) QueryUsage(ctx context.Context, req *connect.Request[meteringv1.QueryUsageRequest]) (*connect.Response[meteringv1.QueryUsageResponse], error) {
	resp, err := g.metering.QueryUsage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
