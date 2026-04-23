package gateway

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	organizationsv1 "github.com/agynio/gateway/gen/agynio/api/organizations/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeOrganizationsClient struct {
	listAccessibleReq   *organizationsv1.ListAccessibleOrganizationsRequest
	listAccessibleResp  *organizationsv1.ListAccessibleOrganizationsResponse
	listAccessibleErr   error
	listAccessibleCalls int
	listMyMembershipsReq   *organizationsv1.ListMyMembershipsRequest
	listMyMembershipsResp  *organizationsv1.ListMyMembershipsResponse
	listMyMembershipsErr   error
	listMyMembershipsCalls int
}

func (f *fakeOrganizationsClient) CreateOrganization(ctx context.Context, in *organizationsv1.CreateOrganizationRequest, opts ...grpc.CallOption) (*organizationsv1.CreateOrganizationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateOrganization not implemented")
}

func (f *fakeOrganizationsClient) GetOrganization(ctx context.Context, in *organizationsv1.GetOrganizationRequest, opts ...grpc.CallOption) (*organizationsv1.GetOrganizationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetOrganization not implemented")
}

func (f *fakeOrganizationsClient) UpdateOrganization(ctx context.Context, in *organizationsv1.UpdateOrganizationRequest, opts ...grpc.CallOption) (*organizationsv1.UpdateOrganizationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "UpdateOrganization not implemented")
}

func (f *fakeOrganizationsClient) DeleteOrganization(ctx context.Context, in *organizationsv1.DeleteOrganizationRequest, opts ...grpc.CallOption) (*organizationsv1.DeleteOrganizationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteOrganization not implemented")
}

func (f *fakeOrganizationsClient) ListOrganizations(ctx context.Context, in *organizationsv1.ListOrganizationsRequest, opts ...grpc.CallOption) (*organizationsv1.ListOrganizationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListOrganizations not implemented")
}

func (f *fakeOrganizationsClient) ListAccessibleOrganizations(ctx context.Context, in *organizationsv1.ListAccessibleOrganizationsRequest, opts ...grpc.CallOption) (*organizationsv1.ListAccessibleOrganizationsResponse, error) {
	f.listAccessibleCalls++
	f.listAccessibleReq = in
	if f.listAccessibleErr != nil {
		return nil, f.listAccessibleErr
	}
	if f.listAccessibleResp == nil {
		f.listAccessibleResp = &organizationsv1.ListAccessibleOrganizationsResponse{}
	}
	return f.listAccessibleResp, nil
}

func (f *fakeOrganizationsClient) CreateMembership(ctx context.Context, in *organizationsv1.CreateMembershipRequest, opts ...grpc.CallOption) (*organizationsv1.CreateMembershipResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateMembership not implemented")
}

func (f *fakeOrganizationsClient) AcceptMembership(ctx context.Context, in *organizationsv1.AcceptMembershipRequest, opts ...grpc.CallOption) (*organizationsv1.AcceptMembershipResponse, error) {
	return nil, status.Error(codes.Unimplemented, "AcceptMembership not implemented")
}

func (f *fakeOrganizationsClient) DeclineMembership(ctx context.Context, in *organizationsv1.DeclineMembershipRequest, opts ...grpc.CallOption) (*organizationsv1.DeclineMembershipResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DeclineMembership not implemented")
}

func (f *fakeOrganizationsClient) RemoveMembership(ctx context.Context, in *organizationsv1.RemoveMembershipRequest, opts ...grpc.CallOption) (*organizationsv1.RemoveMembershipResponse, error) {
	return nil, status.Error(codes.Unimplemented, "RemoveMembership not implemented")
}

func (f *fakeOrganizationsClient) UpdateMembershipRole(ctx context.Context, in *organizationsv1.UpdateMembershipRoleRequest, opts ...grpc.CallOption) (*organizationsv1.UpdateMembershipRoleResponse, error) {
	return nil, status.Error(codes.Unimplemented, "UpdateMembershipRole not implemented")
}

func (f *fakeOrganizationsClient) ListMembers(ctx context.Context, in *organizationsv1.ListMembersRequest, opts ...grpc.CallOption) (*organizationsv1.ListMembersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListMembers not implemented")
}

func (f *fakeOrganizationsClient) ListMyMemberships(ctx context.Context, in *organizationsv1.ListMyMembershipsRequest, opts ...grpc.CallOption) (*organizationsv1.ListMyMembershipsResponse, error) {
	f.listMyMembershipsCalls++
	f.listMyMembershipsReq = in
	if f.listMyMembershipsErr != nil {
		return nil, f.listMyMembershipsErr
	}
	if f.listMyMembershipsResp == nil {
		f.listMyMembershipsResp = &organizationsv1.ListMyMembershipsResponse{}
	}
	return f.listMyMembershipsResp, nil
}

func (f *fakeOrganizationsClient) SetMyOrgNickname(ctx context.Context, in *organizationsv1.SetMyOrgNicknameRequest, opts ...grpc.CallOption) (*organizationsv1.SetMyOrgNicknameResponse, error) {
	return nil, status.Error(codes.Unimplemented, "SetMyOrgNickname not implemented")
}

func TestOrganizationsGatewayListAccessibleOrganizations(t *testing.T) {
	resolved := identity.ResolvedIdentity{
		IdentityID:   "identity-1",
		IdentityType: identity.IdentityTypeUser,
	}
	ctx := identity.WithIdentity(context.Background(), resolved)

	client := &fakeOrganizationsClient{
		listAccessibleResp: &organizationsv1.ListAccessibleOrganizationsResponse{},
	}
	gateway := NewOrganizationsGateway(client)

	req := connect.NewRequest(&organizationsv1.ListAccessibleOrganizationsRequest{})
	resp, err := gateway.ListAccessibleOrganizations(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.listAccessibleCalls != 1 {
		t.Fatalf("expected list accessible to be called once, got %d", client.listAccessibleCalls)
	}
	if client.listAccessibleReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if client.listAccessibleReq.IdentityId != resolved.IdentityID {
		t.Fatalf("expected identity_id %q, got %q", resolved.IdentityID, client.listAccessibleReq.IdentityId)
	}
	if resp.Msg != client.listAccessibleResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestOrganizationsGatewayListAccessibleOrganizationsMissingIdentity(t *testing.T) {
	client := &fakeOrganizationsClient{}
	gateway := NewOrganizationsGateway(client)

	req := connect.NewRequest(&organizationsv1.ListAccessibleOrganizationsRequest{})
	resp, err := gateway.ListAccessibleOrganizations(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Fatalf("expected CodeUnauthenticated, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.listAccessibleCalls != 0 {
		t.Fatalf("expected list accessible to not be called, got %d", client.listAccessibleCalls)
	}
}

func TestListMyMemberships_MissingIdentity(t *testing.T) {
	client := &fakeOrganizationsClient{}
	gateway := NewOrganizationsGateway(client)

	req := connect.NewRequest(&organizationsv1.ListMyMembershipsRequest{})
	resp, err := gateway.ListMyMemberships(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Fatalf("expected CodeUnauthenticated, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.listMyMembershipsCalls != 0 {
		t.Fatalf("expected list my memberships to not be called, got %d", client.listMyMembershipsCalls)
	}
}

func TestListMyMemberships_Success(t *testing.T) {
	resolved := identity.ResolvedIdentity{
		IdentityID:   "identity-1",
		IdentityType: identity.IdentityTypeUser,
	}
	ctx := identity.WithIdentity(context.Background(), resolved)

	client := &fakeOrganizationsClient{
		listMyMembershipsResp: &organizationsv1.ListMyMembershipsResponse{},
	}
	gateway := NewOrganizationsGateway(client)

	req := connect.NewRequest(&organizationsv1.ListMyMembershipsRequest{})
	resp, err := gateway.ListMyMemberships(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.listMyMembershipsCalls != 1 {
		t.Fatalf("expected list my memberships to be called once, got %d", client.listMyMembershipsCalls)
	}
	if client.listMyMembershipsReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.listMyMembershipsResp {
		t.Fatalf("expected response to be forwarded")
	}
}
