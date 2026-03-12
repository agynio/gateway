package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"

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

type platformRetryProvider interface {
	Retries() int
	RetriesConfigured() bool
}

type Team struct {
	client PlatformClient

	agentService          *agentService
	agentValidator        *agentValidator
	toolService           *toolService
	toolValidator         *toolValidator
	mcpServerService      *mcpServerService
	mcpServerValidator    *mcpServerValidator
	workspaceService      *workspaceService
	workspaceValidator    *workspaceValidator
	memoryBucketService   *memoryBucketService
	memoryBucketValidator *memoryBucketValidator
	attachmentService     *attachmentService
	attachmentValidator   *attachmentValidator
}

func NewTeam(client PlatformClient, spec *openapi3.T) *Team {
	if client == nil {
		panic("platform client is required")
	}
	if spec == nil {
		panic("team spec is required")
	}
	graphRetries := graphMutationRetries(client)
	agentValidator, err := newAgentValidator(spec)
	if err != nil {
		panic(fmt.Sprintf("initialize agent validator: %v", err))
	}
	toolValidator, err := newToolValidator(spec)
	if err != nil {
		panic(fmt.Sprintf("initialize tool validator: %v", err))
	}
	mcpValidator, err := newMcpServerValidator(spec)
	if err != nil {
		panic(fmt.Sprintf("initialize mcp validator: %v", err))
	}
	workspaceValidator, err := newWorkspaceValidator(spec)
	if err != nil {
		panic(fmt.Sprintf("initialize workspace validator: %v", err))
	}
	memoryValidator, err := newMemoryBucketValidator(spec)
	if err != nil {
		panic(fmt.Sprintf("initialize memory validator: %v", err))
	}
	attachmentValidator, err := newAttachmentValidator(spec)
	if err != nil {
		panic(fmt.Sprintf("initialize attachment validator: %v", err))
	}
	return &Team{
		client:                client,
		agentService:          newAgentService(client, graphRetries),
		agentValidator:        agentValidator,
		toolService:           newToolService(client, graphRetries),
		toolValidator:         toolValidator,
		mcpServerService:      newMcpServerService(client, graphRetries),
		mcpServerValidator:    mcpValidator,
		workspaceService:      newWorkspaceService(client, graphRetries),
		workspaceValidator:    workspaceValidator,
		memoryBucketService:   newMemoryBucketService(client, graphRetries),
		memoryBucketValidator: memoryValidator,
		attachmentService:     newAttachmentService(client, graphRetries),
		attachmentValidator:   attachmentValidator,
	}
}

func (t *Team) GetAgents(ctx context.Context, request gen.GetAgentsRequestObject) (gen.GetAgentsResponseObject, error) {
	params := agentListParams{
		Page:    1,
		PerPage: 20,
	}
	if request.Params.Q != nil {
		params.Query = strings.TrimSpace(*request.Params.Q)
	}
	if request.Params.Page != nil {
		params.Page = *request.Params.Page
	}
	if request.Params.PerPage != nil {
		params.PerPage = *request.Params.PerPage
	}

	result, err := t.agentService.ListAgents(ctx, params)
	if err != nil {
		return nil, t.handleAgentError(err)
	}

	payload := buildAgentsListPayload(result)
	if err := t.agentValidator.ValidatePaginatedAgents(payload); err != nil {
		return nil, t.agentValidationError(err)
	}

	var resp gen.GetAgents200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.agentSerializationError(err)
	}

	return resp, nil
}

func (t *Team) PostAgents(ctx context.Context, request gen.PostAgentsRequestObject) (gen.PostAgentsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	configMap, err := structToMap(request.Body.Config)
	if err != nil {
		return nil, t.agentSerializationError(err)
	}

	created, err := t.agentService.CreateAgent(ctx, agentCreateInput{
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      configMap,
	})
	if err != nil {
		return nil, t.handleAgentError(err)
	}

	payload := buildAgentPayload(created)
	if err := t.agentValidator.ValidateAgent(payload); err != nil {
		return nil, t.agentValidationError(err)
	}

	var resp gen.PostAgents201JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.agentSerializationError(err)
	}

	return resp, nil
}

func (t *Team) DeleteAgentsId(ctx context.Context, request gen.DeleteAgentsIdRequestObject) (gen.DeleteAgentsIdResponseObject, error) {
	if err := t.agentService.DeleteAgent(ctx, uuid.UUID(request.Id)); err != nil {
		return nil, t.handleAgentError(err)
	}
	return gen.DeleteAgentsId204Response{}, nil
}

func (t *Team) GetAgentsId(ctx context.Context, request gen.GetAgentsIdRequestObject) (gen.GetAgentsIdResponseObject, error) {
	agent, err := t.agentService.GetAgent(ctx, uuid.UUID(request.Id))
	if err != nil {
		return nil, t.handleAgentError(err)
	}

	payload := buildAgentPayload(agent)
	if err := t.agentValidator.ValidateAgent(payload); err != nil {
		return nil, t.agentValidationError(err)
	}

	var resp gen.GetAgentsId200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.agentSerializationError(err)
	}

	return resp, nil
}

func (t *Team) PatchAgentsId(ctx context.Context, request gen.PatchAgentsIdRequestObject) (gen.PatchAgentsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	var configMapPtr *map[string]any
	if request.Body.Config != nil {
		configMap, err := structToMap(*request.Body.Config)
		if err != nil {
			return nil, t.agentSerializationError(err)
		}
		configMapPtr = &configMap
	}

	updated, err := t.agentService.UpdateAgent(ctx, uuid.UUID(request.Id), agentUpdateInput{
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      configMapPtr,
	})
	if err != nil {
		return nil, t.handleAgentError(err)
	}

	payload := buildAgentPayload(updated)
	if err := t.agentValidator.ValidateAgent(payload); err != nil {
		return nil, t.agentValidationError(err)
	}

	var resp gen.PatchAgentsId200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.agentSerializationError(err)
	}

	return resp, nil
}

func (t *Team) handleAgentError(err error) error {
	if err == nil {
		return nil
	}

	var agentErr *agentError
	if errors.As(err, &agentErr) {
		status := agentErr.status
		if status <= 0 {
			status = http.StatusBadGateway
		}
		title := http.StatusText(status)
		if title == "" {
			title = http.StatusText(http.StatusBadGateway)
		}
		detail := strings.TrimSpace(agentErr.detail)
		problem := NewProblem(status, title, detail)
		return NewProblemError(problem, agentErr)
	}

	return t.wrapError(err)
}

func (t *Team) handleResourceError(err error) error {
	if err == nil {
		return nil
	}

	var resErr *resourceError
	if errors.As(err, &resErr) {
		status := resErr.status
		if status <= 0 {
			status = http.StatusBadGateway
		}
		title := http.StatusText(status)
		if title == "" {
			title = http.StatusText(http.StatusBadGateway)
		}
		detail := strings.TrimSpace(resErr.detail)
		problem := NewProblem(status, title, detail)
		return NewProblemError(problem, resErr)
	}

	return t.wrapError(err)
}

func (t *Team) resourceValidationError(resource string, err error) error {
	detail := strings.TrimSpace(err.Error())
	if detail == "" {
		detail = "response validation failed"
	}
	title := fmt.Sprintf("%s response validation failed", titleCase(resource))
	problem := NewProblem(http.StatusInternalServerError, title, detail)
	return NewProblemError(problem, err)
}

func (t *Team) resourceSerializationError(resource string, err error) error {
	title := fmt.Sprintf("%s serialization failed", titleCase(resource))
	detail := fmt.Sprintf("failed to serialize %s payload", resource)
	problem := NewProblem(http.StatusInternalServerError, title, detail)
	return NewProblemError(problem, err)
}

func (t *Team) agentValidationError(err error) error {
	detail := strings.TrimSpace(err.Error())
	if detail == "" {
		detail = "response validation failed"
	}
	problem := NewProblem(http.StatusInternalServerError, "Response validation failed", detail)
	return NewProblemError(problem, err)
}

func (t *Team) agentSerializationError(err error) error {
	problem := NewProblem(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "failed to serialize agent payload")
	return NewProblemError(problem, err)
}

type agentsListPayload struct {
	Items   []agentPayload `json:"items"`
	Page    int            `json:"page"`
	PerPage int            `json:"perPage"`
	Total   int            `json:"total"`
}

type agentPayload struct {
	ID          string         `json:"id"`
	Title       *string        `json:"title,omitempty"`
	Description *string        `json:"description,omitempty"`
	Config      map[string]any `json:"config"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
}

func buildAgentsListPayload(result agentListResult) agentsListPayload {
	items := make([]agentPayload, 0, len(result.Items))
	for _, agent := range result.Items {
		items = append(items, buildAgentPayload(agent))
	}
	return agentsListPayload{
		Items:   items,
		Page:    result.Page,
		PerPage: result.PerPage,
		Total:   result.Total,
	}
}

func buildAgentPayload(agent agentResource) agentPayload {
	clonedConfig := cloneToInterfaceMap(agent.Config)
	return agentPayload{
		ID:          agent.ID.String(),
		Title:       agent.Title,
		Description: agent.Description,
		Config:      clonedConfig,
		CreatedAt:   agent.CreatedAt,
		UpdatedAt:   agent.UpdatedAt,
	}
}

type toolsListPayload struct {
	Items   []toolPayload `json:"items"`
	Page    int           `json:"page"`
	PerPage int           `json:"perPage"`
	Total   int           `json:"total"`
}

type toolPayload struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
}

func buildToolsListPayload(result toolListResult) toolsListPayload {
	items := make([]toolPayload, 0, len(result.Items))
	for _, tool := range result.Items {
		items = append(items, buildToolPayload(tool))
	}
	return toolsListPayload{
		Items:   items,
		Page:    result.Page,
		PerPage: result.PerPage,
		Total:   result.Total,
	}
}

func buildToolPayload(tool toolResource) toolPayload {
	config := cloneToInterfaceMap(tool.Config)
	if len(config) == 0 {
		config = nil
	}
	return toolPayload{
		ID:          tool.ID.String(),
		Type:        tool.Type,
		Name:        tool.Name,
		Description: tool.Description,
		Config:      config,
		CreatedAt:   tool.CreatedAt,
		UpdatedAt:   tool.UpdatedAt,
	}
}

type mcpServersListPayload struct {
	Items   []mcpServerPayload `json:"items"`
	Page    int                `json:"page"`
	PerPage int                `json:"perPage"`
	Total   int                `json:"total"`
}

type mcpServerPayload struct {
	ID          string         `json:"id"`
	Title       *string        `json:"title,omitempty"`
	Description *string        `json:"description,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
}

func buildMcpServersListPayload(result mcpServerListResult) mcpServersListPayload {
	items := make([]mcpServerPayload, 0, len(result.Items))
	for _, server := range result.Items {
		items = append(items, buildMcpServerPayload(server))
	}
	return mcpServersListPayload{
		Items:   items,
		Page:    result.Page,
		PerPage: result.PerPage,
		Total:   result.Total,
	}
}

func buildMcpServerPayload(server mcpServerResource) mcpServerPayload {
	config := cloneToInterfaceMap(server.Config)
	if len(config) == 0 {
		config = nil
	}
	return mcpServerPayload{
		ID:          server.ID.String(),
		Title:       server.Title,
		Description: server.Description,
		Config:      config,
		CreatedAt:   server.CreatedAt,
		UpdatedAt:   server.UpdatedAt,
	}
}

type workspacesListPayload struct {
	Items   []workspacePayload `json:"items"`
	Page    int                `json:"page"`
	PerPage int                `json:"perPage"`
	Total   int                `json:"total"`
}

type workspacePayload struct {
	ID          string         `json:"id"`
	Title       *string        `json:"title,omitempty"`
	Description *string        `json:"description,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
}

func buildWorkspacesListPayload(result workspaceListResult) workspacesListPayload {
	items := make([]workspacePayload, 0, len(result.Items))
	for _, workspace := range result.Items {
		items = append(items, buildWorkspacePayload(workspace))
	}
	return workspacesListPayload{
		Items:   items,
		Page:    result.Page,
		PerPage: result.PerPage,
		Total:   result.Total,
	}
}

func buildWorkspacePayload(workspace workspaceResource) workspacePayload {
	config := cloneToInterfaceMap(workspace.Config)
	if len(config) == 0 {
		config = nil
	}
	return workspacePayload{
		ID:          workspace.ID.String(),
		Title:       workspace.Title,
		Description: workspace.Description,
		Config:      config,
		CreatedAt:   workspace.CreatedAt,
		UpdatedAt:   workspace.UpdatedAt,
	}
}

type memoryBucketsListPayload struct {
	Items   []memoryBucketPayload `json:"items"`
	Page    int                   `json:"page"`
	PerPage int                   `json:"perPage"`
	Total   int                   `json:"total"`
}

type memoryBucketPayload struct {
	ID          string         `json:"id"`
	Title       *string        `json:"title,omitempty"`
	Description *string        `json:"description,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
}

func buildMemoryBucketsListPayload(result memoryBucketListResult) memoryBucketsListPayload {
	items := make([]memoryBucketPayload, 0, len(result.Items))
	for _, bucket := range result.Items {
		items = append(items, buildMemoryBucketPayload(bucket))
	}
	return memoryBucketsListPayload{
		Items:   items,
		Page:    result.Page,
		PerPage: result.PerPage,
		Total:   result.Total,
	}
}

func buildMemoryBucketPayload(bucket memoryBucketResource) memoryBucketPayload {
	config := cloneToInterfaceMap(bucket.Config)
	if len(config) == 0 {
		config = nil
	}
	return memoryBucketPayload{
		ID:          bucket.ID.String(),
		Title:       bucket.Title,
		Description: bucket.Description,
		Config:      config,
		CreatedAt:   bucket.CreatedAt,
		UpdatedAt:   bucket.UpdatedAt,
	}
}

type attachmentsListPayload struct {
	Items   []attachmentPayload `json:"items"`
	Page    int                 `json:"page"`
	PerPage int                 `json:"perPage"`
	Total   int                 `json:"total"`
}

type attachmentPayload struct {
	ID         string     `json:"id"`
	Kind       string     `json:"kind"`
	SourceType string     `json:"sourceType"`
	SourceID   string     `json:"sourceId"`
	TargetType string     `json:"targetType"`
	TargetID   string     `json:"targetId"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

func buildAttachmentsListPayload(result attachmentListResult) attachmentsListPayload {
	items := make([]attachmentPayload, 0, len(result.Items))
	for _, attachment := range result.Items {
		items = append(items, buildAttachmentPayload(attachment))
	}
	return attachmentsListPayload{
		Items:   items,
		Page:    result.Page,
		PerPage: result.PerPage,
		Total:   result.Total,
	}
}

func buildAttachmentPayload(attachment attachmentResource) attachmentPayload {
	return attachmentPayload{
		ID:         attachment.ID.String(),
		Kind:       attachment.Kind,
		SourceType: attachment.SourceType,
		SourceID:   attachment.SourceID.String(),
		TargetType: attachment.TargetType,
		TargetID:   attachment.TargetID.String(),
		CreatedAt:  attachment.CreatedAt,
		UpdatedAt:  attachment.UpdatedAt,
	}
}

func decodePayload(payload any, target any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func structToMap(value any) (map[string]any, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 || string(data) == "null" {
		return make(map[string]any), nil
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if result == nil {
		result = make(map[string]any)
	}
	return result, nil
}

type agentListParams struct {
	Query   string
	Page    int
	PerPage int
}

type agentCreateInput struct {
	Title       *string
	Description *string
	Config      map[string]any
}

type agentUpdateInput struct {
	Title       *string
	Description *string
	Config      *map[string]any
}

type agentResource struct {
	ID          uuid.UUID
	Title       *string
	Description *string
	Config      map[string]any
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type agentListResult struct {
	Items   []agentResource
	Page    int
	PerPage int
	Total   int
}

type agentService struct {
	client       PlatformClient
	graphRetries int
}

func newAgentService(client PlatformClient, graphRetries int) *agentService {
	return &agentService{client: client, graphRetries: normalizeGraphRetries(graphRetries)}
}

func (s *agentService) ListAgents(ctx context.Context, params agentListParams) (agentListResult, error) {
	graph, err := s.fetchGraph(ctx)
	if err != nil {
		return agentListResult{}, err
	}

	agents := make([]agentResource, 0, len(graph.Nodes))
	for _, node := range graph.Nodes {
		agent, ok, convErr := agentFromNode(node, graph.UpdatedAt)
		if convErr != nil {
			return agentListResult{}, newAgentGraphError(convErr)
		}
		if !ok {
			continue
		}
		agents = append(agents, agent)
	}

	filtered := filterAgents(agents, params.Query)
	sortAgents(filtered)

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

	items := append([]agentResource(nil), filtered[start:end]...)

	return agentListResult{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}, nil
}

func (s *agentService) GetAgent(ctx context.Context, id uuid.UUID) (agentResource, error) {
	graph, err := s.fetchGraph(ctx)
	if err != nil {
		return agentResource{}, err
	}

	agent, ok, convErr := agentFromGraphByID(graph, id)
	if convErr != nil {
		return agentResource{}, newAgentGraphError(convErr)
	}
	if !ok {
		return agentResource{}, newAgentNotFoundError(id)
	}

	return agent, nil
}

func (s *agentService) CreateAgent(ctx context.Context, input agentCreateInput) (agentResource, error) {
	graph, err := s.fetchGraph(ctx)
	if err != nil {
		return agentResource{}, err
	}

	newID := uuid.New()
	now := time.Now().UTC()
	maxAttempts := graphMutationAttempts(s.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		node := buildAgentNode(newID, input, now)
		working.Nodes = append(working.Nodes, node)

		copyGraph := working
		updated, persistErr := s.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.fetchGraph(ctx)
					if err != nil {
						return agentResource{}, err
					}
					continue
				}
				return agentResource{}, newAgentConflictError(conflict)
			}
			return agentResource{}, persistErr
		}

		agent, ok, convErr := agentFromGraphByID(updated, newID)
		if convErr != nil {
			return agentResource{}, newAgentGraphError(convErr)
		}
		if !ok {
			return agentResource{}, newAgentGraphError(fmt.Errorf("agent %s missing after create", newID))
		}
		return agent, nil
	}

	return agentResource{}, newAgentConflictError(nil)
}

func (s *agentService) UpdateAgent(ctx context.Context, id uuid.UUID, input agentUpdateInput) (agentResource, error) {
	graph, err := s.fetchGraph(ctx)
	if err != nil {
		return agentResource{}, err
	}

	maxAttempts := graphMutationAttempts(s.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		idx := findNodeIndex(working.Nodes, id.String())
		if idx == -1 {
			return agentResource{}, newAgentNotFoundError(id)
		}

		node := working.Nodes[idx].Clone()
		if node.Config == nil {
			node.Config = make(map[string]any)
		}

		if input.Title != nil {
			node.Config["title"] = *input.Title
		}
		if input.Description != nil {
			node.Config["description"] = *input.Description
		}
		if input.Config != nil {
			node.Config["config"] = sanitizeAgentConfigForStorage(*input.Config)
		}

		ensureAgentMetadata(node.Config, working.UpdatedAt)
		node.Config[agentMetadataUpdatedAtKey] = time.Now().UTC().Format(time.RFC3339Nano)

		working.Nodes[idx] = node

		copyGraph := working
		updated, persistErr := s.persistGraph(ctx, &copyGraph)
		if persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.fetchGraph(ctx)
					if err != nil {
						return agentResource{}, err
					}
					continue
				}
				return agentResource{}, newAgentConflictError(conflict)
			}
			return agentResource{}, persistErr
		}

		agent, ok, convErr := agentFromGraphByID(updated, id)
		if convErr != nil {
			return agentResource{}, newAgentGraphError(convErr)
		}
		if !ok {
			return agentResource{}, newAgentNotFoundError(id)
		}
		return agent, nil
	}

	return agentResource{}, newAgentConflictError(nil)
}

func (s *agentService) DeleteAgent(ctx context.Context, id uuid.UUID) error {
	graph, err := s.fetchGraph(ctx)
	if err != nil {
		return err
	}

	maxAttempts := graphMutationAttempts(s.graphRetries)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		working := graph.Clone()
		idx := findNodeIndex(working.Nodes, id.String())
		if idx == -1 {
			return newAgentNotFoundError(id)
		}

		working.Nodes = append(working.Nodes[:idx], working.Nodes[idx+1:]...)
		filteredEdges, removed := removeEdgesForNode(working.Edges, id.String())
		working.Edges = filteredEdges
		if len(removed) > 0 {
			working.Variables = removeAttachmentVariables(working.Variables, removed)
		}

		copyGraph := working
		if _, persistErr := s.persistGraph(ctx, &copyGraph); persistErr != nil {
			if conflict, ok := isConflictError(persistErr); ok {
				if isVersionConflict(conflict) && attempt < maxAttempts-1 {
					graph, err = s.fetchGraph(ctx)
					if err != nil {
						return err
					}
					continue
				}
				return newAgentConflictError(conflict)
			}
			return persistErr
		}
		return nil
	}

	return newAgentConflictError(nil)
}

func (s *agentService) fetchGraph(ctx context.Context) (*graphDocument, error) {
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

func (s *agentService) persistGraph(ctx context.Context, graph *graphDocument) (*graphDocument, error) {
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

const (
	agentGraphPath             = "/api/graph"
	agentTemplateName          = "agent"
	agentMetadataResourceKey   = "_teamResource"
	agentMetadataResourceValue = "agent"
	agentMetadataVersionKey    = "_teamVersion"
	agentMetadataVersionValue  = 1
	agentMetadataCreatedAtKey  = "_teamCreatedAt"
	agentMetadataUpdatedAtKey  = "_teamUpdatedAt"
)

var allowedAgentConfigKeys = map[string]struct{}{
	"model":                     {},
	"systemPrompt":              {},
	"debounceMs":                {},
	"whenBusy":                  {},
	"processBuffer":             {},
	"sendFinalResponseToThread": {},
	"summarizationKeepTokens":   {},
	"summarizationMaxTokens":    {},
	"restrictOutput":            {},
	"restrictionMessage":        {},
	"restrictionMaxInjections":  {},
	"name":                      {},
	"role":                      {},
}

type graphDocument struct {
	Name      string          `json:"name"`
	Version   int             `json:"version"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Nodes     []graphNode     `json:"nodes"`
	Edges     []graphEdge     `json:"edges"`
	Variables []graphVariable `json:"variables,omitempty"`
}

type graphWriteDocument struct {
	Name      string          `json:"name"`
	Version   int             `json:"version"`
	Nodes     []graphNode     `json:"nodes"`
	Edges     []graphEdge     `json:"edges"`
	Variables []graphVariable `json:"variables"`
}

func (g *graphDocument) Clone() graphDocument {
	clone := graphDocument{
		Name:      g.Name,
		Version:   g.Version,
		UpdatedAt: g.UpdatedAt,
	}
	if len(g.Nodes) > 0 {
		clone.Nodes = make([]graphNode, len(g.Nodes))
		for i, node := range g.Nodes {
			clone.Nodes[i] = node.Clone()
		}
	}
	if len(g.Edges) > 0 {
		clone.Edges = make([]graphEdge, len(g.Edges))
		copy(clone.Edges, g.Edges)
	}
	if len(g.Variables) > 0 {
		clone.Variables = make([]graphVariable, len(g.Variables))
		copy(clone.Variables, g.Variables)
	}
	return clone
}

type graphNode struct {
	ID       string         `json:"id"`
	Template string         `json:"template"`
	Config   map[string]any `json:"config,omitempty"`
	State    map[string]any `json:"state,omitempty"`
	Position *graphPosition `json:"position,omitempty"`
}

func (n graphNode) Clone() graphNode {
	clone := graphNode{
		ID:       n.ID,
		Template: n.Template,
	}
	if len(n.Config) > 0 {
		clone.Config = cloneToInterfaceMap(n.Config)
	}
	if len(n.State) > 0 {
		clone.State = cloneToInterfaceMap(n.State)
	}
	if n.Position != nil {
		pos := *n.Position
		clone.Position = &pos
	}
	return clone
}

type graphEdge struct {
	ID           *string `json:"id,omitempty"`
	Source       string  `json:"source"`
	SourceHandle string  `json:"sourceHandle"`
	Target       string  `json:"target"`
	TargetHandle string  `json:"targetHandle"`
}

type graphVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type graphPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func buildAgentNode(id uuid.UUID, input agentCreateInput, now time.Time) graphNode {
	config := make(map[string]any)
	if input.Title != nil {
		config["title"] = *input.Title
	}
	if input.Description != nil {
		config["description"] = *input.Description
	}
	config["config"] = sanitizeAgentConfigForStorage(input.Config)
	ensureAgentMetadata(config, now)
	config[agentMetadataCreatedAtKey] = now.Format(time.RFC3339Nano)
	config[agentMetadataUpdatedAtKey] = now.Format(time.RFC3339Nano)

	return graphNode{
		ID:       id.String(),
		Template: agentTemplateName,
		Config:   config,
	}
}

func agentFromNode(node graphNode, graphUpdatedAt time.Time) (agentResource, bool, error) {
	if node.Template != agentTemplateName {
		return agentResource{}, false, nil
	}

	agentID, err := uuid.Parse(strings.TrimSpace(node.ID))
	if err != nil {
		return agentResource{}, false, fmt.Errorf("invalid agent id: %w", err)
	}

	config := cloneToInterfaceMap(node.Config)
	ensureAgentMetadata(config, graphUpdatedAt)

	var title *string
	if value, ok := config["title"].(string); ok {
		trimmed := strings.TrimSpace(value)
		title = &trimmed
	}

	var description *string
	if value, ok := config["description"].(string); ok {
		desc := value
		description = &desc
	}

	configValue, _ := config["config"].(map[string]any)
	sanitizedConfig := sanitizeAgentConfig(configValue)

	createdAt := parseMetadataTime(config, agentMetadataCreatedAtKey, graphUpdatedAt)
	updatedAtTime := parseMetadataTimePtr(config, agentMetadataUpdatedAtKey)

	return agentResource{
		ID:          agentID,
		Title:       title,
		Description: description,
		Config:      sanitizedConfig,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAtTime,
	}, true, nil
}

func agentFromGraphByID(graph *graphDocument, id uuid.UUID) (agentResource, bool, error) {
	for _, node := range graph.Nodes {
		if strings.EqualFold(node.ID, id.String()) {
			agent, ok, err := agentFromNode(node, graph.UpdatedAt)
			return agent, ok, err
		}
	}
	return agentResource{}, false, nil
}

func findNodeIndex(nodes []graphNode, id string) int {
	for idx, node := range nodes {
		if strings.EqualFold(node.ID, id) {
			return idx
		}
	}
	return -1
}

func removeEdgesForNode(edges []graphEdge, id string) ([]graphEdge, []string) {
	filtered := edges[:0]
	removed := make([]string, 0)
	for _, edge := range edges {
		if strings.EqualFold(edge.Source, id) || strings.EqualFold(edge.Target, id) {
			if edge.ID != nil {
				removed = append(removed, *edge.ID)
			}
			continue
		}
		filtered = append(filtered, edge)
	}
	return filtered, removed
}

func filterAgents(agents []agentResource, query string) []agentResource {
	trimmed := strings.TrimSpace(strings.ToLower(query))
	if trimmed == "" {
		return append([]agentResource(nil), agents...)
	}

	filtered := make([]agentResource, 0, len(agents))
	for _, agent := range agents {
		if agentMatchesQuery(agent, trimmed) {
			filtered = append(filtered, agent)
		}
	}
	return filtered
}

func agentMatchesQuery(agent agentResource, query string) bool {
	if agent.Title != nil && strings.Contains(strings.ToLower(*agent.Title), query) {
		return true
	}
	if agent.Description != nil && strings.Contains(strings.ToLower(*agent.Description), query) {
		return true
	}
	return false
}

func sortAgents(agents []agentResource) {
	sort.Slice(agents, func(i, j int) bool {
		left := ""
		if agents[i].Title != nil {
			left = strings.ToLower(*agents[i].Title)
		}
		right := ""
		if agents[j].Title != nil {
			right = strings.ToLower(*agents[j].Title)
		}
		if left == right {
			return agents[i].ID.String() < agents[j].ID.String()
		}
		return left < right
	})
}

func sanitizeAgentConfig(source map[string]any) map[string]any {
	if len(source) == 0 {
		return make(map[string]any)
	}

	cloned := make(map[string]any)
	for key, value := range source {
		if _, ok := allowedAgentConfigKeys[key]; ok {
			cloned[key] = value
		}
	}
	return cloned
}

func sanitizeAgentConfigForStorage(source map[string]any) map[string]any {
	if len(source) == 0 {
		return make(map[string]any)
	}
	return sanitizeAgentConfig(source)
}

func ensureAgentMetadata(config map[string]any, fallback time.Time) {
	if config == nil {
		return
	}
	if _, ok := config[agentMetadataResourceKey]; !ok {
		config[agentMetadataResourceKey] = agentMetadataResourceValue
	}
	if _, ok := config[agentMetadataVersionKey]; !ok {
		config[agentMetadataVersionKey] = agentMetadataVersionValue
	}
	if _, ok := config[agentMetadataCreatedAtKey]; !ok {
		base := fallback
		if base.IsZero() {
			base = time.Now().UTC()
		}
		config[agentMetadataCreatedAtKey] = base.Format(time.RFC3339Nano)
	}
}

func parseMetadataTime(config map[string]any, key string, fallback time.Time) time.Time {
	if value, ok := config[key].(string); ok {
		if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
			return parsed.UTC()
		}
	}
	if !fallback.IsZero() {
		return fallback.UTC()
	}
	return time.Now().UTC()
}

func parseMetadataTimePtr(config map[string]any, key string) *time.Time {
	if value, ok := config[key].(string); ok {
		if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
			parsed = parsed.UTC()
			return &parsed
		}
	}
	return nil
}

func cloneToInterfaceMap(source map[string]any) map[string]any {
	if len(source) == 0 {
		return make(map[string]any)
	}
	data, err := json.Marshal(source)
	if err != nil {
		return make(map[string]any)
	}
	var cloned map[string]any
	if err := json.Unmarshal(data, &cloned); err != nil {
		return make(map[string]any)
	}
	if cloned == nil {
		cloned = make(map[string]any)
	}
	return cloned
}

type agentError struct {
	status int
	detail string
	err    error
}

func (e *agentError) Error() string {
	if e == nil {
		return ""
	}
	detail := strings.TrimSpace(e.detail)
	if detail != "" {
		return detail
	}
	if e.err != nil {
		return e.err.Error()
	}
	return http.StatusText(e.status)
}

func (e *agentError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func newAgentNotFoundError(id uuid.UUID) error {
	detail := fmt.Sprintf("agent %s not found", id)
	return &agentError{status: http.StatusNotFound, detail: detail}
}

func newAgentConflictError(conflict *platform.Error) error {
	detail := "graph conflict"
	if conflict != nil {
		if conflict.Problem != nil && conflict.Problem.Detail != nil {
			detail = strings.TrimSpace(*conflict.Problem.Detail)
		} else if len(conflict.Body) > 0 {
			detail = strings.TrimSpace(string(conflict.Body))
		}
	}
	return &agentError{status: http.StatusConflict, detail: detail, err: conflict}
}

func newAgentGraphError(err error) error {
	detail := "invalid graph response"
	if err != nil {
		detail = err.Error()
	}
	return &agentError{status: http.StatusBadGateway, detail: detail, err: err}
}

func isConflictError(err error) (*platform.Error, bool) {
	var platformErr *platform.Error
	if errors.As(err, &platformErr) && platformErr.Status == http.StatusConflict {
		return platformErr, true
	}
	return nil, false
}

func isVersionConflict(conflict *platform.Error) bool {
	if conflict == nil || conflict.Problem == nil {
		return false
	}
	if strings.EqualFold(conflict.Problem.Type, "VERSION_CONFLICT") {
		return true
	}
	if strings.EqualFold(conflict.Problem.Title, "VERSION_CONFLICT") {
		return true
	}
	if conflict.Problem.Detail != nil {
		return strings.EqualFold(strings.TrimSpace(*conflict.Problem.Detail), "VERSION_CONFLICT")
	}
	return false
}

const defaultGraphMutationRetries = 2

func graphMutationRetries(client PlatformClient) int {
	if client == nil {
		return defaultGraphMutationRetries
	}
	provider, ok := client.(platformRetryProvider)
	if !ok {
		return defaultGraphMutationRetries
	}
	if !provider.RetriesConfigured() {
		return defaultGraphMutationRetries
	}
	retries := provider.Retries()
	if retries < 0 {
		return defaultGraphMutationRetries
	}
	return retries
}

func normalizeGraphRetries(retries int) int {
	if retries < 0 {
		return 0
	}
	return retries
}

func graphMutationAttempts(retries int) int {
	normalized := normalizeGraphRetries(retries)
	return normalized + 1
}

func newGraphWriteDocument(graph *graphDocument) graphWriteDocument {
	if graph == nil {
		return graphWriteDocument{
			Name:      "main",
			Nodes:     []graphNode{},
			Edges:     []graphEdge{},
			Variables: []graphVariable{},
		}
	}
	name := strings.TrimSpace(graph.Name)
	if name == "" {
		name = "main"
	}
	nodes := graph.Nodes
	if nodes == nil {
		nodes = []graphNode{}
	}
	edges := graph.Edges
	if edges == nil {
		edges = []graphEdge{}
	}
	variables := graph.Variables
	if variables == nil {
		variables = []graphVariable{}
	}
	return graphWriteDocument{
		Name:      name,
		Version:   graph.Version,
		Nodes:     nodes,
		Edges:     edges,
		Variables: variables,
	}
}

type agentValidator struct {
	agentSchema     *openapi3.SchemaRef
	paginatedSchema *openapi3.SchemaRef
}

func newAgentValidator(spec *openapi3.T) (*agentValidator, error) {
	components := spec.Components
	if components == nil {
		return nil, fmt.Errorf("team spec missing components section")
	}

	agentSchema, ok := components.Schemas["Agent"]
	if !ok {
		return nil, fmt.Errorf("team spec missing Agent schema")
	}
	paginatedSchema, ok := components.Schemas["PaginatedAgents"]
	if !ok {
		return nil, fmt.Errorf("team spec missing PaginatedAgents schema")
	}

	return &agentValidator{
		agentSchema:     agentSchema,
		paginatedSchema: paginatedSchema,
	}, nil
}

func (v *agentValidator) ValidateAgent(value any) error {
	return v.validate(v.agentSchema, value)
}

func (v *agentValidator) ValidatePaginatedAgents(value any) error {
	return v.validate(v.paginatedSchema, value)
}

func (v *agentValidator) validate(schema *openapi3.SchemaRef, value any) error {
	if schema == nil || schema.Value == nil {
		return fmt.Errorf("schema not initialized")
	}
	normalized, err := normalizeForValidation(value)
	if err != nil {
		return err
	}
	return schema.Value.VisitJSON(normalized)
}

func normalizeForValidation(value any) (any, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var normalized any
	if err := json.Unmarshal(data, &normalized); err != nil {
		return nil, err
	}
	return normalized, nil
}

func (t *Team) GetAttachments(ctx context.Context, request gen.GetAttachmentsRequestObject) (gen.GetAttachmentsResponseObject, error) {
	params := attachmentListParams{
		Page:    1,
		PerPage: 20,
	}
	if request.Params.SourceType != nil {
		params.SourceType = string(*request.Params.SourceType)
	}
	if request.Params.SourceId != nil {
		params.SourceID = request.Params.SourceId.String()
	}
	if request.Params.TargetType != nil {
		params.TargetType = string(*request.Params.TargetType)
	}
	if request.Params.TargetId != nil {
		params.TargetID = request.Params.TargetId.String()
	}
	if request.Params.Kind != nil {
		params.Kind = string(*request.Params.Kind)
	}
	if request.Params.Page != nil {
		params.Page = *request.Params.Page
	}
	if request.Params.PerPage != nil {
		params.PerPage = *request.Params.PerPage
	}

	result, err := t.attachmentService.ListAttachments(ctx, params)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildAttachmentsListPayload(result)
	if err := t.attachmentValidator.ValidatePaginatedAttachments(payload); err != nil {
		return nil, t.resourceValidationError("attachment", err)
	}

	var resp gen.GetAttachments200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("attachment", err)
	}

	return resp, nil
}

func (t *Team) PostAttachments(ctx context.Context, request gen.PostAttachmentsRequestObject) (gen.PostAttachmentsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	sourceID, err := uuid.Parse(request.Body.SourceId.String())
	if err != nil {
		return nil, t.resourceSerializationError("attachment", fmt.Errorf("invalid sourceId: %w", err))
	}
	targetID, err := uuid.Parse(request.Body.TargetId.String())
	if err != nil {
		return nil, t.resourceSerializationError("attachment", fmt.Errorf("invalid targetId: %w", err))
	}

	input := attachmentCreateInput{
		Kind:     string(request.Body.Kind),
		SourceID: sourceID,
		TargetID: targetID,
	}

	created, err := t.attachmentService.CreateAttachment(ctx, input)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildAttachmentPayload(created)
	if err := t.attachmentValidator.ValidateAttachment(payload); err != nil {
		return nil, t.resourceValidationError("attachment", err)
	}

	var resp gen.PostAttachments201JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("attachment", err)
	}

	return resp, nil
}

func (t *Team) DeleteAttachmentsId(ctx context.Context, request gen.DeleteAttachmentsIdRequestObject) (gen.DeleteAttachmentsIdResponseObject, error) {
	attachmentID, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("attachment", fmt.Errorf("invalid attachment id: %w", err))
	}
	if err := t.attachmentService.DeleteAttachment(ctx, attachmentID); err != nil {
		return nil, t.handleResourceError(err)
	}
	return gen.DeleteAttachmentsId204Response{}, nil
}

func (t *Team) GetMcpServers(ctx context.Context, request gen.GetMcpServersRequestObject) (gen.GetMcpServersResponseObject, error) {
	params := mcpServerListParams{
		Page:    1,
		PerPage: 20,
	}
	if request.Params.Page != nil {
		params.Page = *request.Params.Page
	}
	if request.Params.PerPage != nil {
		params.PerPage = *request.Params.PerPage
	}

	result, err := t.mcpServerService.ListMcpServers(ctx, params)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildMcpServersListPayload(result)
	if err := t.mcpServerValidator.ValidatePaginatedMcpServers(payload); err != nil {
		return nil, t.resourceValidationError("MCP server", err)
	}

	var resp gen.GetMcpServers200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("MCP server", err)
	}

	return resp, nil
}

func (t *Team) PostMcpServers(ctx context.Context, request gen.PostMcpServersRequestObject) (gen.PostMcpServersResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	config, err := structToMap(request.Body.Config)
	if err != nil {
		return nil, t.resourceSerializationError("MCP server", fmt.Errorf("invalid config: %w", err))
	}

	input := mcpServerCreateInput{
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      config,
	}

	created, err := t.mcpServerService.CreateMcpServer(ctx, input)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildMcpServerPayload(created)
	if err := t.mcpServerValidator.ValidateMcpServer(payload); err != nil {
		return nil, t.resourceValidationError("MCP server", err)
	}

	var resp gen.PostMcpServers201JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("MCP server", err)
	}

	return resp, nil
}

func (t *Team) DeleteMcpServersId(ctx context.Context, request gen.DeleteMcpServersIdRequestObject) (gen.DeleteMcpServersIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("MCP server", fmt.Errorf("invalid MCP server id: %w", err))
	}
	if err := t.mcpServerService.DeleteMcpServer(ctx, id); err != nil {
		return nil, t.handleResourceError(err)
	}
	return gen.DeleteMcpServersId204Response{}, nil
}

func (t *Team) GetMcpServersId(ctx context.Context, request gen.GetMcpServersIdRequestObject) (gen.GetMcpServersIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("MCP server", fmt.Errorf("invalid MCP server id: %w", err))
	}

	server, err := t.mcpServerService.GetMcpServer(ctx, id)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildMcpServerPayload(server)
	if err := t.mcpServerValidator.ValidateMcpServer(payload); err != nil {
		return nil, t.resourceValidationError("MCP server", err)
	}

	var resp gen.GetMcpServersId200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("MCP server", err)
	}

	return resp, nil
}

func (t *Team) PatchMcpServersId(ctx context.Context, request gen.PatchMcpServersIdRequestObject) (gen.PatchMcpServersIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("MCP server", fmt.Errorf("invalid MCP server id: %w", err))
	}

	var configPtr *map[string]any
	if request.Body.Config != nil {
		converted, err := structToMap(*request.Body.Config)
		if err != nil {
			return nil, t.resourceSerializationError("MCP server", fmt.Errorf("invalid config: %w", err))
		}
		configPtr = &converted
	}

	input := mcpServerUpdateInput{
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      configPtr,
	}

	updated, err := t.mcpServerService.UpdateMcpServer(ctx, id, input)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildMcpServerPayload(updated)
	if err := t.mcpServerValidator.ValidateMcpServer(payload); err != nil {
		return nil, t.resourceValidationError("MCP server", err)
	}

	var resp gen.PatchMcpServersId200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("MCP server", err)
	}

	return resp, nil
}

func (t *Team) GetMemoryBuckets(ctx context.Context, request gen.GetMemoryBucketsRequestObject) (gen.GetMemoryBucketsResponseObject, error) {
	params := memoryBucketListParams{
		Page:    1,
		PerPage: 20,
	}
	if request.Params.Page != nil {
		params.Page = *request.Params.Page
	}
	if request.Params.PerPage != nil {
		params.PerPage = *request.Params.PerPage
	}

	result, err := t.memoryBucketService.ListMemoryBuckets(ctx, params)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildMemoryBucketsListPayload(result)
	if err := t.memoryBucketValidator.ValidatePaginatedMemoryBuckets(payload); err != nil {
		return nil, t.resourceValidationError("memory bucket", err)
	}

	var resp gen.GetMemoryBuckets200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("memory bucket", err)
	}

	return resp, nil
}

func (t *Team) PostMemoryBuckets(ctx context.Context, request gen.PostMemoryBucketsRequestObject) (gen.PostMemoryBucketsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	config, err := structToMap(request.Body.Config)
	if err != nil {
		return nil, t.resourceSerializationError("memory bucket", fmt.Errorf("invalid config: %w", err))
	}

	input := memoryBucketCreateInput{
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      config,
	}

	created, err := t.memoryBucketService.CreateMemoryBucket(ctx, input)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildMemoryBucketPayload(created)
	if err := t.memoryBucketValidator.ValidateMemoryBucket(payload); err != nil {
		return nil, t.resourceValidationError("memory bucket", err)
	}

	var resp gen.PostMemoryBuckets201JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("memory bucket", err)
	}

	return resp, nil
}

func (t *Team) DeleteMemoryBucketsId(ctx context.Context, request gen.DeleteMemoryBucketsIdRequestObject) (gen.DeleteMemoryBucketsIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("memory bucket", fmt.Errorf("invalid memory bucket id: %w", err))
	}
	if err := t.memoryBucketService.DeleteMemoryBucket(ctx, id); err != nil {
		return nil, t.handleResourceError(err)
	}
	return gen.DeleteMemoryBucketsId204Response{}, nil
}

func (t *Team) GetMemoryBucketsId(ctx context.Context, request gen.GetMemoryBucketsIdRequestObject) (gen.GetMemoryBucketsIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("memory bucket", fmt.Errorf("invalid memory bucket id: %w", err))
	}

	bucket, err := t.memoryBucketService.GetMemoryBucket(ctx, id)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildMemoryBucketPayload(bucket)
	if err := t.memoryBucketValidator.ValidateMemoryBucket(payload); err != nil {
		return nil, t.resourceValidationError("memory bucket", err)
	}

	var resp gen.GetMemoryBucketsId200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("memory bucket", err)
	}

	return resp, nil
}

func (t *Team) PatchMemoryBucketsId(ctx context.Context, request gen.PatchMemoryBucketsIdRequestObject) (gen.PatchMemoryBucketsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("memory bucket", fmt.Errorf("invalid memory bucket id: %w", err))
	}

	var configPtr *map[string]any
	if request.Body.Config != nil {
		converted, err := structToMap(*request.Body.Config)
		if err != nil {
			return nil, t.resourceSerializationError("memory bucket", fmt.Errorf("invalid config: %w", err))
		}
		configPtr = &converted
	}

	input := memoryBucketUpdateInput{
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      configPtr,
	}

	updated, err := t.memoryBucketService.UpdateMemoryBucket(ctx, id, input)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildMemoryBucketPayload(updated)
	if err := t.memoryBucketValidator.ValidateMemoryBucket(payload); err != nil {
		return nil, t.resourceValidationError("memory bucket", err)
	}

	var resp gen.PatchMemoryBucketsId200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("memory bucket", err)
	}

	return resp, nil
}

func (t *Team) GetTools(ctx context.Context, request gen.GetToolsRequestObject) (gen.GetToolsResponseObject, error) {
	params := toolListParams{
		Page:    1,
		PerPage: 20,
	}
	if request.Params.Type != nil {
		params.Type = string(*request.Params.Type)
	}
	if request.Params.Page != nil {
		params.Page = *request.Params.Page
	}
	if request.Params.PerPage != nil {
		params.PerPage = *request.Params.PerPage
	}

	result, err := t.toolService.ListTools(ctx, params)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildToolsListPayload(result)
	if err := t.toolValidator.ValidatePaginatedTools(payload); err != nil {
		return nil, t.resourceValidationError("tool", err)
	}

	var resp gen.GetTools200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("tool", err)
	}

	return resp, nil
}

func (t *Team) PostTools(ctx context.Context, request gen.PostToolsRequestObject) (gen.PostToolsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	config := make(map[string]any)
	if request.Body.Config != nil {
		config = cloneToInterfaceMap(*request.Body.Config)
	}

	input := toolCreateInput{
		Type:        string(request.Body.Type),
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Config:      config,
	}

	created, err := t.toolService.CreateTool(ctx, input)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildToolPayload(created)
	if err := t.toolValidator.ValidateTool(payload); err != nil {
		return nil, t.resourceValidationError("tool", err)
	}

	var resp gen.PostTools201JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("tool", err)
	}

	return resp, nil
}

func (t *Team) DeleteToolsId(ctx context.Context, request gen.DeleteToolsIdRequestObject) (gen.DeleteToolsIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("tool", fmt.Errorf("invalid tool id: %w", err))
	}
	if err := t.toolService.DeleteTool(ctx, id); err != nil {
		return nil, t.handleResourceError(err)
	}
	return gen.DeleteToolsId204Response{}, nil
}

func (t *Team) GetToolsId(ctx context.Context, request gen.GetToolsIdRequestObject) (gen.GetToolsIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("tool", fmt.Errorf("invalid tool id: %w", err))
	}

	tool, err := t.toolService.GetTool(ctx, id)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildToolPayload(tool)
	if err := t.toolValidator.ValidateTool(payload); err != nil {
		return nil, t.resourceValidationError("tool", err)
	}

	var resp gen.GetToolsId200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("tool", err)
	}

	return resp, nil
}

func (t *Team) PatchToolsId(ctx context.Context, request gen.PatchToolsIdRequestObject) (gen.PatchToolsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("tool", fmt.Errorf("invalid tool id: %w", err))
	}

	var configPtr *map[string]any
	if request.Body.Config != nil {
		cloned := cloneToInterfaceMap(*request.Body.Config)
		configPtr = &cloned
	}

	input := toolUpdateInput{
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Config:      configPtr,
	}

	updated, err := t.toolService.UpdateTool(ctx, id, input)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildToolPayload(updated)
	if err := t.toolValidator.ValidateTool(payload); err != nil {
		return nil, t.resourceValidationError("tool", err)
	}

	var resp gen.PatchToolsId200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("tool", err)
	}

	return resp, nil
}

func (t *Team) GetWorkspaceConfigurations(ctx context.Context, request gen.GetWorkspaceConfigurationsRequestObject) (gen.GetWorkspaceConfigurationsResponseObject, error) {
	params := workspaceListParams{
		Page:    1,
		PerPage: 20,
	}
	if request.Params.Page != nil {
		params.Page = *request.Params.Page
	}
	if request.Params.PerPage != nil {
		params.PerPage = *request.Params.PerPage
	}

	result, err := t.workspaceService.ListWorkspaceConfigurations(ctx, params)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildWorkspacesListPayload(result)
	if err := t.workspaceValidator.ValidatePaginatedWorkspaces(payload); err != nil {
		return nil, t.resourceValidationError("workspace configuration", err)
	}

	var resp gen.GetWorkspaceConfigurations200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("workspace configuration", err)
	}

	return resp, nil
}

func (t *Team) PostWorkspaceConfigurations(ctx context.Context, request gen.PostWorkspaceConfigurationsRequestObject) (gen.PostWorkspaceConfigurationsResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	config, err := structToMap(request.Body.Config)
	if err != nil {
		return nil, t.resourceSerializationError("workspace configuration", fmt.Errorf("invalid config: %w", err))
	}

	input := workspaceCreateInput{
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      config,
	}

	created, err := t.workspaceService.CreateWorkspaceConfiguration(ctx, input)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildWorkspacePayload(created)
	if err := t.workspaceValidator.ValidateWorkspace(payload); err != nil {
		return nil, t.resourceValidationError("workspace configuration", err)
	}

	var resp gen.PostWorkspaceConfigurations201JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("workspace configuration", err)
	}

	return resp, nil
}

func (t *Team) DeleteWorkspaceConfigurationsId(ctx context.Context, request gen.DeleteWorkspaceConfigurationsIdRequestObject) (gen.DeleteWorkspaceConfigurationsIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("workspace configuration", fmt.Errorf("invalid workspace id: %w", err))
	}
	if err := t.workspaceService.DeleteWorkspaceConfiguration(ctx, id); err != nil {
		return nil, t.handleResourceError(err)
	}
	return gen.DeleteWorkspaceConfigurationsId204Response{}, nil
}

func (t *Team) GetWorkspaceConfigurationsId(ctx context.Context, request gen.GetWorkspaceConfigurationsIdRequestObject) (gen.GetWorkspaceConfigurationsIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("workspace configuration", fmt.Errorf("invalid workspace id: %w", err))
	}

	workspace, err := t.workspaceService.GetWorkspaceConfiguration(ctx, id)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildWorkspacePayload(workspace)
	if err := t.workspaceValidator.ValidateWorkspace(payload); err != nil {
		return nil, t.resourceValidationError("workspace configuration", err)
	}

	var resp gen.GetWorkspaceConfigurationsId200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("workspace configuration", err)
	}

	return resp, nil
}

func (t *Team) PatchWorkspaceConfigurationsId(ctx context.Context, request gen.PatchWorkspaceConfigurationsIdRequestObject) (gen.PatchWorkspaceConfigurationsIdResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return nil, t.resourceSerializationError("workspace configuration", fmt.Errorf("invalid workspace id: %w", err))
	}

	var configPtr *map[string]any
	if request.Body.Config != nil {
		converted, err := structToMap(*request.Body.Config)
		if err != nil {
			return nil, t.resourceSerializationError("workspace configuration", fmt.Errorf("invalid config: %w", err))
		}
		configPtr = &converted
	}

	input := workspaceUpdateInput{
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Config:      configPtr,
	}

	updated, err := t.workspaceService.UpdateWorkspaceConfiguration(ctx, id, input)
	if err != nil {
		return nil, t.handleResourceError(err)
	}

	payload := buildWorkspacePayload(updated)
	if err := t.workspaceValidator.ValidateWorkspace(payload); err != nil {
		return nil, t.resourceValidationError("workspace configuration", err)
	}

	var resp gen.PatchWorkspaceConfigurationsId200JSONResponse
	if err := decodePayload(payload, &resp); err != nil {
		return nil, t.resourceSerializationError("workspace configuration", err)
	}

	return resp, nil
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
	problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), detail)
	return NewProblemError(problem, err)
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
		title := http.StatusText(status)
		if title == "" {
			title = http.StatusText(http.StatusBadGateway)
		}
		return NewProblem(status, title, detail)
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
		Status: status,
		Title:  title,
		Type:   ptr(typ),
	}

	if upstream.Status > 0 {
		problem.Status = upstream.Status
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

	return problem
}

func titleCase(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "Resource"
	}
	runes := []rune(trimmed)
	first := strings.ToUpper(string(runes[0]))
	if len(runes) == 1 {
		return first
	}
	return first + string(runes[1:])
}
