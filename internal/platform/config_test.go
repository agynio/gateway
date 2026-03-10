package platform

import (
	"testing"
	"time"
)

func TestLoadConfigFromEnv(t *testing.T) {
	t.Setenv("PLATFORM_BASE_URL", "https://api.example.com/team")
	t.Setenv("FILES_BASE_URL", "https://files.example.com")
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

	if cfg.FilesBaseURL == nil {
		t.Fatal("expected files base URL to be set")
	}
	if got := cfg.FilesBaseURL.String(); got != "https://files.example.com" {
		t.Fatalf("unexpected files baseURL: %s", got)
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

	if cfg.FilesBaseURL != nil {
		t.Fatalf("expected files base URL to be unset")
	}
}

func TestLoadConfigFromEnvInvalidFilesBaseURL(t *testing.T) {
	t.Setenv("PLATFORM_BASE_URL", "https://api.example.com/team")
	t.Setenv("FILES_BASE_URL", "files.example.com")

	if _, err := LoadConfigFromEnv(); err == nil {
		t.Fatalf("expected error when FILES_BASE_URL is invalid")
	}
}

func TestLoadConfigFromEnvMissingBaseURL(t *testing.T) {
	t.Setenv("PLATFORM_BASE_URL", "")

	if _, err := LoadConfigFromEnv(); err == nil {
		t.Fatalf("expected error when PLATFORM_BASE_URL is missing")
	}
}
