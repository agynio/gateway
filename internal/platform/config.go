package platform

import (
	"os"
	"strings"
)

const (
	defaultAgentsGRPCTarget         = "agents:50051"
	defaultThreadsGRPCTarget        = "threads:50051"
	defaultChatGRPCTarget           = "chat:50051"
	defaultNotificationsGRPCTarget  = "notifications:50051"
	defaultFilesGRPCTarget          = "files:50051"
	defaultAgentStateGRPCTarget     = "agent-state:50051"
	defaultTokenCountingGRPCTarget  = "token-counting:50051"
	defaultLLMGRPCTarget            = "llm:50051"
	defaultSecretsGRPCTarget        = "secrets:50051"
	defaultZitiManagementGRPCTarget = "ziti-management:50051"
)

// Config holds the runtime configuration for communicating with upstream services.
type Config struct {
	AgentsGRPCTarget         string
	ThreadsGRPCTarget        string
	ChatGRPCTarget           string
	NotificationsGRPCTarget  string
	FilesGRPCTarget          string
	AgentStateGRPCTarget     string
	TokenCountingGRPCTarget  string
	LLMGRPCTarget            string
	SecretsGRPCTarget        string
	ZitiIdentityFile         string
	ZitiManagementGRPCTarget string
}

// LoadConfigFromEnv constructs a Config instance from environment variables.
func LoadConfigFromEnv() (*Config, error) {
	return &Config{
		AgentsGRPCTarget:         envOrDefault("AGENTS_GRPC_TARGET", defaultAgentsGRPCTarget),
		ThreadsGRPCTarget:        envOrDefault("THREADS_GRPC_TARGET", defaultThreadsGRPCTarget),
		ChatGRPCTarget:           envOrDefault("CHAT_GRPC_TARGET", defaultChatGRPCTarget),
		NotificationsGRPCTarget:  envOrDefault("NOTIFICATIONS_GRPC_TARGET", defaultNotificationsGRPCTarget),
		FilesGRPCTarget:          envOrDefault("FILES_GRPC_TARGET", defaultFilesGRPCTarget),
		AgentStateGRPCTarget:     envOrDefault("AGENT_STATE_GRPC_TARGET", defaultAgentStateGRPCTarget),
		TokenCountingGRPCTarget:  envOrDefault("TOKEN_COUNTING_GRPC_TARGET", defaultTokenCountingGRPCTarget),
		LLMGRPCTarget:            envOrDefault("LLM_GRPC_TARGET", defaultLLMGRPCTarget),
		SecretsGRPCTarget:        envOrDefault("SECRETS_GRPC_TARGET", defaultSecretsGRPCTarget),
		ZitiIdentityFile:         strings.TrimSpace(os.Getenv("ZITI_IDENTITY_FILE")),
		ZitiManagementGRPCTarget: envOrDefault("ZITI_MANAGEMENT_GRPC_TARGET", defaultZitiManagementGRPCTarget),
	}, nil
}

func envOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}
