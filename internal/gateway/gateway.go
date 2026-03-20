package gateway

import (
	agentstatev1 "github.com/agynio/gateway/gen/agynio/api/agent_state/v1"
	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	chatv1 "github.com/agynio/gateway/gen/agynio/api/chat/v1"
	filesv1 "github.com/agynio/gateway/gen/agynio/api/files/v1"
	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
	notificationsv1 "github.com/agynio/gateway/gen/agynio/api/notifications/v1"
	secretsv1 "github.com/agynio/gateway/gen/agynio/api/secrets/v1"
	threadsv1 "github.com/agynio/gateway/gen/agynio/api/threads/v1"
	tokencountingv1 "github.com/agynio/gateway/gen/agynio/api/token_counting/v1"
)

// Gateway forwards ConnectRPC requests to internal gRPC services.
type Gateway struct {
	agents        agentsv1.AgentsServiceClient
	threads       threadsv1.ThreadsServiceClient
	chat          chatv1.ChatServiceClient
	notifications notificationsv1.NotificationsServiceClient
	files         filesv1.FilesServiceClient
	agentState    agentstatev1.AgentStateServiceClient
	tokenCounting tokencountingv1.TokenCountingServiceClient
	llm           llmv1.LLMServiceClient
	secrets       secretsv1.SecretsServiceClient
}

func New(
	agents agentsv1.AgentsServiceClient,
	threads threadsv1.ThreadsServiceClient,
	chat chatv1.ChatServiceClient,
	notifications notificationsv1.NotificationsServiceClient,
	files filesv1.FilesServiceClient,
	agentState agentstatev1.AgentStateServiceClient,
	tokenCounting tokencountingv1.TokenCountingServiceClient,
	llm llmv1.LLMServiceClient,
	secrets secretsv1.SecretsServiceClient,
) *Gateway {
	return &Gateway{
		agents:        agents,
		threads:       threads,
		chat:          chat,
		notifications: notifications,
		files:         files,
		agentState:    agentState,
		tokenCounting: tokenCounting,
		llm:           llm,
		secrets:       secrets,
	}
}
