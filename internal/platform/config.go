package platform

import (
	"os"
	"strings"
)

const (
	defaultTeamsGRPCTarget   = "teams:50051"
	defaultFilesGRPCTarget   = "files:50051"
	defaultLLMGRPCTarget     = "llm:50051"
	defaultSecretsGRPCTarget = "secrets:50051"
)

// Config holds the runtime configuration for communicating with upstream services.
type Config struct {
	TeamsGRPCTarget   string
	FilesGRPCTarget   string
	LLMGRPCTarget     string
	SecretsGRPCTarget string
}

// LoadConfigFromEnv constructs a Config instance from environment variables.
func LoadConfigFromEnv() (*Config, error) {
	return &Config{
		TeamsGRPCTarget:   envOrDefault("TEAMS_GRPC_TARGET", defaultTeamsGRPCTarget),
		FilesGRPCTarget:   envOrDefault("FILES_GRPC_TARGET", defaultFilesGRPCTarget),
		LLMGRPCTarget:     envOrDefault("LLM_GRPC_TARGET", defaultLLMGRPCTarget),
		SecretsGRPCTarget: envOrDefault("SECRETS_GRPC_TARGET", defaultSecretsGRPCTarget),
	}, nil
}

func envOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}
