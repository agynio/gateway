package gateway

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeUsersClient struct {
	batchGetUsersReq   *usersv1.BatchGetUsersRequest
	batchGetUsersResp  *usersv1.BatchGetUsersResponse
	batchGetUsersErr   error
	batchGetUsersCalls int
}

func (f *fakeUsersClient) ResolveOrCreateUser(ctx context.Context, in *usersv1.ResolveOrCreateUserRequest, opts ...grpc.CallOption) (*usersv1.ResolveOrCreateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ResolveOrCreateUser not implemented")
}

func (f *fakeUsersClient) GetUser(ctx context.Context, in *usersv1.GetUserRequest, opts ...grpc.CallOption) (*usersv1.GetUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetUser not implemented")
}

func (f *fakeUsersClient) GetUserByOIDCSubject(ctx context.Context, in *usersv1.GetUserByOIDCSubjectRequest, opts ...grpc.CallOption) (*usersv1.GetUserByOIDCSubjectResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetUserByOIDCSubject not implemented")
}

func (f *fakeUsersClient) GetMe(ctx context.Context, in *usersv1.GetMeRequest, opts ...grpc.CallOption) (*usersv1.GetMeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetMe not implemented")
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
