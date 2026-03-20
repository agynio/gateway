package llmclient

import (
	"fmt"
	"strings"

	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the LLM gRPC connection and client.
type Client struct {
	conn   *grpc.ClientConn
	client llmv1.LLMServiceClient
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
		client: llmv1.NewLLMServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) LLMServiceClient() llmv1.LLMServiceClient {
	return c.client
}
