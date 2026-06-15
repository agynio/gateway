package gateway

import (
	"context"

	"connectrpc.com/connect"
	groupsv1 "github.com/agynio/gateway/gen/agynio/api/groups/v1"
)

type GroupsGateway struct {
	groups groupsv1.GroupsServiceClient
}

func NewGroupsGateway(groups groupsv1.GroupsServiceClient) *GroupsGateway {
	return &GroupsGateway{groups: groups}
}

func (g *GroupsGateway) CreateGroup(ctx context.Context, req *connect.Request[groupsv1.CreateGroupRequest]) (*connect.Response[groupsv1.CreateGroupResponse], error) {
	resp, err := g.groups.CreateGroup(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *GroupsGateway) GetGroup(ctx context.Context, req *connect.Request[groupsv1.GetGroupRequest]) (*connect.Response[groupsv1.GetGroupResponse], error) {
	resp, err := g.groups.GetGroup(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *GroupsGateway) ListGroups(ctx context.Context, req *connect.Request[groupsv1.ListGroupsRequest]) (*connect.Response[groupsv1.ListGroupsResponse], error) {
	resp, err := g.groups.ListGroups(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *GroupsGateway) UpdateGroup(ctx context.Context, req *connect.Request[groupsv1.UpdateGroupRequest]) (*connect.Response[groupsv1.UpdateGroupResponse], error) {
	resp, err := g.groups.UpdateGroup(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *GroupsGateway) DeleteGroup(ctx context.Context, req *connect.Request[groupsv1.DeleteGroupRequest]) (*connect.Response[groupsv1.DeleteGroupResponse], error) {
	resp, err := g.groups.DeleteGroup(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *GroupsGateway) AddMember(ctx context.Context, req *connect.Request[groupsv1.AddMemberRequest]) (*connect.Response[groupsv1.AddMemberResponse], error) {
	resp, err := g.groups.AddMember(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *GroupsGateway) RemoveMember(ctx context.Context, req *connect.Request[groupsv1.RemoveMemberRequest]) (*connect.Response[groupsv1.RemoveMemberResponse], error) {
	resp, err := g.groups.RemoveMember(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *GroupsGateway) ListMembers(ctx context.Context, req *connect.Request[groupsv1.ListMembersRequest]) (*connect.Response[groupsv1.ListMembersResponse], error) {
	resp, err := g.groups.ListMembers(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *GroupsGateway) ListMemberGroups(ctx context.Context, req *connect.Request[groupsv1.ListMemberGroupsRequest]) (*connect.Response[groupsv1.ListMemberGroupsResponse], error) {
	resp, err := g.groups.ListMemberGroups(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
