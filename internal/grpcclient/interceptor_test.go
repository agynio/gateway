package grpcclient

import (
	"context"
	"testing"

	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestIdentityUnaryInterceptorAppendsMetadata(t *testing.T) {
	resolved := identity.ResolvedIdentity{
		IdentityID:   "identity-42",
		IdentityType: identity.IdentityTypeAgent,
	}
	ctx := identity.WithIdentity(context.Background(), resolved)

	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		assertOutgoingIdentity(t, ctx, resolved)
		return nil
	}

	err := IdentityUnaryInterceptor(ctx, "test", nil, nil, nil, invoker)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIdentityStreamInterceptorAppendsMetadata(t *testing.T) {
	resolved := identity.ResolvedIdentity{
		IdentityID:   "identity-99",
		IdentityType: identity.IdentityTypeUser,
	}
	ctx := identity.WithIdentity(context.Background(), resolved)

	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		assertOutgoingIdentity(t, ctx, resolved)
		return noopClientStream{}, nil
	}

	stream, err := IdentityStreamInterceptor(ctx, &grpc.StreamDesc{StreamName: "test"}, nil, "test", streamer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stream == nil {
		t.Fatalf("expected stream")
	}
}

func assertOutgoingIdentity(t *testing.T, ctx context.Context, resolved identity.ResolvedIdentity) {
	t.Helper()
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("expected outgoing metadata")
	}
	assertMetadataValue(t, md, identity.MetadataKeyIdentityID, resolved.IdentityID)
	assertMetadataValue(t, md, identity.MetadataKeyIdentityType, string(resolved.IdentityType))
}

func assertMetadataValue(t *testing.T, md metadata.MD, key, expected string) {
	t.Helper()
	values := md.Get(key)
	if len(values) != 1 {
		t.Fatalf("expected 1 value for %s, got %v", key, values)
	}
	if values[0] != expected {
		t.Fatalf("expected %s=%q, got %q", key, expected, values[0])
	}
}

type noopClientStream struct{}

func (noopClientStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (noopClientStream) Trailer() metadata.MD {
	return nil
}

func (noopClientStream) CloseSend() error {
	return nil
}

func (noopClientStream) Context() context.Context {
	return context.Background()
}

func (noopClientStream) SendMsg(m any) error {
	return nil
}

func (noopClientStream) RecvMsg(m any) error {
	return nil
}
