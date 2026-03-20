package chatclient

import (
	"fmt"
	"strings"

	chatv1 "github.com/agynio/gateway/gen/agynio/api/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the Chat gRPC connection and client.
type Client struct {
	conn   *grpc.ClientConn
	client chatv1.ChatServiceClient
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
		client: chatv1.NewChatServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) ChatServiceClient() chatv1.ChatServiceClient {
	return c.client
}
