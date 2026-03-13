package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParsePaginationDefaults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers", nil)
	pageSize, pageToken, err := parsePagination(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pageSize != defaultPageSize {
		t.Fatalf("unexpected page size: got %d want %d", pageSize, defaultPageSize)
	}
	if pageToken != "" {
		t.Fatalf("unexpected page token: %q", pageToken)
	}
}

func TestParsePaginationPageSizeOverMax(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers?pageSize=500&pageToken=next", nil)
	pageSize, pageToken, err := parsePagination(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pageSize != maxPageSize {
		t.Fatalf("unexpected page size: got %d want %d", pageSize, maxPageSize)
	}
	if pageToken != "next" {
		t.Fatalf("unexpected page token: %q", pageToken)
	}
}

func TestParsePaginationPageSizeZero(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers?pageSize=0", nil)
	if _, _, err := parsePagination(req); err == nil {
		t.Fatalf("expected error")
	}
}

func TestParsePaginationPageSizeNegative(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/llm/v1/providers?pageSize=-5", nil)
	if _, _, err := parsePagination(req); err == nil {
		t.Fatalf("expected error")
	}
}
