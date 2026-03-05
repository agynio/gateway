package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

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
