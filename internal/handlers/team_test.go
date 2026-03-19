package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"google.golang.org/grpc"
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

func (s *stubTeamsClient) CreateAgent(ctx context.Context, req *teamsv1.CreateAgentRequest, _ ...grpc.CallOption) (*teamsv1.CreateAgentResponse, error) {
	call := s.nextCall("CreateAgent", req)
	resp, ok := call.resp.(*teamsv1.CreateAgentResponse)
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

func (s *stubTeamsClient) DeleteAgent(ctx context.Context, req *teamsv1.DeleteAgentRequest, _ ...grpc.CallOption) (*teamsv1.DeleteAgentResponse, error) {
	call := s.nextCall("DeleteAgent", req)
	resp, ok := call.resp.(*teamsv1.DeleteAgentResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListAgents(ctx context.Context, req *teamsv1.ListAgentsRequest, _ ...grpc.CallOption) (*teamsv1.ListAgentsResponse, error) {
	call := s.nextCall("ListAgents", req)
	resp, ok := call.resp.(*teamsv1.ListAgentsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateVolume(ctx context.Context, req *teamsv1.CreateVolumeRequest, _ ...grpc.CallOption) (*teamsv1.CreateVolumeResponse, error) {
	call := s.nextCall("CreateVolume", req)
	resp, ok := call.resp.(*teamsv1.CreateVolumeResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetVolume(ctx context.Context, req *teamsv1.GetVolumeRequest, _ ...grpc.CallOption) (*teamsv1.GetVolumeResponse, error) {
	call := s.nextCall("GetVolume", req)
	resp, ok := call.resp.(*teamsv1.GetVolumeResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateVolume(ctx context.Context, req *teamsv1.UpdateVolumeRequest, _ ...grpc.CallOption) (*teamsv1.UpdateVolumeResponse, error) {
	call := s.nextCall("UpdateVolume", req)
	resp, ok := call.resp.(*teamsv1.UpdateVolumeResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteVolume(ctx context.Context, req *teamsv1.DeleteVolumeRequest, _ ...grpc.CallOption) (*teamsv1.DeleteVolumeResponse, error) {
	call := s.nextCall("DeleteVolume", req)
	resp, ok := call.resp.(*teamsv1.DeleteVolumeResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListVolumes(ctx context.Context, req *teamsv1.ListVolumesRequest, _ ...grpc.CallOption) (*teamsv1.ListVolumesResponse, error) {
	call := s.nextCall("ListVolumes", req)
	resp, ok := call.resp.(*teamsv1.ListVolumesResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateVolumeAttachment(ctx context.Context, req *teamsv1.CreateVolumeAttachmentRequest, _ ...grpc.CallOption) (*teamsv1.CreateVolumeAttachmentResponse, error) {
	call := s.nextCall("CreateVolumeAttachment", req)
	resp, ok := call.resp.(*teamsv1.CreateVolumeAttachmentResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetVolumeAttachment(ctx context.Context, req *teamsv1.GetVolumeAttachmentRequest, _ ...grpc.CallOption) (*teamsv1.GetVolumeAttachmentResponse, error) {
	call := s.nextCall("GetVolumeAttachment", req)
	resp, ok := call.resp.(*teamsv1.GetVolumeAttachmentResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteVolumeAttachment(ctx context.Context, req *teamsv1.DeleteVolumeAttachmentRequest, _ ...grpc.CallOption) (*teamsv1.DeleteVolumeAttachmentResponse, error) {
	call := s.nextCall("DeleteVolumeAttachment", req)
	resp, ok := call.resp.(*teamsv1.DeleteVolumeAttachmentResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListVolumeAttachments(ctx context.Context, req *teamsv1.ListVolumeAttachmentsRequest, _ ...grpc.CallOption) (*teamsv1.ListVolumeAttachmentsResponse, error) {
	call := s.nextCall("ListVolumeAttachments", req)
	resp, ok := call.resp.(*teamsv1.ListVolumeAttachmentsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateMcp(ctx context.Context, req *teamsv1.CreateMcpRequest, _ ...grpc.CallOption) (*teamsv1.CreateMcpResponse, error) {
	call := s.nextCall("CreateMcp", req)
	resp, ok := call.resp.(*teamsv1.CreateMcpResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetMcp(ctx context.Context, req *teamsv1.GetMcpRequest, _ ...grpc.CallOption) (*teamsv1.GetMcpResponse, error) {
	call := s.nextCall("GetMcp", req)
	resp, ok := call.resp.(*teamsv1.GetMcpResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateMcp(ctx context.Context, req *teamsv1.UpdateMcpRequest, _ ...grpc.CallOption) (*teamsv1.UpdateMcpResponse, error) {
	call := s.nextCall("UpdateMcp", req)
	resp, ok := call.resp.(*teamsv1.UpdateMcpResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteMcp(ctx context.Context, req *teamsv1.DeleteMcpRequest, _ ...grpc.CallOption) (*teamsv1.DeleteMcpResponse, error) {
	call := s.nextCall("DeleteMcp", req)
	resp, ok := call.resp.(*teamsv1.DeleteMcpResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListMcps(ctx context.Context, req *teamsv1.ListMcpsRequest, _ ...grpc.CallOption) (*teamsv1.ListMcpsResponse, error) {
	call := s.nextCall("ListMcps", req)
	resp, ok := call.resp.(*teamsv1.ListMcpsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateSkill(ctx context.Context, req *teamsv1.CreateSkillRequest, _ ...grpc.CallOption) (*teamsv1.CreateSkillResponse, error) {
	call := s.nextCall("CreateSkill", req)
	resp, ok := call.resp.(*teamsv1.CreateSkillResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetSkill(ctx context.Context, req *teamsv1.GetSkillRequest, _ ...grpc.CallOption) (*teamsv1.GetSkillResponse, error) {
	call := s.nextCall("GetSkill", req)
	resp, ok := call.resp.(*teamsv1.GetSkillResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateSkill(ctx context.Context, req *teamsv1.UpdateSkillRequest, _ ...grpc.CallOption) (*teamsv1.UpdateSkillResponse, error) {
	call := s.nextCall("UpdateSkill", req)
	resp, ok := call.resp.(*teamsv1.UpdateSkillResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteSkill(ctx context.Context, req *teamsv1.DeleteSkillRequest, _ ...grpc.CallOption) (*teamsv1.DeleteSkillResponse, error) {
	call := s.nextCall("DeleteSkill", req)
	resp, ok := call.resp.(*teamsv1.DeleteSkillResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListSkills(ctx context.Context, req *teamsv1.ListSkillsRequest, _ ...grpc.CallOption) (*teamsv1.ListSkillsResponse, error) {
	call := s.nextCall("ListSkills", req)
	resp, ok := call.resp.(*teamsv1.ListSkillsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateHook(ctx context.Context, req *teamsv1.CreateHookRequest, _ ...grpc.CallOption) (*teamsv1.CreateHookResponse, error) {
	call := s.nextCall("CreateHook", req)
	resp, ok := call.resp.(*teamsv1.CreateHookResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetHook(ctx context.Context, req *teamsv1.GetHookRequest, _ ...grpc.CallOption) (*teamsv1.GetHookResponse, error) {
	call := s.nextCall("GetHook", req)
	resp, ok := call.resp.(*teamsv1.GetHookResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateHook(ctx context.Context, req *teamsv1.UpdateHookRequest, _ ...grpc.CallOption) (*teamsv1.UpdateHookResponse, error) {
	call := s.nextCall("UpdateHook", req)
	resp, ok := call.resp.(*teamsv1.UpdateHookResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteHook(ctx context.Context, req *teamsv1.DeleteHookRequest, _ ...grpc.CallOption) (*teamsv1.DeleteHookResponse, error) {
	call := s.nextCall("DeleteHook", req)
	resp, ok := call.resp.(*teamsv1.DeleteHookResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListHooks(ctx context.Context, req *teamsv1.ListHooksRequest, _ ...grpc.CallOption) (*teamsv1.ListHooksResponse, error) {
	call := s.nextCall("ListHooks", req)
	resp, ok := call.resp.(*teamsv1.ListHooksResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateEnv(ctx context.Context, req *teamsv1.CreateEnvRequest, _ ...grpc.CallOption) (*teamsv1.CreateEnvResponse, error) {
	call := s.nextCall("CreateEnv", req)
	resp, ok := call.resp.(*teamsv1.CreateEnvResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetEnv(ctx context.Context, req *teamsv1.GetEnvRequest, _ ...grpc.CallOption) (*teamsv1.GetEnvResponse, error) {
	call := s.nextCall("GetEnv", req)
	resp, ok := call.resp.(*teamsv1.GetEnvResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateEnv(ctx context.Context, req *teamsv1.UpdateEnvRequest, _ ...grpc.CallOption) (*teamsv1.UpdateEnvResponse, error) {
	call := s.nextCall("UpdateEnv", req)
	resp, ok := call.resp.(*teamsv1.UpdateEnvResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteEnv(ctx context.Context, req *teamsv1.DeleteEnvRequest, _ ...grpc.CallOption) (*teamsv1.DeleteEnvResponse, error) {
	call := s.nextCall("DeleteEnv", req)
	resp, ok := call.resp.(*teamsv1.DeleteEnvResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListEnvs(ctx context.Context, req *teamsv1.ListEnvsRequest, _ ...grpc.CallOption) (*teamsv1.ListEnvsResponse, error) {
	call := s.nextCall("ListEnvs", req)
	resp, ok := call.resp.(*teamsv1.ListEnvsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) CreateInitScript(ctx context.Context, req *teamsv1.CreateInitScriptRequest, _ ...grpc.CallOption) (*teamsv1.CreateInitScriptResponse, error) {
	call := s.nextCall("CreateInitScript", req)
	resp, ok := call.resp.(*teamsv1.CreateInitScriptResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) GetInitScript(ctx context.Context, req *teamsv1.GetInitScriptRequest, _ ...grpc.CallOption) (*teamsv1.GetInitScriptResponse, error) {
	call := s.nextCall("GetInitScript", req)
	resp, ok := call.resp.(*teamsv1.GetInitScriptResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) UpdateInitScript(ctx context.Context, req *teamsv1.UpdateInitScriptRequest, _ ...grpc.CallOption) (*teamsv1.UpdateInitScriptResponse, error) {
	call := s.nextCall("UpdateInitScript", req)
	resp, ok := call.resp.(*teamsv1.UpdateInitScriptResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) DeleteInitScript(ctx context.Context, req *teamsv1.DeleteInitScriptRequest, _ ...grpc.CallOption) (*teamsv1.DeleteInitScriptResponse, error) {
	call := s.nextCall("DeleteInitScript", req)
	resp, ok := call.resp.(*teamsv1.DeleteInitScriptResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubTeamsClient) ListInitScripts(ctx context.Context, req *teamsv1.ListInitScriptsRequest, _ ...grpc.CallOption) (*teamsv1.ListInitScriptsResponse, error) {
	call := s.nextCall("ListInitScripts", req)
	resp, ok := call.resp.(*teamsv1.ListInitScriptsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func TestTeamGetAgents(t *testing.T) {
	client := &stubTeamsClient{t: t}
	team := NewTeam(client)

	createdAt := time.Date(2024, 5, 10, 9, 30, 0, 0, time.UTC)
	updatedAt := createdAt.Add(2 * time.Hour)
	agentID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	modelID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	agent := &teamsv1.Agent{
		Meta: &teamsv1.EntityMeta{
			Id:        agentID.String(),
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		},
		Name:          "Athena",
		Role:          "assistant",
		Model:         modelID.String(),
		Configuration: "{}",
		Image:         "ghcr.io/example/agent:latest",
		Description:   "primary agent",
		Resources: &teamsv1.ComputeResources{
			RequestsCpu: "500m",
			LimitsCpu:   "1",
		},
	}

	client.Expect(teamsCall{
		method: "ListAgents",
		assert: func(req any) {
			reqValue := req.(*teamsv1.ListAgentsRequest)
			if reqValue.GetPageSize() != 10 {
				t.Fatalf("page size mismatch: got %d", reqValue.GetPageSize())
			}
			if reqValue.GetPageToken() != "token" {
				t.Fatalf("page token mismatch: got %s", reqValue.GetPageToken())
			}
		},
		resp: &teamsv1.ListAgentsResponse{
			Agents:        []*teamsv1.Agent{agent},
			NextPageToken: "next",
		},
	})

	resp, err := team.GetAgents(context.Background(), gen.GetAgentsRequestObject{
		Params: gen.GetAgentsParams{
			PageSize:  intPtr(10),
			PageToken: stringPtr("token"),
		},
	})
	if err != nil {
		t.Fatalf("GetAgents error: %v", err)
	}

	payload, ok := resp.(gen.GetAgents200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if payload.NextPageToken == nil || *payload.NextPageToken != "next" {
		t.Fatalf("unexpected next page token")
	}
	if len(payload.Items) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(payload.Items))
	}
	item := payload.Items[0]
	if item.Id != openapi_types.UUID(agentID) {
		t.Fatalf("agent id mismatch")
	}
	if item.Model != openapi_types.UUID(modelID) {
		t.Fatalf("model id mismatch")
	}
	if item.Description == nil || *item.Description != "primary agent" {
		t.Fatalf("description mismatch")
	}
	if item.Resources == nil || item.Resources.RequestsCpu == nil || *item.Resources.RequestsCpu != "500m" {
		t.Fatalf("resources mismatch")
	}
	if item.UpdatedAt == nil || !item.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("updated_at mismatch")
	}

	client.AssertDone()
}

func TestTeamPostEnvs(t *testing.T) {
	client := &stubTeamsClient{t: t}
	team := NewTeam(client)

	createdAt := time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC)
	agentID := openapiUUID(t, "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	envID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	client.Expect(teamsCall{
		method: "CreateEnv",
		assert: func(req any) {
			request := req.(*teamsv1.CreateEnvRequest)
			if request.GetAgentId() != agentID.String() {
				t.Fatalf("agent id mismatch: got %s", request.GetAgentId())
			}
			if request.GetValue() != "top-secret" {
				t.Fatalf("env value mismatch: got %s", request.GetValue())
			}
		},
		resp: &teamsv1.CreateEnvResponse{
			Env: &teamsv1.Env{
				Meta: &teamsv1.EntityMeta{
					Id:        envID.String(),
					CreatedAt: timestamppb.New(createdAt),
				},
				Name:        "TOKEN",
				Description: "service token",
				Target:      &teamsv1.Env_AgentId{AgentId: agentID.String()},
				Source:      &teamsv1.Env_Value{Value: "top-secret"},
			},
		},
	})

	resp, err := team.PostEnvs(context.Background(), gen.PostEnvsRequestObject{
		Body: &gen.EnvCreateRequest{
			Name:        "TOKEN",
			Description: stringPtr("service token"),
			AgentId:     &agentID,
			Value:       stringPtr("top-secret"),
		},
	})
	if err != nil {
		t.Fatalf("PostEnvs error: %v", err)
	}

	payload, ok := resp.(gen.PostEnvs201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if payload.AgentId == nil || *payload.AgentId != agentID {
		t.Fatalf("agent id mismatch")
	}
	if payload.Value == nil || *payload.Value != "top-secret" {
		t.Fatalf("env value mismatch")
	}
	if payload.Description == nil || *payload.Description != "service token" {
		t.Fatalf("description mismatch")
	}

	client.AssertDone()
}

func TestTeamGetInitScripts(t *testing.T) {
	client := &stubTeamsClient{t: t}
	team := NewTeam(client)

	createdAt := time.Date(2024, 7, 1, 9, 0, 0, 0, time.UTC)
	mcpID := openapiUUID(t, "cccccccc-cccc-cccc-cccc-cccccccccccc")
	scriptID := uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")

	client.Expect(teamsCall{
		method: "ListInitScripts",
		assert: func(req any) {
			request := req.(*teamsv1.ListInitScriptsRequest)
			if request.GetMcpId() != mcpID.String() {
				t.Fatalf("mcp id mismatch: got %s", request.GetMcpId())
			}
		},
		resp: &teamsv1.ListInitScriptsResponse{
			InitScripts: []*teamsv1.InitScript{
				{
					Meta: &teamsv1.EntityMeta{
						Id:        scriptID.String(),
						CreatedAt: timestamppb.New(createdAt),
					},
					Script:      "echo boot",
					Description: "boot",
					Target:      &teamsv1.InitScript_McpId{McpId: mcpID.String()},
				},
			},
			NextPageToken: "next",
		},
	})

	resp, err := team.GetInitScripts(context.Background(), gen.GetInitScriptsRequestObject{
		Params: gen.GetInitScriptsParams{
			McpId:     &mcpID,
			PageSize:  intPtr(5),
			PageToken: stringPtr("cursor"),
		},
	})
	if err != nil {
		t.Fatalf("GetInitScripts error: %v", err)
	}

	payload, ok := resp.(gen.GetInitScripts200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if payload.NextPageToken == nil || *payload.NextPageToken != "next" {
		t.Fatalf("next page token mismatch")
	}
	if len(payload.Items) != 1 {
		t.Fatalf("expected 1 init script, got %d", len(payload.Items))
	}
	item := payload.Items[0]
	if item.McpId == nil || *item.McpId != mcpID {
		t.Fatalf("mcp id mismatch")
	}

	client.AssertDone()
}

func TestTeamPostVolumeAttachments(t *testing.T) {
	client := &stubTeamsClient{t: t}
	team := NewTeam(client)

	createdAt := time.Date(2024, 8, 1, 10, 0, 0, 0, time.UTC)
	hookID := openapiUUID(t, "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	volumeID := openapiUUID(t, "ffffffff-ffff-ffff-ffff-ffffffffffff")
	attachmentID := uuid.MustParse("12121212-1212-1212-1212-121212121212")

	client.Expect(teamsCall{
		method: "CreateVolumeAttachment",
		assert: func(req any) {
			request := req.(*teamsv1.CreateVolumeAttachmentRequest)
			if request.GetVolumeId() != volumeID.String() {
				t.Fatalf("volume id mismatch: got %s", request.GetVolumeId())
			}
			if request.GetHookId() != hookID.String() {
				t.Fatalf("hook id mismatch: got %s", request.GetHookId())
			}
		},
		resp: &teamsv1.CreateVolumeAttachmentResponse{
			VolumeAttachment: &teamsv1.VolumeAttachment{
				Meta: &teamsv1.EntityMeta{
					Id:        attachmentID.String(),
					CreatedAt: timestamppb.New(createdAt),
				},
				VolumeId: volumeID.String(),
				Target:   &teamsv1.VolumeAttachment_HookId{HookId: hookID.String()},
			},
		},
	})

	resp, err := team.PostVolumeAttachments(context.Background(), gen.PostVolumeAttachmentsRequestObject{
		Body: &gen.VolumeAttachmentCreateRequest{
			VolumeId: volumeID,
			HookId:   &hookID,
		},
	})
	if err != nil {
		t.Fatalf("PostVolumeAttachments error: %v", err)
	}

	payload, ok := resp.(gen.PostVolumeAttachments201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if payload.HookId == nil || *payload.HookId != hookID {
		t.Fatalf("hook id mismatch")
	}
	if payload.VolumeId != volumeID {
		t.Fatalf("volume id mismatch")
	}

	client.AssertDone()
}

func TestTeamPatchVolumesId(t *testing.T) {
	client := &stubTeamsClient{t: t}
	team := NewTeam(client)

	createdAt := time.Date(2024, 9, 1, 11, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(30 * time.Minute)
	volumeID := openapiUUID(t, "abababab-abab-abab-abab-abababababab")
	volumeIDValue := uuid.UUID(volumeID)

	client.Expect(teamsCall{
		method: "UpdateVolume",
		assert: func(req any) {
			request := req.(*teamsv1.UpdateVolumeRequest)
			if request.GetId() != volumeID.String() {
				t.Fatalf("volume id mismatch: got %s", request.GetId())
			}
			if request.GetSize() != "5Gi" {
				t.Fatalf("size mismatch: got %s", request.GetSize())
			}
			if request.GetDescription() != "cache" {
				t.Fatalf("description mismatch: got %s", request.GetDescription())
			}
			if request.Persistent == nil || !*request.Persistent {
				t.Fatalf("persistent mismatch")
			}
		},
		resp: &teamsv1.UpdateVolumeResponse{
			Volume: &teamsv1.Volume{
				Meta: &teamsv1.EntityMeta{
					Id:        volumeIDValue.String(),
					CreatedAt: timestamppb.New(createdAt),
					UpdatedAt: timestamppb.New(updatedAt),
				},
				Persistent:  true,
				MountPath:   "/cache",
				Size:        "5Gi",
				Description: "cache",
			},
		},
	})

	resp, err := team.PatchVolumesId(context.Background(), gen.PatchVolumesIdRequestObject{
		Id: volumeID,
		Body: &gen.VolumeUpdateRequest{
			Persistent:  boolPtr(true),
			Size:        stringPtr("5Gi"),
			Description: stringPtr("cache"),
		},
	})
	if err != nil {
		t.Fatalf("PatchVolumesId error: %v", err)
	}

	payload, ok := resp.(gen.PatchVolumesId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if payload.Size == nil || *payload.Size != "5Gi" {
		t.Fatalf("size mismatch")
	}
	if payload.Description == nil || *payload.Description != "cache" {
		t.Fatalf("description mismatch")
	}
	if !payload.Persistent {
		t.Fatalf("persistent mismatch")
	}
	if payload.UpdatedAt == nil || !payload.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("updated_at mismatch")
	}

	client.AssertDone()
}

func openapiUUID(t *testing.T, value string) openapi_types.UUID {
	t.Helper()
	parsed, err := uuid.Parse(value)
	if err != nil {
		t.Fatalf("parse uuid: %v", err)
	}
	return openapi_types.UUID(parsed)
}

func intPtr(value int) *int {
	return &value
}

func boolPtr(value bool) *bool {
	return &value
}

var _ teamsv1.TeamsServiceClient = (*stubTeamsClient)(nil)
