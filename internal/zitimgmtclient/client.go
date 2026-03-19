package zitimgmtclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	zitimgmtv1 "github.com/agynio/gateway/gen/agynio/api/ziti_management/v1"
)

type Resolver interface {
	ResolveIdentity(ctx context.Context, sourceIdentity string) (identity.ResolvedIdentity, error)
}

type Client struct {
	conn   *grpc.ClientConn
	client zitimgmtv1.ZitiManagementServiceClient
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
		client: zitimgmtv1.NewZitiManagementServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) ResolveIdentity(ctx context.Context, sourceIdentity string) (identity.ResolvedIdentity, error) {
	trimmed := strings.TrimSpace(sourceIdentity)
	if trimmed == "" {
		return identity.ResolvedIdentity{}, fmt.Errorf("source identity is required")
	}

	response, err := c.client.ResolveIdentity(ctx, &zitimgmtv1.ResolveIdentityRequest{IdentityId: trimmed})
	if err != nil {
		return identity.ResolvedIdentity{}, err
	}

	resolved := response.GetIdentity()
	if resolved == nil {
		return identity.ResolvedIdentity{}, fmt.Errorf("resolved identity missing")
	}

	return identity.ResolvedIdentity{
		IdentityID:   resolved.GetIdentityId(),
		IdentityType: resolved.GetIdentityType(),
		TenantID:     resolved.GetTenantId(),
		AuthMethod:   resolved.GetAuthMethod(),
	}, nil
}
