package gateway

import (
	"context"

	"connectrpc.com/connect"
	appsv1 "github.com/agynio/gateway/gen/agynio/api/apps/v1"
)

func (g *Gateway) RegisterApp(ctx context.Context, req *connect.Request[appsv1.RegisterAppRequest]) (*connect.Response[appsv1.RegisterAppResponse], error) {
	resp, err := g.apps.RegisterApp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) EnrollApp(ctx context.Context, req *connect.Request[appsv1.EnrollAppRequest]) (*connect.Response[appsv1.EnrollAppResponse], error) {
	resp, err := g.apps.EnrollApp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetApp(ctx context.Context, req *connect.Request[appsv1.GetAppRequest]) (*connect.Response[appsv1.GetAppResponse], error) {
	resp, err := g.apps.GetApp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListApps(ctx context.Context, req *connect.Request[appsv1.ListAppsRequest]) (*connect.Response[appsv1.ListAppsResponse], error) {
	resp, err := g.apps.ListApps(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteApp(ctx context.Context, req *connect.Request[appsv1.DeleteAppRequest]) (*connect.Response[appsv1.DeleteAppResponse], error) {
	resp, err := g.apps.DeleteApp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
