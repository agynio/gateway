package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	teamsv1 "github.com/agynio/gateway/gen/agynio/api/teams/v1"
	"github.com/agynio/gateway/internal/gen"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	teamBasePath   = "/team/v1"
	defaultPerPage = 20
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
	page, perPage := normalizePagination(request.Params.Page, request.Params.PerPage)
	query := ""
	if request.Params.Q != nil {
		query = strings.TrimSpace(*request.Params.Q)
	}

	agents, err := listAll(ctx, int32(perPage), func(ctx context.Context, token string, pageSize int32) ([]*teamsv1.Agent, string, error) {
		resp, err := t.client.ListAgents(ctx, &teamsv1.ListAgentsRequest{
			PageSize:  pageSize,
			PageToken: token,
			Query:     query,
		})
		if err != nil {
			return nil, "", err
		}
		return resp.GetAgents(), resp.GetNextPageToken(), nil
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	pageItems := paginateSlice(agents, page, perPage)
	items := make([]gen.Agent, 0, len(pageItems))
	for _, agent := range pageItems {
		converted, err := agentFromProto(agent)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedAgents{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   len(agents),
	}

	return gen.GetAgents200JSONResponse(payload), nil
}

func (t *Team) PostAgents(ctx context.Context, request gen.PostAgentsRequestObject) (gen.PostAgentsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	config, err := agentConfigToProto(request.Body.Config)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := t.client.CreateAgent(ctx, &teamsv1.CreateAgentRequest{
		Title:       stringValue(request.Body.Title),
		Description: stringValue(request.Body.Description),
		Config:      config,
	})
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

	var config *teamsv1.AgentConfig
	if request.Body.Config != nil {
		converted, err := agentConfigToProto(*request.Body.Config)
		if err != nil {
			return nil, requestProblem(err)
		}
		config = converted
	}

	resp, err := t.client.UpdateAgent(ctx, &teamsv1.UpdateAgentRequest{
		Id:          request.Id.String(),
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      config,
	})
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

func (t *Team) GetAttachments(ctx context.Context, request gen.GetAttachmentsRequestObject) (gen.GetAttachmentsResponseObject, error) {
	page, perPage := normalizePagination(request.Params.Page, request.Params.PerPage)

	baseRequest := &teamsv1.ListAttachmentsRequest{PageSize: int32(perPage)}

	if request.Params.SourceType != nil {
		value, err := entityTypeToProto(*request.Params.SourceType)
		if err != nil {
			return nil, requestProblem(err)
		}
		baseRequest.SourceType = value
	}
	if request.Params.SourceId != nil {
		baseRequest.SourceId = request.Params.SourceId.String()
	}
	if request.Params.TargetType != nil {
		value, err := entityTypeToProto(*request.Params.TargetType)
		if err != nil {
			return nil, requestProblem(err)
		}
		baseRequest.TargetType = value
	}
	if request.Params.TargetId != nil {
		baseRequest.TargetId = request.Params.TargetId.String()
	}
	if request.Params.Kind != nil {
		value, err := attachmentKindToProto(*request.Params.Kind)
		if err != nil {
			return nil, requestProblem(err)
		}
		baseRequest.Kind = value
	}

	attachments, err := listAll(ctx, int32(perPage), func(ctx context.Context, token string, pageSize int32) ([]*teamsv1.Attachment, string, error) {
		req := &teamsv1.ListAttachmentsRequest{
			PageSize:   pageSize,
			PageToken:  token,
			SourceType: baseRequest.SourceType,
			SourceId:   baseRequest.SourceId,
			TargetType: baseRequest.TargetType,
			TargetId:   baseRequest.TargetId,
			Kind:       baseRequest.Kind,
		}
		resp, err := t.client.ListAttachments(ctx, req)
		if err != nil {
			return nil, "", err
		}
		return resp.GetAttachments(), resp.GetNextPageToken(), nil
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	pageItems := paginateSlice(attachments, page, perPage)
	items := make([]gen.Attachment, 0, len(pageItems))
	for _, attachment := range pageItems {
		converted, err := attachmentFromProto(attachment)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedAttachments{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   len(attachments),
	}

	return gen.GetAttachments200JSONResponse(payload), nil
}

func (t *Team) PostAttachments(ctx context.Context, request gen.PostAttachmentsRequestObject) (gen.PostAttachmentsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	kind, err := attachmentKindToProto(request.Body.Kind)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := t.client.CreateAttachment(ctx, &teamsv1.CreateAttachmentRequest{
		Kind:     kind,
		SourceId: request.Body.SourceId.String(),
		TargetId: request.Body.TargetId.String(),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	attachment := resp.GetAttachment()
	if attachment == nil {
		return nil, responseProblem(fmt.Errorf("create attachment response missing attachment"))
	}

	converted, err := attachmentFromProto(attachment)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostAttachments201JSONResponse(converted), nil
}

func (t *Team) DeleteAttachmentsId(ctx context.Context, request gen.DeleteAttachmentsIdRequestObject) (gen.DeleteAttachmentsIdResponseObject, error) {
	_, err := t.client.DeleteAttachment(ctx, &teamsv1.DeleteAttachmentRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteAttachmentsId204Response{}, nil
}

func (t *Team) GetMcpServers(ctx context.Context, request gen.GetMcpServersRequestObject) (gen.GetMcpServersResponseObject, error) {
	page, perPage := normalizePagination(request.Params.Page, request.Params.PerPage)

	servers, err := listAll(ctx, int32(perPage), func(ctx context.Context, token string, pageSize int32) ([]*teamsv1.McpServer, string, error) {
		resp, err := t.client.ListMcpServers(ctx, &teamsv1.ListMcpServersRequest{
			PageSize:  pageSize,
			PageToken: token,
		})
		if err != nil {
			return nil, "", err
		}
		return resp.GetMcpServers(), resp.GetNextPageToken(), nil
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	pageItems := paginateSlice(servers, page, perPage)
	items := make([]gen.McpServer, 0, len(pageItems))
	for _, server := range pageItems {
		converted, err := mcpServerFromProto(server)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedMcpServers{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   len(servers),
	}

	return gen.GetMcpServers200JSONResponse(payload), nil
}

func (t *Team) PostMcpServers(ctx context.Context, request gen.PostMcpServersRequestObject) (gen.PostMcpServersResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	config, err := mcpServerConfigToProto(request.Body.Config)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := t.client.CreateMcpServer(ctx, &teamsv1.CreateMcpServerRequest{
		Title:       stringValue(request.Body.Title),
		Description: stringValue(request.Body.Description),
		Config:      config,
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	server := resp.GetMcpServer()
	if server == nil {
		return nil, responseProblem(fmt.Errorf("create mcp server response missing server"))
	}

	converted, err := mcpServerFromProto(server)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostMcpServers201JSONResponse(converted), nil
}

func (t *Team) DeleteMcpServersId(ctx context.Context, request gen.DeleteMcpServersIdRequestObject) (gen.DeleteMcpServersIdResponseObject, error) {
	_, err := t.client.DeleteMcpServer(ctx, &teamsv1.DeleteMcpServerRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteMcpServersId204Response{}, nil
}

func (t *Team) GetMcpServersId(ctx context.Context, request gen.GetMcpServersIdRequestObject) (gen.GetMcpServersIdResponseObject, error) {
	resp, err := t.client.GetMcpServer(ctx, &teamsv1.GetMcpServerRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	server := resp.GetMcpServer()
	if server == nil {
		return nil, responseProblem(fmt.Errorf("get mcp server response missing server"))
	}

	converted, err := mcpServerFromProto(server)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetMcpServersId200JSONResponse(converted), nil
}

func (t *Team) PatchMcpServersId(ctx context.Context, request gen.PatchMcpServersIdRequestObject) (gen.PatchMcpServersIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var config *teamsv1.McpServerConfig
	if request.Body.Config != nil {
		converted, err := mcpServerConfigToProto(*request.Body.Config)
		if err != nil {
			return nil, requestProblem(err)
		}
		config = converted
	}

	resp, err := t.client.UpdateMcpServer(ctx, &teamsv1.UpdateMcpServerRequest{
		Id:          request.Id.String(),
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      config,
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	server := resp.GetMcpServer()
	if server == nil {
		return nil, responseProblem(fmt.Errorf("update mcp server response missing server"))
	}

	converted, err := mcpServerFromProto(server)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchMcpServersId200JSONResponse(converted), nil
}

func (t *Team) GetMemoryBuckets(ctx context.Context, request gen.GetMemoryBucketsRequestObject) (gen.GetMemoryBucketsResponseObject, error) {
	page, perPage := normalizePagination(request.Params.Page, request.Params.PerPage)

	buckets, err := listAll(ctx, int32(perPage), func(ctx context.Context, token string, pageSize int32) ([]*teamsv1.MemoryBucket, string, error) {
		resp, err := t.client.ListMemoryBuckets(ctx, &teamsv1.ListMemoryBucketsRequest{
			PageSize:  pageSize,
			PageToken: token,
		})
		if err != nil {
			return nil, "", err
		}
		return resp.GetMemoryBuckets(), resp.GetNextPageToken(), nil
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	pageItems := paginateSlice(buckets, page, perPage)
	items := make([]gen.MemoryBucket, 0, len(pageItems))
	for _, bucket := range pageItems {
		converted, err := memoryBucketFromProto(bucket)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedMemoryBuckets{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   len(buckets),
	}

	return gen.GetMemoryBuckets200JSONResponse(payload), nil
}

func (t *Team) PostMemoryBuckets(ctx context.Context, request gen.PostMemoryBucketsRequestObject) (gen.PostMemoryBucketsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	config, err := memoryBucketConfigToProto(request.Body.Config)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := t.client.CreateMemoryBucket(ctx, &teamsv1.CreateMemoryBucketRequest{
		Title:       stringValue(request.Body.Title),
		Description: stringValue(request.Body.Description),
		Config:      config,
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	bucket := resp.GetMemoryBucket()
	if bucket == nil {
		return nil, responseProblem(fmt.Errorf("create memory bucket response missing bucket"))
	}

	converted, err := memoryBucketFromProto(bucket)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostMemoryBuckets201JSONResponse(converted), nil
}

func (t *Team) DeleteMemoryBucketsId(ctx context.Context, request gen.DeleteMemoryBucketsIdRequestObject) (gen.DeleteMemoryBucketsIdResponseObject, error) {
	_, err := t.client.DeleteMemoryBucket(ctx, &teamsv1.DeleteMemoryBucketRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteMemoryBucketsId204Response{}, nil
}

func (t *Team) GetMemoryBucketsId(ctx context.Context, request gen.GetMemoryBucketsIdRequestObject) (gen.GetMemoryBucketsIdResponseObject, error) {
	resp, err := t.client.GetMemoryBucket(ctx, &teamsv1.GetMemoryBucketRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	bucket := resp.GetMemoryBucket()
	if bucket == nil {
		return nil, responseProblem(fmt.Errorf("get memory bucket response missing bucket"))
	}

	converted, err := memoryBucketFromProto(bucket)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetMemoryBucketsId200JSONResponse(converted), nil
}

func (t *Team) PatchMemoryBucketsId(ctx context.Context, request gen.PatchMemoryBucketsIdRequestObject) (gen.PatchMemoryBucketsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var config *teamsv1.MemoryBucketConfig
	if request.Body.Config != nil {
		converted, err := memoryBucketConfigToProto(*request.Body.Config)
		if err != nil {
			return nil, requestProblem(err)
		}
		config = converted
	}

	resp, err := t.client.UpdateMemoryBucket(ctx, &teamsv1.UpdateMemoryBucketRequest{
		Id:          request.Id.String(),
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      config,
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	bucket := resp.GetMemoryBucket()
	if bucket == nil {
		return nil, responseProblem(fmt.Errorf("update memory bucket response missing bucket"))
	}

	converted, err := memoryBucketFromProto(bucket)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchMemoryBucketsId200JSONResponse(converted), nil
}

func (t *Team) GetTools(ctx context.Context, request gen.GetToolsRequestObject) (gen.GetToolsResponseObject, error) {
	page, perPage := normalizePagination(request.Params.Page, request.Params.PerPage)

	toolType := teamsv1.ToolType_TOOL_TYPE_UNSPECIFIED
	if request.Params.Type != nil {
		value, err := toolTypeToProto(*request.Params.Type)
		if err != nil {
			return nil, requestProblem(err)
		}
		toolType = value
	}

	tools, err := listAll(ctx, int32(perPage), func(ctx context.Context, token string, pageSize int32) ([]*teamsv1.Tool, string, error) {
		resp, err := t.client.ListTools(ctx, &teamsv1.ListToolsRequest{
			PageSize:  pageSize,
			PageToken: token,
			Type:      toolType,
		})
		if err != nil {
			return nil, "", err
		}
		return resp.GetTools(), resp.GetNextPageToken(), nil
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	pageItems := paginateSlice(tools, page, perPage)
	items := make([]gen.Tool, 0, len(pageItems))
	for _, tool := range pageItems {
		converted, err := toolFromProto(tool)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedTools{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   len(tools),
	}

	return gen.GetTools200JSONResponse(payload), nil
}

func (t *Team) PostTools(ctx context.Context, request gen.PostToolsRequestObject) (gen.PostToolsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	toolType, err := toolTypeToProto(request.Body.Type)
	if err != nil {
		return nil, requestProblem(err)
	}

	config, err := mapToStruct(request.Body.Config)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := t.client.CreateTool(ctx, &teamsv1.CreateToolRequest{
		Type:        toolType,
		Name:        stringValue(request.Body.Name),
		Description: stringValue(request.Body.Description),
		Config:      config,
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	tool := resp.GetTool()
	if tool == nil {
		return nil, responseProblem(fmt.Errorf("create tool response missing tool"))
	}

	converted, err := toolFromProto(tool)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostTools201JSONResponse(converted), nil
}

func (t *Team) DeleteToolsId(ctx context.Context, request gen.DeleteToolsIdRequestObject) (gen.DeleteToolsIdResponseObject, error) {
	_, err := t.client.DeleteTool(ctx, &teamsv1.DeleteToolRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteToolsId204Response{}, nil
}

func (t *Team) GetToolsId(ctx context.Context, request gen.GetToolsIdRequestObject) (gen.GetToolsIdResponseObject, error) {
	resp, err := t.client.GetTool(ctx, &teamsv1.GetToolRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	tool := resp.GetTool()
	if tool == nil {
		return nil, responseProblem(fmt.Errorf("get tool response missing tool"))
	}

	converted, err := toolFromProto(tool)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetToolsId200JSONResponse(converted), nil
}

func (t *Team) PatchToolsId(ctx context.Context, request gen.PatchToolsIdRequestObject) (gen.PatchToolsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var config *structpb.Struct
	if request.Body.Config != nil {
		converted, err := mapToStruct(request.Body.Config)
		if err != nil {
			return nil, requestProblem(err)
		}
		config = converted
	}

	resp, err := t.client.UpdateTool(ctx, &teamsv1.UpdateToolRequest{
		Id:          request.Id.String(),
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Config:      config,
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	tool := resp.GetTool()
	if tool == nil {
		return nil, responseProblem(fmt.Errorf("update tool response missing tool"))
	}

	converted, err := toolFromProto(tool)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchToolsId200JSONResponse(converted), nil
}

func (t *Team) GetWorkspaceConfigurations(ctx context.Context, request gen.GetWorkspaceConfigurationsRequestObject) (gen.GetWorkspaceConfigurationsResponseObject, error) {
	page, perPage := normalizePagination(request.Params.Page, request.Params.PerPage)

	configs, err := listAll(ctx, int32(perPage), func(ctx context.Context, token string, pageSize int32) ([]*teamsv1.WorkspaceConfiguration, string, error) {
		resp, err := t.client.ListWorkspaceConfigurations(ctx, &teamsv1.ListWorkspaceConfigurationsRequest{
			PageSize:  pageSize,
			PageToken: token,
		})
		if err != nil {
			return nil, "", err
		}
		return resp.GetWorkspaceConfigurations(), resp.GetNextPageToken(), nil
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	pageItems := paginateSlice(configs, page, perPage)
	items := make([]gen.WorkspaceConfiguration, 0, len(pageItems))
	for _, cfg := range pageItems {
		converted, err := workspaceConfigurationFromProto(cfg)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := gen.PaginatedWorkspaceConfigurations{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   len(configs),
	}

	return gen.GetWorkspaceConfigurations200JSONResponse(payload), nil
}

func (t *Team) PostWorkspaceConfigurations(ctx context.Context, request gen.PostWorkspaceConfigurationsRequestObject) (gen.PostWorkspaceConfigurationsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	config, err := workspaceConfigToProto(request.Body.Config)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := t.client.CreateWorkspaceConfiguration(ctx, &teamsv1.CreateWorkspaceConfigurationRequest{
		Title:       stringValue(request.Body.Title),
		Description: stringValue(request.Body.Description),
		Config:      config,
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	configuration := resp.GetWorkspaceConfiguration()
	if configuration == nil {
		return nil, responseProblem(fmt.Errorf("create workspace configuration response missing configuration"))
	}

	converted, err := workspaceConfigurationFromProto(configuration)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PostWorkspaceConfigurations201JSONResponse(converted), nil
}

func (t *Team) DeleteWorkspaceConfigurationsId(ctx context.Context, request gen.DeleteWorkspaceConfigurationsIdRequestObject) (gen.DeleteWorkspaceConfigurationsIdResponseObject, error) {
	_, err := t.client.DeleteWorkspaceConfiguration(ctx, &teamsv1.DeleteWorkspaceConfigurationRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}
	return gen.DeleteWorkspaceConfigurationsId204Response{}, nil
}

func (t *Team) GetWorkspaceConfigurationsId(ctx context.Context, request gen.GetWorkspaceConfigurationsIdRequestObject) (gen.GetWorkspaceConfigurationsIdResponseObject, error) {
	resp, err := t.client.GetWorkspaceConfiguration(ctx, &teamsv1.GetWorkspaceConfigurationRequest{Id: request.Id.String()})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	configuration := resp.GetWorkspaceConfiguration()
	if configuration == nil {
		return nil, responseProblem(fmt.Errorf("get workspace configuration response missing configuration"))
	}

	converted, err := workspaceConfigurationFromProto(configuration)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.GetWorkspaceConfigurationsId200JSONResponse(converted), nil
}

func (t *Team) PatchWorkspaceConfigurationsId(ctx context.Context, request gen.PatchWorkspaceConfigurationsIdRequestObject) (gen.PatchWorkspaceConfigurationsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var config *teamsv1.WorkspaceConfig
	if request.Body.Config != nil {
		converted, err := workspaceConfigToProto(*request.Body.Config)
		if err != nil {
			return nil, requestProblem(err)
		}
		config = converted
	}

	resp, err := t.client.UpdateWorkspaceConfiguration(ctx, &teamsv1.UpdateWorkspaceConfigurationRequest{
		Id:          request.Id.String(),
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      config,
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	configuration := resp.GetWorkspaceConfiguration()
	if configuration == nil {
		return nil, responseProblem(fmt.Errorf("update workspace configuration response missing configuration"))
	}

	converted, err := workspaceConfigurationFromProto(configuration)
	if err != nil {
		return nil, responseProblem(err)
	}

	return gen.PatchWorkspaceConfigurationsId200JSONResponse(converted), nil
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
	for {
		batch, nextToken, err := fetch(ctx, pageToken, pageSize)
		if err != nil {
			return nil, err
		}
		items = append(items, batch...)
		if nextToken == "" {
			break
		}
		pageToken = nextToken
	}
	if items == nil {
		return []T{}, nil
	}
	return items, nil
}

func requestProblem(err error) *ProblemError {
	problem := NewProblem(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
	return NewProblemError(problem, err)
}

func responseProblem(err error) *ProblemError {
	problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), err.Error())
	return NewProblemError(problem, err)
}
