package notificationsclient

import (
	"fmt"
	"strings"

	notificationsv1 "github.com/agynio/gateway/gen/agynio/api/notifications/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the Notifications gRPC connection and client.
type Client struct {
	conn   *grpc.ClientConn
	client notificationsv1.NotificationsServiceClient
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
		client: notificationsv1.NewNotificationsServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) NotificationsServiceClient() notificationsv1.NotificationsServiceClient {
	return c.client
}
