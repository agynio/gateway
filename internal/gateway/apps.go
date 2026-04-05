package gateway

import (
	"context"

	"connectrpc.com/connect"
	appsv1 "github.com/agynio/gateway/gen/agynio/api/apps/v1"
)

func (g *Gateway) CreateApp(ctx context.Context, req *connect.Request[appsv1.CreateAppRequest]) (*connect.Response[appsv1.CreateAppResponse], error) {
	resp, err := g.apps.CreateApp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateApp(ctx context.Context, req *connect.Request[appsv1.UpdateAppRequest]) (*connect.Response[appsv1.UpdateAppResponse], error) {
	resp, err := g.apps.UpdateApp(ctx, req.Msg)
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

func (g *Gateway) GetAppBySlug(ctx context.Context, req *connect.Request[appsv1.GetAppBySlugRequest]) (*connect.Response[appsv1.GetAppBySlugResponse], error) {
	resp, err := g.apps.GetAppBySlug(ctx, req.Msg)
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

func (g *Gateway) EnrollApp(ctx context.Context, req *connect.Request[appsv1.EnrollAppRequest]) (*connect.Response[appsv1.EnrollAppResponse], error) {
	resp, err := g.apps.EnrollApp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) InstallApp(ctx context.Context, req *connect.Request[appsv1.InstallAppRequest]) (*connect.Response[appsv1.InstallAppResponse], error) {
	resp, err := g.apps.InstallApp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetInstallation(ctx context.Context, req *connect.Request[appsv1.GetInstallationRequest]) (*connect.Response[appsv1.GetInstallationResponse], error) {
	resp, err := g.apps.GetInstallation(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetInstallationBySlug(ctx context.Context, req *connect.Request[appsv1.GetInstallationBySlugRequest]) (*connect.Response[appsv1.GetInstallationBySlugResponse], error) {
	resp, err := g.apps.GetInstallationBySlug(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListInstallations(ctx context.Context, req *connect.Request[appsv1.ListInstallationsRequest]) (*connect.Response[appsv1.ListInstallationsResponse], error) {
	resp, err := g.apps.ListInstallations(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateInstallation(ctx context.Context, req *connect.Request[appsv1.UpdateInstallationRequest]) (*connect.Response[appsv1.UpdateInstallationResponse], error) {
	resp, err := g.apps.UpdateInstallation(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UninstallApp(ctx context.Context, req *connect.Request[appsv1.UninstallAppRequest]) (*connect.Response[appsv1.UninstallAppResponse], error) {
	resp, err := g.apps.UninstallApp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
