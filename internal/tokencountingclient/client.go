package tokencountingclient

import (
	"fmt"
	"strings"

	tokencountingv1 "github.com/agynio/gateway/gen/agynio/api/token_counting/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the TokenCounting gRPC connection and client.
type Client struct {
	conn   *grpc.ClientConn
	client tokencountingv1.TokenCountingServiceClient
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
		client: tokencountingv1.NewTokenCountingServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) TokenCountingServiceClient() tokencountingv1.TokenCountingServiceClient {
	return c.client
}
