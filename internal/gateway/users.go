package gateway

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersGateway struct {
	users usersv1.UsersServiceClient
}

func NewUsersGateway(users usersv1.UsersServiceClient) *UsersGateway {
	return &UsersGateway{users: users}
}

func (g *UsersGateway) GetMe(ctx context.Context, req *connect.Request[usersv1.GetMeRequest]) (*connect.Response[usersv1.GetMeResponse], error) {
	_, ok := identity.IdentityFromContext(ctx)
	if !ok {
		return nil, toConnectError(status.Error(codes.Unauthenticated, "identity not available"))
	}

	userResp, err := g.users.GetMe(ctx, &usersv1.GetMeRequest{})
	if err != nil {
		return nil, toConnectError(err)
	}

	clusterRole := userResp.ClusterRole
	if clusterRole == usersv1.ClusterRole_CLUSTER_ROLE_UNSPECIFIED {
		clusterRole = usersv1.ClusterRole_CLUSTER_ROLE_ADMIN
	}

	resp := &usersv1.GetMeResponse{
		User:        userResp.User,
		ClusterRole: clusterRole,
	}
	return connect.NewResponse(resp), nil
}

func (g *UsersGateway) CreateUser(ctx context.Context, req *connect.Request[usersv1.CreateUserRequest]) (*connect.Response[usersv1.CreateUserResponse], error) {
	oidcSubject := strings.TrimSpace(req.Msg.GetOidcSubject())
	if oidcSubject == "" {
		return nil, toConnectError(status.Error(codes.InvalidArgument, "oidc_subject is required"))
	}

	name := req.Msg.GetName()
	createResp, err := g.users.ResolveOrCreateUser(ctx, &usersv1.ResolveOrCreateUserRequest{
		OidcSubject: oidcSubject,
		Name:        name,
		Email:       name,
		PhotoUrl:    req.Msg.GetPhotoUrl(),
	})
	if err != nil {
		return nil, toConnectError(err)
	}
	if createResp == nil || createResp.User == nil {
		return nil, toConnectError(status.Error(codes.Internal, "user missing from response"))
	}

	resp := &usersv1.CreateUserResponse{User: createResp.User}
	return connect.NewResponse(resp), nil
}

func (g *UsersGateway) GetUser(ctx context.Context, req *connect.Request[usersv1.GetUserRequest]) (*connect.Response[usersv1.GetUserResponse], error) {
	resp, err := g.users.GetUser(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *UsersGateway) ListUsers(ctx context.Context, req *connect.Request[usersv1.ListUsersRequest]) (*connect.Response[usersv1.ListUsersResponse], error) {
	resp, err := g.users.ListUsers(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *UsersGateway) UpdateUser(ctx context.Context, req *connect.Request[usersv1.UpdateUserRequest]) (*connect.Response[usersv1.UpdateUserResponse], error) {
	resp, err := g.users.UpdateUser(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *UsersGateway) DeleteUser(ctx context.Context, req *connect.Request[usersv1.DeleteUserRequest]) (*connect.Response[usersv1.DeleteUserResponse], error) {
	resp, err := g.users.DeleteUser(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *UsersGateway) BatchGetUsers(ctx context.Context, req *connect.Request[usersv1.BatchGetUsersRequest]) (*connect.Response[usersv1.BatchGetUsersResponse], error) {
	resp, err := g.users.BatchGetUsers(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	for _, u := range resp.Users {
		u.OidcSubject = ""
	}
	return connect.NewResponse(resp), nil
}

func (g *UsersGateway) CreateAPIToken(ctx context.Context, req *connect.Request[usersv1.CreateAPITokenRequest]) (*connect.Response[usersv1.CreateAPITokenResponse], error) {
	resp, err := g.users.CreateAPIToken(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *UsersGateway) ListAPITokens(ctx context.Context, req *connect.Request[usersv1.ListAPITokensRequest]) (*connect.Response[usersv1.ListAPITokensResponse], error) {
	resp, err := g.users.ListAPITokens(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *UsersGateway) RevokeAPIToken(ctx context.Context, req *connect.Request[usersv1.RevokeAPITokenRequest]) (*connect.Response[usersv1.RevokeAPITokenResponse], error) {
	resp, err := g.users.RevokeAPIToken(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
