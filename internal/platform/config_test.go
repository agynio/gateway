package platform

import "testing"

func TestLoadConfigFromEnv(t *testing.T) {
	t.Setenv("TEAMS_GRPC_TARGET", "teams:50052")
	t.Setenv("CHAT_GRPC_TARGET", "chat:50056")
	t.Setenv("FILES_GRPC_TARGET", "files:50053")
	t.Setenv("LLM_GRPC_TARGET", "llm:50054")
	t.Setenv("SECRETS_GRPC_TARGET", "secrets:50055")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := cfg.TeamsGRPCTarget; got != "teams:50052" {
		t.Fatalf("unexpected teams grpc target: %s", got)
	}

	if got := cfg.ChatGRPCTarget; got != "chat:50056" {
		t.Fatalf("unexpected chat grpc target: %s", got)
	}

	if got := cfg.FilesGRPCTarget; got != "files:50053" {
		t.Fatalf("unexpected files grpc target: %s", got)
	}

	if got := cfg.LLMGRPCTarget; got != "llm:50054" {
		t.Fatalf("unexpected llm grpc target: %s", got)
	}

	if got := cfg.SecretsGRPCTarget; got != "secrets:50055" {
		t.Fatalf("unexpected secrets grpc target: %s", got)
	}
}

func TestLoadConfigFromEnvMissingTeamsGRPC(t *testing.T) {
	t.Setenv("TEAMS_GRPC_TARGET", "")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.TeamsGRPCTarget != defaultTeamsGRPCTarget {
		t.Fatalf("unexpected teams grpc target: %s", cfg.TeamsGRPCTarget)
	}
}

func TestLoadConfigFromEnvAllDefaults(t *testing.T) {
	t.Setenv("TEAMS_GRPC_TARGET", "")
	t.Setenv("CHAT_GRPC_TARGET", "")
	t.Setenv("FILES_GRPC_TARGET", "")
	t.Setenv("LLM_GRPC_TARGET", "")
	t.Setenv("SECRETS_GRPC_TARGET", "")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.TeamsGRPCTarget != defaultTeamsGRPCTarget {
		t.Fatalf("unexpected teams grpc target: %s", cfg.TeamsGRPCTarget)
	}
	if cfg.ChatGRPCTarget != defaultChatGRPCTarget {
		t.Fatalf("unexpected chat grpc target: %s", cfg.ChatGRPCTarget)
	}
	if cfg.FilesGRPCTarget != defaultFilesGRPCTarget {
		t.Fatalf("unexpected files grpc target: %s", cfg.FilesGRPCTarget)
	}
	if cfg.LLMGRPCTarget != defaultLLMGRPCTarget {
		t.Fatalf("unexpected llm grpc target: %s", cfg.LLMGRPCTarget)
	}
	if cfg.SecretsGRPCTarget != defaultSecretsGRPCTarget {
		t.Fatalf("unexpected secrets grpc target: %s", cfg.SecretsGRPCTarget)
	}
}
