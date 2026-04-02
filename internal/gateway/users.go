package gateway

import (
	"context"

	"connectrpc.com/connect"
	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
)

type UsersGateway struct {
	users usersv1.UsersServiceClient
}

func NewUsersGateway(users usersv1.UsersServiceClient) *UsersGateway {
	return &UsersGateway{users: users}
}

func (g *UsersGateway) GetMe(ctx context.Context, req *connect.Request[usersv1.GetMeRequest]) (*connect.Response[usersv1.GetMeResponse], error) {
	resp, err := g.users.GetMe(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *UsersGateway) CreateUser(ctx context.Context, req *connect.Request[usersv1.CreateUserRequest]) (*connect.Response[usersv1.CreateUserResponse], error) {
	resp, err := g.users.CreateUser(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
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
