package handlers

import (
	"context"
	"fmt"
	"net/http"

	teamsv1 "github.com/agynio/gateway/gen/agynio/api/teams/v1"
	"github.com/agynio/gateway/internal/gen"
)

const (
	teamBasePath   = "/team/v1"
	defaultPerPage = 20
	maxPages       = 1000
)

func TeamBasePath() string {
	return teamBasePath
}

type Team struct {
	client teamsv1.TeamsServiceClient
}

func NewTeam(client teamsv1.TeamsServiceClient) *Team {
	if client == nil {
		panic("teams client is required")
	}
	return &Team{client: client}
}

func (t *Team) GetAgents(ctx context.Context, request gen.GetAgentsRequestObject) (gen.GetAgentsResponseObject, error) {
	resp, err := t.client.ListAgents(ctx, &teamsv1.ListAgentsRequest{
		PageSize:  pageSizeFromParam(request.Params.PageSize),
		PageToken: stringValue(request.Params.PageToken),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]gen.Agent, 0, len(resp.GetAgents()))
	for _, agent := range resp.GetAgents() {
		converted, err := agentFromProto(agent)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedAgents{
		Items:         items,
		NextPageToken: stringPtr(resp.GetNextPageToken()),
	}

	return gen.GetAgents200JSONResponse(payload), nil
}

func (t *Team) PostAgents(ctx context.Context, request gen.PostAgentsRequestObject) (gen.PostAgentsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.CreateAgent(ctx, agentCreateToProto(*request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	agent := resp.GetAgent()
	if agent == nil {
		return nil, responseProblem(fmt.Errorf("create agent response missing agent"))
	}

	converted, err := agentFromProto(agent)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostAgents201JSONResponse(converted), nil
}

func (t *Team) DeleteAgentsId(ctx context.Context, request gen.DeleteAgentsIdRequestObject) (gen.DeleteAgentsIdResponseObject, error) {
	_, err := t.client.DeleteAgent(ctx, &teamsv1.DeleteAgentRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteAgentsId204Response{}, nil
}

func (t *Team) GetAgentsId(ctx context.Context, request gen.GetAgentsIdRequestObject) (gen.GetAgentsIdResponseObject, error) {
	resp, err := t.client.GetAgent(ctx, &teamsv1.GetAgentRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	agent := resp.GetAgent()
	if agent == nil {
		return nil, responseProblem(fmt.Errorf("get agent response missing agent"))
	}

	converted, err := agentFromProto(agent)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetAgentsId200JSONResponse(converted), nil
}

func (t *Team) PatchAgentsId(ctx context.Context, request gen.PatchAgentsIdRequestObject) (gen.PatchAgentsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.UpdateAgent(ctx, agentUpdateToProto(request.Id, *request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	agent := resp.GetAgent()
	if agent == nil {
		return nil, responseProblem(fmt.Errorf("update agent response missing agent"))
	}

	converted, err := agentFromProto(agent)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchAgentsId200JSONResponse(converted), nil
}

func (t *Team) GetEnvs(ctx context.Context, request gen.GetEnvsRequestObject) (gen.GetEnvsResponseObject, error) {
	resp, err := t.client.ListEnvs(ctx, &teamsv1.ListEnvsRequest{
		PageSize:  pageSizeFromParam(request.Params.PageSize),
		PageToken: stringValue(request.Params.PageToken),
		AgentId:   uuidStringValue(request.Params.AgentId),
		McpId:     uuidStringValue(request.Params.McpId),
		HookId:    uuidStringValue(request.Params.HookId),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]gen.Env, 0, len(resp.GetEnvs()))
	for _, env := range resp.GetEnvs() {
		converted, err := envFromProto(env)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedEnvs{
		Items:         items,
		NextPageToken: stringPtr(resp.GetNextPageToken()),
	}

	return gen.GetEnvs200JSONResponse(payload), nil
}

func (t *Team) PostEnvs(ctx context.Context, request gen.PostEnvsRequestObject) (gen.PostEnvsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	createRequest, err := envCreateToProto(*request.Body)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := t.client.CreateEnv(ctx, createRequest)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	env := resp.GetEnv()
	if env == nil {
		return nil, responseProblem(fmt.Errorf("create env response missing env"))
	}

	converted, err := envFromProto(env)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostEnvs201JSONResponse(converted), nil
}

func (t *Team) DeleteEnvsId(ctx context.Context, request gen.DeleteEnvsIdRequestObject) (gen.DeleteEnvsIdResponseObject, error) {
	_, err := t.client.DeleteEnv(ctx, &teamsv1.DeleteEnvRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteEnvsId204Response{}, nil
}

func (t *Team) GetEnvsId(ctx context.Context, request gen.GetEnvsIdRequestObject) (gen.GetEnvsIdResponseObject, error) {
	resp, err := t.client.GetEnv(ctx, &teamsv1.GetEnvRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	env := resp.GetEnv()
	if env == nil {
		return nil, responseProblem(fmt.Errorf("get env response missing env"))
	}

	converted, err := envFromProto(env)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetEnvsId200JSONResponse(converted), nil
}

func (t *Team) PatchEnvsId(ctx context.Context, request gen.PatchEnvsIdRequestObject) (gen.PatchEnvsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	updateRequest, err := envUpdateToProto(request.Id, *request.Body)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := t.client.UpdateEnv(ctx, updateRequest)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	env := resp.GetEnv()
	if env == nil {
		return nil, responseProblem(fmt.Errorf("update env response missing env"))
	}

	converted, err := envFromProto(env)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchEnvsId200JSONResponse(converted), nil
}

func (t *Team) GetHooks(ctx context.Context, request gen.GetHooksRequestObject) (gen.GetHooksResponseObject, error) {
	resp, err := t.client.ListHooks(ctx, &teamsv1.ListHooksRequest{
		PageSize:  pageSizeFromParam(request.Params.PageSize),
		PageToken: stringValue(request.Params.PageToken),
		AgentId:   uuidStringValue(request.Params.AgentId),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]gen.Hook, 0, len(resp.GetHooks()))
	for _, hook := range resp.GetHooks() {
		converted, err := hookFromProto(hook)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedHooks{
		Items:         items,
		NextPageToken: stringPtr(resp.GetNextPageToken()),
	}

	return gen.GetHooks200JSONResponse(payload), nil
}

func (t *Team) PostHooks(ctx context.Context, request gen.PostHooksRequestObject) (gen.PostHooksResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.CreateHook(ctx, hookCreateToProto(*request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	hook := resp.GetHook()
	if hook == nil {
		return nil, responseProblem(fmt.Errorf("create hook response missing hook"))
	}

	converted, err := hookFromProto(hook)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostHooks201JSONResponse(converted), nil
}

func (t *Team) DeleteHooksId(ctx context.Context, request gen.DeleteHooksIdRequestObject) (gen.DeleteHooksIdResponseObject, error) {
	_, err := t.client.DeleteHook(ctx, &teamsv1.DeleteHookRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteHooksId204Response{}, nil
}

func (t *Team) GetHooksId(ctx context.Context, request gen.GetHooksIdRequestObject) (gen.GetHooksIdResponseObject, error) {
	resp, err := t.client.GetHook(ctx, &teamsv1.GetHookRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	hook := resp.GetHook()
	if hook == nil {
		return nil, responseProblem(fmt.Errorf("get hook response missing hook"))
	}

	converted, err := hookFromProto(hook)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetHooksId200JSONResponse(converted), nil
}

func (t *Team) PatchHooksId(ctx context.Context, request gen.PatchHooksIdRequestObject) (gen.PatchHooksIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.UpdateHook(ctx, hookUpdateToProto(request.Id, *request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	hook := resp.GetHook()
	if hook == nil {
		return nil, responseProblem(fmt.Errorf("update hook response missing hook"))
	}

	converted, err := hookFromProto(hook)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchHooksId200JSONResponse(converted), nil
}

func (t *Team) GetInitScripts(ctx context.Context, request gen.GetInitScriptsRequestObject) (gen.GetInitScriptsResponseObject, error) {
	resp, err := t.client.ListInitScripts(ctx, &teamsv1.ListInitScriptsRequest{
		PageSize:  pageSizeFromParam(request.Params.PageSize),
		PageToken: stringValue(request.Params.PageToken),
		AgentId:   uuidStringValue(request.Params.AgentId),
		McpId:     uuidStringValue(request.Params.McpId),
		HookId:    uuidStringValue(request.Params.HookId),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]gen.InitScript, 0, len(resp.GetInitScripts()))
	for _, script := range resp.GetInitScripts() {
		converted, err := initScriptFromProto(script)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedInitScripts{
		Items:         items,
		NextPageToken: stringPtr(resp.GetNextPageToken()),
	}

	return gen.GetInitScripts200JSONResponse(payload), nil
}

func (t *Team) PostInitScripts(ctx context.Context, request gen.PostInitScriptsRequestObject) (gen.PostInitScriptsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	createRequest, err := initScriptCreateToProto(*request.Body)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := t.client.CreateInitScript(ctx, createRequest)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	script := resp.GetInitScript()
	if script == nil {
		return nil, responseProblem(fmt.Errorf("create init script response missing init script"))
	}

	converted, err := initScriptFromProto(script)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostInitScripts201JSONResponse(converted), nil
}

func (t *Team) DeleteInitScriptsId(ctx context.Context, request gen.DeleteInitScriptsIdRequestObject) (gen.DeleteInitScriptsIdResponseObject, error) {
	_, err := t.client.DeleteInitScript(ctx, &teamsv1.DeleteInitScriptRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteInitScriptsId204Response{}, nil
}

func (t *Team) GetInitScriptsId(ctx context.Context, request gen.GetInitScriptsIdRequestObject) (gen.GetInitScriptsIdResponseObject, error) {
	resp, err := t.client.GetInitScript(ctx, &teamsv1.GetInitScriptRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	script := resp.GetInitScript()
	if script == nil {
		return nil, responseProblem(fmt.Errorf("get init script response missing init script"))
	}

	converted, err := initScriptFromProto(script)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetInitScriptsId200JSONResponse(converted), nil
}

func (t *Team) PatchInitScriptsId(ctx context.Context, request gen.PatchInitScriptsIdRequestObject) (gen.PatchInitScriptsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.UpdateInitScript(ctx, initScriptUpdateToProto(request.Id, *request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	script := resp.GetInitScript()
	if script == nil {
		return nil, responseProblem(fmt.Errorf("update init script response missing init script"))
	}

	converted, err := initScriptFromProto(script)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchInitScriptsId200JSONResponse(converted), nil
}

func (t *Team) GetMcps(ctx context.Context, request gen.GetMcpsRequestObject) (gen.GetMcpsResponseObject, error) {
	resp, err := t.client.ListMcps(ctx, &teamsv1.ListMcpsRequest{
		PageSize:  pageSizeFromParam(request.Params.PageSize),
		PageToken: stringValue(request.Params.PageToken),
		AgentId:   uuidStringValue(request.Params.AgentId),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]gen.Mcp, 0, len(resp.GetMcps()))
	for _, mcp := range resp.GetMcps() {
		converted, err := mcpFromProto(mcp)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedMcps{
		Items:         items,
		NextPageToken: stringPtr(resp.GetNextPageToken()),
	}

	return gen.GetMcps200JSONResponse(payload), nil
}

func (t *Team) PostMcps(ctx context.Context, request gen.PostMcpsRequestObject) (gen.PostMcpsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.CreateMcp(ctx, mcpCreateToProto(*request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	mcp := resp.GetMcp()
	if mcp == nil {
		return nil, responseProblem(fmt.Errorf("create mcp response missing mcp"))
	}

	converted, err := mcpFromProto(mcp)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostMcps201JSONResponse(converted), nil
}

func (t *Team) DeleteMcpsId(ctx context.Context, request gen.DeleteMcpsIdRequestObject) (gen.DeleteMcpsIdResponseObject, error) {
	_, err := t.client.DeleteMcp(ctx, &teamsv1.DeleteMcpRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteMcpsId204Response{}, nil
}

func (t *Team) GetMcpsId(ctx context.Context, request gen.GetMcpsIdRequestObject) (gen.GetMcpsIdResponseObject, error) {
	resp, err := t.client.GetMcp(ctx, &teamsv1.GetMcpRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	mcp := resp.GetMcp()
	if mcp == nil {
		return nil, responseProblem(fmt.Errorf("get mcp response missing mcp"))
	}

	converted, err := mcpFromProto(mcp)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetMcpsId200JSONResponse(converted), nil
}

func (t *Team) PatchMcpsId(ctx context.Context, request gen.PatchMcpsIdRequestObject) (gen.PatchMcpsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.UpdateMcp(ctx, mcpUpdateToProto(request.Id, *request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	mcp := resp.GetMcp()
	if mcp == nil {
		return nil, responseProblem(fmt.Errorf("update mcp response missing mcp"))
	}

	converted, err := mcpFromProto(mcp)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchMcpsId200JSONResponse(converted), nil
}

func (t *Team) GetSkills(ctx context.Context, request gen.GetSkillsRequestObject) (gen.GetSkillsResponseObject, error) {
	resp, err := t.client.ListSkills(ctx, &teamsv1.ListSkillsRequest{
		PageSize:  pageSizeFromParam(request.Params.PageSize),
		PageToken: stringValue(request.Params.PageToken),
		AgentId:   uuidStringValue(request.Params.AgentId),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]gen.Skill, 0, len(resp.GetSkills()))
	for _, skill := range resp.GetSkills() {
		converted, err := skillFromProto(skill)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedSkills{
		Items:         items,
		NextPageToken: stringPtr(resp.GetNextPageToken()),
	}

	return gen.GetSkills200JSONResponse(payload), nil
}

func (t *Team) PostSkills(ctx context.Context, request gen.PostSkillsRequestObject) (gen.PostSkillsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.CreateSkill(ctx, skillCreateToProto(*request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	skill := resp.GetSkill()
	if skill == nil {
		return nil, responseProblem(fmt.Errorf("create skill response missing skill"))
	}

	converted, err := skillFromProto(skill)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostSkills201JSONResponse(converted), nil
}

func (t *Team) DeleteSkillsId(ctx context.Context, request gen.DeleteSkillsIdRequestObject) (gen.DeleteSkillsIdResponseObject, error) {
	_, err := t.client.DeleteSkill(ctx, &teamsv1.DeleteSkillRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteSkillsId204Response{}, nil
}

func (t *Team) GetSkillsId(ctx context.Context, request gen.GetSkillsIdRequestObject) (gen.GetSkillsIdResponseObject, error) {
	resp, err := t.client.GetSkill(ctx, &teamsv1.GetSkillRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	skill := resp.GetSkill()
	if skill == nil {
		return nil, responseProblem(fmt.Errorf("get skill response missing skill"))
	}

	converted, err := skillFromProto(skill)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetSkillsId200JSONResponse(converted), nil
}

func (t *Team) PatchSkillsId(ctx context.Context, request gen.PatchSkillsIdRequestObject) (gen.PatchSkillsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.UpdateSkill(ctx, skillUpdateToProto(request.Id, *request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	skill := resp.GetSkill()
	if skill == nil {
		return nil, responseProblem(fmt.Errorf("update skill response missing skill"))
	}

	converted, err := skillFromProto(skill)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchSkillsId200JSONResponse(converted), nil
}

func (t *Team) GetVolumeAttachments(ctx context.Context, request gen.GetVolumeAttachmentsRequestObject) (gen.GetVolumeAttachmentsResponseObject, error) {
	resp, err := t.client.ListVolumeAttachments(ctx, &teamsv1.ListVolumeAttachmentsRequest{
		PageSize:  pageSizeFromParam(request.Params.PageSize),
		PageToken: stringValue(request.Params.PageToken),
		VolumeId:  uuidStringValue(request.Params.VolumeId),
		AgentId:   uuidStringValue(request.Params.AgentId),
		McpId:     uuidStringValue(request.Params.McpId),
		HookId:    uuidStringValue(request.Params.HookId),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]gen.VolumeAttachment, 0, len(resp.GetVolumeAttachments()))
	for _, attachment := range resp.GetVolumeAttachments() {
		converted, err := volumeAttachmentFromProto(attachment)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedVolumeAttachments{
		Items:         items,
		NextPageToken: stringPtr(resp.GetNextPageToken()),
	}

	return gen.GetVolumeAttachments200JSONResponse(payload), nil
}

func (t *Team) PostVolumeAttachments(ctx context.Context, request gen.PostVolumeAttachmentsRequestObject) (gen.PostVolumeAttachmentsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	createRequest, err := volumeAttachmentCreateToProto(*request.Body)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := t.client.CreateVolumeAttachment(ctx, createRequest)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	attachment := resp.GetVolumeAttachment()
	if attachment == nil {
		return nil, responseProblem(fmt.Errorf("create volume attachment response missing attachment"))
	}

	converted, err := volumeAttachmentFromProto(attachment)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostVolumeAttachments201JSONResponse(converted), nil
}

func (t *Team) DeleteVolumeAttachmentsId(ctx context.Context, request gen.DeleteVolumeAttachmentsIdRequestObject) (gen.DeleteVolumeAttachmentsIdResponseObject, error) {
	_, err := t.client.DeleteVolumeAttachment(ctx, &teamsv1.DeleteVolumeAttachmentRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteVolumeAttachmentsId204Response{}, nil
}

func (t *Team) GetVolumeAttachmentsId(ctx context.Context, request gen.GetVolumeAttachmentsIdRequestObject) (gen.GetVolumeAttachmentsIdResponseObject, error) {
	resp, err := t.client.GetVolumeAttachment(ctx, &teamsv1.GetVolumeAttachmentRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	attachment := resp.GetVolumeAttachment()
	if attachment == nil {
		return nil, responseProblem(fmt.Errorf("get volume attachment response missing attachment"))
	}

	converted, err := volumeAttachmentFromProto(attachment)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetVolumeAttachmentsId200JSONResponse(converted), nil
}

func (t *Team) GetVolumes(ctx context.Context, request gen.GetVolumesRequestObject) (gen.GetVolumesResponseObject, error) {
	resp, err := t.client.ListVolumes(ctx, &teamsv1.ListVolumesRequest{
		PageSize:  pageSizeFromParam(request.Params.PageSize),
		PageToken: stringValue(request.Params.PageToken),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]gen.Volume, 0, len(resp.GetVolumes()))
	for _, volume := range resp.GetVolumes() {
		converted, err := volumeFromProto(volume)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedVolumes{
		Items:         items,
		NextPageToken: stringPtr(resp.GetNextPageToken()),
	}

	return gen.GetVolumes200JSONResponse(payload), nil
}

func (t *Team) PostVolumes(ctx context.Context, request gen.PostVolumesRequestObject) (gen.PostVolumesResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.CreateVolume(ctx, volumeCreateToProto(*request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	volume := resp.GetVolume()
	if volume == nil {
		return nil, responseProblem(fmt.Errorf("create volume response missing volume"))
	}

	converted, err := volumeFromProto(volume)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostVolumes201JSONResponse(converted), nil
}

func (t *Team) DeleteVolumesId(ctx context.Context, request gen.DeleteVolumesIdRequestObject) (gen.DeleteVolumesIdResponseObject, error) {
	_, err := t.client.DeleteVolume(ctx, &teamsv1.DeleteVolumeRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteVolumesId204Response{}, nil
}

func (t *Team) GetVolumesId(ctx context.Context, request gen.GetVolumesIdRequestObject) (gen.GetVolumesIdResponseObject, error) {
	resp, err := t.client.GetVolume(ctx, &teamsv1.GetVolumeRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	volume := resp.GetVolume()
	if volume == nil {
		return nil, responseProblem(fmt.Errorf("get volume response missing volume"))
	}

	converted, err := volumeFromProto(volume)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetVolumesId200JSONResponse(converted), nil
}

func (t *Team) PatchVolumesId(ctx context.Context, request gen.PatchVolumesIdRequestObject) (gen.PatchVolumesIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := t.client.UpdateVolume(ctx, volumeUpdateToProto(request.Id, *request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	volume := resp.GetVolume()
	if volume == nil {
		return nil, responseProblem(fmt.Errorf("update volume response missing volume"))
	}

	converted, err := volumeFromProto(volume)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchVolumesId200JSONResponse(converted), nil
}

func pageSizeFromParam(pageSize *int) int32 {
	if pageSize == nil || *pageSize < 1 {
		return int32(defaultPerPage)
	}
	return int32(*pageSize)
}

func normalizePagination(pagePtr, perPagePtr *int) (int, int) {
	page := 1
	perPage := defaultPerPage
	if pagePtr != nil {
		page = *pagePtr
	}
	if perPagePtr != nil {
		perPage = *perPagePtr
	}
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = defaultPerPage
	}
	return page, perPage
}

func paginateSlice[T any](items []T, page, perPage int) []T {
	if len(items) == 0 {
		return []T{}
	}
	start := (page - 1) * perPage
	if start >= len(items) {
		return []T{}
	}
	end := start + perPage
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func listAll[T any](ctx context.Context, pageSize int32, fetch func(context.Context, string, int32) ([]T, string, error)) ([]T, error) {
	var items []T
	pageToken := ""
	for page := 0; page < maxPages; page++ {
		batch, nextToken, err := fetch(ctx, pageToken, pageSize)
		if err != nil {
			return nil, err
		}
		items = append(items, batch...)
		if nextToken == "" {
			if len(items) == 0 {
				return []T{}, nil
			}
			return items, nil
		}
		if nextToken == pageToken {
			return nil, fmt.Errorf("pagination token did not advance")
		}
		pageToken = nextToken
	}
	return nil, fmt.Errorf("pagination exceeded %d pages", maxPages)
}

func requestProblem(err error) *ProblemError {
	problem := NewProblem(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
	return NewProblemError(problem, err)
}

func responseProblem(err error) *ProblemError {
	problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), err.Error())
	return NewProblemError(problem, err)
}
