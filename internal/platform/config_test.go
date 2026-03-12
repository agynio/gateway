package platform

import (
	"testing"
	"time"
)

func TestLoadConfigFromEnv(t *testing.T) {
	t.Setenv("PLATFORM_BASE_URL", "https://api.example.com/team")
	t.Setenv("TEAMS_GRPC_TARGET", "teams:50052")
	t.Setenv("FILES_GRPC_TARGET", "files:50051")
	t.Setenv("LLM_GRPC_TARGET", "llm:50053")
	t.Setenv("LLM_HTTP_BASE_URL", "https://llm.example.com")
	t.Setenv("PLATFORM_TIMEOUT_MS", "1500")
	t.Setenv("PLATFORM_RETRIES", "3")
	t.Setenv("PLATFORM_AUTH_TOKEN", "secret")
	t.Setenv("PLATFORM_REQUEST_HEADERS_JSON", `{"X-Test":"ok"}`)

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := cfg.BaseURL.String(); got != "https://api.example.com/team" {
		t.Fatalf("unexpected baseURL: %s", got)
	}

	if got := cfg.TeamsGRPCTarget; got != "teams:50052" {
		t.Fatalf("unexpected teams grpc target: %s", got)
	}

	if got := cfg.FilesGRPCTarget; got != "files:50051" {
		t.Fatalf("unexpected files grpc target: %s", got)
	}

	if got := cfg.LLMGRPCTarget; got != "llm:50053" {
		t.Fatalf("unexpected llm grpc target: %s", got)
	}
	if cfg.LLMHTTPBaseURL == nil {
		t.Fatalf("expected llm http base url")
	}
	if got := cfg.LLMHTTPBaseURL.String(); got != "https://llm.example.com" {
		t.Fatalf("unexpected llm http base url: %s", got)
	}

	if cfg.Timeout != 1500*time.Millisecond {
		t.Fatalf("unexpected timeout: %s", cfg.Timeout)
	}

	if cfg.Retries != 3 {
		t.Fatalf("unexpected retries: %d", cfg.Retries)
	}
	if !cfg.RetriesConfigured {
		t.Fatalf("expected retries to be configured")
	}

	if cfg.AuthToken != "secret" {
		t.Fatalf("unexpected auth token: %s", cfg.AuthToken)
	}

	if cfg.Headers.Get("X-Test") != "ok" {
		t.Fatalf("unexpected header: %v", cfg.Headers)
	}
}

func TestLoadConfigFromEnvRetriesZero(t *testing.T) {
	t.Setenv("PLATFORM_BASE_URL", "https://api.example.com/team")
	t.Setenv("TEAMS_GRPC_TARGET", "teams:50052")
	t.Setenv("PLATFORM_RETRIES", "0")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Retries != 0 {
		t.Fatalf("unexpected retries: %d", cfg.Retries)
	}
	if !cfg.RetriesConfigured {
		t.Fatalf("expected retries to be configured")
	}
}

func TestLoadConfigFromEnvRetriesUnset(t *testing.T) {
	t.Setenv("PLATFORM_BASE_URL", "https://api.example.com/team")
	t.Setenv("TEAMS_GRPC_TARGET", "teams:50052")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Retries != defaultRetries {
		t.Fatalf("unexpected retries: %d", cfg.Retries)
	}
	if cfg.RetriesConfigured {
		t.Fatalf("expected retries to be unset")
	}
}

func TestLoadConfigFromEnvMissingBaseURL(t *testing.T) {
	t.Setenv("PLATFORM_BASE_URL", "")
	t.Setenv("TEAMS_GRPC_TARGET", "teams:50052")

	if _, err := LoadConfigFromEnv(); err == nil {
		t.Fatalf("expected error when PLATFORM_BASE_URL is missing")
	}
}

func TestLoadConfigFromEnvMissingTeamsGRPC(t *testing.T) {
	t.Setenv("PLATFORM_BASE_URL", "https://api.example.com/team")
	t.Setenv("TEAMS_GRPC_TARGET", "")

	if _, err := LoadConfigFromEnv(); err == nil {
		t.Fatalf("expected error when TEAMS_GRPC_TARGET is missing")
	}
}

func TestLoadConfigFromEnvLLMHTTPBaseURLEmpty(t *testing.T) {
	t.Setenv("PLATFORM_BASE_URL", "https://api.example.com/team")
	t.Setenv("TEAMS_GRPC_TARGET", "teams:50052")
	t.Setenv("LLM_GRPC_TARGET", "")
	t.Setenv("LLM_HTTP_BASE_URL", "")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.LLMGRPCTarget != "" {
		t.Fatalf("expected empty llm grpc target")
	}
	if cfg.LLMHTTPBaseURL != nil {
		t.Fatalf("expected nil llm http base url")
	}
}

func TestLoadConfigFromEnvLLMHTTPBaseURLInvalid(t *testing.T) {
	t.Setenv("PLATFORM_BASE_URL", "https://api.example.com/team")
	t.Setenv("TEAMS_GRPC_TARGET", "teams:50052")
	t.Setenv("LLM_HTTP_BASE_URL", "example.com")

	if _, err := LoadConfigFromEnv(); err == nil {
		t.Fatalf("expected error when LLM_HTTP_BASE_URL is invalid")
	}
}
