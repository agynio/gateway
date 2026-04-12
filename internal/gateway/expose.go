package gateway

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	exposev1 "github.com/agynio/gateway/gen/agynio/api/expose/v1"
	gatewayv1 "github.com/agynio/gateway/gen/agynio/api/gateway/v1"
	"github.com/agynio/gateway/internal/identity"
)

type ExposeGateway struct {
	expose exposev1.ExposeServiceClient
}

func NewExposeGateway(expose exposev1.ExposeServiceClient) *ExposeGateway {
	return &ExposeGateway{expose: expose}
}

func (g *ExposeGateway) AddExposure(ctx context.Context, req *connect.Request[exposev1.AddExposureRequest]) (*connect.Response[exposev1.AddExposureResponse], error) {
	resp, err := g.expose.AddExposure(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ExposeGateway) RemoveExposure(ctx context.Context, req *connect.Request[exposev1.RemoveExposureRequest]) (*connect.Response[exposev1.RemoveExposureResponse], error) {
	resp, err := g.expose.RemoveExposure(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ExposeGateway) ListExposures(ctx context.Context, req *connect.Request[exposev1.ListExposuresRequest]) (*connect.Response[exposev1.ListExposuresResponse], error) {
	resp, err := g.expose.ListExposures(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ExposeGateway) AddExposureForCaller(ctx context.Context, req *connect.Request[gatewayv1.AddExposureForCallerRequest]) (*connect.Response[gatewayv1.AddExposureForCallerResponse], error) {
	resolved, err := requireAgentIdentity(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := g.expose.AddExposure(ctx, &exposev1.AddExposureRequest{
		WorkloadId: resolved.WorkloadID,
		AgentId:    resolved.IdentityID,
		Port:       req.Msg.GetPort(),
	})
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(&gatewayv1.AddExposureForCallerResponse{Exposure: resp.GetExposure()}), nil
}

func (g *ExposeGateway) RemoveExposureForCaller(ctx context.Context, req *connect.Request[gatewayv1.RemoveExposureForCallerRequest]) (*connect.Response[gatewayv1.RemoveExposureForCallerResponse], error) {
	resolved, err := requireAgentIdentity(ctx)
	if err != nil {
		return nil, err
	}
	_, err = g.expose.RemoveExposure(ctx, &exposev1.RemoveExposureRequest{
		WorkloadId: resolved.WorkloadID,
		Port:       req.Msg.GetPort(),
	})
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(&gatewayv1.RemoveExposureForCallerResponse{}), nil
}

func (g *ExposeGateway) ListExposuresForCaller(ctx context.Context, req *connect.Request[gatewayv1.ListExposuresForCallerRequest]) (*connect.Response[gatewayv1.ListExposuresForCallerResponse], error) {
	resolved, err := requireAgentIdentity(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := g.expose.ListExposures(ctx, &exposev1.ListExposuresRequest{
		WorkloadId: resolved.WorkloadID,
		PageSize:   req.Msg.GetPageSize(),
		PageToken:  req.Msg.GetPageToken(),
	})
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(&gatewayv1.ListExposuresForCallerResponse{
		Exposures:     resp.GetExposures(),
		NextPageToken: resp.GetNextPageToken(),
	}), nil
}

func requireAgentIdentity(ctx context.Context) (identity.ResolvedIdentity, error) {
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok {
		return identity.ResolvedIdentity{}, connect.NewError(connect.CodeUnauthenticated, errors.New("identity not available"))
	}
	if resolved.IdentityType != identity.IdentityTypeAgent {
		return identity.ResolvedIdentity{}, connect.NewError(connect.CodePermissionDenied, errors.New("identity is not an agent"))
	}
	identityID := strings.TrimSpace(resolved.IdentityID)
	if identityID == "" {
		return identity.ResolvedIdentity{}, connect.NewError(connect.CodeInternal, errors.New("agent id missing for agent identity"))
	}
	workloadID := strings.TrimSpace(resolved.WorkloadID)
	if workloadID == "" {
		return identity.ResolvedIdentity{}, connect.NewError(connect.CodeInternal, errors.New("workload id missing for agent identity"))
	}
	resolved.IdentityID = identityID
	resolved.WorkloadID = workloadID
	return resolved, nil
}
