package httpauth

import (
	"context"
	"testing"
)

func TestExtractBearerToken(t *testing.T) {
	token, ok := ExtractBearerToken("Bearer token-123")
	if !ok {
		t.Fatalf("expected token to be extracted")
	}
	if token != "token-123" {
		t.Fatalf("expected token-123, got %q", token)
	}
}

func TestExtractBearerTokenCaseInsensitive(t *testing.T) {
	token, ok := ExtractBearerToken("bearer token-456")
	if !ok {
		t.Fatalf("expected token to be extracted")
	}
	if token != "token-456" {
		t.Fatalf("expected token-456, got %q", token)
	}
}

func TestExtractBearerTokenInvalid(t *testing.T) {
	cases := []string{
		"",
		"Basic abc",
		"Bearer",
		"Bearer   ",
		"Bearer token extra",
	}

	for _, input := range cases {
		if token, ok := ExtractBearerToken(input); ok {
			t.Fatalf("expected no token for %q, got %q", input, token)
		}
	}
}

func TestBearerTokenFromContext(t *testing.T) {
	ctx := WithBearerToken(context.Background(), " token-789 ")
	token, ok := BearerTokenFromContext(ctx)
	if !ok {
		t.Fatalf("expected token in context")
	}
	if token != "token-789" {
		t.Fatalf("expected token-789, got %q", token)
	}
}

func TestBearerTokenFromContextMissing(t *testing.T) {
	if token, ok := BearerTokenFromContext(context.Background()); ok {
		t.Fatalf("expected no token, got %q", token)
	}
}
