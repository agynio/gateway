package platform

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
	defaultZitiLeaseRenewalInterval = 2 * time.Minute
	defaultUsersGRPCTarget          = "users:50051"
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
	ZitiEnabled              bool
	ZitiLeaseRenewalInterval time.Duration
	ZitiManagementGRPCTarget string
	OIDCIssuerURL            string
	OIDCClientID             string
	UsersGRPCTarget          string
}

// LoadConfigFromEnv constructs a Config instance from environment variables.
func LoadConfigFromEnv() (*Config, error) {
	zitiEnabled, err := envBool("ZITI_ENABLED")
	if err != nil {
		return nil, err
	}

	zitiLeaseRenewalInterval, err := envDuration("ZITI_LEASE_RENEWAL_INTERVAL", defaultZitiLeaseRenewalInterval)
	if err != nil {
		return nil, err
	}
	if zitiLeaseRenewalInterval <= 0 {
		return nil, fmt.Errorf("ZITI_LEASE_RENEWAL_INTERVAL must be positive")
	}

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
		ZitiEnabled:              zitiEnabled,
		ZitiLeaseRenewalInterval: zitiLeaseRenewalInterval,
		ZitiManagementGRPCTarget: envOrDefault("ZITI_MANAGEMENT_GRPC_TARGET", defaultZitiManagementGRPCTarget),
		OIDCIssuerURL:            strings.TrimSpace(os.Getenv("OIDC_ISSUER_URL")),
		OIDCClientID:             strings.TrimSpace(os.Getenv("OIDC_CLIENT_ID")),
		UsersGRPCTarget:          envOrDefault("USERS_GRPC_TARGET", defaultUsersGRPCTarget),
	}, nil
}

func envOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func envBool(name string) (bool, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return false, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", name, err)
	}

	return parsed, nil
}

func envDuration(name string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", name, err)
	}

	return parsed, nil
}
