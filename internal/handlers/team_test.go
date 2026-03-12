package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	teamsv1 "github.com/agynio/gateway/gen/agynio/api/teams/v1"
	"github.com/agynio/gateway/internal/gen"
)

type teamsCall struct {
	method string
	assert func(any)
	resp   any
	err    error
}

type stubTeamsClient struct {
	t     *testing.T
	calls []teamsCall
	idx   int
}

func (s *stubTeamsClient) Expect(call teamsCall) {
	s.calls = append(s.calls, call)
}

func (s *stubTeamsClient) AssertDone() {
	if s.idx != len(s.calls) {
		s.t.Fatalf("expected %d calls, got %d", len(s.calls), s.idx)
	}
}

func (s *stubTeamsClient) nextCall(method string, req any) teamsCall {
	if s.idx >= len(s.calls) {
		s.t.Fatalf("unexpected call %s", method)
	}
	call := s.calls[s.idx]
	s.idx++
	if call.method != method {
		s.t.Fatalf("unexpected method: got %s want %s", method, call.method)
	}
	if call.assert != nil {
		call.assert(req)
	}
	return call
}

func (s *stubTeamsClient) ListAgents(ctx context.Context, req *teamsv1.ListAgentsRequest, _ ...grpc.CallOption) (*teamsv1.ListAgentsResponse, error) {
	call := s.nextCall("ListAgents", req)
	if call.resp == nil {
		return nil, call.err
	}
	resp, ok := call.resp.(*teamsv1.ListAgentsResponse)
	if !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateAgent(ctx context.Context, req *teamsv1.CreateAgentRequest, _ ...grpc.CallOption) (*teamsv1.CreateAgentResponse, error) {
	call := s.nextCall("CreateAgent", req)
	if call.resp == nil {
		return nil, call.err
	}
	resp, ok := call.resp.(*teamsv1.CreateAgentResponse)
	if !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteAgent(ctx context.Context, req *teamsv1.DeleteAgentRequest, _ ...grpc.CallOption) (*teamsv1.DeleteAgentResponse, error) {
	call := s.nextCall("DeleteAgent", req)
	resp, ok := call.resp.(*teamsv1.DeleteAgentResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetAgent(ctx context.Context, req *teamsv1.GetAgentRequest, _ ...grpc.CallOption) (*teamsv1.GetAgentResponse, error) {
	call := s.nextCall("GetAgent", req)
	resp, ok := call.resp.(*teamsv1.GetAgentResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateAgent(ctx context.Context, req *teamsv1.UpdateAgentRequest, _ ...grpc.CallOption) (*teamsv1.UpdateAgentResponse, error) {
	call := s.nextCall("UpdateAgent", req)
	resp, ok := call.resp.(*teamsv1.UpdateAgentResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListTools(ctx context.Context, req *teamsv1.ListToolsRequest, _ ...grpc.CallOption) (*teamsv1.ListToolsResponse, error) {
	call := s.nextCall("ListTools", req)
	resp, ok := call.resp.(*teamsv1.ListToolsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateTool(ctx context.Context, req *teamsv1.CreateToolRequest, _ ...grpc.CallOption) (*teamsv1.CreateToolResponse, error) {
	call := s.nextCall("CreateTool", req)
	resp, ok := call.resp.(*teamsv1.CreateToolResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteTool(ctx context.Context, req *teamsv1.DeleteToolRequest, _ ...grpc.CallOption) (*teamsv1.DeleteToolResponse, error) {
	call := s.nextCall("DeleteTool", req)
	resp, ok := call.resp.(*teamsv1.DeleteToolResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetTool(ctx context.Context, req *teamsv1.GetToolRequest, _ ...grpc.CallOption) (*teamsv1.GetToolResponse, error) {
	call := s.nextCall("GetTool", req)
	resp, ok := call.resp.(*teamsv1.GetToolResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateTool(ctx context.Context, req *teamsv1.UpdateToolRequest, _ ...grpc.CallOption) (*teamsv1.UpdateToolResponse, error) {
	call := s.nextCall("UpdateTool", req)
	resp, ok := call.resp.(*teamsv1.UpdateToolResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListMcpServers(ctx context.Context, req *teamsv1.ListMcpServersRequest, _ ...grpc.CallOption) (*teamsv1.ListMcpServersResponse, error) {
	call := s.nextCall("ListMcpServers", req)
	resp, ok := call.resp.(*teamsv1.ListMcpServersResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateMcpServer(ctx context.Context, req *teamsv1.CreateMcpServerRequest, _ ...grpc.CallOption) (*teamsv1.CreateMcpServerResponse, error) {
	call := s.nextCall("CreateMcpServer", req)
	resp, ok := call.resp.(*teamsv1.CreateMcpServerResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteMcpServer(ctx context.Context, req *teamsv1.DeleteMcpServerRequest, _ ...grpc.CallOption) (*teamsv1.DeleteMcpServerResponse, error) {
	call := s.nextCall("DeleteMcpServer", req)
	resp, ok := call.resp.(*teamsv1.DeleteMcpServerResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetMcpServer(ctx context.Context, req *teamsv1.GetMcpServerRequest, _ ...grpc.CallOption) (*teamsv1.GetMcpServerResponse, error) {
	call := s.nextCall("GetMcpServer", req)
	resp, ok := call.resp.(*teamsv1.GetMcpServerResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateMcpServer(ctx context.Context, req *teamsv1.UpdateMcpServerRequest, _ ...grpc.CallOption) (*teamsv1.UpdateMcpServerResponse, error) {
	call := s.nextCall("UpdateMcpServer", req)
	resp, ok := call.resp.(*teamsv1.UpdateMcpServerResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListWorkspaceConfigurations(ctx context.Context, req *teamsv1.ListWorkspaceConfigurationsRequest, _ ...grpc.CallOption) (*teamsv1.ListWorkspaceConfigurationsResponse, error) {
	call := s.nextCall("ListWorkspaceConfigurations", req)
	resp, ok := call.resp.(*teamsv1.ListWorkspaceConfigurationsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateWorkspaceConfiguration(ctx context.Context, req *teamsv1.CreateWorkspaceConfigurationRequest, _ ...grpc.CallOption) (*teamsv1.CreateWorkspaceConfigurationResponse, error) {
	call := s.nextCall("CreateWorkspaceConfiguration", req)
	resp, ok := call.resp.(*teamsv1.CreateWorkspaceConfigurationResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteWorkspaceConfiguration(ctx context.Context, req *teamsv1.DeleteWorkspaceConfigurationRequest, _ ...grpc.CallOption) (*teamsv1.DeleteWorkspaceConfigurationResponse, error) {
	call := s.nextCall("DeleteWorkspaceConfiguration", req)
	resp, ok := call.resp.(*teamsv1.DeleteWorkspaceConfigurationResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetWorkspaceConfiguration(ctx context.Context, req *teamsv1.GetWorkspaceConfigurationRequest, _ ...grpc.CallOption) (*teamsv1.GetWorkspaceConfigurationResponse, error) {
	call := s.nextCall("GetWorkspaceConfiguration", req)
	resp, ok := call.resp.(*teamsv1.GetWorkspaceConfigurationResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateWorkspaceConfiguration(ctx context.Context, req *teamsv1.UpdateWorkspaceConfigurationRequest, _ ...grpc.CallOption) (*teamsv1.UpdateWorkspaceConfigurationResponse, error) {
	call := s.nextCall("UpdateWorkspaceConfiguration", req)
	resp, ok := call.resp.(*teamsv1.UpdateWorkspaceConfigurationResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListMemoryBuckets(ctx context.Context, req *teamsv1.ListMemoryBucketsRequest, _ ...grpc.CallOption) (*teamsv1.ListMemoryBucketsResponse, error) {
	call := s.nextCall("ListMemoryBuckets", req)
	resp, ok := call.resp.(*teamsv1.ListMemoryBucketsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateMemoryBucket(ctx context.Context, req *teamsv1.CreateMemoryBucketRequest, _ ...grpc.CallOption) (*teamsv1.CreateMemoryBucketResponse, error) {
	call := s.nextCall("CreateMemoryBucket", req)
	resp, ok := call.resp.(*teamsv1.CreateMemoryBucketResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteMemoryBucket(ctx context.Context, req *teamsv1.DeleteMemoryBucketRequest, _ ...grpc.CallOption) (*teamsv1.DeleteMemoryBucketResponse, error) {
	call := s.nextCall("DeleteMemoryBucket", req)
	resp, ok := call.resp.(*teamsv1.DeleteMemoryBucketResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetMemoryBucket(ctx context.Context, req *teamsv1.GetMemoryBucketRequest, _ ...grpc.CallOption) (*teamsv1.GetMemoryBucketResponse, error) {
	call := s.nextCall("GetMemoryBucket", req)
	resp, ok := call.resp.(*teamsv1.GetMemoryBucketResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateMemoryBucket(ctx context.Context, req *teamsv1.UpdateMemoryBucketRequest, _ ...grpc.CallOption) (*teamsv1.UpdateMemoryBucketResponse, error) {
	call := s.nextCall("UpdateMemoryBucket", req)
	resp, ok := call.resp.(*teamsv1.UpdateMemoryBucketResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListAttachments(ctx context.Context, req *teamsv1.ListAttachmentsRequest, _ ...grpc.CallOption) (*teamsv1.ListAttachmentsResponse, error) {
	call := s.nextCall("ListAttachments", req)
	resp, ok := call.resp.(*teamsv1.ListAttachmentsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateAttachment(ctx context.Context, req *teamsv1.CreateAttachmentRequest, _ ...grpc.CallOption) (*teamsv1.CreateAttachmentResponse, error) {
	call := s.nextCall("CreateAttachment", req)
	resp, ok := call.resp.(*teamsv1.CreateAttachmentResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetAttachment(ctx context.Context, req *teamsv1.GetAttachmentRequest, _ ...grpc.CallOption) (*teamsv1.GetAttachmentResponse, error) {
	call := s.nextCall("GetAttachment", req)
	resp, ok := call.resp.(*teamsv1.GetAttachmentResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteAttachment(ctx context.Context, req *teamsv1.DeleteAttachmentRequest, _ ...grpc.CallOption) (*teamsv1.DeleteAttachmentResponse, error) {
	call := s.nextCall("DeleteAttachment", req)
	resp, ok := call.resp.(*teamsv1.DeleteAttachmentResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

var _ teamsv1.TeamsServiceClient = (*stubTeamsClient)(nil)

func TestTeamGetAgentsPagination(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	agentID1 := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	agentID2 := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	agentOne := &teamsv1.Agent{
		Meta:   testMeta(agentID1.String()),
		Title:  "Agent One",
		Config: &teamsv1.AgentConfig{Model: "gpt"},
	}
	agentTwo := &teamsv1.Agent{
		Meta:   testMeta(agentID2.String()),
		Title:  "Agent Two",
		Config: &teamsv1.AgentConfig{Model: "gpt-4"},
	}

	stub.Expect(teamsCall{
		method: "ListAgents",
		assert: func(req any) {
			input := req.(*teamsv1.ListAgentsRequest)
			if input.PageSize != 1 {
				t.Fatalf("unexpected page size: %d", input.PageSize)
			}
			if input.PageToken != "" {
				t.Fatalf("unexpected page token: %s", input.PageToken)
			}
			if input.Query != "search" {
				t.Fatalf("unexpected query: %s", input.Query)
			}
		},
		resp: &teamsv1.ListAgentsResponse{Agents: []*teamsv1.Agent{agentOne}, NextPageToken: "next"},
	})
	stub.Expect(teamsCall{
		method: "ListAgents",
		assert: func(req any) {
			input := req.(*teamsv1.ListAgentsRequest)
			if input.PageToken != "next" {
				t.Fatalf("unexpected page token: %s", input.PageToken)
			}
		},
		resp: &teamsv1.ListAgentsResponse{Agents: []*teamsv1.Agent{agentTwo}},
	})

	page := 2
	perPage := 1
	query := " search "
	h := NewTeam(stub)
	resp, err := h.GetAgents(context.Background(), gen.GetAgentsRequestObject{Params: gen.GetAgentsParams{
		Q:       &query,
		Page:    &page,
		PerPage: &perPage,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := resp.(gen.GetAgents200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if list.Total != 2 {
		t.Fatalf("unexpected total: %d", list.Total)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(list.Items))
	}
	if list.Items[0].Title == nil || *list.Items[0].Title != "Agent Two" {
		t.Fatalf("unexpected agent title: %v", list.Items[0].Title)
	}
	if list.Items[0].Config.Model == nil || *list.Items[0].Config.Model != "gpt-4" {
		t.Fatalf("unexpected agent config: %+v", list.Items[0].Config)
	}

	stub.AssertDone()
}

func TestTeamPostAgents(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	agentID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	whenBusy := gen.InjectAfterTools
	processBuffer := gen.AllTogether
	debounce := 42
	restrict := true
	role := "planner"

	stub.Expect(teamsCall{
		method: "CreateAgent",
		assert: func(req any) {
			input := req.(*teamsv1.CreateAgentRequest)
			if input.Title != "Support" {
				t.Fatalf("unexpected title: %s", input.Title)
			}
			if input.Description != "Assist" {
				t.Fatalf("unexpected description: %s", input.Description)
			}
			if input.Config == nil {
				t.Fatalf("missing config")
			}
			if input.Config.Model != "gpt-4" {
				t.Fatalf("unexpected model: %s", input.Config.Model)
			}
			if input.Config.WhenBusy != teamsv1.AgentWhenBusy_AGENT_WHEN_BUSY_INJECT_AFTER_TOOLS {
				t.Fatalf("unexpected when_busy: %s", input.Config.WhenBusy.String())
			}
			if input.Config.ProcessBuffer != teamsv1.AgentProcessBuffer_AGENT_PROCESS_BUFFER_ALL_TOGETHER {
				t.Fatalf("unexpected process buffer: %s", input.Config.ProcessBuffer.String())
			}
			if input.Config.DebounceMs != 42 {
				t.Fatalf("unexpected debounce: %d", input.Config.DebounceMs)
			}
			if !input.Config.RestrictOutput {
				t.Fatalf("expected restrict output true")
			}
			if input.Config.Role != "planner" {
				t.Fatalf("unexpected role: %s", input.Config.Role)
			}
		},
		resp: &teamsv1.CreateAgentResponse{Agent: &teamsv1.Agent{
			Meta:        testMeta(agentID.String()),
			Title:       "Support",
			Description: "Assist",
			Config: &teamsv1.AgentConfig{
				Model:          "gpt-4",
				WhenBusy:       teamsv1.AgentWhenBusy_AGENT_WHEN_BUSY_INJECT_AFTER_TOOLS,
				ProcessBuffer:  teamsv1.AgentProcessBuffer_AGENT_PROCESS_BUFFER_ALL_TOGETHER,
				DebounceMs:     42,
				RestrictOutput: true,
				Role:           "planner",
			},
		}},
	})

	request := gen.PostAgentsRequestObject{Body: &gen.AgentCreateRequest{
		Title:       strPtr("Support"),
		Description: strPtr("Assist"),
		Config: gen.AgentConfig{
			Model:          strPtr("gpt-4"),
			WhenBusy:       &whenBusy,
			ProcessBuffer:  &processBuffer,
			DebounceMs:     &debounce,
			RestrictOutput: &restrict,
			Role:           &role,
		},
	}}

	h := NewTeam(stub)
	resp, err := h.PostAgents(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostAgents201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Title == nil || *created.Title != "Support" {
		t.Fatalf("unexpected title: %v", created.Title)
	}
	if created.Config.WhenBusy == nil || *created.Config.WhenBusy != gen.InjectAfterTools {
		t.Fatalf("unexpected whenBusy: %v", created.Config.WhenBusy)
	}
	if created.Config.ProcessBuffer == nil || *created.Config.ProcessBuffer != gen.AllTogether {
		t.Fatalf("unexpected process buffer: %v", created.Config.ProcessBuffer)
	}
	if created.Config.DebounceMs == nil || *created.Config.DebounceMs != 42 {
		t.Fatalf("unexpected debounce: %v", created.Config.DebounceMs)
	}

	stub.AssertDone()
}

func TestTeamGetTools(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	toolID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	config, err := structpb.NewStruct(map[string]any{"enabled": true})
	if err != nil {
		t.Fatalf("structpb: %v", err)
	}

	stub.Expect(teamsCall{
		method: "ListTools",
		assert: func(req any) {
			input := req.(*teamsv1.ListToolsRequest)
			if input.Type != teamsv1.ToolType_TOOL_TYPE_MANAGE {
				t.Fatalf("unexpected tool type: %s", input.Type.String())
			}
		},
		resp: &teamsv1.ListToolsResponse{Tools: []*teamsv1.Tool{{
			Meta:        testMeta(toolID.String()),
			Type:        teamsv1.ToolType_TOOL_TYPE_MANAGE,
			Name:        "Runner",
			Description: "Runs tasks",
			Config:      config,
		}}},
	})

	toolType := gen.Manage
	h := NewTeam(stub)
	resp, err := h.GetTools(context.Background(), gen.GetToolsRequestObject{Params: gen.GetToolsParams{Type: &toolType}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := resp.(gen.GetTools200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(list.Items))
	}
	if list.Items[0].Name == nil || *list.Items[0].Name != "Runner" {
		t.Fatalf("unexpected name: %v", list.Items[0].Name)
	}
	if list.Items[0].Config == nil || (*list.Items[0].Config)["enabled"] != true {
		t.Fatalf("unexpected config: %v", list.Items[0].Config)
	}

	stub.AssertDone()
}

func TestTeamPostTools(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	toolID := uuid.MustParse("55555555-5555-5555-5555-555555555555")

	stub.Expect(teamsCall{
		method: "CreateTool",
		assert: func(req any) {
			input := req.(*teamsv1.CreateToolRequest)
			if input.Type != teamsv1.ToolType_TOOL_TYPE_MEMORY {
				t.Fatalf("unexpected type: %s", input.Type.String())
			}
			if input.Name != "Memory" {
				t.Fatalf("unexpected name: %s", input.Name)
			}
			if input.Config == nil {
				t.Fatalf("missing config")
			}
			if input.Config.AsMap()["ttl"] != float64(10) {
				t.Fatalf("unexpected config map: %v", input.Config.AsMap())
			}
		},
		resp: &teamsv1.CreateToolResponse{Tool: &teamsv1.Tool{
			Meta:        testMeta(toolID.String()),
			Type:        teamsv1.ToolType_TOOL_TYPE_MEMORY,
			Name:        "Memory",
			Description: "Store",
			Config:      mustStruct(t, map[string]any{"ttl": 10}),
		}},
	})

	config := map[string]any{"ttl": 10}
	request := gen.PostToolsRequestObject{Body: &gen.ToolCreateRequest{
		Type:        gen.Memory,
		Name:        strPtr("Memory"),
		Description: strPtr("Store"),
		Config:      &config,
	}}

	h := NewTeam(stub)
	resp, err := h.PostTools(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostTools201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Config == nil || (*created.Config)["ttl"] != float64(10) {
		t.Fatalf("unexpected config: %v", created.Config)
	}

	stub.AssertDone()
}

func TestTeamPatchTools(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	toolID := uuid.MustParse("66666666-6666-6666-6666-666666666666")
	stub.Expect(teamsCall{
		method: "UpdateTool",
		assert: func(req any) {
			input := req.(*teamsv1.UpdateToolRequest)
			if input.Id != toolID.String() {
				t.Fatalf("unexpected id: %s", input.Id)
			}
			if input.Name == nil || *input.Name != "Updated" {
				t.Fatalf("unexpected name: %v", input.Name)
			}
			if input.Config == nil || input.Config.AsMap()["enabled"] != true {
				t.Fatalf("unexpected config: %v", input.Config)
			}
		},
		resp: &teamsv1.UpdateToolResponse{Tool: &teamsv1.Tool{
			Meta:   testMeta(toolID.String()),
			Type:   teamsv1.ToolType_TOOL_TYPE_SEND_MESSAGE,
			Name:   "Updated",
			Config: mustStruct(t, map[string]any{"enabled": true}),
		}},
	})

	config := map[string]any{"enabled": true}
	name := "Updated"
	request := gen.PatchToolsIdRequestObject{
		Id: toolID,
		Body: &gen.ToolUpdateRequest{
			Name:   &name,
			Config: &config,
		},
	}

	h := NewTeam(stub)
	resp, err := h.PatchToolsId(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, ok := resp.(gen.PatchToolsId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if updated.Name == nil || *updated.Name != "Updated" {
		t.Fatalf("unexpected name: %v", updated.Name)
	}

	stub.AssertDone()
}

func TestTeamPostAttachmentsMcpServerWorkspaceConfiguration(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	attachmentID := uuid.MustParse("77777777-7777-7777-7777-777777777777")
	sourceID := uuid.MustParse("88888888-8888-8888-8888-888888888888")
	targetID := uuid.MustParse("99999999-9999-9999-9999-999999999999")

	stub.Expect(teamsCall{
		method: "CreateAttachment",
		assert: func(req any) {
			input := req.(*teamsv1.CreateAttachmentRequest)
			if input.Kind != teamsv1.AttachmentKind_ATTACHMENT_KIND_MCP_SERVER_WORKSPACE_CONFIGURATION {
				t.Fatalf("unexpected kind: %s", input.Kind.String())
			}
			if input.SourceId != sourceID.String() || input.TargetId != targetID.String() {
				t.Fatalf("unexpected ids: %s %s", input.SourceId, input.TargetId)
			}
		},
		resp: &teamsv1.CreateAttachmentResponse{Attachment: &teamsv1.Attachment{
			Meta:       testMeta(attachmentID.String()),
			Kind:       teamsv1.AttachmentKind_ATTACHMENT_KIND_MCP_SERVER_WORKSPACE_CONFIGURATION,
			SourceType: teamsv1.EntityType_ENTITY_TYPE_MCP_SERVER,
			SourceId:   sourceID.String(),
			TargetType: teamsv1.EntityType_ENTITY_TYPE_WORKSPACE_CONFIGURATION,
			TargetId:   targetID.String(),
		}},
	})

	h := NewTeam(stub)
	resp, err := h.PostAttachments(context.Background(), gen.PostAttachmentsRequestObject{Body: &gen.AttachmentCreateRequest{
		Kind:     gen.McpServerWorkspaceConfiguration,
		SourceId: openapi_types.UUID(sourceID),
		TargetId: openapi_types.UUID(targetID),
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostAttachments201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Kind != gen.McpServerWorkspaceConfiguration {
		t.Fatalf("unexpected kind: %s", created.Kind)
	}

	stub.AssertDone()
}

func TestTeamGetAttachments(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	attachmentID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	sourceID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	targetID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	stub.Expect(teamsCall{
		method: "ListAttachments",
		assert: func(req any) {
			input := req.(*teamsv1.ListAttachmentsRequest)
			if input.SourceType != teamsv1.EntityType_ENTITY_TYPE_AGENT {
				t.Fatalf("unexpected source type: %s", input.SourceType.String())
			}
			if input.TargetType != teamsv1.EntityType_ENTITY_TYPE_TOOL {
				t.Fatalf("unexpected target type: %s", input.TargetType.String())
			}
			if input.Kind != teamsv1.AttachmentKind_ATTACHMENT_KIND_AGENT_TOOL {
				t.Fatalf("unexpected kind: %s", input.Kind.String())
			}
		},
		resp: &teamsv1.ListAttachmentsResponse{Attachments: []*teamsv1.Attachment{{
			Meta:       testMeta(attachmentID.String()),
			Kind:       teamsv1.AttachmentKind_ATTACHMENT_KIND_AGENT_TOOL,
			SourceType: teamsv1.EntityType_ENTITY_TYPE_AGENT,
			SourceId:   sourceID.String(),
			TargetType: teamsv1.EntityType_ENTITY_TYPE_TOOL,
			TargetId:   targetID.String(),
		}}},
	})

	page := 1
	perPage := 10
	sourceType := gen.EntityTypeAgent
	targetType := gen.EntityTypeTool
	kind := gen.AgentTool
	h := NewTeam(stub)
	resp, err := h.GetAttachments(context.Background(), gen.GetAttachmentsRequestObject{Params: gen.GetAttachmentsParams{
		SourceType: &sourceType,
		SourceId:   (*openapi_types.UUID)(&sourceID),
		TargetType: &targetType,
		TargetId:   (*openapi_types.UUID)(&targetID),
		Kind:       &kind,
		Page:       &page,
		PerPage:    &perPage,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := resp.(gen.GetAttachments200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(list.Items))
	}
	if list.Items[0].Kind != gen.AgentTool {
		t.Fatalf("unexpected kind: %s", list.Items[0].Kind)
	}

	stub.AssertDone()
}

func TestTeamPostMcpServers(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	serverID := uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")
	stub.Expect(teamsCall{
		method: "CreateMcpServer",
		assert: func(req any) {
			input := req.(*teamsv1.CreateMcpServerRequest)
			if input.Title != "Test" {
				t.Fatalf("unexpected title: %s", input.Title)
			}
			if input.Config == nil || input.Config.Command != "bash" {
				t.Fatalf("unexpected config: %+v", input.Config)
			}
			if len(input.Config.Env) != 1 || input.Config.Env[0].Name != "FOO" {
				t.Fatalf("unexpected env: %+v", input.Config.Env)
			}
			if input.Config.Restart == nil || input.Config.Restart.BackoffMs != 100 {
				t.Fatalf("unexpected restart config: %+v", input.Config.Restart)
			}
		},
		resp: &teamsv1.CreateMcpServerResponse{McpServer: &teamsv1.McpServer{
			Meta:        testMeta(serverID.String()),
			Title:       "Test",
			Description: "Runner",
			Config: &teamsv1.McpServerConfig{
				Command: "bash",
				Env:     []*teamsv1.McpEnvItem{{Name: "FOO", Value: "bar"}},
				Restart: &teamsv1.McpServerRestartConfig{BackoffMs: 100, MaxAttempts: 2},
			},
		}},
	})

	restart := &struct {
		BackoffMs   *int `json:"backoffMs,omitempty"`
		MaxAttempts *int `json:"maxAttempts,omitempty"`
	}{BackoffMs: intPtr(100), MaxAttempts: intPtr(2)}
	config := gen.McpServerConfig{
		Command: strPtr("bash"),
		Env:     &[]gen.McpEnvItem{{Name: "FOO", Value: "bar"}},
		Restart: restart,
	}
	request := gen.PostMcpServersRequestObject{Body: &gen.McpServerCreateRequest{
		Title:       strPtr("Test"),
		Description: strPtr("Runner"),
		Config:      config,
	}}

	h := NewTeam(stub)
	resp, err := h.PostMcpServers(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostMcpServers201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Config.Command == nil || *created.Config.Command != "bash" {
		t.Fatalf("unexpected command: %v", created.Config.Command)
	}

	stub.AssertDone()
}

func TestTeamPatchMcpServers(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	serverID := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	name := "Updated"
	stub.Expect(teamsCall{
		method: "UpdateMcpServer",
		assert: func(req any) {
			input := req.(*teamsv1.UpdateMcpServerRequest)
			if input.Id != serverID.String() {
				t.Fatalf("unexpected id: %s", input.Id)
			}
			if input.Title == nil || *input.Title != "Updated" {
				t.Fatalf("unexpected title: %v", input.Title)
			}
			if input.Config != nil {
				t.Fatalf("expected nil config")
			}
		},
		resp: &teamsv1.UpdateMcpServerResponse{McpServer: &teamsv1.McpServer{
			Meta:   testMeta(serverID.String()),
			Title:  "Updated",
			Config: &teamsv1.McpServerConfig{},
		}},
	})

	request := gen.PatchMcpServersIdRequestObject{
		Id: serverID,
		Body: &gen.McpServerUpdateRequest{
			Title: &name,
		},
	}

	h := NewTeam(stub)
	resp, err := h.PatchMcpServersId(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, ok := resp.(gen.PatchMcpServersId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if updated.Title == nil || *updated.Title != "Updated" {
		t.Fatalf("unexpected title: %v", updated.Title)
	}

	stub.AssertDone()
}

func TestTeamPostWorkspaceConfigurations(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	configID := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	cpuLimit := gen.WorkspaceConfig_CpuLimit{}
	if err := cpuLimit.FromWorkspaceConfigCpuLimit1("2"); err != nil {
		t.Fatalf("cpu limit: %v", err)
	}
	memoryLimit := gen.WorkspaceConfig_MemoryLimit{}
	if err := memoryLimit.FromWorkspaceConfigMemoryLimit0(4.5); err != nil {
		t.Fatalf("memory limit: %v", err)
	}
	platform := gen.Linuxamd64

	stub.Expect(teamsCall{
		method: "CreateWorkspaceConfiguration",
		assert: func(req any) {
			input := req.(*teamsv1.CreateWorkspaceConfigurationRequest)
			if input.Config == nil {
				t.Fatalf("missing config")
			}
			if input.Config.CpuLimit != "2" || input.Config.MemoryLimit != "4.5" {
				t.Fatalf("unexpected limits: %s %s", input.Config.CpuLimit, input.Config.MemoryLimit)
			}
			if input.Config.Platform != teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_LINUX_AMD64 {
				t.Fatalf("unexpected platform: %s", input.Config.Platform.String())
			}
		},
		resp: &teamsv1.CreateWorkspaceConfigurationResponse{WorkspaceConfiguration: &teamsv1.WorkspaceConfiguration{
			Meta:        testMeta(configID.String()),
			Title:       "Config",
			Description: "Desc",
			Config: &teamsv1.WorkspaceConfig{
				CpuLimit:    "2",
				MemoryLimit: "4.5",
				Platform:    teamsv1.WorkspacePlatform_WORKSPACE_PLATFORM_LINUX_AMD64,
			},
		}},
	})

	request := gen.PostWorkspaceConfigurationsRequestObject{Body: &gen.WorkspaceConfigurationCreateRequest{
		Title:       strPtr("Config"),
		Description: strPtr("Desc"),
		Config: gen.WorkspaceConfig{
			CpuLimit:    &cpuLimit,
			MemoryLimit: &memoryLimit,
			Platform:    &platform,
		},
	}}

	h := NewTeam(stub)
	resp, err := h.PostWorkspaceConfigurations(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostWorkspaceConfigurations201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Config.Platform == nil || *created.Config.Platform != gen.Linuxamd64 {
		t.Fatalf("unexpected platform: %v", created.Config.Platform)
	}

	stub.AssertDone()
}

func TestTeamPatchWorkspaceConfigurations(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	configID := uuid.MustParse("abababab-abab-abab-abab-abababababab")
	name := "Updated"
	stub.Expect(teamsCall{
		method: "UpdateWorkspaceConfiguration",
		assert: func(req any) {
			input := req.(*teamsv1.UpdateWorkspaceConfigurationRequest)
			if input.Id != configID.String() {
				t.Fatalf("unexpected id: %s", input.Id)
			}
			if input.Title == nil || *input.Title != "Updated" {
				t.Fatalf("unexpected title: %v", input.Title)
			}
			if input.Config != nil {
				t.Fatalf("expected nil config")
			}
		},
		resp: &teamsv1.UpdateWorkspaceConfigurationResponse{WorkspaceConfiguration: &teamsv1.WorkspaceConfiguration{
			Meta:   testMeta(configID.String()),
			Title:  "Updated",
			Config: &teamsv1.WorkspaceConfig{},
		}},
	})

	request := gen.PatchWorkspaceConfigurationsIdRequestObject{
		Id: configID,
		Body: &gen.WorkspaceConfigurationUpdateRequest{
			Title: &name,
		},
	}

	h := NewTeam(stub)
	resp, err := h.PatchWorkspaceConfigurationsId(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, ok := resp.(gen.PatchWorkspaceConfigurationsId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if updated.Title == nil || *updated.Title != "Updated" {
		t.Fatalf("unexpected title: %v", updated.Title)
	}

	stub.AssertDone()
}

func TestTeamPostMemoryBuckets(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	bucketID := uuid.MustParse("cdcdcdcd-cdcd-cdcd-cdcd-cdcdcdcdcdcd")
	scope := gen.Global

	stub.Expect(teamsCall{
		method: "CreateMemoryBucket",
		assert: func(req any) {
			input := req.(*teamsv1.CreateMemoryBucketRequest)
			if input.Config == nil || input.Config.Scope != teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_GLOBAL {
				t.Fatalf("unexpected scope: %+v", input.Config)
			}
		},
		resp: &teamsv1.CreateMemoryBucketResponse{MemoryBucket: &teamsv1.MemoryBucket{
			Meta:        testMeta(bucketID.String()),
			Title:       "Bucket",
			Description: "Desc",
			Config:      &teamsv1.MemoryBucketConfig{Scope: teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_GLOBAL},
		}},
	})

	request := gen.PostMemoryBucketsRequestObject{Body: &gen.MemoryBucketCreateRequest{
		Title:       strPtr("Bucket"),
		Description: strPtr("Desc"),
		Config: gen.MemoryBucketConfig{
			Scope: &scope,
		},
	}}

	h := NewTeam(stub)
	resp, err := h.PostMemoryBuckets(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostMemoryBuckets201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Config.Scope == nil || *created.Config.Scope != gen.Global {
		t.Fatalf("unexpected scope: %v", created.Config.Scope)
	}

	stub.AssertDone()
}

func TestTeamPatchMemoryBuckets(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	bucketID := uuid.MustParse("efefefef-efef-efef-efef-efefefefefef")
	scope := gen.PerThread
	stub.Expect(teamsCall{
		method: "UpdateMemoryBucket",
		assert: func(req any) {
			input := req.(*teamsv1.UpdateMemoryBucketRequest)
			if input.Config == nil || input.Config.Scope != teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_PER_THREAD {
				t.Fatalf("unexpected scope: %+v", input.Config)
			}
		},
		resp: &teamsv1.UpdateMemoryBucketResponse{MemoryBucket: &teamsv1.MemoryBucket{
			Meta:   testMeta(bucketID.String()),
			Title:  "Updated",
			Config: &teamsv1.MemoryBucketConfig{Scope: teamsv1.MemoryBucketScope_MEMORY_BUCKET_SCOPE_PER_THREAD},
		}},
	})

	request := gen.PatchMemoryBucketsIdRequestObject{
		Id: bucketID,
		Body: &gen.MemoryBucketUpdateRequest{
			Config: &gen.MemoryBucketConfig{Scope: &scope},
		},
	}

	h := NewTeam(stub)
	resp, err := h.PatchMemoryBucketsId(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, ok := resp.(gen.PatchMemoryBucketsId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if updated.Config.Scope == nil || *updated.Config.Scope != gen.PerThread {
		t.Fatalf("unexpected scope: %v", updated.Config.Scope)
	}

	stub.AssertDone()
}

func TestTeamGrpcErrorMapping(t *testing.T) {
	stub := &stubTeamsClient{t: t}
	stub.Expect(teamsCall{
		method: "ListTools",
		err:    status.Error(codes.NotFound, "missing"),
	})

	h := NewTeam(stub)
	_, err := h.GetTools(context.Background(), gen.GetToolsRequestObject{})
	if err == nil {
		t.Fatalf("expected error")
	}

	var problemErr *ProblemError
	if !errors.As(err, &problemErr) {
		t.Fatalf("expected ProblemError, got %T", err)
	}
	if problemErr.Problem.Status != 404 {
		t.Fatalf("unexpected status: %d", problemErr.Problem.Status)
	}

	stub.AssertDone()
}

func testMeta(id string) *teamsv1.EntityMeta {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	return &teamsv1.EntityMeta{
		Id:        id,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}
}

func mustStruct(t *testing.T, value map[string]any) *structpb.Struct {
	t.Helper()
	result, err := structpb.NewStruct(value)
	if err != nil {
		t.Fatalf("structpb: %v", err)
	}
	return result
}

func strPtr(value string) *string {
	return &value
}
