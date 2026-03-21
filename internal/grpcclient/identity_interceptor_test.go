package grpcclient

import (
	"context"
	"testing"

	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type noopClientStream struct {
	ctx context.Context
}

func (s *noopClientStream) Header() (metadata.MD, error) {
	return metadata.MD{}, nil
}

func (s *noopClientStream) Trailer() metadata.MD {
	return metadata.MD{}
}

func (s *noopClientStream) CloseSend() error {
	return nil
}

func (s *noopClientStream) Context() context.Context {
	return s.ctx
}

func (s *noopClientStream) SendMsg(any) error {
	return nil
}

func (s *noopClientStream) RecvMsg(any) error {
	return nil
}

func TestIdentityPropagationUnaryInterceptor(t *testing.T) {
	interceptor := IdentityPropagationUnaryInterceptor()
	resolved := identity.ResolvedIdentity{IdentityID: "identity-1", IdentityType: identity.IdentityTypeUser}
	ctx := identity.WithIdentity(context.Background(), resolved)

	var got metadata.MD
	var gotOK bool
	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		got, gotOK = metadata.FromOutgoingContext(ctx)
		return nil
	}

	if err := interceptor(ctx, "/test", nil, nil, nil, invoker); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !gotOK {
		t.Fatalf("expected metadata to be present")
	}
	assertIdentityMetadata(t, got, resolved)
}

func TestIdentityPropagationUnaryInterceptorMissingIdentity(t *testing.T) {
	interceptor := IdentityPropagationUnaryInterceptor()

	var got metadata.MD
	var gotOK bool
	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		got, gotOK = metadata.FromOutgoingContext(ctx)
		return nil
	}

	if err := interceptor(context.Background(), "/test", nil, nil, nil, invoker); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertNoIdentityMetadata(t, got, gotOK)
}

func TestIdentityPropagationStreamInterceptor(t *testing.T) {
	interceptor := IdentityPropagationStreamInterceptor()
	resolved := identity.ResolvedIdentity{IdentityID: "identity-2", IdentityType: identity.IdentityTypeAgent}
	ctx := identity.WithIdentity(context.Background(), resolved)

	var got metadata.MD
	var gotOK bool
	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		got, gotOK = metadata.FromOutgoingContext(ctx)
		return &noopClientStream{ctx: ctx}, nil
	}

	stream, err := interceptor(ctx, &grpc.StreamDesc{}, nil, "/test", streamer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stream == nil {
		t.Fatalf("expected stream")
	}
	if !gotOK {
		t.Fatalf("expected metadata to be present")
	}
	assertIdentityMetadata(t, got, resolved)
}

func TestIdentityPropagationStreamInterceptorMissingIdentity(t *testing.T) {
	interceptor := IdentityPropagationStreamInterceptor()

	var got metadata.MD
	var gotOK bool
	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		got, gotOK = metadata.FromOutgoingContext(ctx)
		return &noopClientStream{ctx: ctx}, nil
	}

	stream, err := interceptor(context.Background(), &grpc.StreamDesc{}, nil, "/test", streamer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stream == nil {
		t.Fatalf("expected stream")
	}
	assertNoIdentityMetadata(t, got, gotOK)
}

func assertIdentityMetadata(t *testing.T, md metadata.MD, resolved identity.ResolvedIdentity) {
	t.Helper()
	if got := md.Get(identityIDHeader); len(got) != 1 || got[0] != resolved.IdentityID {
		t.Fatalf("expected identity id %q, got %v", resolved.IdentityID, got)
	}
	if got := md.Get(identityTypeHeader); len(got) != 1 || got[0] != string(resolved.IdentityType) {
		t.Fatalf("expected identity type %q, got %v", resolved.IdentityType, got)
	}
}

func assertNoIdentityMetadata(t *testing.T, md metadata.MD, ok bool) {
	t.Helper()
	if !ok {
		return
	}
	if got := md.Get(identityIDHeader); len(got) != 0 {
		t.Fatalf("expected no identity id metadata, got %v", got)
	}
	if got := md.Get(identityTypeHeader); len(got) != 0 {
		t.Fatalf("expected no identity type metadata, got %v", got)
	}
}
