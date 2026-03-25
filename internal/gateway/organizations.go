package gateway

import (
	"context"

	"connectrpc.com/connect"
	organizationsv1 "github.com/agynio/gateway/gen/agynio/api/organizations/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationsGateway struct {
	organizations organizationsv1.OrganizationsServiceClient
}

func NewOrganizationsGateway(organizations organizationsv1.OrganizationsServiceClient) *OrganizationsGateway {
	return &OrganizationsGateway{organizations: organizations}
}

func (g *OrganizationsGateway) ListAccessibleOrganizations(ctx context.Context, req *connect.Request[organizationsv1.ListAccessibleOrganizationsRequest]) (*connect.Response[organizationsv1.ListAccessibleOrganizationsResponse], error) {
	resolvedIdentity, ok := identity.IdentityFromContext(ctx)
	if !ok {
		return nil, toConnectError(status.Error(codes.Unauthenticated, "identity not available"))
	}

	req.Msg.IdentityId = resolvedIdentity.IdentityID
	resp, err := g.organizations.ListAccessibleOrganizations(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) CreateOrganization(ctx context.Context, req *connect.Request[organizationsv1.CreateOrganizationRequest]) (*connect.Response[organizationsv1.CreateOrganizationResponse], error) {
	resp, err := g.organizations.CreateOrganization(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) GetOrganization(ctx context.Context, req *connect.Request[organizationsv1.GetOrganizationRequest]) (*connect.Response[organizationsv1.GetOrganizationResponse], error) {
	resp, err := g.organizations.GetOrganization(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
