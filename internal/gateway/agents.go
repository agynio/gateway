package gateway

import (
	"context"

	"connectrpc.com/connect"
	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
)

func (g *Gateway) CreateAgent(ctx context.Context, req *connect.Request[agentsv1.CreateAgentRequest]) (*connect.Response[agentsv1.CreateAgentResponse], error) {
	resp, err := g.agents.CreateAgent(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetAgent(ctx context.Context, req *connect.Request[agentsv1.GetAgentRequest]) (*connect.Response[agentsv1.GetAgentResponse], error) {
	resp, err := g.agents.GetAgent(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateAgent(ctx context.Context, req *connect.Request[agentsv1.UpdateAgentRequest]) (*connect.Response[agentsv1.UpdateAgentResponse], error) {
	resp, err := g.agents.UpdateAgent(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteAgent(ctx context.Context, req *connect.Request[agentsv1.DeleteAgentRequest]) (*connect.Response[agentsv1.DeleteAgentResponse], error) {
	resp, err := g.agents.DeleteAgent(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListAgents(ctx context.Context, req *connect.Request[agentsv1.ListAgentsRequest]) (*connect.Response[agentsv1.ListAgentsResponse], error) {
	resp, err := g.agents.ListAgents(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) SetAgentRole(ctx context.Context, req *connect.Request[agentsv1.SetAgentRoleRequest]) (*connect.Response[agentsv1.SetAgentRoleResponse], error) {
	resp, err := g.agents.SetAgentRole(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) RemoveAgentRole(ctx context.Context, req *connect.Request[agentsv1.RemoveAgentRoleRequest]) (*connect.Response[agentsv1.RemoveAgentRoleResponse], error) {
	resp, err := g.agents.RemoveAgentRole(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListAgentRoles(ctx context.Context, req *connect.Request[agentsv1.ListAgentRolesRequest]) (*connect.Response[agentsv1.ListAgentRolesResponse], error) {
	resp, err := g.agents.ListAgentRoles(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListMyAgentRoles(ctx context.Context, req *connect.Request[agentsv1.ListMyAgentRolesRequest]) (*connect.Response[agentsv1.ListMyAgentRolesResponse], error) {
	resp, err := g.agents.ListMyAgentRoles(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateEnvironment(ctx context.Context, req *connect.Request[agentsv1.CreateEnvironmentRequest]) (*connect.Response[agentsv1.CreateEnvironmentResponse], error) {
	resp, err := g.agents.CreateEnvironment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetEnvironment(ctx context.Context, req *connect.Request[agentsv1.GetEnvironmentRequest]) (*connect.Response[agentsv1.GetEnvironmentResponse], error) {
	resp, err := g.agents.GetEnvironment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateEnvironment(ctx context.Context, req *connect.Request[agentsv1.UpdateEnvironmentRequest]) (*connect.Response[agentsv1.UpdateEnvironmentResponse], error) {
	resp, err := g.agents.UpdateEnvironment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteEnvironment(ctx context.Context, req *connect.Request[agentsv1.DeleteEnvironmentRequest]) (*connect.Response[agentsv1.DeleteEnvironmentResponse], error) {
	resp, err := g.agents.DeleteEnvironment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListEnvironments(ctx context.Context, req *connect.Request[agentsv1.ListEnvironmentsRequest]) (*connect.Response[agentsv1.ListEnvironmentsResponse], error) {
	resp, err := g.agents.ListEnvironments(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateSandbox(ctx context.Context, req *connect.Request[agentsv1.CreateSandboxRequest]) (*connect.Response[agentsv1.CreateSandboxResponse], error) {
	resp, err := g.agents.CreateSandbox(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetSandbox(ctx context.Context, req *connect.Request[agentsv1.GetSandboxRequest]) (*connect.Response[agentsv1.GetSandboxResponse], error) {
	resp, err := g.agents.GetSandbox(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListSandboxes(ctx context.Context, req *connect.Request[agentsv1.ListSandboxesRequest]) (*connect.Response[agentsv1.ListSandboxesResponse], error) {
	resp, err := g.agents.ListSandboxes(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) StopSandbox(ctx context.Context, req *connect.Request[agentsv1.StopSandboxRequest]) (*connect.Response[agentsv1.StopSandboxResponse], error) {
	resp, err := g.agents.StopSandbox(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteSandbox(ctx context.Context, req *connect.Request[agentsv1.DeleteSandboxRequest]) (*connect.Response[agentsv1.DeleteSandboxResponse], error) {
	resp, err := g.agents.DeleteSandbox(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) EnsureSandboxRunning(ctx context.Context, req *connect.Request[agentsv1.EnsureSandboxRunningRequest]) (*connect.Response[agentsv1.EnsureSandboxRunningResponse], error) {
	resp, err := g.agents.EnsureSandboxRunning(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateInstance(ctx context.Context, req *connect.Request[agentsv1.CreateInstanceRequest]) (*connect.Response[agentsv1.CreateInstanceResponse], error) {
	resp, err := g.agents.CreateInstance(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetInstance(ctx context.Context, req *connect.Request[agentsv1.GetInstanceRequest]) (*connect.Response[agentsv1.GetInstanceResponse], error) {
	resp, err := g.agents.GetInstance(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListInstances(ctx context.Context, req *connect.Request[agentsv1.ListInstancesRequest]) (*connect.Response[agentsv1.ListInstancesResponse], error) {
	resp, err := g.agents.ListInstances(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) PauseInstance(ctx context.Context, req *connect.Request[agentsv1.PauseInstanceRequest]) (*connect.Response[agentsv1.PauseInstanceResponse], error) {
	resp, err := g.agents.PauseInstance(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ResumeInstance(ctx context.Context, req *connect.Request[agentsv1.ResumeInstanceRequest]) (*connect.Response[agentsv1.ResumeInstanceResponse], error) {
	resp, err := g.agents.ResumeInstance(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteInstance(ctx context.Context, req *connect.Request[agentsv1.DeleteInstanceRequest]) (*connect.Response[agentsv1.DeleteInstanceResponse], error) {
	resp, err := g.agents.DeleteInstance(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) WriteInboxItem(ctx context.Context, req *connect.Request[agentsv1.WriteInboxItemRequest]) (*connect.Response[agentsv1.WriteInboxItemResponse], error) {
	resp, err := g.agents.WriteInboxItem(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetUnackedInboxItems(ctx context.Context, req *connect.Request[agentsv1.GetUnackedInboxItemsRequest]) (*connect.Response[agentsv1.GetUnackedInboxItemsResponse], error) {
	resp, err := g.agents.GetUnackedInboxItems(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) AckInboxItems(ctx context.Context, req *connect.Request[agentsv1.AckInboxItemsRequest]) (*connect.Response[agentsv1.AckInboxItemsResponse], error) {
	resp, err := g.agents.AckInboxItems(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetUnackedInboxCount(ctx context.Context, req *connect.Request[agentsv1.GetUnackedInboxCountRequest]) (*connect.Response[agentsv1.GetUnackedInboxCountResponse], error) {
	resp, err := g.agents.GetUnackedInboxCount(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateVolume(ctx context.Context, req *connect.Request[agentsv1.CreateVolumeRequest]) (*connect.Response[agentsv1.CreateVolumeResponse], error) {
	resp, err := g.agents.CreateVolume(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetVolume(ctx context.Context, req *connect.Request[agentsv1.GetVolumeRequest]) (*connect.Response[agentsv1.GetVolumeResponse], error) {
	resp, err := g.agents.GetVolume(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateVolume(ctx context.Context, req *connect.Request[agentsv1.UpdateVolumeRequest]) (*connect.Response[agentsv1.UpdateVolumeResponse], error) {
	resp, err := g.agents.UpdateVolume(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteVolume(ctx context.Context, req *connect.Request[agentsv1.DeleteVolumeRequest]) (*connect.Response[agentsv1.DeleteVolumeResponse], error) {
	resp, err := g.agents.DeleteVolume(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListVolumes(ctx context.Context, req *connect.Request[agentsv1.ListVolumesRequest]) (*connect.Response[agentsv1.ListVolumesResponse], error) {
	resp, err := g.agents.ListVolumes(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateVolumeAttachment(ctx context.Context, req *connect.Request[agentsv1.CreateVolumeAttachmentRequest]) (*connect.Response[agentsv1.CreateVolumeAttachmentResponse], error) {
	resp, err := g.agents.CreateVolumeAttachment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetVolumeAttachment(ctx context.Context, req *connect.Request[agentsv1.GetVolumeAttachmentRequest]) (*connect.Response[agentsv1.GetVolumeAttachmentResponse], error) {
	resp, err := g.agents.GetVolumeAttachment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteVolumeAttachment(ctx context.Context, req *connect.Request[agentsv1.DeleteVolumeAttachmentRequest]) (*connect.Response[agentsv1.DeleteVolumeAttachmentResponse], error) {
	resp, err := g.agents.DeleteVolumeAttachment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListVolumeAttachments(ctx context.Context, req *connect.Request[agentsv1.ListVolumeAttachmentsRequest]) (*connect.Response[agentsv1.ListVolumeAttachmentsResponse], error) {
	resp, err := g.agents.ListVolumeAttachments(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateMcp(ctx context.Context, req *connect.Request[agentsv1.CreateMcpRequest]) (*connect.Response[agentsv1.CreateMcpResponse], error) {
	resp, err := g.agents.CreateMcp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetMcp(ctx context.Context, req *connect.Request[agentsv1.GetMcpRequest]) (*connect.Response[agentsv1.GetMcpResponse], error) {
	resp, err := g.agents.GetMcp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateMcp(ctx context.Context, req *connect.Request[agentsv1.UpdateMcpRequest]) (*connect.Response[agentsv1.UpdateMcpResponse], error) {
	resp, err := g.agents.UpdateMcp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteMcp(ctx context.Context, req *connect.Request[agentsv1.DeleteMcpRequest]) (*connect.Response[agentsv1.DeleteMcpResponse], error) {
	resp, err := g.agents.DeleteMcp(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListMcps(ctx context.Context, req *connect.Request[agentsv1.ListMcpsRequest]) (*connect.Response[agentsv1.ListMcpsResponse], error) {
	resp, err := g.agents.ListMcps(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateSkill(ctx context.Context, req *connect.Request[agentsv1.CreateSkillRequest]) (*connect.Response[agentsv1.CreateSkillResponse], error) {
	resp, err := g.agents.CreateSkill(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetSkill(ctx context.Context, req *connect.Request[agentsv1.GetSkillRequest]) (*connect.Response[agentsv1.GetSkillResponse], error) {
	resp, err := g.agents.GetSkill(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateSkill(ctx context.Context, req *connect.Request[agentsv1.UpdateSkillRequest]) (*connect.Response[agentsv1.UpdateSkillResponse], error) {
	resp, err := g.agents.UpdateSkill(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteSkill(ctx context.Context, req *connect.Request[agentsv1.DeleteSkillRequest]) (*connect.Response[agentsv1.DeleteSkillResponse], error) {
	resp, err := g.agents.DeleteSkill(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListSkills(ctx context.Context, req *connect.Request[agentsv1.ListSkillsRequest]) (*connect.Response[agentsv1.ListSkillsResponse], error) {
	resp, err := g.agents.ListSkills(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateEnv(ctx context.Context, req *connect.Request[agentsv1.CreateEnvRequest]) (*connect.Response[agentsv1.CreateEnvResponse], error) {
	resp, err := g.agents.CreateEnv(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetEnv(ctx context.Context, req *connect.Request[agentsv1.GetEnvRequest]) (*connect.Response[agentsv1.GetEnvResponse], error) {
	resp, err := g.agents.GetEnv(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateEnv(ctx context.Context, req *connect.Request[agentsv1.UpdateEnvRequest]) (*connect.Response[agentsv1.UpdateEnvResponse], error) {
	resp, err := g.agents.UpdateEnv(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteEnv(ctx context.Context, req *connect.Request[agentsv1.DeleteEnvRequest]) (*connect.Response[agentsv1.DeleteEnvResponse], error) {
	resp, err := g.agents.DeleteEnv(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListEnvs(ctx context.Context, req *connect.Request[agentsv1.ListEnvsRequest]) (*connect.Response[agentsv1.ListEnvsResponse], error) {
	resp, err := g.agents.ListEnvs(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateInitScript(ctx context.Context, req *connect.Request[agentsv1.CreateInitScriptRequest]) (*connect.Response[agentsv1.CreateInitScriptResponse], error) {
	resp, err := g.agents.CreateInitScript(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetInitScript(ctx context.Context, req *connect.Request[agentsv1.GetInitScriptRequest]) (*connect.Response[agentsv1.GetInitScriptResponse], error) {
	resp, err := g.agents.GetInitScript(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateInitScript(ctx context.Context, req *connect.Request[agentsv1.UpdateInitScriptRequest]) (*connect.Response[agentsv1.UpdateInitScriptResponse], error) {
	resp, err := g.agents.UpdateInitScript(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteInitScript(ctx context.Context, req *connect.Request[agentsv1.DeleteInitScriptRequest]) (*connect.Response[agentsv1.DeleteInitScriptResponse], error) {
	resp, err := g.agents.DeleteInitScript(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListInitScripts(ctx context.Context, req *connect.Request[agentsv1.ListInitScriptsRequest]) (*connect.Response[agentsv1.ListInitScriptsResponse], error) {
	resp, err := g.agents.ListInitScripts(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

