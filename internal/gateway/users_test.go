package gateway

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeUsersClient struct {
	resolveReq         *usersv1.ResolveOrCreateUserRequest
	resolveResp        *usersv1.ResolveOrCreateUserResponse
	resolveErr         error
	resolveCalls       int
	getUserReq         *usersv1.GetUserRequest
	getUserResp        *usersv1.GetUserResponse
	getUserErr         error
	getUserCalls       int
	getMeReq           *usersv1.GetMeRequest
	getMeResp          *usersv1.GetMeResponse
	getMeErr           error
	getMeCalls         int
	batchGetUsersReq   *usersv1.BatchGetUsersRequest
	batchGetUsersResp  *usersv1.BatchGetUsersResponse
	batchGetUsersErr   error
	batchGetUsersCalls int
}

func (f *fakeUsersClient) ResolveOrCreateUser(ctx context.Context, in *usersv1.ResolveOrCreateUserRequest, opts ...grpc.CallOption) (*usersv1.ResolveOrCreateUserResponse, error) {
	f.resolveCalls++
	f.resolveReq = in
	if f.resolveErr != nil {
		return nil, f.resolveErr
	}
	if f.resolveResp == nil {
		f.resolveResp = &usersv1.ResolveOrCreateUserResponse{}
	}
	return f.resolveResp, nil
}

func (f *fakeUsersClient) GetUser(ctx context.Context, in *usersv1.GetUserRequest, opts ...grpc.CallOption) (*usersv1.GetUserResponse, error) {
	f.getUserCalls++
	f.getUserReq = in
	if f.getUserErr != nil {
		return nil, f.getUserErr
	}
	if f.getUserResp == nil {
		f.getUserResp = &usersv1.GetUserResponse{}
	}
	return f.getUserResp, nil
}

func (f *fakeUsersClient) GetUserByOIDCSubject(ctx context.Context, in *usersv1.GetUserByOIDCSubjectRequest, opts ...grpc.CallOption) (*usersv1.GetUserByOIDCSubjectResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetUserByOIDCSubject not implemented")
}

func (f *fakeUsersClient) GetMe(ctx context.Context, in *usersv1.GetMeRequest, opts ...grpc.CallOption) (*usersv1.GetMeResponse, error) {
	f.getMeCalls++
	f.getMeReq = in
	if f.getMeErr != nil {
		return nil, f.getMeErr
	}
	if f.getMeResp == nil {
		f.getMeResp = &usersv1.GetMeResponse{}
	}
	return f.getMeResp, nil
}

func (f *fakeUsersClient) BatchGetUsers(ctx context.Context, in *usersv1.BatchGetUsersRequest, opts ...grpc.CallOption) (*usersv1.BatchGetUsersResponse, error) {
	f.batchGetUsersCalls++
	f.batchGetUsersReq = in
	if f.batchGetUsersErr != nil {
		return nil, f.batchGetUsersErr
	}
	if f.batchGetUsersResp == nil {
		f.batchGetUsersResp = &usersv1.BatchGetUsersResponse{}
	}
	return f.batchGetUsersResp, nil
}

func (f *fakeUsersClient) UpdateUser(ctx context.Context, in *usersv1.UpdateUserRequest, opts ...grpc.CallOption) (*usersv1.UpdateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "UpdateUser not implemented")
}

func (f *fakeUsersClient) ListUsers(ctx context.Context, in *usersv1.ListUsersRequest, opts ...grpc.CallOption) (*usersv1.ListUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListUsers not implemented")
}

func (f *fakeUsersClient) CreateUser(ctx context.Context, in *usersv1.CreateUserRequest, opts ...grpc.CallOption) (*usersv1.CreateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateUser not implemented")
}

func (f *fakeUsersClient) DeleteUser(ctx context.Context, in *usersv1.DeleteUserRequest, opts ...grpc.CallOption) (*usersv1.DeleteUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteUser not implemented")
}

func (f *fakeUsersClient) CreateAPIToken(ctx context.Context, in *usersv1.CreateAPITokenRequest, opts ...grpc.CallOption) (*usersv1.CreateAPITokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateAPIToken not implemented")
}

func (f *fakeUsersClient) ListAPITokens(ctx context.Context, in *usersv1.ListAPITokensRequest, opts ...grpc.CallOption) (*usersv1.ListAPITokensResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListAPITokens not implemented")
}

func (f *fakeUsersClient) RevokeAPIToken(ctx context.Context, in *usersv1.RevokeAPITokenRequest, opts ...grpc.CallOption) (*usersv1.RevokeAPITokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "RevokeAPIToken not implemented")
}

func (f *fakeUsersClient) ResolveAPIToken(ctx context.Context, in *usersv1.ResolveAPITokenRequest, opts ...grpc.CallOption) (*usersv1.ResolveAPITokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ResolveAPIToken not implemented")
}

func TestGetMe_Success(t *testing.T) {
	resolved := identity.ResolvedIdentity{
		IdentityID:   "identity-1",
		IdentityType: identity.IdentityTypeUser,
	}
	ctx := identity.WithIdentity(context.Background(), resolved)

	client := &fakeUsersClient{
		getMeResp: &usersv1.GetMeResponse{
			User:        &usersv1.User{Meta: &usersv1.EntityMeta{Id: "user-1"}, Name: "Ada"},
			ClusterRole: usersv1.ClusterRole_CLUSTER_ROLE_ADMIN,
		},
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.GetMeRequest{})
	resp, err := gateway.GetMe(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.getMeCalls != 1 {
		t.Fatalf("expected get me to be called once, got %d", client.getMeCalls)
	}
	if client.getMeReq == nil {
		t.Fatalf("expected get me request")
	}
	if resp.Msg.User != client.getMeResp.User {
		t.Fatalf("expected user to be forwarded")
	}
	if resp.Msg.ClusterRole != client.getMeResp.ClusterRole {
		t.Fatalf("expected cluster role %v, got %v", client.getMeResp.ClusterRole, resp.Msg.ClusterRole)
	}
}

func TestGetMe_MissingIdentity(t *testing.T) {
	client := &fakeUsersClient{}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.GetMeRequest{})
	resp, err := gateway.GetMe(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Fatalf("expected CodeUnauthenticated, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.getMeCalls != 0 {
		t.Fatalf("expected get me to not be called, got %d", client.getMeCalls)
	}
}

func TestGetMe_ClusterRoleDefault(t *testing.T) {
	resolved := identity.ResolvedIdentity{
		IdentityID:   "identity-1",
		IdentityType: identity.IdentityTypeUser,
	}
	ctx := identity.WithIdentity(context.Background(), resolved)

	client := &fakeUsersClient{
		getMeResp: &usersv1.GetMeResponse{
			User:        &usersv1.User{Meta: &usersv1.EntityMeta{Id: "user-2"}, Name: "Linus"},
			ClusterRole: usersv1.ClusterRole_CLUSTER_ROLE_UNSPECIFIED,
		},
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.GetMeRequest{})
	resp, err := gateway.GetMe(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if resp.Msg.ClusterRole != usersv1.ClusterRole_CLUSTER_ROLE_ADMIN {
		t.Fatalf("expected cluster role default admin, got %v", resp.Msg.ClusterRole)
	}
}

func TestCreateUser_Success(t *testing.T) {
	client := &fakeUsersClient{
		resolveResp: &usersv1.ResolveOrCreateUserResponse{
			User: &usersv1.User{Meta: &usersv1.EntityMeta{Id: "user-1"}, Name: "Ada"},
		},
	}
	gateway := NewUsersGateway(client)

	name := "Ada"
	photoURL := "https://example.com/avatar.png"
	req := connect.NewRequest(&usersv1.CreateUserRequest{
		OidcSubject: " subject-1 ",
		Name:        &name,
		PhotoUrl:    &photoURL,
	})
	resp, err := gateway.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.resolveCalls != 1 {
		t.Fatalf("expected resolve-or-create to be called once, got %d", client.resolveCalls)
	}
	if client.resolveReq == nil {
		t.Fatalf("expected resolve-or-create request")
	}
	if client.resolveReq.OidcSubject != "subject-1" {
		t.Fatalf("expected oidc subject %q, got %q", "subject-1", client.resolveReq.OidcSubject)
	}
	if client.resolveReq.Name != name {
		t.Fatalf("expected name %q, got %q", name, client.resolveReq.Name)
	}
	if client.resolveReq.Email != name {
		t.Fatalf("expected email %q, got %q", name, client.resolveReq.Email)
	}
	if client.resolveReq.PhotoUrl != photoURL {
		t.Fatalf("expected photo url %q, got %q", photoURL, client.resolveReq.PhotoUrl)
	}
	if resp.Msg.User != client.resolveResp.User {
		t.Fatalf("expected user to be forwarded")
	}
}

func TestCreateUser_EmptyOidcSubject(t *testing.T) {
	client := &fakeUsersClient{}
	gateway := NewUsersGateway(client)

	name := "Ada"
	req := connect.NewRequest(&usersv1.CreateUserRequest{OidcSubject: "  ", Name: &name})
	resp, err := gateway.CreateUser(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeInvalidArgument {
		t.Fatalf("expected CodeInvalidArgument, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.resolveCalls != 0 {
		t.Fatalf("expected resolve-or-create to not be called, got %d", client.resolveCalls)
	}
}

func TestCreateUser_NilUserResponse(t *testing.T) {
	client := &fakeUsersClient{resolveResp: &usersv1.ResolveOrCreateUserResponse{}}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.CreateUserRequest{OidcSubject: "subject-1"})
	resp, err := gateway.CreateUser(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeInternal {
		t.Fatalf("expected CodeInternal, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.resolveCalls != 1 {
		t.Fatalf("expected resolve-or-create to be called once, got %d", client.resolveCalls)
	}
}

func TestUsersGatewayBatchGetUsers(t *testing.T) {
	client := &fakeUsersClient{
		batchGetUsersResp: &usersv1.BatchGetUsersResponse{
			Users: []*usersv1.User{
				{OidcSubject: "subject-1", Name: "Ada"},
				{OidcSubject: "subject-2", Name: "Linus"},
			},
		},
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.BatchGetUsersRequest{IdentityIds: []string{"identity-1", "identity-2"}})
	resp, err := gateway.BatchGetUsers(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.batchGetUsersCalls != 1 {
		t.Fatalf("expected batch get users to be called once, got %d", client.batchGetUsersCalls)
	}
	if client.batchGetUsersReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.batchGetUsersResp {
		t.Fatalf("expected response to be forwarded")
	}
	for i, user := range resp.Msg.Users {
		if user.OidcSubject != "" {
			t.Fatalf("expected oidc_subject to be stripped at index %d", i)
		}
	}
	if resp.Msg.Users[0].Name != "Ada" {
		t.Fatalf("expected name Ada, got %q", resp.Msg.Users[0].Name)
	}
	if resp.Msg.Users[1].Name != "Linus" {
		t.Fatalf("expected name Linus, got %q", resp.Msg.Users[1].Name)
	}
}

func TestUsersGatewayBatchGetUsersError(t *testing.T) {
	client := &fakeUsersClient{
		batchGetUsersErr: status.Error(codes.NotFound, "missing"),
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.BatchGetUsersRequest{})
	resp, err := gateway.BatchGetUsers(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeNotFound {
		t.Fatalf("expected CodeNotFound, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.batchGetUsersCalls != 1 {
		t.Fatalf("expected batch get users to be called once, got %d", client.batchGetUsersCalls)
	}
	if client.batchGetUsersReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
}
