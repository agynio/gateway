package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/agynio/gateway/internal/gen"
	"github.com/agynio/gateway/internal/platform"
)

type stubCall struct {
	Method       string
	Path         string
	Query        url.Values
	BodyAssert   func(any)
	Status       int
	ResponseJSON string
	Err          error
	Responder    func(any) (int, string, error)
}

type stubPlatformClient struct {
	t     *testing.T
	calls []stubCall
	idx   int
}

func (s *stubPlatformClient) Expect(call stubCall) {
	s.calls = append(s.calls, call)
}

func (s *stubPlatformClient) Do(_ context.Context, method, path string, query url.Values, body any, out any) (int, error) {
	if s.idx >= len(s.calls) {
		s.t.Fatalf("unexpected call: %s %s", method, path)
	}
	call := s.calls[s.idx]
	s.idx++

	if method != call.Method {
		s.t.Fatalf("unexpected method: got %s want %s", method, call.Method)
	}
	if path != call.Path {
		s.t.Fatalf("unexpected path: got %s want %s", path, call.Path)
	}
	if call.Query != nil {
		if query == nil {
			query = url.Values{}
		}
		if len(query) != len(call.Query) {
			s.t.Fatalf("unexpected query: got %v want %v", query, call.Query)
		}
		for key, expected := range call.Query {
			values := query[key]
			if len(values) != len(expected) {
				s.t.Fatalf("unexpected query[%s]: got %v want %v", key, values, expected)
			}
			for i := range values {
				if values[i] != expected[i] {
					s.t.Fatalf("unexpected query[%s]: got %s want %s", key, values[i], expected[i])
				}
			}
		}
	}

	if call.BodyAssert != nil {
		call.BodyAssert(body)
	}

	status := call.Status
	response := call.ResponseJSON
	err := call.Err
	if call.Responder != nil {
		status, response, err = call.Responder(body)
	}

	if response != "" && out != nil {
		if err := json.Unmarshal([]byte(response), out); err != nil {
			s.t.Fatalf("unmarshal response: %v", err)
		}
	}

	if status == 0 {
		status = http.StatusOK
	}
	return status, err
}

func (s *stubPlatformClient) AssertDone() {
	if s.idx != len(s.calls) {
		s.t.Fatalf("expected %d calls, got %d", len(s.calls), s.idx)
	}
}

func mustMarshalGraph(t *testing.T, graph graphDocument) string {
	t.Helper()
	data, err := json.Marshal(graph)
	if err != nil {
		t.Fatalf("marshal graph: %v", err)
	}
	return string(data)
}

func newGraphDocument(nodes []graphNode, edges []graphEdge, variables []graphVariable) graphDocument {
	return graphDocument{
		Name:      "main",
		Version:   1,
		UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Nodes:     nodes,
		Edges:     edges,
		Variables: variables,
	}
}

func metadataMap(t *testing.T, config map[string]any) map[string]any {
	t.Helper()
	raw, ok := config[teamMetadataKey]
	if !ok {
		t.Fatalf("metadata missing")
	}
	meta, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("metadata not map: %T", raw)
	}
	return meta
}

func TestTeamGetAgents(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	stub.Expect(stubCall{
		Method: http.MethodGet,
		Path:   "/api/graph",
		Status: http.StatusOK,
		ResponseJSON: `{
			"name":"main",
			"version":1,
			"updatedAt":"2024-01-01T00:00:00Z",
			"nodes":[
				{
					"id":"11111111-1111-1111-1111-111111111111",
					"template":"agent",
					"config":{
						"title":"Demo Agent",
						"description":"Primary",
						"config":{"model":"gpt-5","systemPrompt":"hello"},
						"_teamResource":"agent",
						"_teamVersion":1,
						"_teamCreatedAt":"2024-01-01T00:00:00Z",
						"_teamUpdatedAt":"2024-01-02T00:00:00Z"
					}
				},
				{
					"id":"22222222-2222-2222-2222-222222222222",
					"template":"tool",
					"config":{}
				}
			],
			"edges":[]
		}`,
	})

	h := NewTeam(stub)

	resp, err := h.GetAgents(context.Background(), gen.GetAgentsRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := resp.(gen.GetAgents200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}

	if len(list.Items) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(list.Items))
	}

	item := list.Items[0]
	if item.Title == nil || *item.Title != "Demo Agent" {
		t.Fatalf("unexpected title: %v", item.Title)
	}
	if item.Description == nil || *item.Description != "Primary" {
		t.Fatalf("unexpected description: %v", item.Description)
	}
	if item.Config.Model == nil || *item.Config.Model != "gpt-5" {
		t.Fatalf("unexpected config model: %+v", item.Config)
	}
	if item.Config.SystemPrompt == nil || *item.Config.SystemPrompt != "hello" {
		t.Fatalf("unexpected config prompt: %+v", item.Config)
	}
	if item.Config.RestrictOutput != nil {
		t.Fatalf("unexpected metadata in config: %+v", item.Config)
	}

	stub.AssertDone()
}

func TestTeamPostAgents(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	stub.Expect(stubCall{
		Method: http.MethodGet,
		Path:   "/api/graph",
		ResponseJSON: `{
			"name":"main",
			"version":3,
			"updatedAt":"2024-01-01T00:00:00Z",
			"nodes":[],
			"edges":[]
		}`,
	})
	var createdID string
	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		BodyAssert: func(body any) {
			doc, ok := body.(*graphDocument)
			if !ok {
				t.Fatalf("unexpected body type: %T", body)
			}
			if len(doc.Nodes) != 1 {
				t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
			}
			node := doc.Nodes[0]
			createdID = node.ID
			if node.Template != agentTemplateName {
				t.Fatalf("unexpected template: %s", node.Template)
			}
			if node.Config[agentMetadataResourceKey] != agentMetadataResourceValue {
				t.Fatalf("missing resource metadata: %+v", node.Config)
			}
			if cfg, ok := node.Config["config"].(map[string]any); !ok || cfg["model"] != "gpt-4" {
				t.Fatalf("unexpected inner config: %+v", node.Config["config"])
			}
		},
		Responder: func(any) (int, string, error) {
			if createdID == "" {
				t.Fatalf("create call did not capture node id")
			}
			response := fmt.Sprintf(`{
				"name":"main",
				"version":4,
				"updatedAt":"2024-01-02T00:00:00Z",
				"nodes":[{
					"id":"%s",
					"template":"agent",
					"config":{
						"title":"New Agent",
						"config":{"model":"gpt-4"},
						"_teamResource":"agent",
						"_teamVersion":1,
						"_teamCreatedAt":"2024-01-02T00:00:00Z",
						"_teamUpdatedAt":"2024-01-02T00:00:00Z"
					}
				}],
				"edges":[]
			}`, createdID)
			return http.StatusOK, response, nil
		},
	})

	h := NewTeam(stub)
	body := &gen.PostAgentsJSONRequestBody{}
	body.Config.Model = ptr("gpt-4")
	body.Title = ptr("New Agent")

	resp, err := h.PostAgents(context.Background(), gen.PostAgentsRequestObject{Body: body})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostAgents201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Id.String() != createdID {
		t.Fatalf("unexpected id: %s", created.Id)
	}

	stub.AssertDone()
}

func TestTeamPostAgentsConflictRetry(t *testing.T) {
	conflictErr := &platform.Error{Status: http.StatusConflict, Problem: &platform.Problem{Detail: ptr("conflict")}}
	stub := &stubPlatformClient{t: t}
	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		ResponseJSON: `{"name":"main","version":1,"updatedAt":"2024-01-01T00:00:00Z","nodes":[],"edges":[]}`,
	})
	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Err:    conflictErr,
		Status: http.StatusConflict,
	})
	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		ResponseJSON: `{"name":"main","version":2,"updatedAt":"2024-01-01T00:01:00Z","nodes":[],"edges":[]}`,
	})
	var retriedID string
	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		BodyAssert: func(body any) {
			doc, ok := body.(*graphDocument)
			if !ok {
				t.Fatalf("unexpected body type: %T", body)
			}
			if len(doc.Nodes) != 1 {
				t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
			}
			retriedID = doc.Nodes[0].ID
		},
		Responder: func(any) (int, string, error) {
			if retriedID == "" {
				t.Fatalf("retry body missing id")
			}
			response := fmt.Sprintf(`{"name":"main","version":3,"updatedAt":"2024-01-01T00:02:00Z","nodes":[{"id":"%s","template":"agent","config":{"config":{"model":"gpt"},"_teamResource":"agent","_teamVersion":1,"_teamCreatedAt":"2024-01-01T00:02:00Z","_teamUpdatedAt":"2024-01-01T00:02:00Z"}}],"edges":[]}`,
				retriedID)
			return http.StatusOK, response, nil
		},
	})

	h := NewTeam(stub)
	body := &gen.PostAgentsJSONRequestBody{}
	body.Config.Model = ptr("gpt")

	if _, err := h.PostAgents(context.Background(), gen.PostAgentsRequestObject{Body: body}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stub.AssertDone()
}

func TestTeamGetAgentsPlatformError(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	stub.Expect(stubCall{
		Method: http.MethodGet,
		Path:   "/api/graph",
		Err: &platform.Error{
			Status:  http.StatusUnprocessableEntity,
			Problem: &platform.Problem{Detail: ptr("bad data"), Status: http.StatusUnprocessableEntity},
		},
		Status: http.StatusUnprocessableEntity,
	})

	h := NewTeam(stub)

	_, err := h.GetAgents(context.Background(), gen.GetAgentsRequestObject{})
	if err == nil {
		t.Fatalf("expected error")
	}

	var problemErr *ProblemError
	if !errors.As(err, &problemErr) {
		t.Fatalf("expected ProblemError, got %T", err)
	}
	if problemErr.Problem.Status != http.StatusUnprocessableEntity {
		t.Fatalf("unexpected status: %d", problemErr.Problem.Status)
	}
	if problemErr.Problem.Detail == nil || *problemErr.Problem.Detail != "bad data" {
		t.Fatalf("unexpected detail: %v", problemErr.Problem.Detail)
	}

	stub.AssertDone()
}

func TestTeamGetAgentsUnexpectedError(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	stub.Expect(stubCall{
		Method: http.MethodGet,
		Path:   "/api/graph",
		Err:    errors.New("boom"),
	})

	h := NewTeam(stub)

	_, err := h.GetAgents(context.Background(), gen.GetAgentsRequestObject{})
	if err == nil {
		t.Fatalf("expected error")
	}

	var problemErr *ProblemError
	if !errors.As(err, &problemErr) {
		t.Fatalf("expected ProblemError, got %T", err)
	}
	if problemErr.Problem.Status != http.StatusBadGateway {
		t.Fatalf("unexpected status: %d", problemErr.Problem.Status)
	}

	stub.AssertDone()
}

func TestTeamGetTools(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	toolID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	graph := newGraphDocument([]graphNode{
		{
			ID:       toolID.String(),
			Template: "manageTool",
			Config: map[string]any{
				"type":        "manage",
				"name":        "Manage",
				"description": "Primary tool",
				"config": map[string]any{
					"foo": "bar",
				},
				teamMetadataKey: map[string]any{
					teamMetadataResourceField: toolMetadataResource,
					teamMetadataVersionField:  teamMetadataVersion,
					teamMetadataCreatedAt:     "2024-01-01T00:00:00Z",
					teamMetadataUpdatedAt:     "2024-01-02T00:00:00Z",
				},
			},
		},
	}, nil, nil)
	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, graph),
	})

	h := NewTeam(stub)

	resp, err := h.GetTools(context.Background(), gen.GetToolsRequestObject{})
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

	item := list.Items[0]
	if string(item.Type) != "manage" {
		t.Fatalf("unexpected type: %s", item.Type)
	}
	if item.Name == nil || *item.Name != "Manage" {
		t.Fatalf("unexpected name: %v", item.Name)
	}
	if item.Description == nil || *item.Description != "Primary tool" {
		t.Fatalf("unexpected description: %v", item.Description)
	}
	if item.Config == nil {
		t.Fatalf("expected config")
	}
	if _, ok := (*item.Config)["foo"]; !ok {
		t.Fatalf("expected foo in config: %v", item.Config)
	}
	if _, exists := (*item.Config)[teamMetadataKey]; exists {
		t.Fatalf("metadata leaked in config: %v", item.Config)
	}

	stub.AssertDone()
}

func TestTeamPostTools(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, newGraphDocument(nil, nil, nil)),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			graphBody, ok := body.(*graphDocument)
			if !ok {
				t.Fatalf("unexpected body type: %T", body)
			}
			if len(graphBody.Nodes) != 1 {
				t.Fatalf("expected 1 node, got %d", len(graphBody.Nodes))
			}
			node := graphBody.Nodes[0]
			if node.Template != "manageTool" {
				t.Fatalf("unexpected template: %s", node.Template)
			}
			if node.Config["type"] != "manage" {
				t.Fatalf("unexpected type config: %v", node.Config["type"])
			}
			if node.Config["name"] != "Runner" {
				t.Fatalf("unexpected name config: %v", node.Config["name"])
			}
			meta := metadataMap(t, node.Config)
			if meta[teamMetadataResourceField] != toolMetadataResource {
				t.Fatalf("unexpected resource: %v", meta)
			}
			switch version := meta[teamMetadataVersionField].(type) {
			case float64:
				if int(version) != teamMetadataVersion {
					t.Fatalf("unexpected version: %v", version)
				}
			case int:
				if version != teamMetadataVersion {
					t.Fatalf("unexpected version: %v", version)
				}
			default:
				t.Fatalf("unexpected version type: %T", version)
			}
			if _, ok := meta[teamMetadataCreatedAt].(string); !ok {
				t.Fatalf("createdAt missing: %v", meta)
			}
			if _, ok := meta[teamMetadataUpdatedAt].(string); !ok {
				t.Fatalf("updatedAt missing: %v", meta)
			}
			configValue, ok := node.Config["config"].(map[string]any)
			if !ok {
				t.Fatalf("config missing: %v", node.Config)
			}
			if configValue["enabled"] != true {
				t.Fatalf("unexpected config payload: %v", configValue)
			}
			return http.StatusOK, mustMarshalGraph(t, graphBody.Clone()), nil
		},
	})

	h := NewTeam(stub)

	resp, err := h.PostTools(context.Background(), gen.PostToolsRequestObject{
		Body: &gen.PostToolsJSONRequestBody{
			Type:        gen.PostToolsJSONBodyType("manage"),
			Name:        ptr("Runner"),
			Description: ptr("Plan executor"),
			Config: &map[string]interface{}{
				"enabled": true,
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostTools201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if string(created.Type) != "manage" {
		t.Fatalf("unexpected type: %s", created.Type)
	}
	if created.Name == nil || *created.Name != "Runner" {
		t.Fatalf("unexpected name: %v", created.Name)
	}
	if created.Config == nil || (*created.Config)["enabled"] != true {
		t.Fatalf("unexpected config: %v", created.Config)
	}
	if created.CreatedAt.IsZero() {
		t.Fatalf("expected createdAt")
	}

	stub.AssertDone()
}

func TestTeamPatchTools(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	toolID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	originalUpdated := "2024-01-02T00:00:00Z"
	graph := newGraphDocument([]graphNode{
		{
			ID:       toolID.String(),
			Template: "manageTool",
			Config: map[string]any{
				"type":        "manage",
				"name":        "Runner",
				"description": "Initial",
				"config": map[string]any{
					"enabled": true,
				},
				teamMetadataKey: map[string]any{
					teamMetadataResourceField: toolMetadataResource,
					teamMetadataVersionField:  teamMetadataVersion,
					teamMetadataCreatedAt:     "2024-01-01T00:00:00Z",
					teamMetadataUpdatedAt:     originalUpdated,
				},
			},
		},
	}, nil, nil)

	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, graph),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			graphBody := body.(*graphDocument)
			if len(graphBody.Nodes) != 1 {
				t.Fatalf("expected 1 node, got %d", len(graphBody.Nodes))
			}
			node := graphBody.Nodes[0]
			if node.Config["name"] != "Updated" {
				t.Fatalf("name not updated: %v", node.Config["name"])
			}
			if node.Config["description"] != "Revised" {
				t.Fatalf("description not updated: %v", node.Config["description"])
			}
			configValue := node.Config["config"].(map[string]any)
			if configValue["enabled"] != false {
				t.Fatalf("config not updated: %v", configValue)
			}
			meta := metadataMap(t, node.Config)
			if meta[teamMetadataCreatedAt] != "2024-01-01T00:00:00Z" {
				t.Fatalf("createdAt changed: %v", meta[teamMetadataCreatedAt])
			}
			updated, ok := meta[teamMetadataUpdatedAt].(string)
			if !ok || updated == "" {
				t.Fatalf("updatedAt missing: %v", meta)
			}
			if updated == originalUpdated {
				t.Fatalf("updatedAt not refreshed")
			}
			return http.StatusOK, mustMarshalGraph(t, graphBody.Clone()), nil
		},
	})

	h := NewTeam(stub)

	resp, err := h.PatchToolsId(context.Background(), gen.PatchToolsIdRequestObject{
		Id: toolID,
		Body: &gen.PatchToolsIdJSONRequestBody{
			Name:        ptr("Updated"),
			Description: ptr("Revised"),
			Config: &map[string]interface{}{
				"enabled": false,
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedResp, ok := resp.(gen.PatchToolsId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if updatedResp.Name == nil || *updatedResp.Name != "Updated" {
		t.Fatalf("unexpected response name: %v", updatedResp.Name)
	}
	if updatedResp.Description == nil || *updatedResp.Description != "Revised" {
		t.Fatalf("unexpected response description: %v", updatedResp.Description)
	}
	if updatedResp.Config == nil || (*updatedResp.Config)["enabled"] != false {
		t.Fatalf("unexpected response config: %v", updatedResp.Config)
	}

	stub.AssertDone()
}

func TestTeamDeleteTools(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	toolID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	edgeID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	graph := newGraphDocument([]graphNode{
		{
			ID:       toolID.String(),
			Template: "manageTool",
			Config: map[string]any{
				"type":   "manage",
				"config": map[string]any{},
				teamMetadataKey: map[string]any{
					teamMetadataResourceField: toolMetadataResource,
					teamMetadataVersionField:  teamMetadataVersion,
					teamMetadataCreatedAt:     "2024-01-01T00:00:00Z",
					teamMetadataUpdatedAt:     "2024-01-02T00:00:00Z",
				},
			},
		},
	}, []graphEdge{
		{
			ID:           ptr(edgeID),
			Source:       "44444444-4444-4444-4444-444444444444",
			SourceHandle: "tools",
			Target:       toolID.String(),
			TargetHandle: "$self",
		},
	}, []graphVariable{
		{Key: attachmentCreatedAtKey(edgeID), Value: "2024-01-01T00:00:00Z"},
		{Key: attachmentUpdatedAtKey(edgeID), Value: "2024-01-02T00:00:00Z"},
	})

	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, graph),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			graphBody := body.(*graphDocument)
			if len(graphBody.Nodes) != 0 {
				t.Fatalf("expected tool node removed")
			}
			if len(graphBody.Edges) != 0 {
				t.Fatalf("expected edges removed")
			}
			for _, variable := range graphBody.Variables {
				if strings.Contains(variable.Key, edgeID) {
					t.Fatalf("attachment variable not removed: %v", variable)
				}
			}
			return http.StatusOK, mustMarshalGraph(t, graphBody.Clone()), nil
		},
	})

	h := NewTeam(stub)

	_, err := h.DeleteToolsId(context.Background(), gen.DeleteToolsIdRequestObject{Id: toolID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stub.AssertDone()
}

func TestTeamPostAttachments(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	agentID := uuid.MustParse("55555555-5555-5555-5555-555555555555")
	toolID := uuid.MustParse("66666666-6666-6666-6666-666666666666")
	graph := newGraphDocument([]graphNode{
		{
			ID:       agentID.String(),
			Template: "agent",
			Config: map[string]any{
				"config":                  map[string]any{},
				agentMetadataResourceKey:  agentMetadataResourceValue,
				agentMetadataVersionKey:   agentMetadataVersionValue,
				agentMetadataCreatedAtKey: "2024-01-01T00:00:00Z",
			},
		},
		{
			ID:       toolID.String(),
			Template: "manageTool",
			Config: map[string]any{
				"type": "manage",
				teamMetadataKey: map[string]any{
					teamMetadataResourceField: toolMetadataResource,
					teamMetadataVersionField:  teamMetadataVersion,
					teamMetadataCreatedAt:     "2024-01-01T00:00:00Z",
					teamMetadataUpdatedAt:     "2024-01-01T00:00:00Z",
				},
			},
		},
	}, nil, nil)

	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, graph),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			graphBody := body.(*graphDocument)
			if len(graphBody.Edges) != 1 {
				t.Fatalf("expected 1 edge, got %d", len(graphBody.Edges))
			}
			edge := graphBody.Edges[0]
			if edge.Source != agentID.String() || edge.Target != toolID.String() {
				t.Fatalf("unexpected edge endpoints: %+v", edge)
			}
			if edge.SourceHandle != "tools" || edge.TargetHandle != "$self" {
				t.Fatalf("unexpected handles: %+v", edge)
			}
			if edge.ID == nil {
				t.Fatalf("edge ID missing")
			}
			createdKey := attachmentCreatedAtKey(*edge.ID)
			updatedKey := attachmentUpdatedAtKey(*edge.ID)
			foundCreated := false
			foundUpdated := false
			for _, variable := range graphBody.Variables {
				if variable.Key == createdKey {
					foundCreated = true
				}
				if variable.Key == updatedKey {
					foundUpdated = true
				}
			}
			if !foundCreated || !foundUpdated {
				t.Fatalf("attachment timestamps missing: %v", graphBody.Variables)
			}
			return http.StatusOK, mustMarshalGraph(t, graphBody.Clone()), nil
		},
	})

	h := NewTeam(stub)

	resp, err := h.PostAttachments(context.Background(), gen.PostAttachmentsRequestObject{
		Body: &gen.PostAttachmentsJSONRequestBody{
			Kind:     gen.PostAttachmentsJSONBodyKindAgentTool,
			SourceId: agentID,
			TargetId: toolID,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostAttachments201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if string(created.Kind) != "agent_tool" {
		t.Fatalf("unexpected kind: %s", created.Kind)
	}
	if created.SourceType != gen.PostAttachments201JSONResponseSourceType("agent") {
		t.Fatalf("unexpected source type: %s", created.SourceType)
	}
	if created.TargetType != gen.PostAttachments201JSONResponseTargetType("tool") {
		t.Fatalf("unexpected target type: %s", created.TargetType)
	}
	if created.CreatedAt.IsZero() {
		t.Fatalf("expected createdAt")
	}

	stub.AssertDone()
}

func TestTeamGetAttachments(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	agentID := uuid.MustParse("77777777-7777-7777-7777-777777777777")
	toolID := uuid.MustParse("88888888-8888-8888-8888-888888888888")
	edgeID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	graph := newGraphDocument([]graphNode{
		{ID: agentID.String(), Template: "agent"},
		{ID: toolID.String(), Template: "manageTool", Config: map[string]any{
			"type": "manage",
			teamMetadataKey: map[string]any{
				teamMetadataResourceField: toolMetadataResource,
				teamMetadataVersionField:  teamMetadataVersion,
				teamMetadataCreatedAt:     "2024-01-01T00:00:00Z",
			},
		}},
	}, []graphEdge{
		{
			ID:           ptr(edgeID),
			Source:       agentID.String(),
			SourceHandle: "tools",
			Target:       toolID.String(),
			TargetHandle: "$self",
		},
	}, []graphVariable{
		{Key: attachmentCreatedAtKey(edgeID), Value: "2024-01-01T00:00:00Z"},
		{Key: attachmentUpdatedAtKey(edgeID), Value: "2024-01-02T00:00:00Z"},
	})

	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, graph),
	})

	h := NewTeam(stub)

	resp, err := h.GetAttachments(context.Background(), gen.GetAttachmentsRequestObject{})
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
	item := list.Items[0]
	if string(item.Kind) != "agent_tool" {
		t.Fatalf("unexpected kind: %s", item.Kind)
	}
	if item.SourceId.String() != agentID.String() || item.TargetId.String() != toolID.String() {
		t.Fatalf("unexpected endpoints: %s -> %s", item.SourceId, item.TargetId)
	}
	if item.CreatedAt.IsZero() {
		t.Fatalf("expected createdAt")
	}

	stub.AssertDone()
}

func TestTeamDeleteAttachments(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	edgeID := "cccccccc-cccc-cccc-cccc-cccccccccccc"
	graph := newGraphDocument([]graphNode{
		{ID: "99999999-9999-9999-9999-999999999999", Template: "agent"},
		{ID: "00000000-0000-0000-0000-000000000000", Template: "manageTool", Config: map[string]any{
			"type": "manage",
			teamMetadataKey: map[string]any{
				teamMetadataResourceField: toolMetadataResource,
				teamMetadataVersionField:  teamMetadataVersion,
				teamMetadataCreatedAt:     "2024-01-01T00:00:00Z",
			},
		}},
	}, []graphEdge{
		{
			ID:           ptr(edgeID),
			Source:       "99999999-9999-9999-9999-999999999999",
			SourceHandle: "tools",
			Target:       "00000000-0000-0000-0000-000000000000",
			TargetHandle: "$self",
		},
	}, []graphVariable{
		{Key: attachmentCreatedAtKey(edgeID), Value: "2024-01-01T00:00:00Z"},
		{Key: attachmentUpdatedAtKey(edgeID), Value: "2024-01-02T00:00:00Z"},
	})

	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, graph),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			graphBody := body.(*graphDocument)
			if len(graphBody.Edges) != 0 {
				t.Fatalf("expected edges removed")
			}
			for _, variable := range graphBody.Variables {
				if strings.Contains(variable.Key, edgeID) {
					t.Fatalf("attachment variable not removed: %v", variable)
				}
			}
			return http.StatusOK, mustMarshalGraph(t, graphBody.Clone()), nil
		},
	})

	h := NewTeam(stub)

	attachmentID := uuid.MustParse(edgeID)
	_, err := h.DeleteAttachmentsId(context.Background(), gen.DeleteAttachmentsIdRequestObject{Id: attachmentID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stub.AssertDone()
}

func TestTeamPostMcpServers(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, newGraphDocument(nil, nil, nil)),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			graphBody := body.(*graphDocument)
			if len(graphBody.Nodes) != 1 {
				t.Fatalf("expected node appended")
			}
			node := graphBody.Nodes[0]
			if node.Template != mcpServerTemplateName {
				t.Fatalf("unexpected template: %s", node.Template)
			}
			config, ok := node.Config["config"].(map[string]any)
			if !ok {
				t.Fatalf("config missing: %v", node.Config)
			}
			if config["command"] != "./run" {
				t.Fatalf("command not set: %v", config)
			}
			meta := metadataMap(t, node.Config)
			if meta[teamMetadataResourceField] != mcpServerMetadataResource {
				t.Fatalf("metadata resource mismatch: %v", meta)
			}
			return http.StatusOK, mustMarshalGraph(t, graphBody.Clone()), nil
		},
	})

	h := NewTeam(stub)

	body := gen.PostMcpServersJSONRequestBody{}
	body.Title = ptr("Responder")
	body.Description = ptr("Handles requests")
	body.Config.Command = ptr("./run")
	resp, err := h.PostMcpServers(context.Background(), gen.PostMcpServersRequestObject{
		Body: &body,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostMcpServers201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Title == nil || *created.Title != "Responder" {
		t.Fatalf("unexpected title: %v", created.Title)
	}

	stub.AssertDone()
}

func TestTeamPatchMcpServers(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	serverID := uuid.MustParse("12121212-1212-1212-1212-121212121212")
	graph := newGraphDocument([]graphNode{
		{
			ID:       serverID.String(),
			Template: mcpServerTemplateName,
			Config: map[string]any{
				"title":       "Original",
				"description": "First",
				"config":      map[string]any{"command": "./run"},
				teamMetadataKey: map[string]any{
					teamMetadataResourceField: mcpServerMetadataResource,
					teamMetadataVersionField:  teamMetadataVersion,
					teamMetadataCreatedAt:     "2024-01-01T00:00:00Z",
					teamMetadataUpdatedAt:     "2024-01-02T00:00:00Z",
				},
			},
		},
	}, nil, nil)

	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, graph),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			node := body.(*graphDocument).Nodes[0]
			if node.Config["title"] != "Updated" {
				t.Fatalf("title not updated: %v", node.Config["title"])
			}
			config := node.Config["config"].(map[string]any)
			if config["command"] != "./start" {
				t.Fatalf("config not updated: %v", config)
			}
			return http.StatusOK, mustMarshalGraph(t, body.(*graphDocument).Clone()), nil
		},
	})

	h := NewTeam(stub)
	config := struct {
		Command *string `json:"command,omitempty"`
		Env     *[]struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"env,omitempty"`
		HeartbeatIntervalMs *int    `json:"heartbeatIntervalMs,omitempty"`
		Namespace           *string `json:"namespace,omitempty"`
		RequestTimeoutMs    *int    `json:"requestTimeoutMs,omitempty"`
		Restart             *struct {
			BackoffMs   *int `json:"backoffMs,omitempty"`
			MaxAttempts *int `json:"maxAttempts,omitempty"`
		} `json:"restart,omitempty"`
		StaleTimeoutMs   *int    `json:"staleTimeoutMs,omitempty"`
		StartupTimeoutMs *int    `json:"startupTimeoutMs,omitempty"`
		Workdir          *string `json:"workdir,omitempty"`
	}{}
	config.Command = ptr("./start")
	body := gen.PatchMcpServersIdJSONRequestBody{Title: ptr("Updated"), Config: &config}
	resp, err := h.PatchMcpServersId(context.Background(), gen.PatchMcpServersIdRequestObject{
		Id:   serverID,
		Body: &body,
	})
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
	stub := &stubPlatformClient{t: t}
	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, newGraphDocument(nil, nil, nil)),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			node := body.(*graphDocument).Nodes[0]
			if node.Template != workspaceTemplateName {
				t.Fatalf("unexpected template: %s", node.Template)
			}
			config := node.Config["config"].(map[string]any)
			if config["image"] != "workspace:latest" {
				t.Fatalf("unexpected config: %v", config)
			}
			meta := metadataMap(t, node.Config)
			if meta[teamMetadataResourceField] != workspaceMetadataResource {
				t.Fatalf("metadata mismatch: %v", meta)
			}
			return http.StatusOK, mustMarshalGraph(t, body.(*graphDocument).Clone()), nil
		},
	})

	h := NewTeam(stub)
	body := gen.PostWorkspaceConfigurationsJSONRequestBody{}
	body.Title = ptr("Dev")
	body.Description = ptr("Development")
	body.Config.Image = ptr("workspace:latest")
	resp, err := h.PostWorkspaceConfigurations(context.Background(), gen.PostWorkspaceConfigurationsRequestObject{
		Body: &body,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostWorkspaceConfigurations201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Title == nil || *created.Title != "Dev" {
		t.Fatalf("unexpected title: %v", created.Title)
	}

	stub.AssertDone()
}

func TestTeamPatchWorkspaceConfigurations(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	workspaceID := uuid.MustParse("34343434-3434-3434-3434-343434343434")
	graph := newGraphDocument([]graphNode{
		{
			ID:       workspaceID.String(),
			Template: workspaceTemplateName,
			Config: map[string]any{
				"title":  "Dev",
				"config": map[string]any{"image": "workspace:latest"},
				teamMetadataKey: map[string]any{
					teamMetadataResourceField: workspaceMetadataResource,
					teamMetadataVersionField:  teamMetadataVersion,
					teamMetadataCreatedAt:     "2024-01-01T00:00:00Z",
					teamMetadataUpdatedAt:     "2024-01-02T00:00:00Z",
				},
			},
		},
	}, nil, nil)

	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, graph),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			node := body.(*graphDocument).Nodes[0]
			if node.Config["title"] != "Prod" {
				t.Fatalf("title not updated")
			}
			return http.StatusOK, mustMarshalGraph(t, body.(*graphDocument).Clone()), nil
		},
	})

	h := NewTeam(stub)
	resp, err := h.PatchWorkspaceConfigurationsId(context.Background(), gen.PatchWorkspaceConfigurationsIdRequestObject{
		Id: workspaceID,
		Body: &gen.PatchWorkspaceConfigurationsIdJSONRequestBody{
			Title: ptr("Prod"),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	patched, ok := resp.(gen.PatchWorkspaceConfigurationsId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if patched.Title == nil || *patched.Title != "Prod" {
		t.Fatalf("unexpected title: %v", patched.Title)
	}

	stub.AssertDone()
}

func TestTeamPostMemoryBuckets(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, newGraphDocument(nil, nil, nil)),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			node := body.(*graphDocument).Nodes[0]
			if node.Template != memoryBucketTemplateName {
				t.Fatalf("unexpected template: %s", node.Template)
			}
			meta := metadataMap(t, node.Config)
			if meta[teamMetadataResourceField] != memoryBucketMetadataResource {
				t.Fatalf("metadata mismatch: %v", meta)
			}
			return http.StatusOK, mustMarshalGraph(t, body.(*graphDocument).Clone()), nil
		},
	})

	h := NewTeam(stub)
	resp, err := h.PostMemoryBuckets(context.Background(), gen.PostMemoryBucketsRequestObject{
		Body: &gen.PostMemoryBucketsJSONRequestBody{
			Title:       ptr("User State"),
			Description: ptr("Stores user preferences"),
			Config: struct {
				CollectionPrefix *string                                   `json:"collectionPrefix,omitempty"`
				Scope            *gen.PostMemoryBucketsJSONBodyConfigScope `json:"scope,omitempty"`
			}{
				CollectionPrefix: ptr("users"),
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(gen.PostMemoryBuckets201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Title == nil || *created.Title != "User State" {
		t.Fatalf("unexpected title: %v", created.Title)
	}

	stub.AssertDone()
}

func TestTeamPatchMemoryBuckets(t *testing.T) {
	stub := &stubPlatformClient{t: t}
	bucketID := uuid.MustParse("45454545-4545-4545-4545-454545454545")
	graph := newGraphDocument([]graphNode{
		{
			ID:       bucketID.String(),
			Template: memoryBucketTemplateName,
			Config: map[string]any{
				"title":  "User State",
				"config": map[string]any{"collectionPrefix": "users"},
				teamMetadataKey: map[string]any{
					teamMetadataResourceField: memoryBucketMetadataResource,
					teamMetadataVersionField:  teamMetadataVersion,
					teamMetadataCreatedAt:     "2024-01-01T00:00:00Z",
					teamMetadataUpdatedAt:     "2024-01-02T00:00:00Z",
				},
			},
		},
	}, nil, nil)

	stub.Expect(stubCall{
		Method:       http.MethodGet,
		Path:         "/api/graph",
		Status:       http.StatusOK,
		ResponseJSON: mustMarshalGraph(t, graph),
	})

	stub.Expect(stubCall{
		Method: http.MethodPost,
		Path:   "/api/graph",
		Responder: func(body any) (int, string, error) {
			node := body.(*graphDocument).Nodes[0]
			config := node.Config["config"].(map[string]any)
			if config["collectionPrefix"] != "sessions" {
				t.Fatalf("config not updated: %v", config)
			}
			return http.StatusOK, mustMarshalGraph(t, body.(*graphDocument).Clone()), nil
		},
	})

	h := NewTeam(stub)
	resp, err := h.PatchMemoryBucketsId(context.Background(), gen.PatchMemoryBucketsIdRequestObject{
		Id: bucketID,
		Body: &gen.PatchMemoryBucketsIdJSONRequestBody{
			Config: &struct {
				CollectionPrefix *string                                      `json:"collectionPrefix,omitempty"`
				Scope            *gen.PatchMemoryBucketsIdJSONBodyConfigScope `json:"scope,omitempty"`
			}{
				CollectionPrefix: ptr("sessions"),
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, ok := resp.(gen.PatchMemoryBucketsId200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if updated.Config.CollectionPrefix == nil || *updated.Config.CollectionPrefix != "sessions" {
		t.Fatalf("unexpected config: %+v", updated.Config)
	}

	stub.AssertDone()
}
