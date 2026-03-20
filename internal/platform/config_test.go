package platform

import "testing"

func TestLoadConfigFromEnv(t *testing.T) {
	t.Setenv("AGENTS_GRPC_TARGET", "agents:50052")
	t.Setenv("THREADS_GRPC_TARGET", "threads:50053")
	t.Setenv("CHAT_GRPC_TARGET", "chat:50054")
	t.Setenv("NOTIFICATIONS_GRPC_TARGET", "notifications:50055")
	t.Setenv("FILES_GRPC_TARGET", "files:50053")
	t.Setenv("AGENT_STATE_GRPC_TARGET", "agent-state:50056")
	t.Setenv("TOKEN_COUNTING_GRPC_TARGET", "token-counting:50057")
	t.Setenv("LLM_GRPC_TARGET", "llm:50058")
	t.Setenv("SECRETS_GRPC_TARGET", "secrets:50059")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := cfg.AgentsGRPCTarget; got != "agents:50052" {
		t.Fatalf("unexpected agents grpc target: %s", got)
	}

	if got := cfg.ThreadsGRPCTarget; got != "threads:50053" {
		t.Fatalf("unexpected threads grpc target: %s", got)
	}

	if got := cfg.ChatGRPCTarget; got != "chat:50054" {
		t.Fatalf("unexpected chat grpc target: %s", got)
	}

	if got := cfg.NotificationsGRPCTarget; got != "notifications:50055" {
		t.Fatalf("unexpected notifications grpc target: %s", got)
	}

	if got := cfg.FilesGRPCTarget; got != "files:50053" {
		t.Fatalf("unexpected files grpc target: %s", got)
	}

	if got := cfg.AgentStateGRPCTarget; got != "agent-state:50056" {
		t.Fatalf("unexpected agent state grpc target: %s", got)
	}

	if got := cfg.TokenCountingGRPCTarget; got != "token-counting:50057" {
		t.Fatalf("unexpected token counting grpc target: %s", got)
	}

	if got := cfg.LLMGRPCTarget; got != "llm:50058" {
		t.Fatalf("unexpected llm grpc target: %s", got)
	}

	if got := cfg.SecretsGRPCTarget; got != "secrets:50059" {
		t.Fatalf("unexpected secrets grpc target: %s", got)
	}
}

func TestLoadConfigFromEnvMissingAgentsGRPC(t *testing.T) {
	t.Setenv("AGENTS_GRPC_TARGET", "")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.AgentsGRPCTarget != defaultAgentsGRPCTarget {
		t.Fatalf("unexpected agents grpc target: %s", cfg.AgentsGRPCTarget)
	}
}

func TestLoadConfigFromEnvAllDefaults(t *testing.T) {
	t.Setenv("AGENTS_GRPC_TARGET", "")
	t.Setenv("THREADS_GRPC_TARGET", "")
	t.Setenv("CHAT_GRPC_TARGET", "")
	t.Setenv("NOTIFICATIONS_GRPC_TARGET", "")
	t.Setenv("FILES_GRPC_TARGET", "")
	t.Setenv("AGENT_STATE_GRPC_TARGET", "")
	t.Setenv("TOKEN_COUNTING_GRPC_TARGET", "")
	t.Setenv("LLM_GRPC_TARGET", "")
	t.Setenv("SECRETS_GRPC_TARGET", "")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.AgentsGRPCTarget != defaultAgentsGRPCTarget {
		t.Fatalf("unexpected agents grpc target: %s", cfg.AgentsGRPCTarget)
	}
	if cfg.ThreadsGRPCTarget != defaultThreadsGRPCTarget {
		t.Fatalf("unexpected threads grpc target: %s", cfg.ThreadsGRPCTarget)
	}
	if cfg.ChatGRPCTarget != defaultChatGRPCTarget {
		t.Fatalf("unexpected chat grpc target: %s", cfg.ChatGRPCTarget)
	}
	if cfg.NotificationsGRPCTarget != defaultNotificationsGRPCTarget {
		t.Fatalf("unexpected notifications grpc target: %s", cfg.NotificationsGRPCTarget)
	}
	if cfg.FilesGRPCTarget != defaultFilesGRPCTarget {
		t.Fatalf("unexpected files grpc target: %s", cfg.FilesGRPCTarget)
	}
	if cfg.AgentStateGRPCTarget != defaultAgentStateGRPCTarget {
		t.Fatalf("unexpected agent state grpc target: %s", cfg.AgentStateGRPCTarget)
	}
	if cfg.TokenCountingGRPCTarget != defaultTokenCountingGRPCTarget {
		t.Fatalf("unexpected token counting grpc target: %s", cfg.TokenCountingGRPCTarget)
	}
	if cfg.LLMGRPCTarget != defaultLLMGRPCTarget {
		t.Fatalf("unexpected llm grpc target: %s", cfg.LLMGRPCTarget)
	}
	if cfg.SecretsGRPCTarget != defaultSecretsGRPCTarget {
		t.Fatalf("unexpected secrets grpc target: %s", cfg.SecretsGRPCTarget)
	}
}
