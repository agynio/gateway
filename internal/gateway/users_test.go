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
	createUserReq      *usersv1.CreateUserRequest
	createUserResp     *usersv1.CreateUserResponse
	createUserErr      error
	createUserCalls    int
	updateMeReq        *usersv1.UpdateMeRequest
	updateMeResp       *usersv1.UpdateMeResponse
	updateMeErr        error
	updateMeCalls      int
	searchUsersReq     *usersv1.SearchUsersRequest
	searchUsersResp    *usersv1.SearchUsersResponse
	searchUsersErr     error
	searchUsersCalls   int
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
	createDeviceReq    *usersv1.CreateDeviceRequest
	createDeviceResp   *usersv1.CreateDeviceResponse
	createDeviceErr    error
	createDeviceCalls  int
	listDevicesReq     *usersv1.ListDevicesRequest
	listDevicesResp    *usersv1.ListDevicesResponse
	listDevicesErr     error
	listDevicesCalls   int
	deleteDeviceReq    *usersv1.DeleteDeviceRequest
	deleteDeviceResp   *usersv1.DeleteDeviceResponse
	deleteDeviceErr    error
	deleteDeviceCalls  int
}

func (f *fakeUsersClient) ResolveOrCreateUser(ctx context.Context, in *usersv1.ResolveOrCreateUserRequest, opts ...grpc.CallOption) (*usersv1.ResolveOrCreateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ResolveOrCreateUser not implemented")
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

func (f *fakeUsersClient) UpdateMe(ctx context.Context, in *usersv1.UpdateMeRequest, opts ...grpc.CallOption) (*usersv1.UpdateMeResponse, error) {
	f.updateMeCalls++
	f.updateMeReq = in
	if f.updateMeErr != nil {
		return nil, f.updateMeErr
	}
	if f.updateMeResp == nil {
		f.updateMeResp = &usersv1.UpdateMeResponse{}
	}
	return f.updateMeResp, nil
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

func (f *fakeUsersClient) SearchUsers(ctx context.Context, in *usersv1.SearchUsersRequest, opts ...grpc.CallOption) (*usersv1.SearchUsersResponse, error) {
	f.searchUsersCalls++
	f.searchUsersReq = in
	if f.searchUsersErr != nil {
		return nil, f.searchUsersErr
	}
	if f.searchUsersResp == nil {
		f.searchUsersResp = &usersv1.SearchUsersResponse{}
	}
	return f.searchUsersResp, nil
}

func (f *fakeUsersClient) CreateUser(ctx context.Context, in *usersv1.CreateUserRequest, opts ...grpc.CallOption) (*usersv1.CreateUserResponse, error) {
	f.createUserCalls++
	f.createUserReq = in
	if f.createUserErr != nil {
		return nil, f.createUserErr
	}
	if f.createUserResp == nil {
		f.createUserResp = &usersv1.CreateUserResponse{}
	}
	return f.createUserResp, nil
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

func (f *fakeUsersClient) CreateDevice(ctx context.Context, in *usersv1.CreateDeviceRequest, opts ...grpc.CallOption) (*usersv1.CreateDeviceResponse, error) {
	f.createDeviceCalls++
	f.createDeviceReq = in
	if f.createDeviceErr != nil {
		return nil, f.createDeviceErr
	}
	if f.createDeviceResp == nil {
		f.createDeviceResp = &usersv1.CreateDeviceResponse{}
	}
	return f.createDeviceResp, nil
}

func (f *fakeUsersClient) ListDevices(ctx context.Context, in *usersv1.ListDevicesRequest, opts ...grpc.CallOption) (*usersv1.ListDevicesResponse, error) {
	f.listDevicesCalls++
	f.listDevicesReq = in
	if f.listDevicesErr != nil {
		return nil, f.listDevicesErr
	}
	if f.listDevicesResp == nil {
		f.listDevicesResp = &usersv1.ListDevicesResponse{}
	}
	return f.listDevicesResp, nil
}

func (f *fakeUsersClient) DeleteDevice(ctx context.Context, in *usersv1.DeleteDeviceRequest, opts ...grpc.CallOption) (*usersv1.DeleteDeviceResponse, error) {
	f.deleteDeviceCalls++
	f.deleteDeviceReq = in
	if f.deleteDeviceErr != nil {
		return nil, f.deleteDeviceErr
	}
	if f.deleteDeviceResp == nil {
		f.deleteDeviceResp = &usersv1.DeleteDeviceResponse{}
	}
	return f.deleteDeviceResp, nil
}

func (f *fakeUsersClient) ResolveAPIToken(ctx context.Context, in *usersv1.ResolveAPITokenRequest, opts ...grpc.CallOption) (*usersv1.ResolveAPITokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ResolveAPIToken not implemented")
}

func TestGetMe_Success(t *testing.T) {
	client := &fakeUsersClient{
		getMeResp: &usersv1.GetMeResponse{
			User:        &usersv1.User{Meta: &usersv1.EntityMeta{Id: "user-1"}, Name: "Ada"},
			ClusterRole: usersv1.ClusterRole_CLUSTER_ROLE_UNSPECIFIED,
		},
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.GetMeRequest{})
	resp, err := gateway.GetMe(context.Background(), req)
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
	if client.getMeReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.getMeResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestUpdateMe_Success(t *testing.T) {
	client := &fakeUsersClient{
		updateMeResp: &usersv1.UpdateMeResponse{
			User: &usersv1.User{Meta: &usersv1.EntityMeta{Id: "user-1"}, Name: "Ada"},
		},
	}
	gateway := NewUsersGateway(client)

	name := "Ada"
	username := "ada"
	req := connect.NewRequest(&usersv1.UpdateMeRequest{Name: &name, Username: &username})
	resp, err := gateway.UpdateMe(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.updateMeCalls != 1 {
		t.Fatalf("expected update me to be called once, got %d", client.updateMeCalls)
	}
	if client.updateMeReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.updateMeResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestUpdateMe_Error(t *testing.T) {
	client := &fakeUsersClient{updateMeErr: status.Error(codes.PermissionDenied, "denied")}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.UpdateMeRequest{})
	resp, err := gateway.UpdateMe(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected CodePermissionDenied, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.updateMeCalls != 1 {
		t.Fatalf("expected update me to be called once, got %d", client.updateMeCalls)
	}
}

func TestCreateUser_Success(t *testing.T) {
	client := &fakeUsersClient{
		createUserResp: &usersv1.CreateUserResponse{
			User: &usersv1.User{Meta: &usersv1.EntityMeta{Id: "user-1"}, Name: "Ada"},
		},
	}
	gateway := NewUsersGateway(client)

	name := "Ada"
	nickname := "Ada Lovelace"
	username := "ada"
	photoURL := "https://example.com/avatar.png"
	req := connect.NewRequest(&usersv1.CreateUserRequest{
		OidcSubject: " subject-1 ",
		Name:        &name,
		Nickname:    &nickname,
		Username:    &username,
		PhotoUrl:    &photoURL,
		ClusterRole: usersv1.ClusterRole_CLUSTER_ROLE_ADMIN,
	})
	resp, err := gateway.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.createUserCalls != 1 {
		t.Fatalf("expected create user to be called once, got %d", client.createUserCalls)
	}
	if client.createUserReq == nil {
		t.Fatalf("expected create user request")
	}
	if client.createUserReq.OidcSubject != "subject-1" {
		t.Fatalf("expected oidc subject %q, got %q", "subject-1", client.createUserReq.OidcSubject)
	}
	if client.createUserReq.Name == nil || *client.createUserReq.Name != name {
		t.Fatalf("expected name %q, got %q", name, client.createUserReq.GetName())
	}
	if client.createUserReq.Nickname == nil || *client.createUserReq.Nickname != nickname {
		t.Fatalf("expected nickname %q, got %q", nickname, client.createUserReq.GetNickname())
	}
	if client.createUserReq.Username == nil || *client.createUserReq.Username != username {
		t.Fatalf("expected username %q, got %q", username, client.createUserReq.GetUsername())
	}
	if client.createUserReq.PhotoUrl == nil || *client.createUserReq.PhotoUrl != photoURL {
		t.Fatalf("expected photo url %q, got %q", photoURL, client.createUserReq.GetPhotoUrl())
	}
	if client.createUserReq.ClusterRole != usersv1.ClusterRole_CLUSTER_ROLE_ADMIN {
		t.Fatalf("expected cluster role %v, got %v", usersv1.ClusterRole_CLUSTER_ROLE_ADMIN, client.createUserReq.ClusterRole)
	}
	if resp.Msg != client.createUserResp {
		t.Fatalf("expected response to be forwarded")
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
	if client.createUserCalls != 0 {
		t.Fatalf("expected create user to not be called, got %d", client.createUserCalls)
	}
}

func TestCreateUser_NilUserResponse(t *testing.T) {
	client := &fakeUsersClient{createUserResp: &usersv1.CreateUserResponse{}}
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
	if client.createUserCalls != 1 {
		t.Fatalf("expected create user to be called once, got %d", client.createUserCalls)
	}
}

func TestUsersGatewaySearchUsersSuccess(t *testing.T) {
	client := &fakeUsersClient{
		searchUsersResp: &usersv1.SearchUsersResponse{
			Users: []*usersv1.UserDirectoryEntry{{IdentityId: "user-1", Name: "Ada"}},
		},
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.SearchUsersRequest{Prefix: "ad", Limit: 10})
	resp, err := gateway.SearchUsers(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.searchUsersCalls != 1 {
		t.Fatalf("expected search users to be called once, got %d", client.searchUsersCalls)
	}
	if client.searchUsersReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.searchUsersResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestUsersGatewaySearchUsersError(t *testing.T) {
	client := &fakeUsersClient{searchUsersErr: status.Error(codes.NotFound, "missing")}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.SearchUsersRequest{Prefix: "missing"})
	resp, err := gateway.SearchUsers(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeNotFound {
		t.Fatalf("expected CodeNotFound, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.searchUsersCalls != 1 {
		t.Fatalf("expected search users to be called once, got %d", client.searchUsersCalls)
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

func TestUsersGatewayCreateDeviceSuccess(t *testing.T) {
	client := &fakeUsersClient{
		createDeviceResp: &usersv1.CreateDeviceResponse{
			Device: &usersv1.Device{Meta: &usersv1.EntityMeta{Id: "device-1"}, Name: "Laptop"},
		},
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.CreateDeviceRequest{Name: "Laptop"})
	resp, err := gateway.CreateDevice(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.createDeviceCalls != 1 {
		t.Fatalf("expected create device to be called once, got %d", client.createDeviceCalls)
	}
	if client.createDeviceReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.createDeviceResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestUsersGatewayCreateDeviceError(t *testing.T) {
	client := &fakeUsersClient{
		createDeviceErr: status.Error(codes.NotFound, "missing"),
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.CreateDeviceRequest{Name: "Laptop"})
	resp, err := gateway.CreateDevice(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeNotFound {
		t.Fatalf("expected CodeNotFound, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.createDeviceCalls != 1 {
		t.Fatalf("expected create device to be called once, got %d", client.createDeviceCalls)
	}
	if client.createDeviceReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
}

func TestUsersGatewayListDevicesSuccess(t *testing.T) {
	client := &fakeUsersClient{
		listDevicesResp: &usersv1.ListDevicesResponse{
			Devices: []*usersv1.Device{{Meta: &usersv1.EntityMeta{Id: "device-1"}, Name: "Laptop"}},
		},
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.ListDevicesRequest{PageSize: 25})
	resp, err := gateway.ListDevices(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.listDevicesCalls != 1 {
		t.Fatalf("expected list devices to be called once, got %d", client.listDevicesCalls)
	}
	if client.listDevicesReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.listDevicesResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestUsersGatewayListDevicesError(t *testing.T) {
	client := &fakeUsersClient{
		listDevicesErr: status.Error(codes.PermissionDenied, "denied"),
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.ListDevicesRequest{PageToken: "next"})
	resp, err := gateway.ListDevices(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected CodePermissionDenied, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.listDevicesCalls != 1 {
		t.Fatalf("expected list devices to be called once, got %d", client.listDevicesCalls)
	}
	if client.listDevicesReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
}

func TestUsersGatewayDeleteDeviceSuccess(t *testing.T) {
	client := &fakeUsersClient{
		deleteDeviceResp: &usersv1.DeleteDeviceResponse{},
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.DeleteDeviceRequest{Id: "device-1"})
	resp, err := gateway.DeleteDevice(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.deleteDeviceCalls != 1 {
		t.Fatalf("expected delete device to be called once, got %d", client.deleteDeviceCalls)
	}
	if client.deleteDeviceReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.deleteDeviceResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestUsersGatewayDeleteDeviceError(t *testing.T) {
	client := &fakeUsersClient{
		deleteDeviceErr: status.Error(codes.Internal, "boom"),
	}
	gateway := NewUsersGateway(client)

	req := connect.NewRequest(&usersv1.DeleteDeviceRequest{Id: "device-1"})
	resp, err := gateway.DeleteDevice(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeInternal {
		t.Fatalf("expected CodeInternal, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.deleteDeviceCalls != 1 {
		t.Fatalf("expected delete device to be called once, got %d", client.deleteDeviceCalls)
	}
	if client.deleteDeviceReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
}
