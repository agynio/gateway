package gateway

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	exposev1 "github.com/agynio/gateway/gen/agynio/api/expose/v1"
	"github.com/agynio/gateway/internal/clusteradminresolver"
	"github.com/agynio/gateway/internal/httpauth"
	"github.com/agynio/gateway/internal/identity"
)

type ExposeGateway struct {
	expose               exposev1.ExposeServiceClient
	clusterAdminResolver *clusteradminresolver.Resolver
}

func NewExposeGateway(expose exposev1.ExposeServiceClient, clusterAdminResolver *clusteradminresolver.Resolver) *ExposeGateway {
	return &ExposeGateway{expose: expose, clusterAdminResolver: clusterAdminResolver}
}

func (g *ExposeGateway) AddExposure(ctx context.Context, req *connect.Request[exposev1.AddExposureRequest]) (*connect.Response[exposev1.AddExposureResponse], error) {
	caller, err := g.resolveExposureCaller(ctx)
	if err != nil {
		return nil, err
	}
	workloadID, agentID, err := resolveAddExposureIDs(caller, req.Msg.GetWorkloadId(), req.Msg.GetAgentId())
	if err != nil {
		return nil, err
	}
	resp, err := g.expose.AddExposure(ctx, &exposev1.AddExposureRequest{
		WorkloadId: workloadID,
		AgentId:    agentID,
		Port:       req.Msg.GetPort(),
	})
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ExposeGateway) RemoveExposure(ctx context.Context, req *connect.Request[exposev1.RemoveExposureRequest]) (*connect.Response[exposev1.RemoveExposureResponse], error) {
	caller, err := g.resolveExposureCaller(ctx)
	if err != nil {
		return nil, err
	}
	workloadID, err := resolveWorkloadID(caller, req.Msg.GetWorkloadId())
	if err != nil {
		return nil, err
	}
	resp, err := g.expose.RemoveExposure(ctx, &exposev1.RemoveExposureRequest{
		WorkloadId: workloadID,
		Port:       req.Msg.GetPort(),
	})
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *ExposeGateway) ListExposures(ctx context.Context, req *connect.Request[exposev1.ListExposuresRequest]) (*connect.Response[exposev1.ListExposuresResponse], error) {
	caller, err := g.resolveExposureCaller(ctx)
	if err != nil {
		return nil, err
	}
	workloadID, err := resolveWorkloadID(caller, req.Msg.GetWorkloadId())
	if err != nil {
		return nil, err
	}
	resp, err := g.expose.ListExposures(ctx, &exposev1.ListExposuresRequest{
		WorkloadId: workloadID,
		PageSize:   req.Msg.GetPageSize(),
		PageToken:  req.Msg.GetPageToken(),
	})
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

type exposureCaller struct {
	identity       identity.ResolvedIdentity
	isClusterAdmin bool
}

func (g *ExposeGateway) resolveExposureCaller(ctx context.Context) (exposureCaller, error) {
	if token, ok := httpauth.BearerTokenFromContext(ctx); ok {
		if g.clusterAdminResolver != nil && g.clusterAdminResolver.Matches(token) {
			return exposureCaller{isClusterAdmin: true}, nil
		}
	}
	resolved, err := requireAgentIdentity(ctx)
	if err != nil {
		return exposureCaller{}, err
	}
	return exposureCaller{identity: resolved}, nil
}

func resolveAddExposureIDs(caller exposureCaller, workloadID, agentID string) (string, string, error) {
	if caller.isClusterAdmin {
		resolvedWorkloadID, err := requireClusterAdminID(workloadID, "workload")
		if err != nil {
			return "", "", err
		}
		resolvedAgentID, err := requireClusterAdminID(agentID, "agent")
		if err != nil {
			return "", "", err
		}
		return resolvedWorkloadID, resolvedAgentID, nil
	}
	resolvedWorkloadID, err := resolveAgentIDMatch(caller.identity.WorkloadID, workloadID, "workload")
	if err != nil {
		return "", "", err
	}
	resolvedAgentID, err := resolveAgentIDMatch(caller.identity.IdentityID, agentID, "agent")
	if err != nil {
		return "", "", err
	}
	return resolvedWorkloadID, resolvedAgentID, nil
}

func resolveWorkloadID(caller exposureCaller, workloadID string) (string, error) {
	if caller.isClusterAdmin {
		return requireClusterAdminID(workloadID, "workload")
	}
	return resolveAgentIDMatch(caller.identity.WorkloadID, workloadID, "workload")
}

func requireClusterAdminID(value, label string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", connect.NewError(connect.CodePermissionDenied, errors.New(label+" id required for cluster admin"))
	}
	return trimmed, nil
}

func resolveAgentIDMatch(expectedID, providedID, label string) (string, error) {
	trimmed := strings.TrimSpace(providedID)
	if trimmed != "" && trimmed != expectedID {
		return "", connect.NewError(connect.CodePermissionDenied, errors.New(label+" id does not match identity"))
	}
	return expectedID, nil
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
