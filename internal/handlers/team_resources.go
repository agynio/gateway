package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"

	"github.com/agynio/gateway/internal/platform"
)

type resourceError struct {
	status int
	detail string
	err    error
}

func (e *resourceError) Error() string {
	if e == nil {
		return ""
	}
	if strings.TrimSpace(e.detail) != "" {
		return e.detail
	}
	if e.err != nil {
		return e.err.Error()
	}
	return http.StatusText(e.status)
}

func (e *resourceError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func newResourceError(status int, detail string, err error) *resourceError {
	if status <= 0 {
		status = http.StatusInternalServerError
	}
	detail = strings.TrimSpace(detail)
	if detail == "" {
		detail = http.StatusText(status)
	}
	return &resourceError{status: status, detail: detail, err: err}
}

const (
	teamMetadataKey           = "_team"
	teamMetadataResourceField = "resource"
	teamMetadataVersionField  = "version"
	teamMetadataCreatedAt     = "createdAt"
	teamMetadataUpdatedAt     = "updatedAt"
	teamMetadataVersion       = 1
)

const (
	toolMetadataResource         = "tool"
	mcpServerMetadataResource    = "mcpServer"
	workspaceMetadataResource    = "workspaceConfiguration"
	memoryBucketMetadataResource = "memoryBucket"
)

const (
	mcpServerTemplateName    = "mcpServer"
	workspaceTemplateName    = "workspace"
	memoryBucketTemplateName = "memory"
)

type teamMetadata struct {
	Resource  string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func ensureTeamMetadata(config map[string]any, resource string, now time.Time) teamMetadata {
	if config == nil {
		config = make(map[string]any)
	}
	metaValue, _ := config[teamMetadataKey]
	metaMap, _ := metaValue.(map[string]any)
	if metaMap == nil {
		metaMap = make(map[string]any)
	}
	metaMap[teamMetadataResourceField] = resource
	metaMap[teamMetadataVersionField] = teamMetadataVersion
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if _, ok := metaMap[teamMetadataCreatedAt]; !ok {
		metaMap[teamMetadataCreatedAt] = now.Format(time.RFC3339Nano)
	}
	metaMap[teamMetadataUpdatedAt] = now.Format(time.RFC3339Nano)
	config[teamMetadataKey] = metaMap

	parsed, err := parseTeamMetadata(config, resource)
	if err != nil {
		return teamMetadata{Resource: resource, CreatedAt: now, UpdatedAt: ptr(now)}
	}
	return parsed
}

func updateTeamMetadataTimestamp(config map[string]any, resource string, now time.Time) error {
	metaValue, ok := config[teamMetadataKey]
	if !ok {
		return fmt.Errorf("missing %s metadata", teamMetadataKey)
	}
	metaMap, ok := metaValue.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid %s metadata", teamMetadataKey)
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	metaMap[teamMetadataUpdatedAt] = now.Format(time.RFC3339Nano)
	config[teamMetadataKey] = metaMap

	_, err := parseTeamMetadata(config, resource)
	return err
}

func parseTeamMetadata(config map[string]any, expectedResource string) (teamMetadata, error) {
	metaValue, ok := config[teamMetadataKey]
	if !ok {
		return teamMetadata{}, fmt.Errorf("missing %s metadata", teamMetadataKey)
	}
	metaMap, ok := metaValue.(map[string]any)
	if !ok {
		return teamMetadata{}, fmt.Errorf("invalid %s metadata", teamMetadataKey)
	}

	resourceValue, ok := metaMap[teamMetadataResourceField].(string)
	if !ok || strings.TrimSpace(resourceValue) == "" {
		return teamMetadata{}, fmt.Errorf("missing %s.%s", teamMetadataKey, teamMetadataResourceField)
	}
	if resourceValue != expectedResource {
		return teamMetadata{}, fmt.Errorf("unexpected resource %s", resourceValue)
	}

	versionValue, ok := metaMap[teamMetadataVersionField]
	if !ok {
		return teamMetadata{}, fmt.Errorf("missing %s.%s", teamMetadataKey, teamMetadataVersionField)
	}
	var versionNumber int
	switch v := versionValue.(type) {
	case float64:
		versionNumber = int(v)
	case int:
		versionNumber = v
	case json.Number:
		parsed, err := strconv.Atoi(v.String())
		if err != nil {
			return teamMetadata{}, fmt.Errorf("invalid version: %w", err)
		}
		versionNumber = parsed
	default:
		return teamMetadata{}, fmt.Errorf("invalid %s.%s", teamMetadataKey, teamMetadataVersionField)
	}
	if versionNumber != teamMetadataVersion {
		return teamMetadata{}, fmt.Errorf("unsupported metadata version %d", versionNumber)
	}

	createdRaw, ok := metaMap[teamMetadataCreatedAt].(string)
	if !ok {
		return teamMetadata{}, fmt.Errorf("missing %s.%s", teamMetadataKey, teamMetadataCreatedAt)
	}
	createdAt, err := time.Parse(time.RFC3339Nano, createdRaw)
	if err != nil {
		return teamMetadata{}, fmt.Errorf("invalid createdAt: %w", err)
	}
	createdAt = createdAt.UTC()

	var updatedAtPtr *time.Time
	if updatedRaw, ok := metaMap[teamMetadataUpdatedAt].(string); ok && strings.TrimSpace(updatedRaw) != "" {
		parsed, err := time.Parse(time.RFC3339Nano, updatedRaw)
		if err != nil {
			return teamMetadata{}, fmt.Errorf("invalid updatedAt: %w", err)
		}
		parsed = parsed.UTC()
		updatedAtPtr = &parsed
	}

	return teamMetadata{
		Resource:  resourceValue,
		CreatedAt: createdAt,
		UpdatedAt: updatedAtPtr,
	}, nil
}

type graphService struct {
	client       PlatformClient
	graphRetries int
}

func newGraphService(client PlatformClient, graphRetries int) graphService {
	return graphService{client: client, graphRetries: normalizeGraphRetries(graphRetries)}
}

func (s graphService) fetchGraph(ctx context.Context) (*graphDocument, error) {
	var graph graphDocument
	status, err := s.client.Do(ctx, http.MethodGet, agentGraphPath, nil, nil, &graph)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected graph status %d", status)
	}
	return &graph, nil
}

func (s graphService) persistGraph(ctx context.Context, graph *graphDocument) (*graphDocument, error) {
	payload := newGraphWriteDocument(graph)
	var updated graphDocument
	status, err := s.client.Do(ctx, http.MethodPost, agentGraphPath, nil, &payload, &updated)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected graph status %d", status)
	}
	return &updated, nil
}

func cloneWithoutTeamMetadata(config map[string]any) map[string]any {
	cloned := cloneToInterfaceMap(config)
	delete(cloned, teamMetadataKey)
	return cloned
}

func extractConflictDetail(conflict *platform.Error) string {
	if conflict == nil {
		return "graph conflict"
	}
	if conflict.Problem != nil && conflict.Problem.Detail != nil {
		trimmed := strings.TrimSpace(*conflict.Problem.Detail)
		if trimmed != "" {
			return trimmed
		}
	}
	if len(conflict.Body) > 0 {
		trimmed := strings.TrimSpace(string(conflict.Body))
		if trimmed != "" {
			return trimmed
		}
	}
	if conflict.Err != nil {
		return conflict.Err.Error()
	}
	return "graph conflict"
}

var toolTemplateByType = map[string]string{
	"manage":             "manageTool",
	"memory":             "memoryTool",
	"shell_command":      "shellTool",
	"send_message":       "sendMessageTool",
	"send_slack_message": "sendSlackMessageTool",
	"remind_me":          "remindMeTool",
	"github_clone_repo":  "githubCloneRepoTool",
	"call_agent":         "callAgentTool",
	"finish":             "finishTool",
}

var toolTypeByTemplate = map[string]string{
	"manageTool":           "manage",
	"memoryTool":           "memory",
	"shellTool":            "shell_command",
	"sendMessageTool":      "send_message",
	"sendSlackMessageTool": "send_slack_message",
	"remindMeTool":         "remind_me",
	"githubCloneRepoTool":  "github_clone_repo",
	"callAgentTool":        "call_agent",
	"finishTool":           "finish",
}

type toolService struct {
	graph graphService
}

func newToolService(client PlatformClient, graphRetries int) *toolService {
	return &toolService{graph: newGraphService(client, graphRetries)}
}

type toolListParams struct {
	Query   string
	Type    string
	Page    int
	PerPage int
}

type toolCreateInput struct {
	Type        string
	Name        *string
	Description *string
	Config      map[string]any
}

type toolUpdateInput struct {
	Name        *string
	Description *string
	Config      *map[string]any
}

type toolResource struct {
	ID          uuid.UUID
	Type        string
	Name        *string
	Description *string
	Config      map[string]any
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type toolListResult struct {
	Items   []toolResource
	Page    int
	PerPage int
	Total   int
}

func (s *toolService) ListTools(ctx context.Context, params toolListParams) (toolListResult, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return toolListResult{}, err
	}

	tools := make([]toolResource, 0, len(graph.Nodes))
	for _, node := range graph.Nodes {
		tool, ok, convErr := toolFromNode(node)
		if convErr != nil {
			return toolListResult{}, newResourceError(http.StatusInternalServerError, fmt.Sprintf("invalid tool node %s", node.ID), convErr)
		}
		if !ok {
			continue
		}
		tools = append(tools, tool)
	}

	filtered := filterTools(tools, params)
	sortTools(filtered)

	page := params.Page
	if page < 1 {
		page = 1
	}
	perPage := params.PerPage
	if perPage < 1 {
		perPage = 20
	}

	total := len(filtered)
	start := (page - 1) * perPage
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}

	items := append([]toolResource(nil), filtered[start:end]...)

	return toolListResult{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}, nil
}

func (s *toolService) GetTool(ctx context.Context, id uuid.UUID) (toolResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return toolResource{}, err
	}

	tool, ok, convErr := toolFromGraphByID(graph, id)
	if convErr != nil {
		return toolResource{}, newResourceError(http.StatusInternalServerError, "invalid tool data", convErr)
	}
	if !ok {
		return toolResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("tool %s not found", id), nil)
	}
	return tool, nil
}

func (s *toolService) CreateTool(ctx context.Context, input toolCreateInput) (toolResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return toolResource{}, err
	}

	newID := uuid.New()
	now := time.Now().UTC()
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		node, buildErr := buildToolNode(newID, input, now)
		if buildErr != nil {
			return toolResource{}, newResourceError(http.StatusInternalServerError, "failed to build tool node", buildErr)
		}
		working.Nodes = append(working.Nodes, node)

		copyGraph := working
		updated, persistErr := s.graph.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return toolResource{}, err
					}
					continue
				}
				return toolResource{}, newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return toolResource{}, persistErr
		}

		tool, ok, convErr := toolFromGraphByID(updated, newID)
		if convErr != nil {
			return toolResource{}, newResourceError(http.StatusInternalServerError, "invalid tool data", convErr)
		}
		if !ok {
			return toolResource{}, newResourceError(http.StatusInternalServerError, "tool missing after create", nil)
		}
		return tool, nil
	}

	return toolResource{}, newResourceError(http.StatusConflict, "graph conflict", nil)
}

func (s *toolService) UpdateTool(ctx context.Context, id uuid.UUID, input toolUpdateInput) (toolResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return toolResource{}, err
	}
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		idx := findNodeIndex(working.Nodes, id.String())
		if idx == -1 {
			return toolResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("tool %s not found", id), nil)
		}

		node := working.Nodes[idx].Clone()
		if node.Config == nil {
			node.Config = make(map[string]any)
		}

		meta, metaErr := parseTeamMetadata(node.Config, toolMetadataResource)
		if metaErr != nil {
			return toolResource{}, newResourceError(http.StatusInternalServerError, "invalid tool metadata", metaErr)
		}

		if input.Name != nil {
			node.Config["name"] = *input.Name
		}
		if input.Description != nil {
			node.Config["description"] = *input.Description
		}
		if input.Config != nil {
			node.Config["config"] = sanitizeToolConfigForStorage(*input.Config)
		}

		if err := updateTeamMetadataTimestamp(node.Config, toolMetadataResource, time.Now().UTC()); err != nil {
			return toolResource{}, newResourceError(http.StatusInternalServerError, "failed to update metadata", err)
		}
		working.Nodes[idx] = node

		copyGraph := working
		updated, persistErr := s.graph.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return toolResource{}, err
					}
					continue
				}
				return toolResource{}, newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return toolResource{}, persistErr
		}

		tool, ok, convErr := toolFromGraphByID(updated, id)
		if convErr != nil {
			return toolResource{}, newResourceError(http.StatusInternalServerError, "invalid tool data", convErr)
		}
		if !ok {
			return toolResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("tool %s not found", id), nil)
		}
		tool.CreatedAt = meta.CreatedAt
		return tool, nil
	}

	return toolResource{}, newResourceError(http.StatusConflict, "graph conflict", nil)
}

func (s *toolService) DeleteTool(ctx context.Context, id uuid.UUID) error {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return err
	}
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		idx := findNodeIndex(working.Nodes, id.String())
		if idx == -1 {
			return newResourceError(http.StatusNotFound, fmt.Sprintf("tool %s not found", id), nil)
		}

		working.Nodes = append(working.Nodes[:idx], working.Nodes[idx+1:]...)
		filteredEdges, removed := removeEdgesForNode(working.Edges, id.String())
		working.Edges = filteredEdges
		if len(removed) > 0 {
			working.Variables = removeAttachmentVariables(working.Variables, removed)
		}

		copyGraph := working
		if _, persistErr := s.graph.persistGraph(ctx, &copyGraph); persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return err
					}
					continue
				}
				return newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return persistErr
		}
		return nil
	}

	return newResourceError(http.StatusConflict, "graph conflict", nil)
}

func filterTools(tools []toolResource, params toolListParams) []toolResource {
	trimmedQuery := strings.TrimSpace(strings.ToLower(params.Query))
	filtered := make([]toolResource, 0, len(tools))

	for _, tool := range tools {
		if params.Type != "" && tool.Type != params.Type {
			continue
		}
		if trimmedQuery != "" && !toolMatchesQuery(tool, trimmedQuery) {
			continue
		}
		filtered = append(filtered, tool)
	}
	return filtered
}

func toolMatchesQuery(tool toolResource, query string) bool {
	if tool.Name != nil && strings.Contains(strings.ToLower(*tool.Name), query) {
		return true
	}
	if tool.Description != nil && strings.Contains(strings.ToLower(*tool.Description), query) {
		return true
	}
	if strings.Contains(strings.ToLower(tool.ID.String()), query) {
		return true
	}
	return false
}

func sortTools(tools []toolResource) {
	sort.Slice(tools, func(i, j int) bool {
		left := ""
		if tools[i].Name != nil {
			left = strings.ToLower(*tools[i].Name)
		}
		right := ""
		if tools[j].Name != nil {
			right = strings.ToLower(*tools[j].Name)
		}
		if left == right {
			return tools[i].ID.String() < tools[j].ID.String()
		}
		return left < right
	})
}

func toolFromNode(node graphNode) (toolResource, bool, error) {
	toolType, ok := toolTypeByTemplate[node.Template]
	if !ok {
		return toolResource{}, false, nil
	}

	id, err := uuid.Parse(strings.TrimSpace(node.ID))
	if err != nil {
		return toolResource{}, false, fmt.Errorf("invalid tool id: %w", err)
	}

	config := cloneToInterfaceMap(node.Config)
	meta, metaErr := parseTeamMetadata(config, toolMetadataResource)
	if metaErr != nil {
		return toolResource{}, false, metaErr
	}

	clean := cloneWithoutTeamMetadata(node.Config)

	var namePtr *string
	if nameValue, ok := clean["name"].(string); ok {
		trimmed := strings.TrimSpace(nameValue)
		if trimmed != "" {
			namePtr = &trimmed
		}
	}

	var descriptionPtr *string
	if descValue, ok := clean["description"].(string); ok {
		description := descValue
		if strings.TrimSpace(description) != "" {
			descriptionPtr = &description
		}
	}

	var configMap map[string]any
	if rawConfig, ok := clean["config"].(map[string]any); ok {
		configMap = cloneToInterfaceMap(rawConfig)
	}
	if configMap == nil {
		configMap = make(map[string]any)
	}

	return toolResource{
		ID:          id,
		Type:        toolType,
		Name:        namePtr,
		Description: descriptionPtr,
		Config:      configMap,
		CreatedAt:   meta.CreatedAt,
		UpdatedAt:   meta.UpdatedAt,
	}, true, nil
}

func toolFromGraphByID(graph *graphDocument, id uuid.UUID) (toolResource, bool, error) {
	for _, node := range graph.Nodes {
		if strings.EqualFold(node.ID, id.String()) {
			return toolFromNode(node)
		}
	}
	return toolResource{}, false, nil
}

func buildToolNode(id uuid.UUID, input toolCreateInput, now time.Time) (graphNode, error) {
	template, ok := toolTemplateByType[input.Type]
	if !ok {
		return graphNode{}, fmt.Errorf("unsupported tool type %s", input.Type)
	}

	config := make(map[string]any)
	config["type"] = input.Type
	if input.Name != nil {
		config["name"] = *input.Name
	}
	if input.Description != nil {
		config["description"] = *input.Description
	}
	config["config"] = sanitizeToolConfigForStorage(input.Config)

	ensureTeamMetadata(config, toolMetadataResource, now)

	return graphNode{
		ID:       id.String(),
		Template: template,
		Config:   config,
	}, nil
}

func sanitizeToolConfigForStorage(source map[string]any) map[string]any {
	if len(source) == 0 {
		return make(map[string]any)
	}
	return cloneToInterfaceMap(source)
}

type toolValidator struct {
	toolSchema      *openapi3.SchemaRef
	paginatedSchema *openapi3.SchemaRef
}

func newToolValidator(spec *openapi3.T) (*toolValidator, error) {
	components := spec.Components
	if components == nil {
		return nil, fmt.Errorf("team spec missing components section")
	}

	toolSchema, ok := components.Schemas["Tool"]
	if !ok {
		return nil, fmt.Errorf("team spec missing Tool schema")
	}
	paginatedSchema, ok := components.Schemas["PaginatedTools"]
	if !ok {
		return nil, fmt.Errorf("team spec missing PaginatedTools schema")
	}

	return &toolValidator{
		toolSchema:      toolSchema,
		paginatedSchema: paginatedSchema,
	}, nil
}

func (v *toolValidator) ValidateTool(value any) error {
	if v == nil || v.toolSchema == nil || v.toolSchema.Value == nil {
		return fmt.Errorf("tool schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return v.toolSchema.Value.VisitJSON(normalized)
}

func (v *toolValidator) ValidatePaginatedTools(value any) error {
	if v == nil || v.paginatedSchema == nil || v.paginatedSchema.Value == nil {
		return fmt.Errorf("paginated tool schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return v.paginatedSchema.Value.VisitJSON(normalized)
}

type mcpServerService struct {
	graph graphService
}

func newMcpServerService(client PlatformClient, graphRetries int) *mcpServerService {
	return &mcpServerService{graph: newGraphService(client, graphRetries)}
}

type mcpServerListParams struct {
	Query   string
	Page    int
	PerPage int
}

type mcpServerCreateInput struct {
	Title       *string
	Description *string
	Config      map[string]any
}

type mcpServerUpdateInput struct {
	Title       *string
	Description *string
	Config      *map[string]any
}

type mcpServerResource struct {
	ID          uuid.UUID
	Title       *string
	Description *string
	Config      map[string]any
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type mcpServerListResult struct {
	Items   []mcpServerResource
	Page    int
	PerPage int
	Total   int
}

func (s *mcpServerService) ListMcpServers(ctx context.Context, params mcpServerListParams) (mcpServerListResult, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return mcpServerListResult{}, err
	}

	servers := make([]mcpServerResource, 0, len(graph.Nodes))
	for _, node := range graph.Nodes {
		server, ok, convErr := mcpServerFromNode(node)
		if convErr != nil {
			return mcpServerListResult{}, newResourceError(http.StatusInternalServerError, fmt.Sprintf("invalid mcpServer node %s", node.ID), convErr)
		}
		if !ok {
			continue
		}
		servers = append(servers, server)
	}

	filtered := filterMcpServers(servers, params.Query)
	sortMcpServers(filtered)

	page := params.Page
	if page < 1 {
		page = 1
	}
	perPage := params.PerPage
	if perPage < 1 {
		perPage = 20
	}

	total := len(filtered)
	start := (page - 1) * perPage
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}

	items := append([]mcpServerResource(nil), filtered[start:end]...)

	return mcpServerListResult{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}, nil
}

func (s *mcpServerService) GetMcpServer(ctx context.Context, id uuid.UUID) (mcpServerResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return mcpServerResource{}, err
	}

	server, ok, convErr := mcpServerFromGraphByID(graph, id)
	if convErr != nil {
		return mcpServerResource{}, newResourceError(http.StatusInternalServerError, "invalid mcpServer data", convErr)
	}
	if !ok {
		return mcpServerResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("mcpServer %s not found", id), nil)
	}
	return server, nil
}

func (s *mcpServerService) CreateMcpServer(ctx context.Context, input mcpServerCreateInput) (mcpServerResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return mcpServerResource{}, err
	}

	newID := uuid.New()
	now := time.Now().UTC()
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		node := buildMcpServerNode(newID, input, now)
		working.Nodes = append(working.Nodes, node)

		copyGraph := working
		updated, persistErr := s.graph.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return mcpServerResource{}, err
					}
					continue
				}
				return mcpServerResource{}, newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return mcpServerResource{}, persistErr
		}

		server, ok, convErr := mcpServerFromGraphByID(updated, newID)
		if convErr != nil {
			return mcpServerResource{}, newResourceError(http.StatusInternalServerError, "invalid mcpServer data", convErr)
		}
		if !ok {
			return mcpServerResource{}, newResourceError(http.StatusInternalServerError, "mcpServer missing after create", nil)
		}
		return server, nil
	}

	return mcpServerResource{}, newResourceError(http.StatusConflict, "graph conflict", nil)
}

func (s *mcpServerService) UpdateMcpServer(ctx context.Context, id uuid.UUID, input mcpServerUpdateInput) (mcpServerResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return mcpServerResource{}, err
	}
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		idx := findNodeIndex(working.Nodes, id.String())
		if idx == -1 {
			return mcpServerResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("mcpServer %s not found", id), nil)
		}

		node := working.Nodes[idx].Clone()
		if node.Template != mcpServerTemplateName {
			return mcpServerResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("mcpServer %s not found", id), nil)
		}

		if node.Config == nil {
			node.Config = make(map[string]any)
		}

		if _, err := parseTeamMetadata(node.Config, mcpServerMetadataResource); err != nil {
			return mcpServerResource{}, newResourceError(http.StatusInternalServerError, "invalid mcpServer metadata", err)
		}

		if input.Title != nil {
			node.Config["title"] = *input.Title
		}
		if input.Description != nil {
			node.Config["description"] = *input.Description
		}
		if input.Config != nil {
			node.Config["config"] = cloneToInterfaceMap(*input.Config)
		}

		if err := updateTeamMetadataTimestamp(node.Config, mcpServerMetadataResource, time.Now().UTC()); err != nil {
			return mcpServerResource{}, newResourceError(http.StatusInternalServerError, "failed to update metadata", err)
		}

		working.Nodes[idx] = node

		copyGraph := working
		updated, persistErr := s.graph.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return mcpServerResource{}, err
					}
					continue
				}
				return mcpServerResource{}, newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return mcpServerResource{}, persistErr
		}

		server, ok, convErr := mcpServerFromGraphByID(updated, id)
		if convErr != nil {
			return mcpServerResource{}, newResourceError(http.StatusInternalServerError, "invalid mcpServer data", convErr)
		}
		if !ok {
			return mcpServerResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("mcpServer %s not found", id), nil)
		}
		return server, nil
	}

	return mcpServerResource{}, newResourceError(http.StatusConflict, "graph conflict", nil)
}

func (s *mcpServerService) DeleteMcpServer(ctx context.Context, id uuid.UUID) error {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return err
	}
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		idx := findNodeIndex(working.Nodes, id.String())
		if idx == -1 {
			return newResourceError(http.StatusNotFound, fmt.Sprintf("mcpServer %s not found", id), nil)
		}
		if working.Nodes[idx].Template != mcpServerTemplateName {
			return newResourceError(http.StatusNotFound, fmt.Sprintf("mcpServer %s not found", id), nil)
		}

		working.Nodes = append(working.Nodes[:idx], working.Nodes[idx+1:]...)
		filteredEdges, removed := removeEdgesForNode(working.Edges, id.String())
		working.Edges = filteredEdges
		if len(removed) > 0 {
			working.Variables = removeAttachmentVariables(working.Variables, removed)
		}

		copyGraph := working
		if _, persistErr := s.graph.persistGraph(ctx, &copyGraph); persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return err
					}
					continue
				}
				return newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return persistErr
		}
		return nil
	}

	return newResourceError(http.StatusConflict, "graph conflict", nil)
}

func filterMcpServers(servers []mcpServerResource, query string) []mcpServerResource {
	trimmed := strings.TrimSpace(strings.ToLower(query))
	if trimmed == "" {
		return append([]mcpServerResource(nil), servers...)
	}

	filtered := make([]mcpServerResource, 0, len(servers))
	for _, server := range servers {
		if mcpServerMatchesQuery(server, trimmed) {
			filtered = append(filtered, server)
		}
	}
	return filtered
}

func mcpServerMatchesQuery(server mcpServerResource, query string) bool {
	if server.Title != nil && strings.Contains(strings.ToLower(*server.Title), query) {
		return true
	}
	if server.Description != nil && strings.Contains(strings.ToLower(*server.Description), query) {
		return true
	}
	if strings.Contains(strings.ToLower(server.ID.String()), query) {
		return true
	}
	return false
}

func sortMcpServers(servers []mcpServerResource) {
	sort.Slice(servers, func(i, j int) bool {
		left := ""
		if servers[i].Title != nil {
			left = strings.ToLower(*servers[i].Title)
		}
		right := ""
		if servers[j].Title != nil {
			right = strings.ToLower(*servers[j].Title)
		}
		if left == right {
			return servers[i].ID.String() < servers[j].ID.String()
		}
		return left < right
	})
}

func mcpServerFromNode(node graphNode) (mcpServerResource, bool, error) {
	if node.Template != mcpServerTemplateName {
		return mcpServerResource{}, false, nil
	}

	id, err := uuid.Parse(strings.TrimSpace(node.ID))
	if err != nil {
		return mcpServerResource{}, false, fmt.Errorf("invalid mcpServer id: %w", err)
	}

	config := cloneToInterfaceMap(node.Config)
	meta, metaErr := parseTeamMetadata(config, mcpServerMetadataResource)
	if metaErr != nil {
		return mcpServerResource{}, false, metaErr
	}

	clean := cloneWithoutTeamMetadata(node.Config)

	var titlePtr *string
	if titleValue, ok := clean["title"].(string); ok {
		trimmed := strings.TrimSpace(titleValue)
		if trimmed != "" {
			titlePtr = &trimmed
		}
	}

	var descriptionPtr *string
	if descValue, ok := clean["description"].(string); ok {
		desc := descValue
		if strings.TrimSpace(desc) != "" {
			descriptionPtr = &desc
		}
	}

	var configMap map[string]any
	if raw, ok := clean["config"].(map[string]any); ok {
		configMap = cloneToInterfaceMap(raw)
	}
	if configMap == nil {
		configMap = make(map[string]any)
	}

	return mcpServerResource{
		ID:          id,
		Title:       titlePtr,
		Description: descriptionPtr,
		Config:      configMap,
		CreatedAt:   meta.CreatedAt,
		UpdatedAt:   meta.UpdatedAt,
	}, true, nil
}

func mcpServerFromGraphByID(graph *graphDocument, id uuid.UUID) (mcpServerResource, bool, error) {
	for _, node := range graph.Nodes {
		if strings.EqualFold(node.ID, id.String()) {
			return mcpServerFromNode(node)
		}
	}
	return mcpServerResource{}, false, nil
}

func buildMcpServerNode(id uuid.UUID, input mcpServerCreateInput, now time.Time) graphNode {
	config := make(map[string]any)
	if input.Title != nil {
		config["title"] = *input.Title
	}
	if input.Description != nil {
		config["description"] = *input.Description
	}
	config["config"] = cloneToInterfaceMap(input.Config)
	ensureTeamMetadata(config, mcpServerMetadataResource, now)

	return graphNode{
		ID:       id.String(),
		Template: mcpServerTemplateName,
		Config:   config,
	}
}

type mcpServerValidator struct {
	schema          *openapi3.SchemaRef
	paginatedSchema *openapi3.SchemaRef
}

func newMcpServerValidator(spec *openapi3.T) (*mcpServerValidator, error) {
	components := spec.Components
	if components == nil {
		return nil, fmt.Errorf("team spec missing components section")
	}

	schema, ok := components.Schemas["McpServer"]
	if !ok {
		return nil, fmt.Errorf("team spec missing McpServer schema")
	}
	paginatedSchema, ok := components.Schemas["PaginatedMcpServers"]
	if !ok {
		return nil, fmt.Errorf("team spec missing PaginatedMcpServers schema")
	}

	return &mcpServerValidator{
		schema:          schema,
		paginatedSchema: paginatedSchema,
	}, nil
}

func (v *mcpServerValidator) ValidateMcpServer(value any) error {
	if v == nil || v.schema == nil || v.schema.Value == nil {
		return fmt.Errorf("mcpServer schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return v.schema.Value.VisitJSON(normalized)
}

func (v *mcpServerValidator) ValidatePaginatedMcpServers(value any) error {
	if v == nil || v.paginatedSchema == nil || v.paginatedSchema.Value == nil {
		return fmt.Errorf("paginated mcpServer schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return v.paginatedSchema.Value.VisitJSON(normalized)
}

type workspaceService struct {
	graph graphService
}

func newWorkspaceService(client PlatformClient, graphRetries int) *workspaceService {
	return &workspaceService{graph: newGraphService(client, graphRetries)}
}

type workspaceListParams struct {
	Query   string
	Page    int
	PerPage int
}

type workspaceCreateInput struct {
	Title       *string
	Description *string
	Config      map[string]any
}

type workspaceUpdateInput struct {
	Title       *string
	Description *string
	Config      *map[string]any
}

type workspaceResource struct {
	ID          uuid.UUID
	Title       *string
	Description *string
	Config      map[string]any
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type workspaceListResult struct {
	Items   []workspaceResource
	Page    int
	PerPage int
	Total   int
}

func (s *workspaceService) ListWorkspaceConfigurations(ctx context.Context, params workspaceListParams) (workspaceListResult, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return workspaceListResult{}, err
	}

	workspaces := make([]workspaceResource, 0, len(graph.Nodes))
	for _, node := range graph.Nodes {
		workspace, ok, convErr := workspaceFromNode(node)
		if convErr != nil {
			return workspaceListResult{}, newResourceError(http.StatusInternalServerError, fmt.Sprintf("invalid workspace node %s", node.ID), convErr)
		}
		if !ok {
			continue
		}
		workspaces = append(workspaces, workspace)
	}

	filtered := filterWorkspaces(workspaces, params.Query)
	sortWorkspaces(filtered)

	page := params.Page
	if page < 1 {
		page = 1
	}
	perPage := params.PerPage
	if perPage < 1 {
		perPage = 20
	}

	total := len(filtered)
	start := (page - 1) * perPage
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}

	items := append([]workspaceResource(nil), filtered[start:end]...)

	return workspaceListResult{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}, nil
}

func (s *workspaceService) GetWorkspaceConfiguration(ctx context.Context, id uuid.UUID) (workspaceResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return workspaceResource{}, err
	}

	workspace, ok, convErr := workspaceFromGraphByID(graph, id)
	if convErr != nil {
		return workspaceResource{}, newResourceError(http.StatusInternalServerError, "invalid workspace data", convErr)
	}
	if !ok {
		return workspaceResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("workspaceConfiguration %s not found", id), nil)
	}
	return workspace, nil
}

func (s *workspaceService) CreateWorkspaceConfiguration(ctx context.Context, input workspaceCreateInput) (workspaceResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return workspaceResource{}, err
	}

	newID := uuid.New()
	now := time.Now().UTC()
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		node := buildWorkspaceNode(newID, input, now)
		working.Nodes = append(working.Nodes, node)

		copyGraph := working
		updated, persistErr := s.graph.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return workspaceResource{}, err
					}
					continue
				}
				return workspaceResource{}, newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return workspaceResource{}, persistErr
		}

		workspace, ok, convErr := workspaceFromGraphByID(updated, newID)
		if convErr != nil {
			return workspaceResource{}, newResourceError(http.StatusInternalServerError, "invalid workspace data", convErr)
		}
		if !ok {
			return workspaceResource{}, newResourceError(http.StatusInternalServerError, "workspace missing after create", nil)
		}
		return workspace, nil
	}

	return workspaceResource{}, newResourceError(http.StatusConflict, "graph conflict", nil)
}

func (s *workspaceService) UpdateWorkspaceConfiguration(ctx context.Context, id uuid.UUID, input workspaceUpdateInput) (workspaceResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return workspaceResource{}, err
	}
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		idx := findNodeIndex(working.Nodes, id.String())
		if idx == -1 {
			return workspaceResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("workspaceConfiguration %s not found", id), nil)
		}

		node := working.Nodes[idx].Clone()
		if node.Template != workspaceTemplateName {
			return workspaceResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("workspaceConfiguration %s not found", id), nil)
		}

		if node.Config == nil {
			node.Config = make(map[string]any)
		}

		if _, err := parseTeamMetadata(node.Config, workspaceMetadataResource); err != nil {
			return workspaceResource{}, newResourceError(http.StatusInternalServerError, "invalid workspace metadata", err)
		}

		if input.Title != nil {
			node.Config["title"] = *input.Title
		}
		if input.Description != nil {
			node.Config["description"] = *input.Description
		}
		if input.Config != nil {
			node.Config["config"] = cloneToInterfaceMap(*input.Config)
		}

		if err := updateTeamMetadataTimestamp(node.Config, workspaceMetadataResource, time.Now().UTC()); err != nil {
			return workspaceResource{}, newResourceError(http.StatusInternalServerError, "failed to update metadata", err)
		}

		working.Nodes[idx] = node

		copyGraph := working
		updated, persistErr := s.graph.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return workspaceResource{}, err
					}
					continue
				}
				return workspaceResource{}, newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return workspaceResource{}, persistErr
		}

		workspace, ok, convErr := workspaceFromGraphByID(updated, id)
		if convErr != nil {
			return workspaceResource{}, newResourceError(http.StatusInternalServerError, "invalid workspace data", convErr)
		}
		if !ok {
			return workspaceResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("workspaceConfiguration %s not found", id), nil)
		}
		return workspace, nil
	}

	return workspaceResource{}, newResourceError(http.StatusConflict, "graph conflict", nil)
}

func (s *workspaceService) DeleteWorkspaceConfiguration(ctx context.Context, id uuid.UUID) error {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return err
	}
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		idx := findNodeIndex(working.Nodes, id.String())
		if idx == -1 {
			return newResourceError(http.StatusNotFound, fmt.Sprintf("workspaceConfiguration %s not found", id), nil)
		}
		if working.Nodes[idx].Template != workspaceTemplateName {
			return newResourceError(http.StatusNotFound, fmt.Sprintf("workspaceConfiguration %s not found", id), nil)
		}

		working.Nodes = append(working.Nodes[:idx], working.Nodes[idx+1:]...)
		filteredEdges, removed := removeEdgesForNode(working.Edges, id.String())
		working.Edges = filteredEdges
		if len(removed) > 0 {
			working.Variables = removeAttachmentVariables(working.Variables, removed)
		}

		copyGraph := working
		if _, persistErr := s.graph.persistGraph(ctx, &copyGraph); persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return err
					}
					continue
				}
				return newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return persistErr
		}
		return nil
	}

	return newResourceError(http.StatusConflict, "graph conflict", nil)
}

func filterWorkspaces(workspaces []workspaceResource, query string) []workspaceResource {
	trimmed := strings.TrimSpace(strings.ToLower(query))
	if trimmed == "" {
		return append([]workspaceResource(nil), workspaces...)
	}
	filtered := make([]workspaceResource, 0, len(workspaces))
	for _, workspace := range workspaces {
		if workspaceMatchesQuery(workspace, trimmed) {
			filtered = append(filtered, workspace)
		}
	}
	return filtered
}

func workspaceMatchesQuery(workspace workspaceResource, query string) bool {
	if workspace.Title != nil && strings.Contains(strings.ToLower(*workspace.Title), query) {
		return true
	}
	if workspace.Description != nil && strings.Contains(strings.ToLower(*workspace.Description), query) {
		return true
	}
	if strings.Contains(strings.ToLower(workspace.ID.String()), query) {
		return true
	}
	return false
}

func sortWorkspaces(workspaces []workspaceResource) {
	sort.Slice(workspaces, func(i, j int) bool {
		left := ""
		if workspaces[i].Title != nil {
			left = strings.ToLower(*workspaces[i].Title)
		}
		right := ""
		if workspaces[j].Title != nil {
			right = strings.ToLower(*workspaces[j].Title)
		}
		if left == right {
			return workspaces[i].ID.String() < workspaces[j].ID.String()
		}
		return left < right
	})
}

func workspaceFromNode(node graphNode) (workspaceResource, bool, error) {
	if node.Template != workspaceTemplateName {
		return workspaceResource{}, false, nil
	}

	id, err := uuid.Parse(strings.TrimSpace(node.ID))
	if err != nil {
		return workspaceResource{}, false, fmt.Errorf("invalid workspace id: %w", err)
	}

	config := cloneToInterfaceMap(node.Config)
	meta, metaErr := parseTeamMetadata(config, workspaceMetadataResource)
	if metaErr != nil {
		return workspaceResource{}, false, metaErr
	}

	clean := cloneWithoutTeamMetadata(node.Config)

	var titlePtr *string
	if titleValue, ok := clean["title"].(string); ok {
		trimmed := strings.TrimSpace(titleValue)
		if trimmed != "" {
			titlePtr = &trimmed
		}
	}

	var descriptionPtr *string
	if descValue, ok := clean["description"].(string); ok {
		desc := descValue
		if strings.TrimSpace(desc) != "" {
			descriptionPtr = &desc
		}
	}

	var configMap map[string]any
	if raw, ok := clean["config"].(map[string]any); ok {
		configMap = cloneToInterfaceMap(raw)
	}
	if configMap == nil {
		configMap = make(map[string]any)
	}

	return workspaceResource{
		ID:          id,
		Title:       titlePtr,
		Description: descriptionPtr,
		Config:      configMap,
		CreatedAt:   meta.CreatedAt,
		UpdatedAt:   meta.UpdatedAt,
	}, true, nil
}

func workspaceFromGraphByID(graph *graphDocument, id uuid.UUID) (workspaceResource, bool, error) {
	for _, node := range graph.Nodes {
		if strings.EqualFold(node.ID, id.String()) {
			return workspaceFromNode(node)
		}
	}
	return workspaceResource{}, false, nil
}

func buildWorkspaceNode(id uuid.UUID, input workspaceCreateInput, now time.Time) graphNode {
	config := make(map[string]any)
	if input.Title != nil {
		config["title"] = *input.Title
	}
	if input.Description != nil {
		config["description"] = *input.Description
	}
	config["config"] = cloneToInterfaceMap(input.Config)
	ensureTeamMetadata(config, workspaceMetadataResource, now)

	return graphNode{
		ID:       id.String(),
		Template: workspaceTemplateName,
		Config:   config,
	}
}

type workspaceValidator struct {
	schema          *openapi3.SchemaRef
	paginatedSchema *openapi3.SchemaRef
}

func newWorkspaceValidator(spec *openapi3.T) (*workspaceValidator, error) {
	components := spec.Components
	if components == nil {
		return nil, fmt.Errorf("team spec missing components section")
	}

	schema, ok := components.Schemas["WorkspaceConfiguration"]
	if !ok {
		return nil, fmt.Errorf("team spec missing WorkspaceConfiguration schema")
	}
	paginatedSchema, ok := components.Schemas["PaginatedWorkspaceConfigurations"]
	if !ok {
		return nil, fmt.Errorf("team spec missing PaginatedWorkspaceConfigurations schema")
	}

	return &workspaceValidator{
		schema:          schema,
		paginatedSchema: paginatedSchema,
	}, nil
}

func (v *workspaceValidator) ValidateWorkspace(value any) error {
	if v == nil || v.schema == nil || v.schema.Value == nil {
		return fmt.Errorf("workspace schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return v.schema.Value.VisitJSON(normalized)
}

func (v *workspaceValidator) ValidatePaginatedWorkspaces(value any) error {
	if v == nil || v.paginatedSchema == nil || v.paginatedSchema.Value == nil {
		return fmt.Errorf("paginated workspace schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return v.paginatedSchema.Value.VisitJSON(normalized)
}

type memoryBucketService struct {
	graph graphService
}

func newMemoryBucketService(client PlatformClient, graphRetries int) *memoryBucketService {
	return &memoryBucketService{graph: newGraphService(client, graphRetries)}
}

type memoryBucketListParams struct {
	Query   string
	Page    int
	PerPage int
}

type memoryBucketCreateInput struct {
	Title       *string
	Description *string
	Config      map[string]any
}

type memoryBucketUpdateInput struct {
	Title       *string
	Description *string
	Config      *map[string]any
}

type memoryBucketResource struct {
	ID          uuid.UUID
	Title       *string
	Description *string
	Config      map[string]any
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type memoryBucketListResult struct {
	Items   []memoryBucketResource
	Page    int
	PerPage int
	Total   int
}

func (s *memoryBucketService) ListMemoryBuckets(ctx context.Context, params memoryBucketListParams) (memoryBucketListResult, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return memoryBucketListResult{}, err
	}

	buckets := make([]memoryBucketResource, 0, len(graph.Nodes))
	for _, node := range graph.Nodes {
		bucket, ok, convErr := memoryBucketFromNode(node)
		if convErr != nil {
			return memoryBucketListResult{}, newResourceError(http.StatusInternalServerError, fmt.Sprintf("invalid memory bucket node %s", node.ID), convErr)
		}
		if !ok {
			continue
		}
		buckets = append(buckets, bucket)
	}

	filtered := filterMemoryBuckets(buckets, params.Query)
	sortMemoryBuckets(filtered)

	page := params.Page
	if page < 1 {
		page = 1
	}
	perPage := params.PerPage
	if perPage < 1 {
		perPage = 20
	}

	total := len(filtered)
	start := (page - 1) * perPage
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}

	items := append([]memoryBucketResource(nil), filtered[start:end]...)

	return memoryBucketListResult{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}, nil
}

func (s *memoryBucketService) GetMemoryBucket(ctx context.Context, id uuid.UUID) (memoryBucketResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return memoryBucketResource{}, err
	}

	bucket, ok, convErr := memoryBucketFromGraphByID(graph, id)
	if convErr != nil {
		return memoryBucketResource{}, newResourceError(http.StatusInternalServerError, "invalid memory bucket data", convErr)
	}
	if !ok {
		return memoryBucketResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("memoryBucket %s not found", id), nil)
	}
	return bucket, nil
}

func (s *memoryBucketService) CreateMemoryBucket(ctx context.Context, input memoryBucketCreateInput) (memoryBucketResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return memoryBucketResource{}, err
	}

	newID := uuid.New()
	now := time.Now().UTC()
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		node := buildMemoryBucketNode(newID, input, now)
		working.Nodes = append(working.Nodes, node)

		copyGraph := working
		updated, persistErr := s.graph.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return memoryBucketResource{}, err
					}
					continue
				}
				return memoryBucketResource{}, newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return memoryBucketResource{}, persistErr
		}

		bucket, ok, convErr := memoryBucketFromGraphByID(updated, newID)
		if convErr != nil {
			return memoryBucketResource{}, newResourceError(http.StatusInternalServerError, "invalid memory bucket data", convErr)
		}
		if !ok {
			return memoryBucketResource{}, newResourceError(http.StatusInternalServerError, "memory bucket missing after create", nil)
		}
		return bucket, nil
	}

	return memoryBucketResource{}, newResourceError(http.StatusConflict, "graph conflict", nil)
}

func (s *memoryBucketService) UpdateMemoryBucket(ctx context.Context, id uuid.UUID, input memoryBucketUpdateInput) (memoryBucketResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return memoryBucketResource{}, err
	}
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		idx := findNodeIndex(working.Nodes, id.String())
		if idx == -1 {
			return memoryBucketResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("memoryBucket %s not found", id), nil)
		}

		node := working.Nodes[idx].Clone()
		if node.Template != memoryBucketTemplateName {
			return memoryBucketResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("memoryBucket %s not found", id), nil)
		}

		if node.Config == nil {
			node.Config = make(map[string]any)
		}

		if _, err := parseTeamMetadata(node.Config, memoryBucketMetadataResource); err != nil {
			return memoryBucketResource{}, newResourceError(http.StatusInternalServerError, "invalid memory bucket metadata", err)
		}

		if input.Title != nil {
			node.Config["title"] = *input.Title
		}
		if input.Description != nil {
			node.Config["description"] = *input.Description
		}
		if input.Config != nil {
			node.Config["config"] = cloneToInterfaceMap(*input.Config)
		}

		if err := updateTeamMetadataTimestamp(node.Config, memoryBucketMetadataResource, time.Now().UTC()); err != nil {
			return memoryBucketResource{}, newResourceError(http.StatusInternalServerError, "failed to update metadata", err)
		}

		working.Nodes[idx] = node

		copyGraph := working
		updated, persistErr := s.graph.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return memoryBucketResource{}, err
					}
					continue
				}
				return memoryBucketResource{}, newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return memoryBucketResource{}, persistErr
		}

		bucket, ok, convErr := memoryBucketFromGraphByID(updated, id)
		if convErr != nil {
			return memoryBucketResource{}, newResourceError(http.StatusInternalServerError, "invalid memory bucket data", convErr)
		}
		if !ok {
			return memoryBucketResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("memoryBucket %s not found", id), nil)
		}
		return bucket, nil
	}

	return memoryBucketResource{}, newResourceError(http.StatusConflict, "graph conflict", nil)
}

func (s *memoryBucketService) DeleteMemoryBucket(ctx context.Context, id uuid.UUID) error {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return err
	}
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		idx := findNodeIndex(working.Nodes, id.String())
		if idx == -1 {
			return newResourceError(http.StatusNotFound, fmt.Sprintf("memoryBucket %s not found", id), nil)
		}
		if working.Nodes[idx].Template != memoryBucketTemplateName {
			return newResourceError(http.StatusNotFound, fmt.Sprintf("memoryBucket %s not found", id), nil)
		}

		working.Nodes = append(working.Nodes[:idx], working.Nodes[idx+1:]...)
		filteredEdges, removed := removeEdgesForNode(working.Edges, id.String())
		working.Edges = filteredEdges
		if len(removed) > 0 {
			working.Variables = removeAttachmentVariables(working.Variables, removed)
		}

		copyGraph := working
		if _, persistErr := s.graph.persistGraph(ctx, &copyGraph); persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return err
					}
					continue
				}
				return newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return persistErr
		}
		return nil
	}

	return newResourceError(http.StatusConflict, "graph conflict", nil)
}

func filterMemoryBuckets(buckets []memoryBucketResource, query string) []memoryBucketResource {
	trimmed := strings.TrimSpace(strings.ToLower(query))
	if trimmed == "" {
		return append([]memoryBucketResource(nil), buckets...)
	}
	filtered := make([]memoryBucketResource, 0, len(buckets))
	for _, bucket := range buckets {
		if memoryBucketMatchesQuery(bucket, trimmed) {
			filtered = append(filtered, bucket)
		}
	}
	return filtered
}

func memoryBucketMatchesQuery(bucket memoryBucketResource, query string) bool {
	if bucket.Title != nil && strings.Contains(strings.ToLower(*bucket.Title), query) {
		return true
	}
	if bucket.Description != nil && strings.Contains(strings.ToLower(*bucket.Description), query) {
		return true
	}
	if strings.Contains(strings.ToLower(bucket.ID.String()), query) {
		return true
	}
	return false
}

func sortMemoryBuckets(buckets []memoryBucketResource) {
	sort.Slice(buckets, func(i, j int) bool {
		left := ""
		if buckets[i].Title != nil {
			left = strings.ToLower(*buckets[i].Title)
		}
		right := ""
		if buckets[j].Title != nil {
			right = strings.ToLower(*buckets[j].Title)
		}
		if left == right {
			return buckets[i].ID.String() < buckets[j].ID.String()
		}
		return left < right
	})
}

func memoryBucketFromNode(node graphNode) (memoryBucketResource, bool, error) {
	if node.Template != memoryBucketTemplateName {
		return memoryBucketResource{}, false, nil
	}

	id, err := uuid.Parse(strings.TrimSpace(node.ID))
	if err != nil {
		return memoryBucketResource{}, false, fmt.Errorf("invalid memory bucket id: %w", err)
	}

	config := cloneToInterfaceMap(node.Config)
	meta, metaErr := parseTeamMetadata(config, memoryBucketMetadataResource)
	if metaErr != nil {
		return memoryBucketResource{}, false, metaErr
	}

	clean := cloneWithoutTeamMetadata(node.Config)

	var titlePtr *string
	if titleValue, ok := clean["title"].(string); ok {
		trimmed := strings.TrimSpace(titleValue)
		if trimmed != "" {
			titlePtr = &trimmed
		}
	}

	var descriptionPtr *string
	if descValue, ok := clean["description"].(string); ok {
		desc := descValue
		if strings.TrimSpace(desc) != "" {
			descriptionPtr = &desc
		}
	}

	var configMap map[string]any
	if raw, ok := clean["config"].(map[string]any); ok {
		configMap = cloneToInterfaceMap(raw)
	}
	if configMap == nil {
		configMap = make(map[string]any)
	}

	return memoryBucketResource{
		ID:          id,
		Title:       titlePtr,
		Description: descriptionPtr,
		Config:      configMap,
		CreatedAt:   meta.CreatedAt,
		UpdatedAt:   meta.UpdatedAt,
	}, true, nil
}

func memoryBucketFromGraphByID(graph *graphDocument, id uuid.UUID) (memoryBucketResource, bool, error) {
	for _, node := range graph.Nodes {
		if strings.EqualFold(node.ID, id.String()) {
			return memoryBucketFromNode(node)
		}
	}
	return memoryBucketResource{}, false, nil
}

func buildMemoryBucketNode(id uuid.UUID, input memoryBucketCreateInput, now time.Time) graphNode {
	config := make(map[string]any)
	if input.Title != nil {
		config["title"] = *input.Title
	}
	if input.Description != nil {
		config["description"] = *input.Description
	}
	config["config"] = cloneToInterfaceMap(input.Config)
	ensureTeamMetadata(config, memoryBucketMetadataResource, now)

	return graphNode{
		ID:       id.String(),
		Template: memoryBucketTemplateName,
		Config:   config,
	}
}

type memoryBucketValidator struct {
	schema          *openapi3.SchemaRef
	paginatedSchema *openapi3.SchemaRef
}

func newMemoryBucketValidator(spec *openapi3.T) (*memoryBucketValidator, error) {
	components := spec.Components
	if components == nil {
		return nil, fmt.Errorf("team spec missing components section")
	}

	schema, ok := components.Schemas["MemoryBucket"]
	if !ok {
		return nil, fmt.Errorf("team spec missing MemoryBucket schema")
	}
	paginatedSchema, ok := components.Schemas["PaginatedMemoryBuckets"]
	if !ok {
		return nil, fmt.Errorf("team spec missing PaginatedMemoryBuckets schema")
	}

	return &memoryBucketValidator{
		schema:          schema,
		paginatedSchema: paginatedSchema,
	}, nil
}

func (v *memoryBucketValidator) ValidateMemoryBucket(value any) error {
	if v == nil || v.schema == nil || v.schema.Value == nil {
		return fmt.Errorf("memory bucket schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return v.schema.Value.VisitJSON(normalized)
}

func (v *memoryBucketValidator) ValidatePaginatedMemoryBuckets(value any) error {
	if v == nil || v.paginatedSchema == nil || v.paginatedSchema.Value == nil {
		return fmt.Errorf("paginated memory bucket schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return v.paginatedSchema.Value.VisitJSON(normalized)
}

type attachmentService struct {
	graph graphService
}

func newAttachmentService(client PlatformClient, graphRetries int) *attachmentService {
	return &attachmentService{graph: newGraphService(client, graphRetries)}
}

type attachmentListParams struct {
	SourceType string
	SourceID   string
	TargetType string
	TargetID   string
	Kind       string
	Page       int
	PerPage    int
}

type attachmentCreateInput struct {
	Kind     string
	SourceID uuid.UUID
	TargetID uuid.UUID
}

type attachmentResource struct {
	ID         uuid.UUID
	Kind       string
	SourceType string
	SourceID   uuid.UUID
	TargetType string
	TargetID   uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

type attachmentListResult struct {
	Items   []attachmentResource
	Page    int
	PerPage int
	Total   int
}

type attachmentKindConfig struct {
	Kind         string
	SourceType   string
	TargetType   string
	SourceHandle string
	TargetHandle string
}

var supportedAttachmentKinds = map[string]attachmentKindConfig{
	"agent_tool": {
		Kind:         "agent_tool",
		SourceType:   "agent",
		TargetType:   "tool",
		SourceHandle: "tools",
		TargetHandle: "$self",
	},
	"agent_mcpServer": {
		Kind:         "agent_mcpServer",
		SourceType:   "agent",
		TargetType:   "mcpServer",
		SourceHandle: "mcp",
		TargetHandle: "$self",
	},
	"mcpServer_workspaceConfiguration": {
		Kind:         "mcpServer_workspaceConfiguration",
		SourceType:   "mcpServer",
		TargetType:   "workspaceConfiguration",
		SourceHandle: "workspace",
		TargetHandle: "$self",
	},
}

var unsupportedAttachmentKinds = map[string]struct{}{
	"agent_workspaceConfiguration": {},
	"agent_memoryBucket":           {},
}

type attachmentSignature struct {
	SourceType   string
	TargetType   string
	SourceHandle string
	TargetHandle string
}

var attachmentKindBySignature = func() map[attachmentSignature]attachmentKindConfig {
	m := make(map[attachmentSignature]attachmentKindConfig, len(supportedAttachmentKinds))
	for _, cfg := range supportedAttachmentKinds {
		sig := attachmentSignature{
			SourceType:   cfg.SourceType,
			TargetType:   cfg.TargetType,
			SourceHandle: cfg.SourceHandle,
			TargetHandle: cfg.TargetHandle,
		}
		m[sig] = cfg
	}
	return m
}()

func (s *attachmentService) ListAttachments(ctx context.Context, params attachmentListParams) (attachmentListResult, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return attachmentListResult{}, err
	}

	nodesByID := make(map[string]graphNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		nodesByID[node.ID] = node
	}

	varMap := variablesToMap(graph.Variables)
	attachments := make([]attachmentResource, 0, len(graph.Edges))
	for _, edge := range graph.Edges {
		attachment, ok, convErr := attachmentFromEdge(edge, nodesByID, varMap)
		if convErr != nil {
			return attachmentListResult{}, newResourceError(http.StatusInternalServerError, "invalid attachment data", convErr)
		}
		if !ok {
			continue
		}
		attachments = append(attachments, attachment)
	}

	filtered := filterAttachments(attachments, params)
	sortAttachments(filtered)

	page := params.Page
	if page < 1 {
		page = 1
	}
	perPage := params.PerPage
	if perPage < 1 {
		perPage = 20
	}

	total := len(filtered)
	start := (page - 1) * perPage
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}

	items := append([]attachmentResource(nil), filtered[start:end]...)

	return attachmentListResult{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}, nil
}

func (s *attachmentService) GetAttachment(ctx context.Context, id uuid.UUID) (attachmentResource, error) {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return attachmentResource{}, err
	}

	nodesByID := make(map[string]graphNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		nodesByID[node.ID] = node
	}
	varMap := variablesToMap(graph.Variables)

	for _, edge := range graph.Edges {
		if edge.ID == nil || !strings.EqualFold(*edge.ID, id.String()) {
			continue
		}
		attachment, ok, convErr := attachmentFromEdge(edge, nodesByID, varMap)
		if convErr != nil {
			return attachmentResource{}, newResourceError(http.StatusInternalServerError, "invalid attachment data", convErr)
		}
		if !ok {
			break
		}
		return attachment, nil
	}

	return attachmentResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("attachment %s not found", id), nil)
}

func (s *attachmentService) CreateAttachment(ctx context.Context, input attachmentCreateInput) (attachmentResource, error) {
	if _, unsupported := unsupportedAttachmentKinds[input.Kind]; unsupported {
		return attachmentResource{}, newResourceError(http.StatusNotImplemented, fmt.Sprintf("attachment kind %s not supported", input.Kind), nil)
	}
	cfg, ok := supportedAttachmentKinds[input.Kind]
	if !ok {
		return attachmentResource{}, newResourceError(http.StatusBadRequest, fmt.Sprintf("attachment kind %s is invalid", input.Kind), nil)
	}

	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return attachmentResource{}, err
	}

	nodesByID := make(map[string]graphNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		nodesByID[node.ID] = node
	}

	sourceNode, ok := nodesByID[input.SourceID.String()]
	if !ok {
		return attachmentResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("source %s not found", input.SourceID), nil)
	}
	sourceType, err := entityTypeForTemplate(sourceNode.Template)
	if err != nil {
		return attachmentResource{}, newResourceError(http.StatusInternalServerError, "invalid source node", err)
	}
	if sourceType != cfg.SourceType {
		return attachmentResource{}, newResourceError(http.StatusBadRequest, fmt.Sprintf("source must be %s", cfg.SourceType), nil)
	}

	targetNode, ok := nodesByID[input.TargetID.String()]
	if !ok {
		return attachmentResource{}, newResourceError(http.StatusNotFound, fmt.Sprintf("target %s not found", input.TargetID), nil)
	}
	targetType, err := entityTypeForTemplate(targetNode.Template)
	if err != nil {
		return attachmentResource{}, newResourceError(http.StatusInternalServerError, "invalid target node", err)
	}
	if targetType != cfg.TargetType {
		return attachmentResource{}, newResourceError(http.StatusBadRequest, fmt.Sprintf("target must be %s", cfg.TargetType), nil)
	}

	varMap := variablesToMap(graph.Variables)
	existingAttachments := make(map[string]attachmentResource)
	for _, edge := range graph.Edges {
		attachment, ok, convErr := attachmentFromEdge(edge, nodesByID, varMap)
		if convErr != nil {
			return attachmentResource{}, newResourceError(http.StatusInternalServerError, "invalid attachment data", convErr)
		}
		if !ok {
			continue
		}
		key := attachmentDedupKey(attachment.Kind, attachment.SourceID, attachment.TargetID)
		existingAttachments[key] = attachment
	}

	key := attachmentDedupKey(cfg.Kind, input.SourceID, input.TargetID)
	if _, exists := existingAttachments[key]; exists {
		return attachmentResource{}, newResourceError(http.StatusConflict, "attachment already exists", nil)
	}

	newID := uuid.New()
	now := time.Now().UTC()
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()

		edge := graphEdge{
			ID:           ptr(newID.String()),
			Source:       input.SourceID.String(),
			SourceHandle: cfg.SourceHandle,
			Target:       input.TargetID.String(),
			TargetHandle: cfg.TargetHandle,
		}
		working.Edges = append(working.Edges, edge)
		working.Variables = upsertVariable(working.Variables, attachmentCreatedAtKey(newID.String()), now.Format(time.RFC3339Nano))
		working.Variables = upsertVariable(working.Variables, attachmentUpdatedAtKey(newID.String()), now.Format(time.RFC3339Nano))

		copyGraph := working
		updated, persistErr := s.graph.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return attachmentResource{}, err
					}
					// refresh nodes and variables for subsequent iteration
					nodesByID = make(map[string]graphNode, len(graph.Nodes))
					for _, node := range graph.Nodes {
						nodesByID[node.ID] = node
					}
					varMap = variablesToMap(graph.Variables)
					continue
				}
				return attachmentResource{}, newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return attachmentResource{}, persistErr
		}

		nodesByID = make(map[string]graphNode, len(updated.Nodes))
		for _, node := range updated.Nodes {
			nodesByID[node.ID] = node
		}
		varMap = variablesToMap(updated.Variables)

		for _, edge := range updated.Edges {
			if edge.ID == nil || *edge.ID != newID.String() {
				continue
			}
			attachment, ok, convErr := attachmentFromEdge(edge, nodesByID, varMap)
			if convErr != nil {
				return attachmentResource{}, newResourceError(http.StatusInternalServerError, "invalid attachment data", convErr)
			}
			if !ok {
				return attachmentResource{}, newResourceError(http.StatusInternalServerError, "attachment missing after create", nil)
			}
			return attachment, nil
		}

		return attachmentResource{}, newResourceError(http.StatusInternalServerError, "attachment missing after create", nil)
	}

	return attachmentResource{}, newResourceError(http.StatusConflict, "graph conflict", nil)
}

func (s *attachmentService) DeleteAttachment(ctx context.Context, id uuid.UUID) error {
	graph, err := s.graph.fetchGraph(ctx)
	if err != nil {
		return err
	}

	found := false
	for _, edge := range graph.Edges {
		if edge.ID != nil && strings.EqualFold(*edge.ID, id.String()) {
			found = true
			break
		}
	}
	if !found {
		return newResourceError(http.StatusNotFound, fmt.Sprintf("attachment %s not found", id), nil)
	}
	maxAttempts := graphMutationAttempts(s.graph.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		filteredEdges := make([]graphEdge, 0, len(working.Edges))
		for _, edge := range working.Edges {
			if edge.ID != nil && strings.EqualFold(*edge.ID, id.String()) {
				continue
			}
			filteredEdges = append(filteredEdges, edge)
		}
		working.Edges = filteredEdges
		working.Variables = removeAttachmentVariables(working.Variables, []string{id.String()})

		copyGraph := working
		if _, persistErr := s.graph.persistGraph(ctx, &copyGraph); persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.graph.fetchGraph(ctx)
					if err != nil {
						return err
					}
					continue
				}
				return newResourceError(http.StatusConflict, extractConflictDetail(conflict), conflict)
			}
			return persistErr
		}
		return nil
	}

	return newResourceError(http.StatusConflict, "graph conflict", nil)
}

func filterAttachments(attachments []attachmentResource, params attachmentListParams) []attachmentResource {
	filtered := make([]attachmentResource, 0, len(attachments))
	for _, attachment := range attachments {
		if params.Kind != "" && attachment.Kind != params.Kind {
			continue
		}
		if params.SourceType != "" && attachment.SourceType != params.SourceType {
			continue
		}
		if params.SourceID != "" && !strings.EqualFold(attachment.SourceID.String(), params.SourceID) {
			continue
		}
		if params.TargetType != "" && attachment.TargetType != params.TargetType {
			continue
		}
		if params.TargetID != "" && !strings.EqualFold(attachment.TargetID.String(), params.TargetID) {
			continue
		}
		filtered = append(filtered, attachment)
	}
	return filtered
}

func sortAttachments(attachments []attachmentResource) {
	sort.Slice(attachments, func(i, j int) bool {
		left := attachments[i]
		right := attachments[j]
		if !left.CreatedAt.Equal(right.CreatedAt) {
			return left.CreatedAt.Before(right.CreatedAt)
		}
		return left.ID.String() < right.ID.String()
	})
}

func attachmentFromEdge(edge graphEdge, nodes map[string]graphNode, vars map[string]string) (attachmentResource, bool, error) {
	if edge.ID == nil {
		return attachmentResource{}, false, nil
	}
	rawEdgeID := *edge.ID

	sourceNode, ok := nodes[edge.Source]
	if !ok {
		return attachmentResource{}, false, fmt.Errorf("attachment source node %s missing", edge.Source)
	}
	sourceType, err := entityTypeForTemplate(sourceNode.Template)
	if err != nil {
		return attachmentResource{}, false, err
	}

	targetNode, ok := nodes[edge.Target]
	if !ok {
		return attachmentResource{}, false, fmt.Errorf("attachment target node %s missing", edge.Target)
	}
	targetType, err := entityTypeForTemplate(targetNode.Template)
	if err != nil {
		return attachmentResource{}, false, err
	}

	sig := attachmentSignature{
		SourceType:   sourceType,
		TargetType:   targetType,
		SourceHandle: edge.SourceHandle,
		TargetHandle: edge.TargetHandle,
	}
	cfg, ok := attachmentKindBySignature[sig]
	if !ok {
		return attachmentResource{}, false, nil
	}

	createdKey := attachmentCreatedAtKey(rawEdgeID)
	createdValue, ok := vars[createdKey]
	createdAt := time.Time{}
	if ok && strings.TrimSpace(createdValue) != "" {
		parsed, err := time.Parse(time.RFC3339Nano, createdValue)
		if err != nil {
			return attachmentResource{}, false, fmt.Errorf("attachment %s invalid createdAt: %w", rawEdgeID, err)
		}
		createdAt = parsed.UTC()
	}

	var updatedAtPtr *time.Time
	updatedValue, ok := vars[attachmentUpdatedAtKey(rawEdgeID)]
	if ok && strings.TrimSpace(updatedValue) != "" {
		parsed, err := time.Parse(time.RFC3339Nano, updatedValue)
		if err != nil {
			return attachmentResource{}, false, fmt.Errorf("attachment %s invalid updatedAt: %w", rawEdgeID, err)
		}
		parsed = parsed.UTC()
		updatedAtPtr = &parsed
	}

	sourceID, err := uuid.Parse(edge.Source)
	if err != nil {
		return attachmentResource{}, false, fmt.Errorf("invalid source id: %w", err)
	}
	targetID, err := uuid.Parse(edge.Target)
	if err != nil {
		return attachmentResource{}, false, fmt.Errorf("invalid target id: %w", err)
	}

	return attachmentResource{
		ID:         attachmentIDFromEdge(rawEdgeID),
		Kind:       cfg.Kind,
		SourceType: cfg.SourceType,
		SourceID:   sourceID,
		TargetType: cfg.TargetType,
		TargetID:   targetID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAtPtr,
	}, true, nil
}

func attachmentIDFromEdge(rawID string) uuid.UUID {
	parsed, err := uuid.Parse(rawID)
	if err == nil {
		return parsed
	}
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(rawID))
}

func entityTypeForTemplate(template string) (string, error) {
	switch template {
	case agentTemplateName:
		return "agent", nil
	case mcpServerTemplateName:
		return "mcpServer", nil
	case workspaceTemplateName:
		return "workspaceConfiguration", nil
	case memoryBucketTemplateName:
		return "memoryBucket", nil
	default:
		if toolType, ok := toolTypeByTemplate[template]; ok {
			_ = toolType
			return "tool", nil
		}
	}
	return "", fmt.Errorf("unknown template %s", template)
}

func variablesToMap(vars []graphVariable) map[string]string {
	result := make(map[string]string, len(vars))
	for _, v := range vars {
		result[v.Key] = v.Value
	}
	return result
}

func attachmentCreatedAtKey(id string) string {
	return fmt.Sprintf("team.attachment.%s.createdAt", id)
}

func attachmentUpdatedAtKey(id string) string {
	return fmt.Sprintf("team.attachment.%s.updatedAt", id)
}

func upsertVariable(vars []graphVariable, key, value string) []graphVariable {
	for i := range vars {
		if vars[i].Key == key {
			vars[i].Value = value
			return vars
		}
	}
	return append(vars, graphVariable{Key: key, Value: value})
}

func removeAttachmentVariables(vars []graphVariable, ids []string) []graphVariable {
	if len(ids) == 0 {
		return vars
	}
	lookup := make(map[string]struct{}, len(ids)*2)
	for _, id := range ids {
		lookup[attachmentCreatedAtKey(id)] = struct{}{}
		lookup[attachmentUpdatedAtKey(id)] = struct{}{}
	}
	filtered := vars[:0]
	for _, v := range vars {
		if _, drop := lookup[v.Key]; drop {
			continue
		}
		filtered = append(filtered, v)
	}
	return filtered
}

func attachmentDedupKey(kind string, sourceID uuid.UUID, targetID uuid.UUID) string {
	return fmt.Sprintf("%s/%s/%s", kind, sourceID.String(), targetID.String())
}

type attachmentValidator struct {
	schema          *openapi3.SchemaRef
	paginatedSchema *openapi3.SchemaRef
}

func newAttachmentValidator(spec *openapi3.T) (*attachmentValidator, error) {
	components := spec.Components
	if components == nil {
		return nil, fmt.Errorf("team spec missing components section")
	}

	schema, ok := components.Schemas["Attachment"]
	if !ok {
		return nil, fmt.Errorf("team spec missing Attachment schema")
	}
	paginatedSchema, ok := components.Schemas["PaginatedAttachments"]
	if !ok {
		return nil, fmt.Errorf("team spec missing PaginatedAttachments schema")
	}

	return &attachmentValidator{
		schema:          schema,
		paginatedSchema: paginatedSchema,
	}, nil
}

func (v *attachmentValidator) ValidateAttachment(value any) error {
	if v == nil || v.schema == nil || v.schema.Value == nil {
		return fmt.Errorf("attachment schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return v.schema.Value.VisitJSON(normalized)
}

func (v *attachmentValidator) ValidatePaginatedAttachments(value any) error {
	if v == nil || v.paginatedSchema == nil || v.paginatedSchema.Value == nil {
		return fmt.Errorf("paginated attachment schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return v.paginatedSchema.Value.VisitJSON(normalized)
}
