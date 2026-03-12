package teamsclient

import (
	"fmt"
	"strings"

	teamsv1 "github.com/agynio/gateway/gen/agynio/api/teams/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the Teams gRPC connection and client.
type Client struct {
	conn   *grpc.ClientConn
	client teamsv1.TeamsServiceClient
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
		client: teamsv1.NewTeamsServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) TeamsServiceClient() teamsv1.TeamsServiceClient {
	if c == nil {
		return nil
	}
	return c.client
}
