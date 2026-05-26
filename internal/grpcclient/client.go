package grpcclient

import (
	"context"
	"fmt"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps a gRPC connection with a typed service client.
type Client[T any] struct {
	conn   *grpc.ClientConn
	client T
}

type Option func(*clientOptions)

type clientOptions struct {
	dialOptions []grpc.DialOption
}

func WithDialOption(option grpc.DialOption) Option {
	return func(options *clientOptions) {
		options.dialOptions = append(options.dialOptions, option)
	}
}

func WithContextDialer(dialer func(context.Context, string) (net.Conn, error)) Option {
	return WithDialOption(grpc.WithContextDialer(dialer))
}

func New[T any](target string, factory func(grpc.ClientConnInterface) T, opts ...Option) (*Client[T], error) {
	if strings.TrimSpace(target) == "" {
		return nil, fmt.Errorf("target is required")
	}
	if factory == nil {
		return nil, fmt.Errorf("client factory is required")
	}
	options := clientOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(identityUnaryClientInterceptor()),
		grpc.WithChainStreamInterceptor(identityStreamClientInterceptor()),
	}
	dialOptions = append(dialOptions, options.dialOptions...)
	conn, err := grpc.NewClient(target, dialOptions...)
	if err != nil {
		return nil, err
	}

	return &Client[T]{
		conn:   conn,
		client: factory(conn),
	}, nil
}

func (c *Client[T]) Close() error {
	return c.conn.Close()
}

func (c *Client[T]) Service() T {
	return c.client
}
