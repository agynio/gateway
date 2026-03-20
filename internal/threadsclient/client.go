package threadsclient

import (
	"fmt"
	"strings"

	threadsv1 "github.com/agynio/gateway/gen/agynio/api/threads/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the Threads gRPC connection and client.
type Client struct {
	conn   *grpc.ClientConn
	client threadsv1.ThreadsServiceClient
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
		client: threadsv1.NewThreadsServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) ThreadsServiceClient() threadsv1.ThreadsServiceClient {
	return c.client
}
