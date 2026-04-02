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

func (g *OrganizationsGateway) UpdateOrganization(ctx context.Context, req *connect.Request[organizationsv1.UpdateOrganizationRequest]) (*connect.Response[organizationsv1.UpdateOrganizationResponse], error) {
	resp, err := g.organizations.UpdateOrganization(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) DeleteOrganization(ctx context.Context, req *connect.Request[organizationsv1.DeleteOrganizationRequest]) (*connect.Response[organizationsv1.DeleteOrganizationResponse], error) {
	resp, err := g.organizations.DeleteOrganization(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) ListOrganizations(ctx context.Context, req *connect.Request[organizationsv1.ListOrganizationsRequest]) (*connect.Response[organizationsv1.ListOrganizationsResponse], error) {
	resp, err := g.organizations.ListOrganizations(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) CreateMembership(ctx context.Context, req *connect.Request[organizationsv1.CreateMembershipRequest]) (*connect.Response[organizationsv1.CreateMembershipResponse], error) {
	resp, err := g.organizations.CreateMembership(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) AcceptMembership(ctx context.Context, req *connect.Request[organizationsv1.AcceptMembershipRequest]) (*connect.Response[organizationsv1.AcceptMembershipResponse], error) {
	resp, err := g.organizations.AcceptMembership(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) DeclineMembership(ctx context.Context, req *connect.Request[organizationsv1.DeclineMembershipRequest]) (*connect.Response[organizationsv1.DeclineMembershipResponse], error) {
	resp, err := g.organizations.DeclineMembership(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) RemoveMembership(ctx context.Context, req *connect.Request[organizationsv1.RemoveMembershipRequest]) (*connect.Response[organizationsv1.RemoveMembershipResponse], error) {
	resp, err := g.organizations.RemoveMembership(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) UpdateMembershipRole(ctx context.Context, req *connect.Request[organizationsv1.UpdateMembershipRoleRequest]) (*connect.Response[organizationsv1.UpdateMembershipRoleResponse], error) {
	resp, err := g.organizations.UpdateMembershipRole(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) ListMembers(ctx context.Context, req *connect.Request[organizationsv1.ListMembersRequest]) (*connect.Response[organizationsv1.ListMembersResponse], error) {
	resp, err := g.organizations.ListMembers(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *OrganizationsGateway) ListMyMemberships(ctx context.Context, req *connect.Request[organizationsv1.ListMyMembershipsRequest]) (*connect.Response[organizationsv1.ListMyMembershipsResponse], error) {
	resp, err := g.organizations.ListMyMemberships(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
