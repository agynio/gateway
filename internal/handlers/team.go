package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/agynio/gateway/internal/gen"
	"github.com/agynio/gateway/internal/platform"
)

const teamBasePath = "/team/v1"

func TeamBasePath() string {
	return teamBasePath
}

type PlatformClient interface {
	Do(ctx context.Context, method, path string, query url.Values, body any, out any) (int, error)
}

type Team struct {
	client PlatformClient
}

func NewTeam(client PlatformClient) *Team {
	if client == nil {
		panic("platform client is required")
	}
	return &Team{client: client}
}

func (t *Team) GetAgents(ctx context.Context, request gen.GetAgentsRequestObject) (gen.GetAgentsResponseObject, error) {
	query := url.Values{}
	if request.Params.Q != nil {
		query.Set("q", *request.Params.Q)
	}
	if request.Params.Page != nil {
		query.Set("page", strconv.Itoa(*request.Params.Page))
	}
	if request.Params.PerPage != nil {
		query.Set("perPage", strconv.Itoa(*request.Params.PerPage))
	}

	var resp gen.GetAgents200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/agents", query, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PostAgents(ctx context.Context, request gen.PostAgentsRequestObject) (gen.PostAgentsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PostAgents201JSONResponse
	status, err := t.do(ctx, http.MethodPost, "/agents", nil, request.Body, &resp)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusCreated {
		return nil, t.unexpectedStatus(status)
	}
	return resp, nil
}

func (t *Team) DeleteAgentsId(ctx context.Context, request gen.DeleteAgentsIdRequestObject) (gen.DeleteAgentsIdResponseObject, error) {
	status, err := t.do(ctx, http.MethodDelete, "/agents/"+request.Id.String(), nil, nil, nil)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusNoContent {
		return nil, t.unexpectedStatus(status)
	}
	return gen.DeleteAgentsId204Response{}, nil
}

func (t *Team) GetAgentsId(ctx context.Context, request gen.GetAgentsIdRequestObject) (gen.GetAgentsIdResponseObject, error) {
	var resp gen.GetAgentsId200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/agents/"+request.Id.String(), nil, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PatchAgentsId(ctx context.Context, request gen.PatchAgentsIdRequestObject) (gen.PatchAgentsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PatchAgentsId200JSONResponse
	if _, err := t.do(ctx, http.MethodPatch, "/agents/"+request.Id.String(), nil, request.Body, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) GetAttachments(ctx context.Context, request gen.GetAttachmentsRequestObject) (gen.GetAttachmentsResponseObject, error) {
	query := url.Values{}
	if request.Params.SourceType != nil {
		query.Set("sourceType", string(*request.Params.SourceType))
	}
	if request.Params.SourceId != nil {
		query.Set("sourceId", request.Params.SourceId.String())
	}
	if request.Params.TargetType != nil {
		query.Set("targetType", string(*request.Params.TargetType))
	}
	if request.Params.TargetId != nil {
		query.Set("targetId", request.Params.TargetId.String())
	}
	if request.Params.Kind != nil {
		query.Set("kind", string(*request.Params.Kind))
	}
	if request.Params.Page != nil {
		query.Set("page", strconv.Itoa(*request.Params.Page))
	}
	if request.Params.PerPage != nil {
		query.Set("perPage", strconv.Itoa(*request.Params.PerPage))
	}

	var resp gen.GetAttachments200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/attachments", query, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PostAttachments(ctx context.Context, request gen.PostAttachmentsRequestObject) (gen.PostAttachmentsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PostAttachments201JSONResponse
	status, err := t.do(ctx, http.MethodPost, "/attachments", nil, request.Body, &resp)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusCreated {
		return nil, t.unexpectedStatus(status)
	}
	return resp, nil
}

func (t *Team) DeleteAttachmentsId(ctx context.Context, request gen.DeleteAttachmentsIdRequestObject) (gen.DeleteAttachmentsIdResponseObject, error) {
	status, err := t.do(ctx, http.MethodDelete, "/attachments/"+request.Id.String(), nil, nil, nil)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusNoContent {
		return nil, t.unexpectedStatus(status)
	}
	return gen.DeleteAttachmentsId204Response{}, nil
}

func (t *Team) GetMcpServers(ctx context.Context, request gen.GetMcpServersRequestObject) (gen.GetMcpServersResponseObject, error) {
	query := url.Values{}
	if request.Params.Page != nil {
		query.Set("page", strconv.Itoa(*request.Params.Page))
	}
	if request.Params.PerPage != nil {
		query.Set("perPage", strconv.Itoa(*request.Params.PerPage))
	}

	var resp gen.GetMcpServers200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/mcp-servers", query, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PostMcpServers(ctx context.Context, request gen.PostMcpServersRequestObject) (gen.PostMcpServersResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PostMcpServers201JSONResponse
	status, err := t.do(ctx, http.MethodPost, "/mcp-servers", nil, request.Body, &resp)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusCreated {
		return nil, t.unexpectedStatus(status)
	}
	return resp, nil
}

func (t *Team) DeleteMcpServersId(ctx context.Context, request gen.DeleteMcpServersIdRequestObject) (gen.DeleteMcpServersIdResponseObject, error) {
	status, err := t.do(ctx, http.MethodDelete, "/mcp-servers/"+request.Id.String(), nil, nil, nil)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusNoContent {
		return nil, t.unexpectedStatus(status)
	}
	return gen.DeleteMcpServersId204Response{}, nil
}

func (t *Team) GetMcpServersId(ctx context.Context, request gen.GetMcpServersIdRequestObject) (gen.GetMcpServersIdResponseObject, error) {
	var resp gen.GetMcpServersId200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/mcp-servers/"+request.Id.String(), nil, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PatchMcpServersId(ctx context.Context, request gen.PatchMcpServersIdRequestObject) (gen.PatchMcpServersIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PatchMcpServersId200JSONResponse
	if _, err := t.do(ctx, http.MethodPatch, "/mcp-servers/"+request.Id.String(), nil, request.Body, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) GetMemoryBuckets(ctx context.Context, request gen.GetMemoryBucketsRequestObject) (gen.GetMemoryBucketsResponseObject, error) {
	query := url.Values{}
	if request.Params.Page != nil {
		query.Set("page", strconv.Itoa(*request.Params.Page))
	}
	if request.Params.PerPage != nil {
		query.Set("perPage", strconv.Itoa(*request.Params.PerPage))
	}

	var resp gen.GetMemoryBuckets200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/memory-buckets", query, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PostMemoryBuckets(ctx context.Context, request gen.PostMemoryBucketsRequestObject) (gen.PostMemoryBucketsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PostMemoryBuckets201JSONResponse
	status, err := t.do(ctx, http.MethodPost, "/memory-buckets", nil, request.Body, &resp)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusCreated {
		return nil, t.unexpectedStatus(status)
	}
	return resp, nil
}

func (t *Team) DeleteMemoryBucketsId(ctx context.Context, request gen.DeleteMemoryBucketsIdRequestObject) (gen.DeleteMemoryBucketsIdResponseObject, error) {
	status, err := t.do(ctx, http.MethodDelete, "/memory-buckets/"+request.Id.String(), nil, nil, nil)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusNoContent {
		return nil, t.unexpectedStatus(status)
	}
	return gen.DeleteMemoryBucketsId204Response{}, nil
}

func (t *Team) GetMemoryBucketsId(ctx context.Context, request gen.GetMemoryBucketsIdRequestObject) (gen.GetMemoryBucketsIdResponseObject, error) {
	var resp gen.GetMemoryBucketsId200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/memory-buckets/"+request.Id.String(), nil, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PatchMemoryBucketsId(ctx context.Context, request gen.PatchMemoryBucketsIdRequestObject) (gen.PatchMemoryBucketsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PatchMemoryBucketsId200JSONResponse
	if _, err := t.do(ctx, http.MethodPatch, "/memory-buckets/"+request.Id.String(), nil, request.Body, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) GetTools(ctx context.Context, request gen.GetToolsRequestObject) (gen.GetToolsResponseObject, error) {
	query := url.Values{}
	if request.Params.Type != nil {
		query.Set("type", string(*request.Params.Type))
	}
	if request.Params.Page != nil {
		query.Set("page", strconv.Itoa(*request.Params.Page))
	}
	if request.Params.PerPage != nil {
		query.Set("perPage", strconv.Itoa(*request.Params.PerPage))
	}

	var resp gen.GetTools200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/tools", query, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PostTools(ctx context.Context, request gen.PostToolsRequestObject) (gen.PostToolsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PostTools201JSONResponse
	status, err := t.do(ctx, http.MethodPost, "/tools", nil, request.Body, &resp)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusCreated {
		return nil, t.unexpectedStatus(status)
	}
	return resp, nil
}

func (t *Team) DeleteToolsId(ctx context.Context, request gen.DeleteToolsIdRequestObject) (gen.DeleteToolsIdResponseObject, error) {
	status, err := t.do(ctx, http.MethodDelete, "/tools/"+request.Id.String(), nil, nil, nil)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusNoContent {
		return nil, t.unexpectedStatus(status)
	}
	return gen.DeleteToolsId204Response{}, nil
}

func (t *Team) GetToolsId(ctx context.Context, request gen.GetToolsIdRequestObject) (gen.GetToolsIdResponseObject, error) {
	var resp gen.GetToolsId200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/tools/"+request.Id.String(), nil, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PatchToolsId(ctx context.Context, request gen.PatchToolsIdRequestObject) (gen.PatchToolsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PatchToolsId200JSONResponse
	if _, err := t.do(ctx, http.MethodPatch, "/tools/"+request.Id.String(), nil, request.Body, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) GetWorkspaceConfigurations(ctx context.Context, request gen.GetWorkspaceConfigurationsRequestObject) (gen.GetWorkspaceConfigurationsResponseObject, error) {
	query := url.Values{}
	if request.Params.Page != nil {
		query.Set("page", strconv.Itoa(*request.Params.Page))
	}
	if request.Params.PerPage != nil {
		query.Set("perPage", strconv.Itoa(*request.Params.PerPage))
	}

	var resp gen.GetWorkspaceConfigurations200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/workspace-configurations", query, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PostWorkspaceConfigurations(ctx context.Context, request gen.PostWorkspaceConfigurationsRequestObject) (gen.PostWorkspaceConfigurationsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PostWorkspaceConfigurations201JSONResponse
	status, err := t.do(ctx, http.MethodPost, "/workspace-configurations", nil, request.Body, &resp)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusCreated {
		return nil, t.unexpectedStatus(status)
	}
	return resp, nil
}

func (t *Team) DeleteWorkspaceConfigurationsId(ctx context.Context, request gen.DeleteWorkspaceConfigurationsIdRequestObject) (gen.DeleteWorkspaceConfigurationsIdResponseObject, error) {
	status, err := t.do(ctx, http.MethodDelete, "/workspace-configurations/"+request.Id.String(), nil, nil, nil)
	if err != nil {
		return nil, t.wrapError(err)
	}
	if status != http.StatusNoContent {
		return nil, t.unexpectedStatus(status)
	}
	return gen.DeleteWorkspaceConfigurationsId204Response{}, nil
}

func (t *Team) GetWorkspaceConfigurationsId(ctx context.Context, request gen.GetWorkspaceConfigurationsIdRequestObject) (gen.GetWorkspaceConfigurationsIdResponseObject, error) {
	var resp gen.GetWorkspaceConfigurationsId200JSONResponse
	if _, err := t.do(ctx, http.MethodGet, "/workspace-configurations/"+request.Id.String(), nil, nil, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) PatchWorkspaceConfigurationsId(ctx context.Context, request gen.PatchWorkspaceConfigurationsIdRequestObject) (gen.PatchWorkspaceConfigurationsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var resp gen.PatchWorkspaceConfigurationsId200JSONResponse
	if _, err := t.do(ctx, http.MethodPatch, "/workspace-configurations/"+request.Id.String(), nil, request.Body, &resp); err != nil {
		return nil, t.wrapError(err)
	}
	return resp, nil
}

func (t *Team) do(ctx context.Context, method, suffix string, query url.Values, body any, out any) (int, error) {
	path := joinPath(teamBasePath, suffix)
	if len(query) == 0 {
		query = nil
	}
	return t.client.Do(ctx, method, path, query, body, out)
}

func (t *Team) wrapError(err error) error {
	if err == nil {
		return nil
	}

	var platformErr *platform.Error
	if errors.As(err, &platformErr) {
		problem := problemFromPlatform(platformErr)
		return NewProblemError(problem, platformErr)
	}

	detail := strings.TrimSpace(err.Error())
	problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), detail, nil)
	return NewProblemError(problem, err)
}

func (t *Team) unexpectedStatus(status int) error {
	detail := fmt.Sprintf("upstream returned unexpected status %d", status)
	problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), detail, nil)
	return NewProblemError(problem, nil)
}

func problemFromPlatform(err *platform.Error) gen.Problem {
	status := err.Status
	if status <= 0 {
		status = http.StatusBadGateway
	}

	if err.Problem == nil {
		detail := strings.TrimSpace(string(err.Body))
		if detail == "" && err.Err != nil {
			detail = err.Err.Error()
		}
		return NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), detail, nil)
	}

	upstream := err.Problem
	title := strings.TrimSpace(upstream.Title)
	if title == "" {
		title = http.StatusText(status)
	}
	typ := strings.TrimSpace(upstream.Type)
	if typ == "" {
		typ = problemTypeDefault
	}

	problem := gen.Problem{
		Status: int32(status),
		Title:  title,
		Type:   typ,
	}

	if upstream.Status > 0 {
		problem.Status = int32(upstream.Status)
	}

	if upstream.Detail != nil {
		trimmed := strings.TrimSpace(*upstream.Detail)
		if trimmed != "" {
			problem.Detail = ptr(trimmed)
		}
	}

	if upstream.Instance != nil {
		trimmed := strings.TrimSpace(*upstream.Instance)
		if trimmed != "" {
			problem.Instance = ptr(trimmed)
		}
	}

	if len(upstream.Errors) > 0 {
		copied := make(map[string][]string, len(upstream.Errors))
		for key, values := range upstream.Errors {
			copied[key] = append([]string(nil), values...)
		}
		problem.Errors = &copied
	}

	return problem
}

func joinPath(prefix, suffix string) string {
	if suffix == "" {
		return prefix
	}
	if strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimRight(prefix, "/")
	}
	if !strings.HasPrefix(suffix, "/") {
		suffix = "/" + suffix
	}
	return prefix + suffix
}
