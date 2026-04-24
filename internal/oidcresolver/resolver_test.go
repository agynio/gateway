package oidcresolver

import (
	"context"
	"net/http"
	"strings"
	"testing"

	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/oidcauth"
	"github.com/agynio/gateway/internal/oidctestutil"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestResolveFromTokenUsesExistingUser(t *testing.T) {
	provider := oidctestutil.NewProvider(t)
	verifier, err := oidcauth.NewVerifier(context.Background(), provider.Issuer, provider.ClientID)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	usersClient := &fakeUsersClient{
		getBySubjectResp: &usersv1.GetUserByOIDCSubjectResponse{User: buildUser("user-1")},
	}
	resolver, err := NewResolver(verifier, usersClient, provider.Server.Client())
	if err != nil {
		t.Fatalf("failed to create resolver: %v", err)
	}

	resolved, err := resolver.ResolveFromToken(context.Background(), provider.Token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertResolvedIdentity(t, resolved, "user-1")
	if usersClient.getBySubjectCalls != 1 {
		t.Fatalf("expected get-by-subject to be called once, got %d", usersClient.getBySubjectCalls)
	}
	if usersClient.lastSubject != provider.Subject {
		t.Fatalf("expected subject %q, got %q", provider.Subject, usersClient.lastSubject)
	}
	if usersClient.resolveCalls != 0 {
		t.Fatalf("expected resolve-or-create not to be called, got %d", usersClient.resolveCalls)
	}
	if provider.UserinfoCalls != 0 {
		t.Fatalf("expected userinfo not to be called, got %d", provider.UserinfoCalls)
	}
}

func TestResolveFromTokenCreatesUserOnFirstLogin(t *testing.T) {
	provider := oidctestutil.NewProvider(t)
	provider.UserInfo.PreferredUsername = " ada "
	verifier, err := oidcauth.NewVerifier(context.Background(), provider.Issuer, provider.ClientID)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	usersClient := &fakeUsersClient{
		getBySubjectErr: status.Error(codes.NotFound, "not found"),
		resolveResp:     &usersv1.ResolveOrCreateUserResponse{User: buildUser("user-2")},
	}
	resolver, err := NewResolver(verifier, usersClient, provider.Server.Client())
	if err != nil {
		t.Fatalf("failed to create resolver: %v", err)
	}

	resolved, err := resolver.ResolveFromToken(context.Background(), provider.Token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertResolvedIdentity(t, resolved, "user-2")
	if usersClient.getBySubjectCalls != 1 {
		t.Fatalf("expected get-by-subject to be called once, got %d", usersClient.getBySubjectCalls)
	}
	if usersClient.resolveCalls != 1 {
		t.Fatalf("expected resolve-or-create to be called once, got %d", usersClient.resolveCalls)
	}
	if usersClient.lastResolve == nil {
		t.Fatalf("expected resolve-or-create request to be captured")
	}
	if usersClient.lastResolve.OidcSubject != provider.UserInfo.Subject {
		t.Fatalf("expected resolve subject %q, got %q", provider.UserInfo.Subject, usersClient.lastResolve.OidcSubject)
	}
	if usersClient.lastResolve.Name != provider.UserInfo.Name {
		t.Fatalf("expected resolve name %q, got %q", provider.UserInfo.Name, usersClient.lastResolve.Name)
	}
	if usersClient.lastResolve.Email != provider.UserInfo.Email {
		t.Fatalf("expected resolve email %q, got %q", provider.UserInfo.Email, usersClient.lastResolve.Email)
	}
	if usersClient.lastResolve.PhotoUrl != provider.UserInfo.Picture {
		t.Fatalf("expected resolve photo %q, got %q", provider.UserInfo.Picture, usersClient.lastResolve.PhotoUrl)
	}
	if usersClient.lastResolve.PreferredUsername == nil {
		t.Fatalf("expected resolve preferred username to be set")
	}
	if *usersClient.lastResolve.PreferredUsername != "ada" {
		t.Fatalf("expected resolve preferred username %q, got %q", "ada", *usersClient.lastResolve.PreferredUsername)
	}
	if provider.UserinfoCalls != 1 {
		t.Fatalf("expected userinfo to be called once, got %d", provider.UserinfoCalls)
	}
	if provider.LastAuthHeader != "Bearer "+provider.Token {
		t.Fatalf("expected bearer token header, got %q", provider.LastAuthHeader)
	}
}

func TestResolveFromTokenGetUserBySubjectError(t *testing.T) {
	provider := oidctestutil.NewProvider(t)
	verifier, err := oidcauth.NewVerifier(context.Background(), provider.Issuer, provider.ClientID)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	usersClient := &fakeUsersClient{
		getBySubjectErr: status.Error(codes.Internal, "boom"),
	}
	resolver, err := NewResolver(verifier, usersClient, provider.Server.Client())
	if err != nil {
		t.Fatalf("failed to create resolver: %v", err)
	}

	_, err = resolver.ResolveFromToken(context.Background(), provider.Token)
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected internal error, got %v", err)
	}
	if provider.UserinfoCalls != 0 {
		t.Fatalf("expected userinfo not to be called, got %d", provider.UserinfoCalls)
	}
	if usersClient.resolveCalls != 0 {
		t.Fatalf("expected resolve-or-create not to be called, got %d", usersClient.resolveCalls)
	}
}

func TestResolveFromTokenUserInfoStatusFailure(t *testing.T) {
	provider := oidctestutil.NewProvider(t, oidctestutil.WithUserInfoResponse(http.StatusBadGateway, "bad gateway"))
	verifier, err := oidcauth.NewVerifier(context.Background(), provider.Issuer, provider.ClientID)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	usersClient := &fakeUsersClient{
		getBySubjectErr: status.Error(codes.NotFound, "not found"),
	}
	resolver, err := NewResolver(verifier, usersClient, provider.Server.Client())
	if err != nil {
		t.Fatalf("failed to create resolver: %v", err)
	}

	_, err = resolver.ResolveFromToken(context.Background(), provider.Token)
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected internal error, got %v", err)
	}
	if !strings.Contains(err.Error(), "failed to fetch user info") {
		t.Fatalf("expected userinfo failure, got %v", err)
	}
	if provider.UserinfoCalls != 1 {
		t.Fatalf("expected userinfo to be called once, got %d", provider.UserinfoCalls)
	}
	if usersClient.resolveCalls != 0 {
		t.Fatalf("expected resolve-or-create not to be called, got %d", usersClient.resolveCalls)
	}
}

func TestResolveFromTokenUserInfoSubjectMismatch(t *testing.T) {
	provider := oidctestutil.NewProvider(t, oidctestutil.WithUserInfo(oidc.UserInfo{Subject: "mismatch"}))
	verifier, err := oidcauth.NewVerifier(context.Background(), provider.Issuer, provider.ClientID)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	usersClient := &fakeUsersClient{
		getBySubjectErr: status.Error(codes.NotFound, "not found"),
	}
	resolver, err := NewResolver(verifier, usersClient, provider.Server.Client())
	if err != nil {
		t.Fatalf("failed to create resolver: %v", err)
	}

	_, err = resolver.ResolveFromToken(context.Background(), provider.Token)
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected internal error, got %v", err)
	}
	if !strings.Contains(err.Error(), "subject does not match expected subject") {
		t.Fatalf("expected subject mismatch error, got %v", err)
	}
	if provider.UserinfoCalls != 1 {
		t.Fatalf("expected userinfo to be called once, got %d", provider.UserinfoCalls)
	}
	if usersClient.resolveCalls != 0 {
		t.Fatalf("expected resolve-or-create not to be called, got %d", usersClient.resolveCalls)
	}
}

type fakeUsersClient struct {
	getBySubjectResp  *usersv1.GetUserByOIDCSubjectResponse
	getBySubjectErr   error
	resolveResp       *usersv1.ResolveOrCreateUserResponse
	resolveErr        error
	getBySubjectCalls int
	resolveCalls      int
	lastSubject       string
	lastResolve       *usersv1.ResolveOrCreateUserRequest
}

func (f *fakeUsersClient) ResolveOrCreateUser(ctx context.Context, in *usersv1.ResolveOrCreateUserRequest, opts ...grpc.CallOption) (*usersv1.ResolveOrCreateUserResponse, error) {
	f.resolveCalls++
	f.lastResolve = in
	return f.resolveResp, f.resolveErr
}

func (f *fakeUsersClient) GetUser(ctx context.Context, in *usersv1.GetUserRequest, opts ...grpc.CallOption) (*usersv1.GetUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) GetUserByOIDCSubject(ctx context.Context, in *usersv1.GetUserByOIDCSubjectRequest, opts ...grpc.CallOption) (*usersv1.GetUserByOIDCSubjectResponse, error) {
	f.getBySubjectCalls++
	f.lastSubject = in.GetOidcSubject()
	return f.getBySubjectResp, f.getBySubjectErr
}

func (f *fakeUsersClient) GetMe(ctx context.Context, in *usersv1.GetMeRequest, opts ...grpc.CallOption) (*usersv1.GetMeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) UpdateMe(ctx context.Context, in *usersv1.UpdateMeRequest, opts ...grpc.CallOption) (*usersv1.UpdateMeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) BatchGetUsers(ctx context.Context, in *usersv1.BatchGetUsersRequest, opts ...grpc.CallOption) (*usersv1.BatchGetUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) UpdateUser(ctx context.Context, in *usersv1.UpdateUserRequest, opts ...grpc.CallOption) (*usersv1.UpdateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) ListUsers(ctx context.Context, in *usersv1.ListUsersRequest, opts ...grpc.CallOption) (*usersv1.ListUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) SearchUsers(ctx context.Context, in *usersv1.SearchUsersRequest, opts ...grpc.CallOption) (*usersv1.SearchUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) CreateUser(ctx context.Context, in *usersv1.CreateUserRequest, opts ...grpc.CallOption) (*usersv1.CreateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) CreateDevice(ctx context.Context, in *usersv1.CreateDeviceRequest, opts ...grpc.CallOption) (*usersv1.CreateDeviceResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) ListDevices(ctx context.Context, in *usersv1.ListDevicesRequest, opts ...grpc.CallOption) (*usersv1.ListDevicesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) DeleteDevice(ctx context.Context, in *usersv1.DeleteDeviceRequest, opts ...grpc.CallOption) (*usersv1.DeleteDeviceResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) DeleteUser(ctx context.Context, in *usersv1.DeleteUserRequest, opts ...grpc.CallOption) (*usersv1.DeleteUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) CreateAPIToken(ctx context.Context, in *usersv1.CreateAPITokenRequest, opts ...grpc.CallOption) (*usersv1.CreateAPITokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) ListAPITokens(ctx context.Context, in *usersv1.ListAPITokensRequest, opts ...grpc.CallOption) (*usersv1.ListAPITokensResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) RevokeAPIToken(ctx context.Context, in *usersv1.RevokeAPITokenRequest, opts ...grpc.CallOption) (*usersv1.RevokeAPITokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeUsersClient) ResolveAPIToken(ctx context.Context, in *usersv1.ResolveAPITokenRequest, opts ...grpc.CallOption) (*usersv1.ResolveAPITokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func buildUser(id string) *usersv1.User {
	return &usersv1.User{Meta: &usersv1.EntityMeta{Id: id}}
}

func assertResolvedIdentity(t *testing.T, resolved identity.ResolvedIdentity, expectedID string) {
	t.Helper()
	if resolved.IdentityID != expectedID {
		t.Fatalf("expected identity id %q, got %q", expectedID, resolved.IdentityID)
	}
	if resolved.IdentityType != identity.IdentityTypeUser {
		t.Fatalf("expected identity type %q, got %q", identity.IdentityTypeUser, resolved.IdentityType)
	}
}
