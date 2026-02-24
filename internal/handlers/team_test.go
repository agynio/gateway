package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/agynio/gateway/internal/gen"
	"github.com/agynio/gateway/internal/platform"
)

type stubPlatformClient struct {
	t              *testing.T
	expectMethod   string
	expectPath     string
	expectQuery    url.Values
	expectBodyJSON string

	responseStatus int
	responseJSON   string
	err            error

	called bool
}

func (s *stubPlatformClient) Do(ctx context.Context, method, path string, query url.Values, body any, out any) (int, error) {
	s.called = true

	if s.expectMethod != "" && method != s.expectMethod {
		s.t.Fatalf("unexpected method: %s", method)
	}

	if s.expectPath != "" && path != s.expectPath {
		s.t.Fatalf("unexpected path: %s", path)
	}

	if s.expectQuery != nil {
		if query == nil {
			query = url.Values{}
		}
		if !reflect.DeepEqual(query, s.expectQuery) {
			s.t.Fatalf("unexpected query: %v", query)
		}
	}

	if s.expectBodyJSON != "" {
		payload, err := json.Marshal(body)
		if err != nil {
			s.t.Fatalf("marshal body: %v", err)
		}
		if !jsonEqual(payload, []byte(s.expectBodyJSON)) {
			s.t.Fatalf("unexpected body: %s", payload)
		}
	}

	if s.responseJSON != "" && out != nil {
		if err := json.Unmarshal([]byte(s.responseJSON), out); err != nil {
			s.t.Fatalf("unmarshal response: %v", err)
		}
	}

	if s.responseStatus == 0 {
		s.responseStatus = http.StatusOK
	}

	return s.responseStatus, s.err
}

func jsonEqual(a, b []byte) bool {
	var left any
	var right any
	if err := json.Unmarshal(a, &left); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &right); err != nil {
		return false
	}
	return reflect.DeepEqual(left, right)
}

func TestTeamGetAgents(t *testing.T) {
	stub := &stubPlatformClient{
		t:            t,
		expectMethod: http.MethodGet,
		expectPath:   "/team/v1/agents",
		expectQuery: url.Values{
			"q":       []string{"demo"},
			"page":    []string{"2"},
			"perPage": []string{"25"},
		},
		responseJSON: `{"items":[],"page":2,"perPage":25,"total":0}`,
	}

	handler := NewTeam(stub)

	resp, err := handler.GetAgents(context.Background(), gen.GetAgentsRequestObject{
		Params: gen.GetAgentsParams{
			Q:       ptr("demo"),
			Page:    ptr(2),
			PerPage: ptr(25),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, ok := resp.(gen.GetAgents200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}

	if result.Page != 2 || result.PerPage != 25 {
		t.Fatalf("unexpected response payload: %+v", result)
	}

	if !stub.called {
		t.Fatalf("platform client was not invoked")
	}
}

func TestTeamPostAgents(t *testing.T) {
	stub := &stubPlatformClient{
		t:              t,
		expectMethod:   http.MethodPost,
		expectPath:     "/team/v1/agents",
		expectBodyJSON: `{"config":{"model":"gpt"},"title":"Demo"}`,
		responseStatus: http.StatusCreated,
		responseJSON:   `{"config":{"model":"gpt"},"createdAt":"2024-01-01T00:00:00Z","id":"9d4cdb42-8c84-4d84-9f39-a740249bca7a"}`,
	}

	handler := NewTeam(stub)

	body := &gen.PostAgentsJSONRequestBody{}
	body.Config.Model = ptr("gpt")
	body.Title = ptr("Demo")

	resp, err := handler.PostAgents(context.Background(), gen.PostAgentsRequestObject{Body: body})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := resp.(gen.PostAgents201JSONResponse); !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}

	if !stub.called {
		t.Fatalf("platform client was not invoked")
	}
}

func TestTeamPlatformProblemError(t *testing.T) {
	stub := &stubPlatformClient{
		t:            t,
		expectMethod: http.MethodGet,
		expectPath:   "/team/v1/agents",
		err: &platform.Error{
			Status: http.StatusUnprocessableEntity,
			Problem: &platform.Problem{
				Title:  "Invalid",
				Status: http.StatusUnprocessableEntity,
				Detail: ptr("bad data"),
			},
		},
	}

	handler := NewTeam(stub)

	_, err := handler.GetAgents(context.Background(), gen.GetAgentsRequestObject{Params: gen.GetAgentsParams{}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var problemErr *ProblemError
	if !errors.As(err, &problemErr) {
		t.Fatalf("expected ProblemError, got %T", err)
	}

	if int(problemErr.Problem.Status) != http.StatusUnprocessableEntity {
		t.Fatalf("unexpected status: %d", problemErr.Problem.Status)
	}

	if problemErr.Problem.Detail == nil || *problemErr.Problem.Detail != "bad data" {
		t.Fatalf("unexpected detail: %+v", problemErr.Problem.Detail)
	}
}

func TestTeamPlatformErrorWithoutProblem(t *testing.T) {
	stub := &stubPlatformClient{
		t:            t,
		expectMethod: http.MethodGet,
		expectPath:   "/team/v1/agents",
		err: &platform.Error{
			Status: http.StatusUnauthorized,
			Body:   []byte("denied"),
		},
	}

	handler := NewTeam(stub)

	_, err := handler.GetAgents(context.Background(), gen.GetAgentsRequestObject{Params: gen.GetAgentsParams{}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var problemErr *ProblemError
	if !errors.As(err, &problemErr) {
		t.Fatalf("expected ProblemError, got %T", err)
	}

	if int(problemErr.Problem.Status) != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", problemErr.Problem.Status)
	}

	if problemErr.Problem.Detail == nil || *problemErr.Problem.Detail != "denied" {
		t.Fatalf("unexpected detail: %+v", problemErr.Problem.Detail)
	}
}

func TestTeamUnexpectedStatus(t *testing.T) {
	stub := &stubPlatformClient{
		t:              t,
		expectMethod:   http.MethodDelete,
		expectPath:     "/team/v1/agents/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		responseStatus: http.StatusOK,
	}

	handler := NewTeam(stub)

	_, err := handler.DeleteAgentsId(context.Background(), gen.DeleteAgentsIdRequestObject{Id: mustUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")})
	if err == nil {
		t.Fatalf("expected error")
	}

	if !stub.called {
		t.Fatalf("platform client was not invoked")
	}

	var problemErr *ProblemError
	if !errors.As(err, &problemErr) {
		t.Fatalf("expected ProblemError, got %T", err)
	}
}

func TestTeamWrapsNonPlatformError(t *testing.T) {
	stub := &stubPlatformClient{
		t:            t,
		expectMethod: http.MethodGet,
		expectPath:   "/team/v1/tools",
		err:          errors.New("boom"),
	}

	handler := NewTeam(stub)

	_, err := handler.GetTools(context.Background(), gen.GetToolsRequestObject{Params: gen.GetToolsParams{}})
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
}

func mustUUID(value string) openapi_types.UUID {
	return openapi_types.UUID(uuid.MustParse(value))
}
