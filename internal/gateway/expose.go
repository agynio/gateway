package gateway

import (
	"context"

	"connectrpc.com/connect"
	exposev1 "github.com/agynio/gateway/gen/agynio/api/expose/v1"
)

type ExposeGateway struct {
	expose exposev1.ExposeServiceClient
}

func NewExposeGateway(expose exposev1.ExposeServiceClient) *ExposeGateway {
	return &ExposeGateway{expose: expose}
}

func (g *ExposeGateway) AddExposure(ctx context.Context, req *connect.Request[exposev1.AddExposureRequest]) (*connect.Response[exposev1.AddExposureResponse], error) {
	resp, err := g.expose.AddExposure(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ExposeGateway) RemoveExposure(ctx context.Context, req *connect.Request[exposev1.RemoveExposureRequest]) (*connect.Response[exposev1.RemoveExposureResponse], error) {
	resp, err := g.expose.RemoveExposure(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ExposeGateway) ListExposures(ctx context.Context, req *connect.Request[exposev1.ListExposuresRequest]) (*connect.Response[exposev1.ListExposuresResponse], error) {
	resp, err := g.expose.ListExposures(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
