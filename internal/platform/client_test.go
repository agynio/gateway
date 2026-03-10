package platform

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientDoSuccess(t *testing.T) {
	var received struct {
		Method string
		Path   string
		Header http.Header
		Body   map[string]any
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received.Method = r.Method
		received.Path = r.URL.Path
		received.Header = r.Header.Clone()
		if err := json.NewDecoder(r.Body).Decode(&received.Body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	cfg := &Config{BaseURL: baseURL, Timeout: time.Second, Retries: 0, Headers: make(http.Header)}
	cfg.Headers.Set("X-Forwarded-For", "gateway")

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	client.retryWait = 0

	var resp struct {
		Ok bool `json:"ok"`
	}

	status, err := client.Do(context.Background(), http.MethodPost, "/team/v1/agents", nil, map[string]string{"name": "demo"}, &resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status != http.StatusCreated {
		t.Fatalf("unexpected status: %d", status)
	}

	if !resp.Ok {
		t.Fatalf("response not decoded")
	}

	if received.Method != http.MethodPost {
		t.Fatalf("unexpected method: %s", received.Method)
	}

	if received.Path != "/team/v1/agents" {
		t.Fatalf("unexpected path: %s", received.Path)
	}

	if received.Header.Get("X-Forwarded-For") != "gateway" {
		t.Fatalf("header not forwarded: %v", received.Header)
	}

	if received.Body["name"] != "demo" {
		t.Fatalf("unexpected body: %v", received.Body)
	}
}

func TestClientDoProblemError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"title":"Invalid","status":422,"detail":"bad"}`))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	cfg := &Config{BaseURL: baseURL, Timeout: time.Second, Retries: 0, Headers: make(http.Header)}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	client.retryWait = 0

	status, err := client.Do(context.Background(), http.MethodGet, "/team/v1/agents", nil, nil, nil)
	if status != http.StatusUnprocessableEntity {
		t.Fatalf("unexpected status: %d", status)
	}

	var pErr *Error
	if ok := errors.As(err, &pErr); !ok {
		t.Fatalf("expected platform error, got: %v", err)
	}

	if pErr.Problem == nil || pErr.Problem.Title != "Invalid" {
		t.Fatalf("unexpected problem: %+v", pErr.Problem)
	}
}

func TestClientRetriesIdempotentErrors(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) == 1 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	cfg := &Config{BaseURL: baseURL, Timeout: time.Second, Retries: 1, Headers: make(http.Header)}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	client.retryWait = 0

	var resp map[string]any
	status, err := client.Do(context.Background(), http.MethodGet, "/team/v1/agents", nil, nil, &resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status != http.StatusOK {
		t.Fatalf("unexpected status: %d", status)
	}

	if attempts.Load() != 2 {
		t.Fatalf("expected two attempts, got %d", attempts.Load())
	}
}
