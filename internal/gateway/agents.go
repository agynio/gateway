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

func (g *Gateway) CreateHook(ctx context.Context, req *connect.Request[agentsv1.CreateHookRequest]) (*connect.Response[agentsv1.CreateHookResponse], error) {
	resp, err := g.agents.CreateHook(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetHook(ctx context.Context, req *connect.Request[agentsv1.GetHookRequest]) (*connect.Response[agentsv1.GetHookResponse], error) {
	resp, err := g.agents.GetHook(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateHook(ctx context.Context, req *connect.Request[agentsv1.UpdateHookRequest]) (*connect.Response[agentsv1.UpdateHookResponse], error) {
	resp, err := g.agents.UpdateHook(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteHook(ctx context.Context, req *connect.Request[agentsv1.DeleteHookRequest]) (*connect.Response[agentsv1.DeleteHookResponse], error) {
	resp, err := g.agents.DeleteHook(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListHooks(ctx context.Context, req *connect.Request[agentsv1.ListHooksRequest]) (*connect.Response[agentsv1.ListHooksResponse], error) {
	resp, err := g.agents.ListHooks(ctx, req.Msg)
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

func (g *Gateway) CreateImagePullSecretAttachment(ctx context.Context, req *connect.Request[agentsv1.CreateImagePullSecretAttachmentRequest]) (*connect.Response[agentsv1.CreateImagePullSecretAttachmentResponse], error) {
	resp, err := g.agents.CreateImagePullSecretAttachment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetImagePullSecretAttachment(ctx context.Context, req *connect.Request[agentsv1.GetImagePullSecretAttachmentRequest]) (*connect.Response[agentsv1.GetImagePullSecretAttachmentResponse], error) {
	resp, err := g.agents.GetImagePullSecretAttachment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteImagePullSecretAttachment(ctx context.Context, req *connect.Request[agentsv1.DeleteImagePullSecretAttachmentRequest]) (*connect.Response[agentsv1.DeleteImagePullSecretAttachmentResponse], error) {
	resp, err := g.agents.DeleteImagePullSecretAttachment(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListImagePullSecretAttachments(ctx context.Context, req *connect.Request[agentsv1.ListImagePullSecretAttachmentsRequest]) (*connect.Response[agentsv1.ListImagePullSecretAttachmentsResponse], error) {
	resp, err := g.agents.ListImagePullSecretAttachments(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
