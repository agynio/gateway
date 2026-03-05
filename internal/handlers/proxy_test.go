package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type stubProxyTarget struct {
	url     *url.URL
	headers http.Header
}

func (s *stubProxyTarget) BaseURL() *url.URL {
	if s.url == nil {
		return nil
	}
	cloned := *s.url
	return &cloned
}

func (s *stubProxyTarget) DefaultHeaders() http.Header {
	cloned := make(http.Header, len(s.headers))
	for key, values := range s.headers {
		cloned[key] = append([]string(nil), values...)
	}
	return cloned
}

func TestNewUpstreamProxy_ForwardsRequestWithHeaders(t *testing.T) {
	t.Helper()

	captured := make(chan *http.Request, 1)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured <- r
		w.WriteHeader(http.StatusTeapot)
	}))
	t.Cleanup(upstream.Close)

	targetURL, err := url.Parse(upstream.URL)
	if err != nil {
		t.Fatalf("parse upstream url: %v", err)
	}

	handler := NewUpstreamProxy(&stubProxyTarget{
		url:     targetURL,
		headers: http.Header{"Authorization": []string{"Bearer test-token"}},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/agents/threads?limit=10", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusTeapot {
		t.Fatalf("expected status %d, got %d", http.StatusTeapot, rr.Result().StatusCode)
	}

	var forwarded *http.Request
	select {
	case forwarded = <-captured:
	default:
		t.Fatal("expected upstream request")
	}

	if forwarded.URL.Path != "/api/agents/threads" {
		t.Fatalf("unexpected path: %s", forwarded.URL.Path)
	}

	if want := "limit=10"; forwarded.URL.RawQuery != want {
		t.Fatalf("unexpected query: %s", forwarded.URL.RawQuery)
	}

	if got := forwarded.Header.Get("Authorization"); got != "Bearer test-token" {
		t.Fatalf("expected authorization header to be set, got %q", got)
	}
}
