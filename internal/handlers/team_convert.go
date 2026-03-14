package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"google.golang.org/protobuf/types/known/structpb"

	teamsv1 "github.com/agynio/gateway/gen/agynio/api/teams/v1"
	"github.com/agynio/gateway/internal/gen"
)

type entityMeta struct {
	id        openapi_types.UUID
	createdAt time.Time
	updatedAt *time.Time
}

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

func agentFromProto(agent *teamsv1.Agent) (gen.Agent, error) {
	if agent == nil {
		return gen.Agent{}, fmt.Errorf("agent missing")
	}

	meta, err := metaFromProto(agent.GetMeta())
	if err != nil {
		return gen.Agent{}, err
	}

	config, err := agentConfigFromProto(agent.GetConfig())
	if err != nil {
		return gen.Agent{}, err
	}

	result := gen.Agent{
		Id:        meta.id,
		CreatedAt: meta.createdAt,
		UpdatedAt: meta.updatedAt,
		Config:    config,
	}

	if agent.Title != "" {
		result.Title = stringPtr(agent.Title)
	}
	if agent.Description != "" {
		result.Description = stringPtr(agent.Description)
	}

	return result, nil
}

func agentConfigFromProto(config *teamsv1.AgentConfig) (gen.AgentConfig, error) {
	if config == nil {
		return gen.AgentConfig{}, fmt.Errorf("agent config missing")
	}

	whenBusy, err := agentWhenBusyFromProto(config.GetWhenBusy())
	if err != nil {
		return gen.AgentConfig{}, err
	}

	processBuffer, err := agentProcessBufferFromProto(config.GetProcessBuffer())
	if err != nil {
		return gen.AgentConfig{}, err
	}

	result := gen.AgentConfig{
		WhenBusy:      whenBusy,
		ProcessBuffer: processBuffer,
	}

	if config.Model != "" {
		result.Model = stringPtr(config.Model)
	}
	if config.SystemPrompt != "" {
		result.SystemPrompt = stringPtr(config.SystemPrompt)
	}
	result.DebounceMs = intPtr(config.DebounceMs)
	result.SendFinalResponseToThread = boolPtr(config.SendFinalResponseToThread)
	result.SummarizationKeepTokens = intPtr(config.SummarizationKeepTokens)
	result.SummarizationMaxTokens = intPtr(config.SummarizationMaxTokens)
	result.RestrictOutput = boolPtr(config.RestrictOutput)
	if config.RestrictionMessage != "" {
		result.RestrictionMessage = stringPtr(config.RestrictionMessage)
	}
	result.RestrictionMaxInjections = intPtr(config.RestrictionMaxInjections)
	if config.Name != "" {
		result.Name = stringPtr(config.Name)
	}
	if config.Role != "" {
		result.Role = stringPtr(config.Role)
	}

	return result, nil
}

func agentConfigToProto(config gen.AgentConfig) (*teamsv1.AgentConfig, error) {
	whenBusy, err := agentWhenBusyToProto(config.WhenBusy)
	if err != nil {
		return nil, err
	}

	processBuffer, err := agentProcessBufferToProto(config.ProcessBuffer)
	if err != nil {
		return nil, err
	}

	return &teamsv1.AgentConfig{
		Model:                     stringValue(config.Model),
		SystemPrompt:              stringValue(config.SystemPrompt),
		DebounceMs:                int32Value(config.DebounceMs),
		WhenBusy:                  whenBusy,
		ProcessBuffer:             processBuffer,
		SendFinalResponseToThread: boolValue(config.SendFinalResponseToThread),
		SummarizationKeepTokens:   int32Value(config.SummarizationKeepTokens),
		SummarizationMaxTokens:    int32Value(config.SummarizationMaxTokens),
		RestrictOutput:            boolValue(config.RestrictOutput),
		RestrictionMessage:        stringValue(config.RestrictionMessage),
		RestrictionMaxInjections:  int32Value(config.RestrictionMaxInjections),
		Name:                      stringValue(config.Name),
		Role:                      stringValue(config.Role),
	}, nil
}

func toolFromProto(tool *teamsv1.Tool) (gen.Tool, error) {
	if tool == nil {
		return gen.Tool{}, fmt.Errorf("tool missing")
	}

	meta, err := metaFromProto(tool.GetMeta())
	if err != nil {
		return gen.Tool{}, err
	}

	toolType, err := toolTypeFromProto(tool.GetType())
	if err != nil {
		return gen.Tool{}, err
	}

	config, err := structToMap(tool.GetConfig())
	if err != nil {
		return gen.Tool{}, err
	}

	result := gen.Tool{
		Id:        meta.id,
		CreatedAt: meta.createdAt,
		UpdatedAt: meta.updatedAt,
		Type:      toolType,
	}

	if tool.Name != "" {
		result.Name = stringPtr(tool.Name)
	}
	if tool.Description != "" {
		result.Description = stringPtr(tool.Description)
	}
	if config != nil {
		result.Config = &config
	}

	return result, nil
}

func mcpServerFromProto(server *teamsv1.McpServer) (gen.McpServer, error) {
	if server == nil {
		return gen.McpServer{}, fmt.Errorf("mcp server missing")
	}

	meta, err := metaFromProto(server.GetMeta())
	if err != nil {
		return gen.McpServer{}, err
	}

	config, err := mcpServerConfigFromProto(server.GetConfig())
	if err != nil {
		return gen.McpServer{}, err
	}

	result := gen.McpServer{
		Id:        meta.id,
		CreatedAt: meta.createdAt,
		UpdatedAt: meta.updatedAt,
		Config:    config,
	}

	if server.Title != "" {
		result.Title = stringPtr(server.Title)
	}
	if server.Description != "" {
		result.Description = stringPtr(server.Description)
	}

	return result, nil
}

func mcpServerConfigFromProto(config *teamsv1.McpServerConfig) (gen.McpServerConfig, error) {
	if config == nil {
		return gen.McpServerConfig{}, fmt.Errorf("mcp config missing")
	}

	result := gen.McpServerConfig{}

	if config.Command != "" {
		result.Command = stringPtr(config.Command)
	}
	if config.Namespace != "" {
		result.Namespace = stringPtr(config.Namespace)
	}
	if config.Workdir != "" {
		result.Workdir = stringPtr(config.Workdir)
	}
	result.RequestTimeoutMs = intPtr(config.RequestTimeoutMs)
	result.StartupTimeoutMs = intPtr(config.StartupTimeoutMs)
	result.HeartbeatIntervalMs = intPtr(config.HeartbeatIntervalMs)
	result.StaleTimeoutMs = intPtr(config.StaleTimeoutMs)

	if env := config.GetEnv(); len(env) > 0 {
		items := make([]gen.McpEnvItem, 0, len(env))
		for _, item := range env {
			if item == nil {
				continue
			}
			items = append(items, gen.McpEnvItem{Name: item.GetName(), Value: item.GetValue()})
		}
		if len(items) > 0 {
			result.Env = &items
		}
	}

	if restart := config.GetRestart(); restart != nil {
		restartConfig := &struct {
			BackoffMs   *int `json:"backoffMs,omitempty"`
			MaxAttempts *int `json:"maxAttempts,omitempty"`
		}{}
		restartConfig.BackoffMs = intPtr(restart.BackoffMs)
		restartConfig.MaxAttempts = intPtr(restart.MaxAttempts)
		result.Restart = restartConfig
	}

	if filter := config.GetToolFilter(); filter != nil {
		mode, err := mcpToolFilterModeFromProto(filter.GetMode())
		if err != nil {
			return gen.McpServerConfig{}, err
		}

		toolFilter := &struct {
			Mode  gen.McpServerConfigToolFilterMode `json:"mode"`
			Rules *[]struct {
				Pattern string `json:"pattern"`
			} `json:"rules,omitempty"`
		}{
			Mode: mode,
		}

		if rules := filter.GetRules(); len(rules) > 0 {
			items := make([]struct {
				Pattern string `json:"pattern"`
			}, 0, len(rules))
			for _, rule := range rules {
				if rule == nil {
					continue
				}
				items = append(items, struct {
					Pattern string `json:"pattern"`
				}{Pattern: rule.GetPattern()})
			}
			if len(items) > 0 {
				toolFilter.Rules = &items
			}
		}

		result.ToolFilter = toolFilter
	}

	return result, nil
}

func mcpServerConfigToProto(config gen.McpServerConfig) (*teamsv1.McpServerConfig, error) {
	var envItems []*teamsv1.McpEnvItem
	if config.Env != nil {
		for _, item := range *config.Env {
			copyItem := item
			envItems = append(envItems, &teamsv1.McpEnvItem{Name: copyItem.Name, Value: copyItem.Value})
		}
	}

	var restart *teamsv1.McpServerRestartConfig
	if config.Restart != nil {
		restart = &teamsv1.McpServerRestartConfig{
			MaxAttempts: int32Value(config.Restart.MaxAttempts),
			BackoffMs:   int32Value(config.Restart.BackoffMs),
		}
	}

	var toolFilter *teamsv1.McpToolFilter
	if config.ToolFilter != nil {
		mode, err := mcpToolFilterModeToProto(config.ToolFilter.Mode)
		if err != nil {
			return nil, err
		}

		var rules []*teamsv1.McpToolFilterRule
		if config.ToolFilter.Rules != nil {
			for _, rule := range *config.ToolFilter.Rules {
				copyRule := rule
				rules = append(rules, &teamsv1.McpToolFilterRule{Pattern: copyRule.Pattern})
			}
		}

		toolFilter = &teamsv1.McpToolFilter{Mode: mode, Rules: rules}
	}

	return &teamsv1.McpServerConfig{
		Namespace:           stringValue(config.Namespace),
		Command:             stringValue(config.Command),
		Workdir:             stringValue(config.Workdir),
		Env:                 envItems,
		RequestTimeoutMs:    int32Value(config.RequestTimeoutMs),
		StartupTimeoutMs:    int32Value(config.StartupTimeoutMs),
		HeartbeatIntervalMs: int32Value(config.HeartbeatIntervalMs),
		StaleTimeoutMs:      int32Value(config.StaleTimeoutMs),
		Restart:             restart,
		ToolFilter:          toolFilter,
	}, nil
}

func workspaceConfigurationFromProto(cfg *teamsv1.WorkspaceConfiguration) (gen.WorkspaceConfiguration, error) {
	if cfg == nil {
		return gen.WorkspaceConfiguration{}, fmt.Errorf("workspace configuration missing")
	}

	meta, err := metaFromProto(cfg.GetMeta())
	if err != nil {
		return gen.WorkspaceConfiguration{}, err
	}

	config, err := workspaceConfigFromProto(cfg.GetConfig())
	if err != nil {
		return gen.WorkspaceConfiguration{}, err
	}

	result := gen.WorkspaceConfiguration{
		Id:        meta.id,
		CreatedAt: meta.createdAt,
		UpdatedAt: meta.updatedAt,
		Config:    config,
	}

	if cfg.Title != "" {
		result.Title = stringPtr(cfg.Title)
	}
	if cfg.Description != "" {
		result.Description = stringPtr(cfg.Description)
	}

	return result, nil
}

func workspaceConfigFromProto(config *teamsv1.WorkspaceConfig) (gen.WorkspaceConfig, error) {
	if config == nil {
		return gen.WorkspaceConfig{}, fmt.Errorf("workspace config missing")
	}

	platform, err := workspacePlatformFromProto(config.GetPlatform())
	if err != nil {
		return gen.WorkspaceConfig{}, err
	}

	cpuLimit, err := cpuLimitFromProto(config.GetCpuLimit())
	if err != nil {
		return gen.WorkspaceConfig{}, err
	}

	memoryLimit, err := memoryLimitFromProto(config.GetMemoryLimit())
	if err != nil {
		return gen.WorkspaceConfig{}, err
	}

	nixConfig, err := structToMap(config.GetNix())
	if err != nil {
		return gen.WorkspaceConfig{}, err
	}

	result := gen.WorkspaceConfig{
		CpuLimit:    cpuLimit,
		MemoryLimit: memoryLimit,
		Platform:    platform,
	}

	if config.Image != "" {
		result.Image = stringPtr(config.Image)
	}
	if config.InitialScript != "" {
		result.InitialScript = stringPtr(config.InitialScript)
	}
	result.EnableDinD = boolPtr(config.EnableDind)
	result.TtlSeconds = intPtr(config.TtlSeconds)
	if env := config.GetEnv(); len(env) > 0 {
		items := make([]gen.WorkspaceEnvItem, 0, len(env))
		for _, item := range env {
			if item == nil {
				continue
			}
			items = append(items, gen.WorkspaceEnvItem{Name: item.GetName(), Value: item.GetValue()})
		}
		if len(items) > 0 {
			result.Env = &items
		}
	}
	if nixConfig != nil {
		result.Nix = &nixConfig
	}
	if volume := config.GetVolumes(); volume != nil {
		vol := gen.WorkspaceVolumeConfig{}
		if volume.MountPath != "" {
			vol.MountPath = stringPtr(volume.MountPath)
		}
		vol.Enabled = boolPtr(volume.Enabled)
		if vol.MountPath != nil || vol.Enabled != nil {
			result.Volumes = &vol
		}
	}

	return result, nil
}

func workspaceConfigToProto(config gen.WorkspaceConfig) (*teamsv1.WorkspaceConfig, error) {
	cpuLimit, err := cpuLimitToString(config.CpuLimit)
	if err != nil {
		return nil, err
	}

	memoryLimit, err := memoryLimitToString(config.MemoryLimit)
	if err != nil {
		return nil, err
	}

	platform, err := workspacePlatformToProto(config.Platform)
	if err != nil {
		return nil, err
	}

	nixConfig, err := mapToStruct(config.Nix)
	if err != nil {
		return nil, err
	}

	var envItems []*teamsv1.WorkspaceEnvItem
	if config.Env != nil {
		for _, item := range *config.Env {
			copyItem := item
			envItems = append(envItems, &teamsv1.WorkspaceEnvItem{Name: copyItem.Name, Value: copyItem.Value})
		}
	}

	var volumes *teamsv1.WorkspaceVolumeConfig
	if config.Volumes != nil {
		volumes = &teamsv1.WorkspaceVolumeConfig{
			Enabled:   boolValue(config.Volumes.Enabled),
			MountPath: stringValue(config.Volumes.MountPath),
		}
	}

	return &teamsv1.WorkspaceConfig{
		Image:         stringValue(config.Image),
		Env:           envItems,
		InitialScript: stringValue(config.InitialScript),
		CpuLimit:      cpuLimit,
		MemoryLimit:   memoryLimit,
		Platform:      platform,
		EnableDind:    boolValue(config.EnableDinD),
		TtlSeconds:    int32Value(config.TtlSeconds),
		Nix:           nixConfig,
		Volumes:       volumes,
	}, nil
}

func memoryBucketFromProto(bucket *teamsv1.MemoryBucket) (gen.MemoryBucket, error) {
	if bucket == nil {
		return gen.MemoryBucket{}, fmt.Errorf("memory bucket missing")
	}

	meta, err := metaFromProto(bucket.GetMeta())
	if err != nil {
		return gen.MemoryBucket{}, err
	}

	config, err := memoryBucketConfigFromProto(bucket.GetConfig())
	if err != nil {
		return gen.MemoryBucket{}, err
	}

	result := gen.MemoryBucket{
		Id:        meta.id,
		CreatedAt: meta.createdAt,
		UpdatedAt: meta.updatedAt,
		Config:    config,
	}

	if bucket.Title != "" {
		result.Title = stringPtr(bucket.Title)
	}
	if bucket.Description != "" {
		result.Description = stringPtr(bucket.Description)
	}

	return result, nil
}

func memoryBucketConfigFromProto(config *teamsv1.MemoryBucketConfig) (gen.MemoryBucketConfig, error) {
	if config == nil {
		return gen.MemoryBucketConfig{}, fmt.Errorf("memory bucket config missing")
	}

	scope, err := memoryBucketScopeFromProto(config.GetScope())
	if err != nil {
		return gen.MemoryBucketConfig{}, err
	}

	result := gen.MemoryBucketConfig{Scope: scope}
	if config.CollectionPrefix != "" {
		result.CollectionPrefix = stringPtr(config.CollectionPrefix)
	}

	return result, nil
}

func memoryBucketConfigToProto(config gen.MemoryBucketConfig) (*teamsv1.MemoryBucketConfig, error) {
	scope, err := memoryBucketScopeToProto(config.Scope)
	if err != nil {
		return nil, err
	}

	return &teamsv1.MemoryBucketConfig{
		Scope:            scope,
		CollectionPrefix: stringValue(config.CollectionPrefix),
	}, nil
}

func variableFromProto(variable *teamsv1.Variable) (gen.Variable, error) {
	if variable == nil {
		return gen.Variable{}, fmt.Errorf("variable missing")
	}

	meta, err := metaFromProto(variable.GetMeta())
	if err != nil {
		return gen.Variable{}, err
	}

	result := gen.Variable{
		Id:        meta.id,
		CreatedAt: meta.createdAt,
		UpdatedAt: meta.updatedAt,
		Key:       variable.GetKey(),
		Value:     variable.GetValue(),
	}

	if variable.Description != "" {
		result.Description = stringPtr(variable.Description)
	}

	return result, nil
}

func variableCreateToProto(request gen.VariableCreateRequest) (*teamsv1.CreateVariableRequest, error) {
	return &teamsv1.CreateVariableRequest{
		Key:         request.Key,
		Value:       request.Value,
		Description: stringValue(request.Description),
	}, nil
}

func variableUpdateToProto(request gen.VariableUpdateRequest) (*teamsv1.UpdateVariableRequest, error) {
	return &teamsv1.UpdateVariableRequest{
		Key:         request.Key,
		Value:       request.Value,
		Description: request.Description,
	}, nil
}

func attachmentFromProto(attachment *teamsv1.Attachment) (gen.Attachment, error) {
	if attachment == nil {
		return gen.Attachment{}, fmt.Errorf("attachment missing")
	}

	meta, err := metaFromProto(attachment.GetMeta())
	if err != nil {
		return gen.Attachment{}, err
	}

	kind, err := attachmentKindFromProto(attachment.GetKind())
	if err != nil {
		return gen.Attachment{}, err
	}

	sourceType, err := entityTypeFromProto(attachment.GetSourceType())
	if err != nil {
		return gen.Attachment{}, err
	}

	targetType, err := entityTypeFromProto(attachment.GetTargetType())
	if err != nil {
		return gen.Attachment{}, err
	}

	sourceID, err := uuid.Parse(attachment.GetSourceId())
	if err != nil {
		return gen.Attachment{}, fmt.Errorf("parse source id: %w", err)
	}

	targetID, err := uuid.Parse(attachment.GetTargetId())
	if err != nil {
		return gen.Attachment{}, fmt.Errorf("parse target id: %w", err)
	}

	return gen.Attachment{
		Id:         meta.id,
		CreatedAt:  meta.createdAt,
		UpdatedAt:  meta.updatedAt,
		Kind:       kind,
		SourceType: sourceType,
		SourceId:   openapi_types.UUID(sourceID),
		TargetType: targetType,
		TargetId:   openapi_types.UUID(targetID),
	}, nil
}

func toolTypeFromProto(toolType teamsv1.ToolType) (gen.ToolType, error) {
	switch toolType {
	case teamsv1.ToolType_TOOL_TYPE_MANAGE:
		return gen.Manage, nil
	case teamsv1.ToolType_TOOL_TYPE_MEMORY:
		return gen.Memory, nil
	case teamsv1.ToolType_TOOL_TYPE_SHELL_COMMAND:
		return gen.ShellCommand, nil
	case teamsv1.ToolType_TOOL_TYPE_SEND_MESSAGE:
		return gen.SendMessage, nil
	case teamsv1.ToolType_TOOL_TYPE_SEND_SLACK_MESSAGE:
		return gen.SendSlackMessage, nil
	case teamsv1.ToolType_TOOL_TYPE_REMIND_ME:
		return gen.RemindMe, nil
	case teamsv1.ToolType_TOOL_TYPE_GITHUB_CLONE_REPO:
		return gen.GithubCloneRepo, nil
	case teamsv1.ToolType_TOOL_TYPE_CALL_AGENT:
		return gen.CallAgent, nil
	default:
		return "", fmt.Errorf("unsupported tool type: %s", toolType.String())
	}
}

func toolTypeToProto(toolType gen.ToolType) (teamsv1.ToolType, error) {
	switch toolType {
	case gen.Manage:
		return teamsv1.ToolType_TOOL_TYPE_MANAGE, nil
	case gen.Memory:
		return teamsv1.ToolType_TOOL_TYPE_MEMORY, nil
	case gen.ShellCommand:
		return teamsv1.ToolType_TOOL_TYPE_SHELL_COMMAND, nil
	case gen.SendMessage:
		return teamsv1.ToolType_TOOL_TYPE_SEND_MESSAGE, nil
	case gen.SendSlackMessage:
		return teamsv1.ToolType_TOOL_TYPE_SEND_SLACK_MESSAGE, nil
	case gen.RemindMe:
		return teamsv1.ToolType_TOOL_TYPE_REMIND_ME, nil
	case gen.GithubCloneRepo:
		return teamsv1.ToolType_TOOL_TYPE_GITHUB_CLONE_REPO, nil
	case gen.CallAgent:
		return teamsv1.ToolType_TOOL_TYPE_CALL_AGENT, nil
	default:
		return teamsv1.ToolType_TOOL_TYPE_UNSPECIFIED, fmt.Errorf("unsupported tool type: %s", toolType)
	}
}

func attachmentKindFromProto(kind teamsv1.AttachmentKind) (gen.AttachmentKind, error) {
	switch kind {
	case teamsv1.AttachmentKind_ATTACHMENT_KIND_AGENT_TOOL:
		return gen.AgentTool, nil
	case teamsv1.AttachmentKind_ATTACHMENT_KIND_AGENT_MEMORY_BUCKET:
		return gen.AgentMemoryBucket, nil
	case teamsv1.AttachmentKind_ATTACHMENT_KIND_AGENT_WORKSPACE_CONFIGURATION:
		return gen.AgentWorkspaceConfiguration, nil
	case teamsv1.AttachmentKind_ATTACHMENT_KIND_AGENT_MCP_SERVER:
		return gen.AgentMcpServer, nil
	case teamsv1.AttachmentKind_ATTACHMENT_KIND_MCP_SERVER_WORKSPACE_CONFIGURATION:
		return gen.McpServerWorkspaceConfiguration, nil
	default:
		return "", fmt.Errorf("unsupported attachment kind: %s", kind.String())
	}
}

func attachmentKindToProto(kind gen.AttachmentKind) (teamsv1.AttachmentKind, error) {
	switch kind {
	case gen.AgentTool:
		return teamsv1.AttachmentKind_ATTACHMENT_KIND_AGENT_TOOL, nil
	case gen.AgentMemoryBucket:
		return teamsv1.AttachmentKind_ATTACHMENT_KIND_AGENT_MEMORY_BUCKET, nil
	case gen.AgentWorkspaceConfiguration:
		return teamsv1.AttachmentKind_ATTACHMENT_KIND_AGENT_WORKSPACE_CONFIGURATION, nil
	case gen.AgentMcpServer:
		return teamsv1.AttachmentKind_ATTACHMENT_KIND_AGENT_MCP_SERVER, nil
	case gen.McpServerWorkspaceConfiguration:
		return teamsv1.AttachmentKind_ATTACHMENT_KIND_MCP_SERVER_WORKSPACE_CONFIGURATION, nil
	default:
		return teamsv1.AttachmentKind_ATTACHMENT_KIND_UNSPECIFIED, fmt.Errorf("unsupported attachment kind: %s", kind)
	}
}

func entityTypeFromProto(entity teamsv1.EntityType) (gen.EntityType, error) {
	switch entity {
	case teamsv1.EntityType_ENTITY_TYPE_AGENT:
		return gen.EntityTypeAgent, nil
	case teamsv1.EntityType_ENTITY_TYPE_TOOL:
		return gen.EntityTypeTool, nil
	case teamsv1.EntityType_ENTITY_TYPE_MCP_SERVER:
		return gen.EntityTypeMcpServer, nil
	case teamsv1.EntityType_ENTITY_TYPE_WORKSPACE_CONFIGURATION:
		return gen.EntityTypeWorkspaceConfiguration, nil
	case teamsv1.EntityType_ENTITY_TYPE_MEMORY_BUCKET:
		return gen.EntityTypeMemoryBucket, nil
	case teamsv1.EntityType_ENTITY_TYPE_VARIABLE:
		return gen.EntityTypeVariable, nil
	default:
		return "", fmt.Errorf("unsupported entity type: %s", entity.String())
	}
}

func entityTypeToProto(entity gen.EntityType) (teamsv1.EntityType, error) {
	switch entity {
	case gen.EntityTypeAgent:
		return teamsv1.EntityType_ENTITY_TYPE_AGENT, nil
	case gen.EntityTypeTool:
		return teamsv1.EntityType_ENTITY_TYPE_TOOL, nil
	case gen.EntityTypeMcpServer:
		return teamsv1.EntityType_ENTITY_TYPE_MCP_SERVER, nil
	case gen.EntityTypeWorkspaceConfiguration:
		return teamsv1.EntityType_ENTITY_TYPE_WORKSPACE_CONFIGURATION, nil
	case gen.EntityTypeMemoryBucket:
		return teamsv1.EntityType_ENTITY_TYPE_MEMORY_BUCKET, nil
	case gen.EntityTypeVariable:
		return teamsv1.EntityType_ENTITY_TYPE_VARIABLE, nil
	default:
		return teamsv1.EntityType_ENTITY_TYPE_UNSPECIFIED, fmt.Errorf("unsupported entity type: %s", entity)
	}
}

func mcpToolFilterModeFromProto(mode teamsv1.McpToolFilterMode) (gen.McpServerConfigToolFilterMode, error) {
	switch mode {
	case teamsv1.McpToolFilterMode_MCP_TOOL_FILTER_MODE_ALLOW:
		return gen.Allow, nil
	case teamsv1.McpToolFilterMode_MCP_TOOL_FILTER_MODE_DENY:
		return gen.Deny, nil
	default:
		return "", fmt.Errorf("unsupported mcp tool filter mode: %s", mode.String())
	}
}

func mcpToolFilterModeToProto(mode gen.McpServerConfigToolFilterMode) (teamsv1.McpToolFilterMode, error) {
	switch mode {
	case gen.Allow:
		return teamsv1.McpToolFilterMode_MCP_TOOL_FILTER_MODE_ALLOW, nil
	case gen.Deny:
		return teamsv1.McpToolFilterMode_MCP_TOOL_FILTER_MODE_DENY, nil
	default:
		return teamsv1.McpToolFilterMode_MCP_TOOL_FILTER_MODE_UNSPECIFIED, fmt.Errorf("unsupported mcp tool filter mode: %s", mode)
	}
}

func memoryBucketScopeFromProto(scope teamsv1.MemoryBucketScope) (*gen.MemoryBucketConfigScope, error) {
	switch scope {
	case teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_UNSPECIFIED:
		return nil, nil
	case teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_GLOBAL:
		value := gen.Global
		return &value, nil
	case teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_PER_THREAD:
		value := gen.PerThread
		return &value, nil
	default:
		return nil, fmt.Errorf("unsupported memory bucket scope: %s", scope.String())
	}
}

func memoryBucketScopeToProto(scope *gen.MemoryBucketConfigScope) (teamsv1.MemoryBucketScope, error) {
	if scope == nil {
		return teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_UNSPECIFIED, nil
	}

	switch *scope {
	case gen.Global:
		return teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_GLOBAL, nil
	case gen.PerThread:
		return teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_PER_THREAD, nil
	default:
		return teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_UNSPECIFIED, fmt.Errorf("unsupported memory bucket scope: %s", *scope)
	}
}

func workspacePlatformFromProto(platform teamsv1.WorkspacePlatform) (*gen.WorkspacePlatform, error) {
	switch platform {
	case teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_UNSPECIFIED:
		return nil, nil
	case teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_LINUX_AMD64:
		value := gen.Linuxamd64
		return &value, nil
	case teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_LINUX_ARM64:
		value := gen.Linuxarm64
		return &value, nil
	case teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_AUTO:
		value := gen.Auto
		return &value, nil
	default:
		return nil, fmt.Errorf("unsupported workspace platform: %s", platform.String())
	}
}

func workspacePlatformToProto(platform *gen.WorkspacePlatform) (teamsv1.WorkspacePlatform, error) {
	if platform == nil {
		return teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_UNSPECIFIED, nil
	}

	switch *platform {
	case gen.Linuxamd64:
		return teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_LINUX_AMD64, nil
	case gen.Linuxarm64:
		return teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_LINUX_ARM64, nil
	case gen.Auto:
		return teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_AUTO, nil
	default:
		return teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_UNSPECIFIED, fmt.Errorf("unsupported workspace platform: %s", *platform)
	}
}

func agentWhenBusyFromProto(value teamsv1.AgentWhenBusy) (*gen.AgentConfigWhenBusy, error) {
	switch value {
	case teamsv1.AgentWhenBusy_AGENT_WHEN_BUSY_UNSPECIFIED:
		return nil, nil
	case teamsv1.AgentWhenBusy_AGENT_WHEN_BUSY_WAIT:
		v := gen.Wait
		return &v, nil
	case teamsv1.AgentWhenBusy_AGENT_WHEN_BUSY_INJECT_AFTER_TOOLS:
		v := gen.InjectAfterTools
		return &v, nil
	default:
		return nil, fmt.Errorf("unsupported when_busy: %s", value.String())
	}
}

func agentWhenBusyToProto(value *gen.AgentConfigWhenBusy) (teamsv1.AgentWhenBusy, error) {
	if value == nil {
		return teamsv1.AgentWhenBusy_AGENT_WHEN_BUSY_UNSPECIFIED, nil
	}

	switch *value {
	case gen.Wait:
		return teamsv1.AgentWhenBusy_AGENT_WHEN_BUSY_WAIT, nil
	case gen.InjectAfterTools:
		return teamsv1.AgentWhenBusy_AGENT_WHEN_BUSY_INJECT_AFTER_TOOLS, nil
	default:
		return teamsv1.AgentWhenBusy_AGENT_WHEN_BUSY_UNSPECIFIED, fmt.Errorf("unsupported when_busy: %s", *value)
	}
}

func agentProcessBufferFromProto(value teamsv1.AgentProcessBuffer) (*gen.AgentConfigProcessBuffer, error) {
	switch value {
	case teamsv1.AgentProcessBuffer_AGENT_PROCESS_BUFFER_UNSPECIFIED:
		return nil, nil
	case teamsv1.AgentProcessBuffer_AGENT_PROCESS_BUFFER_ALL_TOGETHER:
		v := gen.AllTogether
		return &v, nil
	case teamsv1.AgentProcessBuffer_AGENT_PROCESS_BUFFER_ONE_BY_ONE:
		v := gen.OneByOne
		return &v, nil
	default:
		return nil, fmt.Errorf("unsupported process_buffer: %s", value.String())
	}
}

func agentProcessBufferToProto(value *gen.AgentConfigProcessBuffer) (teamsv1.AgentProcessBuffer, error) {
	if value == nil {
		return teamsv1.AgentProcessBuffer_AGENT_PROCESS_BUFFER_UNSPECIFIED, nil
	}

	switch *value {
	case gen.AllTogether:
		return teamsv1.AgentProcessBuffer_AGENT_PROCESS_BUFFER_ALL_TOGETHER, nil
	case gen.OneByOne:
		return teamsv1.AgentProcessBuffer_AGENT_PROCESS_BUFFER_ONE_BY_ONE, nil
	default:
		return teamsv1.AgentProcessBuffer_AGENT_PROCESS_BUFFER_UNSPECIFIED, fmt.Errorf("unsupported process_buffer: %s", *value)
	}
}

func cpuLimitFromProto(value string) (*gen.WorkspaceConfig_CpuLimit, error) {
	if value == "" {
		return nil, nil
	}

	limit := gen.WorkspaceConfig_CpuLimit{}
	if err := limit.FromWorkspaceConfigCpuLimit1(gen.WorkspaceConfigCpuLimit1(value)); err != nil {
		return nil, err
	}
	return &limit, nil
}

func memoryLimitFromProto(value string) (*gen.WorkspaceConfig_MemoryLimit, error) {
	if value == "" {
		return nil, nil
	}

	limit := gen.WorkspaceConfig_MemoryLimit{}
	if err := limit.FromWorkspaceConfigMemoryLimit1(gen.WorkspaceConfigMemoryLimit1(value)); err != nil {
		return nil, err
	}
	return &limit, nil
}

func cpuLimitToString(limit *gen.WorkspaceConfig_CpuLimit) (string, error) {
	if limit == nil {
		return "", nil
	}

	if value, err := limit.AsWorkspaceConfigCpuLimit1(); err == nil {
		return string(value), nil
	}

	if value, err := limit.AsWorkspaceConfigCpuLimit0(); err == nil {
		return strconv.FormatFloat(float64(value), 'f', -1, 32), nil
	}

	return "", fmt.Errorf("unsupported workspace limit")
}

func memoryLimitToString(limit *gen.WorkspaceConfig_MemoryLimit) (string, error) {
	if limit == nil {
		return "", nil
	}

	if value, err := limit.AsWorkspaceConfigMemoryLimit1(); err == nil {
		return string(value), nil
	}

	if value, err := limit.AsWorkspaceConfigMemoryLimit0(); err == nil {
		return strconv.FormatFloat(float64(value), 'f', -1, 32), nil
	}

	return "", fmt.Errorf("unsupported workspace limit")
}

func structToMap(value *structpb.Struct) (map[string]any, error) {
	if value == nil {
		return nil, nil
	}

	result := value.AsMap()
	if result == nil {
		result = map[string]any{}
	}

	return result, nil
}

func mapToStruct(value *map[string]any) (*structpb.Struct, error) {
	if value == nil {
		return nil, nil
	}

	if *value == nil {
		return structpb.NewStruct(map[string]any{})
	}

	return structpb.NewStruct(*value)
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	v := value
	return &v
}

func intPtr(value int32) *int {
	v := int(value)
	return &v
}

func boolPtr(value bool) *bool {
	v := value
	return &v
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func int32Value(value *int) int32 {
	if value == nil {
		return 0
	}
	return int32(*value)
}

func boolValue(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}
