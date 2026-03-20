package secretsclient

import (
	"fmt"
	"strings"

	secretsv1 "github.com/agynio/gateway/gen/agynio/api/secrets/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the Secrets gRPC connection and client.
type Client struct {
	conn   *grpc.ClientConn
	client secretsv1.SecretsServiceClient
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
		client: secretsv1.NewSecretsServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SecretsServiceClient() secretsv1.SecretsServiceClient {
	return c.client
}
