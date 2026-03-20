package agentstateclient

import (
	"fmt"
	"strings"

	agentstatev1 "github.com/agynio/gateway/gen/agynio/api/agent_state/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the AgentState gRPC connection and client.
type Client struct {
	conn   *grpc.ClientConn
	client agentstatev1.AgentStateServiceClient
}

func NewClient(target string) (*Client, error) {
	if strings.TrimSpace(target) == "" {
		return nil, fmt.Errorf("target is required")
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: agentstatev1.NewAgentStateServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) AgentStateServiceClient() agentstatev1.AgentStateServiceClient {
	return c.client
}
