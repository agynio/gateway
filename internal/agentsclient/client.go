package agentsclient

import (
	"fmt"
	"strings"

	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the Agents gRPC connection and client.
type Client struct {
	conn   *grpc.ClientConn
	client agentsv1.AgentsServiceClient
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
		client: agentsv1.NewAgentsServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) AgentsServiceClient() agentsv1.AgentsServiceClient {
	return c.client
}
