package gateway

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	exposev1 "github.com/agynio/gateway/gen/agynio/api/expose/v1"
	"github.com/agynio/gateway/internal/identity"
)

type ExposeGateway struct {
	expose exposev1.ExposeServiceClient
}

func NewExposeGateway(expose exposev1.ExposeServiceClient) *ExposeGateway {
	return &ExposeGateway{expose: expose}
}

func (g *ExposeGateway) AddExposure(ctx context.Context, req *connect.Request[exposev1.AddExposureRequest]) (*connect.Response[exposev1.AddExposureResponse], error) {
	if err := applyWorkloadID(ctx, &req.Msg.WorkloadId); err != nil {
		return nil, err
	}
	resp, err := g.expose.AddExposure(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ExposeGateway) RemoveExposure(ctx context.Context, req *connect.Request[exposev1.RemoveExposureRequest]) (*connect.Response[exposev1.RemoveExposureResponse], error) {
	if err := applyWorkloadID(ctx, &req.Msg.WorkloadId); err != nil {
		return nil, err
	}
	resp, err := g.expose.RemoveExposure(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ExposeGateway) ListExposures(ctx context.Context, req *connect.Request[exposev1.ListExposuresRequest]) (*connect.Response[exposev1.ListExposuresResponse], error) {
	if err := applyWorkloadID(ctx, &req.Msg.WorkloadId); err != nil {
		return nil, err
	}
	resp, err := g.expose.ListExposures(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func applyWorkloadID(ctx context.Context, workloadID *string) error {
	resolved, ok := identity.IdentityFromContext(ctx)
	if !ok || resolved.IdentityType != identity.IdentityTypeAgent {
		return nil
	}

	identityWorkloadID := strings.TrimSpace(resolved.WorkloadID)
	if identityWorkloadID == "" {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("workload id missing for agent identity"))
	}

	requestedWorkloadID := strings.TrimSpace(*workloadID)
	if requestedWorkloadID == "" {
		*workloadID = identityWorkloadID
		return nil
	}
	if requestedWorkloadID != identityWorkloadID {
		return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("workload id does not match identity"))
	}
	*workloadID = requestedWorkloadID
	return nil
}
