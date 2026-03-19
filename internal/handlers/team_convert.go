package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	teamsv1 "github.com/agynio/gateway/gen/agynio/api/teams/v1"
	"github.com/agynio/gateway/internal/gen"
)

type entityMeta struct {
	id        openapi_types.UUID
	createdAt time.Time
	updatedAt *time.Time
}

type targetKind int

const (
	targetNone targetKind = iota
	targetAgent
	targetMcp
	targetHook
)

type sourceKind int

const (
	sourceNone sourceKind = iota
	sourceValue
	sourceSecret
)

func metaFromProto(meta *teamsv1.EntityMeta) (entityMeta, error) {
	if meta == nil {
		return entityMeta{}, fmt.Errorf("metadata missing")
	}

	parsedID, err := uuid.Parse(meta.GetId())
	if err != nil {
		return entityMeta{}, fmt.Errorf("parse id: %w", err)
	}

	createdAt := meta.GetCreatedAt()
	if createdAt == nil {
		return entityMeta{}, fmt.Errorf("created_at missing")
	}

	var updatedAt *time.Time
	if updated := meta.GetUpdatedAt(); updated != nil {
		value := updated.AsTime().UTC()
		updatedAt = &value
	}

	return entityMeta{
		id:        openapi_types.UUID(parsedID),
		createdAt: createdAt.AsTime().UTC(),
		updatedAt: updatedAt,
	}, nil
}

func parseUUID(value string) (openapi_types.UUID, error) {
	if value == "" {
		return openapi_types.UUID{}, fmt.Errorf("uuid missing")
	}
	parsed, err := uuid.Parse(value)
	if err != nil {
		return openapi_types.UUID{}, fmt.Errorf("parse uuid: %w", err)
	}
	return openapi_types.UUID(parsed), nil
}

func uuidPtrFromString(value string) (*openapi_types.UUID, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("parse uuid: %w", err)
	}
	result := openapi_types.UUID(parsed)
	return &result, nil
}

func uuidStringValue(value *openapi_types.UUID) string {
	if value == nil {
		return ""
	}
	return value.String()
}

func computeResourcesFromProto(resources *teamsv1.ComputeResources) *gen.ComputeResources {
	if resources == nil {
		return nil
	}
	return &gen.ComputeResources{
		LimitsCpu:      stringPtr(resources.GetLimitsCpu()),
		LimitsMemory:   stringPtr(resources.GetLimitsMemory()),
		RequestsCpu:    stringPtr(resources.GetRequestsCpu()),
		RequestsMemory: stringPtr(resources.GetRequestsMemory()),
	}
}

func computeResourcesToProto(resources *gen.ComputeResources) *teamsv1.ComputeResources {
	if resources == nil {
		return nil
	}
	return &teamsv1.ComputeResources{
		LimitsCpu:      stringValue(resources.LimitsCpu),
		LimitsMemory:   stringValue(resources.LimitsMemory),
		RequestsCpu:    stringValue(resources.RequestsCpu),
		RequestsMemory: stringValue(resources.RequestsMemory),
	}
}

func agentFromProto(agent *teamsv1.Agent) (gen.Agent, error) {
	if agent == nil {
		return gen.Agent{}, fmt.Errorf("agent missing")
	}

	meta, err := metaFromProto(agent.GetMeta())
	if err != nil {
		return gen.Agent{}, err
	}

	modelID, err := parseUUID(agent.GetModel())
	if err != nil {
		return gen.Agent{}, fmt.Errorf("parse model id: %w", err)
	}

	result := gen.Agent{
		Configuration: agent.GetConfiguration(),
		CreatedAt:     meta.createdAt,
		Id:            meta.id,
		Image:         agent.GetImage(),
		Model:         modelID,
		Name:          agent.GetName(),
		Resources:     computeResourcesFromProto(agent.GetResources()),
		Role:          agent.GetRole(),
		UpdatedAt:     meta.updatedAt,
	}

	if agent.Description != "" {
		result.Description = stringPtr(agent.Description)
	}

	return result, nil
}

func agentCreateToProto(request gen.AgentCreateRequest) *teamsv1.CreateAgentRequest {
	return &teamsv1.CreateAgentRequest{
		Name:          request.Name,
		Role:          request.Role,
		Model:         request.Model.String(),
		Description:   stringValue(request.Description),
		Configuration: stringValue(request.Configuration),
		Image:         request.Image,
		Resources:     computeResourcesToProto(request.Resources),
	}
}

func agentUpdateToProto(id openapi_types.UUID, request gen.AgentUpdateRequest) *teamsv1.UpdateAgentRequest {
	var model *string
	if request.Model != nil {
		value := request.Model.String()
		model = &value
	}

	return &teamsv1.UpdateAgentRequest{
		Id:            id.String(),
		Name:          request.Name,
		Role:          request.Role,
		Model:         model,
		Description:   request.Description,
		Configuration: request.Configuration,
		Image:         request.Image,
		Resources:     computeResourcesToProto(request.Resources),
	}
}

func envFromProto(env *teamsv1.Env) (gen.Env, error) {
	if env == nil {
		return gen.Env{}, fmt.Errorf("env missing")
	}

	meta, err := metaFromProto(env.GetMeta())
	if err != nil {
		return gen.Env{}, err
	}

	agentID, mcpID, hookID, err := parseTargetIDs(env.GetAgentId(), env.GetMcpId(), env.GetHookId())
	if err != nil {
		return gen.Env{}, err
	}

	result := gen.Env{
		AgentId:   agentID,
		CreatedAt: meta.createdAt,
		HookId:    hookID,
		Id:        meta.id,
		McpId:     mcpID,
		Name:      env.GetName(),
		UpdatedAt: meta.updatedAt,
	}

	if env.Description != "" {
		result.Description = stringPtr(env.Description)
	}

	value := env.GetValue()
	secretID := env.GetSecretId()
	if value != "" && secretID != "" {
		return gen.Env{}, fmt.Errorf("env source has both value and secret")
	}
	if value != "" {
		result.Value = stringPtr(value)
	}
	if secretID != "" {
		parsedSecret, err := uuidPtrFromString(secretID)
		if err != nil {
			return gen.Env{}, fmt.Errorf("parse secret id: %w", err)
		}
		result.SecretId = parsedSecret
	}

	return result, nil
}

func envCreateToProto(request gen.EnvCreateRequest) (*teamsv1.CreateEnvRequest, error) {
	createRequest := &teamsv1.CreateEnvRequest{
		Name:        request.Name,
		Description: stringValue(request.Description),
	}

	kind, value, err := resolveTargetKind(request.AgentId, request.McpId, request.HookId)
	if err != nil {
		return nil, err
	}
	if kind == targetNone {
		return nil, fmt.Errorf("env target missing")
	}
	switch kind {
	case targetAgent:
		createRequest.Target = &teamsv1.CreateEnvRequest_AgentId{AgentId: value}
	case targetMcp:
		createRequest.Target = &teamsv1.CreateEnvRequest_McpId{McpId: value}
	case targetHook:
		createRequest.Target = &teamsv1.CreateEnvRequest_HookId{HookId: value}
	}

	sourceKindValue, sourceVal, err := resolveSourceKind(request.Value, request.SecretId)
	if err != nil {
		return nil, err
	}
	if sourceKindValue == sourceNone {
		return nil, fmt.Errorf("env source missing")
	}
	switch sourceKindValue {
	case sourceValue:
		createRequest.Source = &teamsv1.CreateEnvRequest_Value{Value: sourceVal}
	case sourceSecret:
		createRequest.Source = &teamsv1.CreateEnvRequest_SecretId{SecretId: sourceVal}
	}

	return createRequest, nil
}

func envUpdateToProto(id openapi_types.UUID, request gen.EnvUpdateRequest) (*teamsv1.UpdateEnvRequest, error) {
	if request.Value != nil && request.SecretId != nil {
		return nil, fmt.Errorf("env update includes value and secret")
	}

	update := &teamsv1.UpdateEnvRequest{
		Id:          id.String(),
		Name:        request.Name,
		Description: request.Description,
		Value:       request.Value,
	}
	if request.SecretId != nil {
		secret := request.SecretId.String()
		update.SecretId = &secret
	}

	return update, nil
}

func hookFromProto(hook *teamsv1.Hook) (gen.Hook, error) {
	if hook == nil {
		return gen.Hook{}, fmt.Errorf("hook missing")
	}

	meta, err := metaFromProto(hook.GetMeta())
	if err != nil {
		return gen.Hook{}, err
	}

	agentID, err := parseUUID(hook.GetAgentId())
	if err != nil {
		return gen.Hook{}, fmt.Errorf("parse agent id: %w", err)
	}

	result := gen.Hook{
		AgentId:   agentID,
		CreatedAt: meta.createdAt,
		Event:     hook.GetEvent(),
		Function:  hook.GetFunction(),
		Id:        meta.id,
		Image:     hook.GetImage(),
		Resources: computeResourcesFromProto(hook.GetResources()),
		UpdatedAt: meta.updatedAt,
	}

	if hook.Description != "" {
		result.Description = stringPtr(hook.Description)
	}

	return result, nil
}

func hookCreateToProto(request gen.HookCreateRequest) *teamsv1.CreateHookRequest {
	return &teamsv1.CreateHookRequest{
		AgentId:     request.AgentId.String(),
		Event:       request.Event,
		Function:    request.Function,
		Image:       request.Image,
		Resources:   computeResourcesToProto(request.Resources),
		Description: stringValue(request.Description),
	}
}

func hookUpdateToProto(id openapi_types.UUID, request gen.HookUpdateRequest) *teamsv1.UpdateHookRequest {
	return &teamsv1.UpdateHookRequest{
		Id:          id.String(),
		Event:       request.Event,
		Function:    request.Function,
		Image:       request.Image,
		Resources:   computeResourcesToProto(request.Resources),
		Description: request.Description,
	}
}

func initScriptFromProto(script *teamsv1.InitScript) (gen.InitScript, error) {
	if script == nil {
		return gen.InitScript{}, fmt.Errorf("init script missing")
	}

	meta, err := metaFromProto(script.GetMeta())
	if err != nil {
		return gen.InitScript{}, err
	}

	agentID, mcpID, hookID, err := parseTargetIDs(script.GetAgentId(), script.GetMcpId(), script.GetHookId())
	if err != nil {
		return gen.InitScript{}, err
	}

	result := gen.InitScript{
		AgentId:   agentID,
		CreatedAt: meta.createdAt,
		HookId:    hookID,
		Id:        meta.id,
		McpId:     mcpID,
		Script:    script.GetScript(),
		UpdatedAt: meta.updatedAt,
	}

	if script.Description != "" {
		result.Description = stringPtr(script.Description)
	}

	return result, nil
}

func initScriptCreateToProto(request gen.InitScriptCreateRequest) (*teamsv1.CreateInitScriptRequest, error) {
	createRequest := &teamsv1.CreateInitScriptRequest{
		Script:      request.Script,
		Description: stringValue(request.Description),
	}

	kind, value, err := resolveTargetKind(request.AgentId, request.McpId, request.HookId)
	if err != nil {
		return nil, err
	}
	if kind == targetNone {
		return nil, fmt.Errorf("init script target missing")
	}
	switch kind {
	case targetAgent:
		createRequest.Target = &teamsv1.CreateInitScriptRequest_AgentId{AgentId: value}
	case targetMcp:
		createRequest.Target = &teamsv1.CreateInitScriptRequest_McpId{McpId: value}
	case targetHook:
		createRequest.Target = &teamsv1.CreateInitScriptRequest_HookId{HookId: value}
	}

	return createRequest, nil
}

func initScriptUpdateToProto(id openapi_types.UUID, request gen.InitScriptUpdateRequest) *teamsv1.UpdateInitScriptRequest {
	return &teamsv1.UpdateInitScriptRequest{
		Id:          id.String(),
		Script:      request.Script,
		Description: request.Description,
	}
}

func mcpFromProto(mcp *teamsv1.Mcp) (gen.Mcp, error) {
	if mcp == nil {
		return gen.Mcp{}, fmt.Errorf("mcp missing")
	}

	meta, err := metaFromProto(mcp.GetMeta())
	if err != nil {
		return gen.Mcp{}, err
	}

	agentID, err := parseUUID(mcp.GetAgentId())
	if err != nil {
		return gen.Mcp{}, fmt.Errorf("parse agent id: %w", err)
	}

	result := gen.Mcp{
		AgentId:   agentID,
		Command:   mcp.GetCommand(),
		CreatedAt: meta.createdAt,
		Id:        meta.id,
		Image:     mcp.GetImage(),
		Resources: computeResourcesFromProto(mcp.GetResources()),
		UpdatedAt: meta.updatedAt,
	}

	if mcp.Description != "" {
		result.Description = stringPtr(mcp.Description)
	}

	return result, nil
}

func mcpCreateToProto(request gen.McpCreateRequest) *teamsv1.CreateMcpRequest {
	return &teamsv1.CreateMcpRequest{
		AgentId:     request.AgentId.String(),
		Command:     request.Command,
		Description: stringValue(request.Description),
		Image:       request.Image,
		Resources:   computeResourcesToProto(request.Resources),
	}
}

func mcpUpdateToProto(id openapi_types.UUID, request gen.McpUpdateRequest) *teamsv1.UpdateMcpRequest {
	return &teamsv1.UpdateMcpRequest{
		Id:          id.String(),
		Command:     request.Command,
		Description: request.Description,
		Image:       request.Image,
		Resources:   computeResourcesToProto(request.Resources),
	}
}

func skillFromProto(skill *teamsv1.Skill) (gen.Skill, error) {
	if skill == nil {
		return gen.Skill{}, fmt.Errorf("skill missing")
	}

	meta, err := metaFromProto(skill.GetMeta())
	if err != nil {
		return gen.Skill{}, err
	}

	agentID, err := parseUUID(skill.GetAgentId())
	if err != nil {
		return gen.Skill{}, fmt.Errorf("parse agent id: %w", err)
	}

	result := gen.Skill{
		AgentId:   agentID,
		Body:      skill.GetBody(),
		CreatedAt: meta.createdAt,
		Id:        meta.id,
		Name:      skill.GetName(),
		UpdatedAt: meta.updatedAt,
	}

	if skill.Description != "" {
		result.Description = stringPtr(skill.Description)
	}

	return result, nil
}

func skillCreateToProto(request gen.SkillCreateRequest) *teamsv1.CreateSkillRequest {
	return &teamsv1.CreateSkillRequest{
		AgentId:     request.AgentId.String(),
		Name:        request.Name,
		Body:        request.Body,
		Description: stringValue(request.Description),
	}
}

func skillUpdateToProto(id openapi_types.UUID, request gen.SkillUpdateRequest) *teamsv1.UpdateSkillRequest {
	return &teamsv1.UpdateSkillRequest{
		Id:          id.String(),
		Name:        request.Name,
		Body:        request.Body,
		Description: request.Description,
	}
}

func volumeFromProto(volume *teamsv1.Volume) (gen.Volume, error) {
	if volume == nil {
		return gen.Volume{}, fmt.Errorf("volume missing")
	}

	meta, err := metaFromProto(volume.GetMeta())
	if err != nil {
		return gen.Volume{}, err
	}

	result := gen.Volume{
		CreatedAt:  meta.createdAt,
		Id:         meta.id,
		MountPath:  volume.GetMountPath(),
		Persistent: volume.GetPersistent(),
		UpdatedAt:  meta.updatedAt,
	}

	if volume.Description != "" {
		result.Description = stringPtr(volume.Description)
	}
	if volume.Size != "" {
		result.Size = stringPtr(volume.Size)
	}

	return result, nil
}

func volumeCreateToProto(request gen.VolumeCreateRequest) *teamsv1.CreateVolumeRequest {
	return &teamsv1.CreateVolumeRequest{
		Persistent:  request.Persistent,
		MountPath:   request.MountPath,
		Size:        stringValue(request.Size),
		Description: stringValue(request.Description),
	}
}

func volumeUpdateToProto(id openapi_types.UUID, request gen.VolumeUpdateRequest) *teamsv1.UpdateVolumeRequest {
	return &teamsv1.UpdateVolumeRequest{
		Id:          id.String(),
		Persistent:  request.Persistent,
		MountPath:   request.MountPath,
		Size:        request.Size,
		Description: request.Description,
	}
}

func volumeAttachmentFromProto(attachment *teamsv1.VolumeAttachment) (gen.VolumeAttachment, error) {
	if attachment == nil {
		return gen.VolumeAttachment{}, fmt.Errorf("volume attachment missing")
	}

	meta, err := metaFromProto(attachment.GetMeta())
	if err != nil {
		return gen.VolumeAttachment{}, err
	}

	volumeID, err := parseUUID(attachment.GetVolumeId())
	if err != nil {
		return gen.VolumeAttachment{}, fmt.Errorf("parse volume id: %w", err)
	}

	agentID, mcpID, hookID, err := parseTargetIDs(attachment.GetAgentId(), attachment.GetMcpId(), attachment.GetHookId())
	if err != nil {
		return gen.VolumeAttachment{}, err
	}

	return gen.VolumeAttachment{
		AgentId:   agentID,
		CreatedAt: meta.createdAt,
		HookId:    hookID,
		Id:        meta.id,
		McpId:     mcpID,
		UpdatedAt: meta.updatedAt,
		VolumeId:  volumeID,
	}, nil
}

func volumeAttachmentCreateToProto(request gen.VolumeAttachmentCreateRequest) (*teamsv1.CreateVolumeAttachmentRequest, error) {
	createRequest := &teamsv1.CreateVolumeAttachmentRequest{
		VolumeId: request.VolumeId.String(),
	}

	kind, value, err := resolveTargetKind(request.AgentId, request.McpId, request.HookId)
	if err != nil {
		return nil, err
	}
	if kind == targetNone {
		return nil, fmt.Errorf("volume attachment target missing")
	}
	switch kind {
	case targetAgent:
		createRequest.Target = &teamsv1.CreateVolumeAttachmentRequest_AgentId{AgentId: value}
	case targetMcp:
		createRequest.Target = &teamsv1.CreateVolumeAttachmentRequest_McpId{McpId: value}
	case targetHook:
		createRequest.Target = &teamsv1.CreateVolumeAttachmentRequest_HookId{HookId: value}
	}

	return createRequest, nil
}

func parseTargetIDs(agentID, mcpID, hookID string) (*openapi_types.UUID, *openapi_types.UUID, *openapi_types.UUID, error) {
	count := 0
	if agentID != "" {
		count++
	}
	if mcpID != "" {
		count++
	}
	if hookID != "" {
		count++
	}
	if count > 1 {
		return nil, nil, nil, fmt.Errorf("multiple targets set")
	}

	agentUUID, err := uuidPtrFromString(agentID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse agent id: %w", err)
	}
	mcpUUID, err := uuidPtrFromString(mcpID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse mcp id: %w", err)
	}
	hookUUID, err := uuidPtrFromString(hookID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse hook id: %w", err)
	}

	return agentUUID, mcpUUID, hookUUID, nil
}

func resolveTargetKind(agentID, mcpID, hookID *openapi_types.UUID) (targetKind, string, error) {
	count := 0
	if agentID != nil {
		count++
	}
	if mcpID != nil {
		count++
	}
	if hookID != nil {
		count++
	}
	if count == 0 {
		return targetNone, "", nil
	}
	if count > 1 {
		return targetNone, "", fmt.Errorf("multiple targets set")
	}
	if agentID != nil {
		return targetAgent, agentID.String(), nil
	}
	if mcpID != nil {
		return targetMcp, mcpID.String(), nil
	}
	return targetHook, hookID.String(), nil
}

func resolveSourceKind(value *string, secretID *openapi_types.UUID) (sourceKind, string, error) {
	count := 0
	if value != nil {
		count++
	}
	if secretID != nil {
		count++
	}
	if count == 0 {
		return sourceNone, "", nil
	}
	if count > 1 {
		return sourceNone, "", fmt.Errorf("multiple sources set")
	}
	if value != nil {
		return sourceValue, *value, nil
	}
	return sourceSecret, secretID.String(), nil
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	v := value
	return &v
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
