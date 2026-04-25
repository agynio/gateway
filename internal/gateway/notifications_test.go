package gateway

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	gatewayv1connect "github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	notificationsv1 "github.com/agynio/gateway/gen/agynio/api/notifications/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type fakeNotificationsStream struct {
	ctx       context.Context
	responses []*notificationsv1.SubscribeResponse
	err       error
	idx       int
}

func (f *fakeNotificationsStream) Recv() (*notificationsv1.SubscribeResponse, error) {
	if f.idx < len(f.responses) {
		resp := f.responses[f.idx]
		f.idx++
		return resp, nil
	}
	if f.err != nil {
		return nil, f.err
	}
	return nil, io.EOF
}

func (f *fakeNotificationsStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (f *fakeNotificationsStream) Trailer() metadata.MD {
	return nil
}

func (f *fakeNotificationsStream) CloseSend() error {
	return nil
}

func (f *fakeNotificationsStream) Context() context.Context {
	if f.ctx != nil {
		return f.ctx
	}
	return context.Background()
}

func (f *fakeNotificationsStream) SendMsg(any) error {
	return nil
}

func (f *fakeNotificationsStream) RecvMsg(any) error {
	return nil
}

type fakeNotificationsClient struct {
	subscribe      func(ctx context.Context, req *notificationsv1.SubscribeRequest) (grpc.ServerStreamingClient[notificationsv1.SubscribeResponse], error)
	subscribeReq   *notificationsv1.SubscribeRequest
	subscribeCalls int
}

func (f *fakeNotificationsClient) Publish(ctx context.Context, in *notificationsv1.PublishRequest, opts ...grpc.CallOption) (*notificationsv1.PublishResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Publish not implemented")
}

func (f *fakeNotificationsClient) Subscribe(ctx context.Context, in *notificationsv1.SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[notificationsv1.SubscribeResponse], error) {
	f.subscribeCalls++
	f.subscribeReq = in
	if f.subscribe == nil {
		return nil, status.Error(codes.Unimplemented, "Subscribe not implemented")
	}
	return f.subscribe(ctx, in)
}

type identityInterceptor struct {
	resolved identity.ResolvedIdentity
}

func (i identityInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return next(identity.WithIdentity(ctx, i.resolved), req)
	}
}

func (i identityInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i identityInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		return next(identity.WithIdentity(ctx, i.resolved), conn)
	}
}

func newNotificationsGatewayClient(t *testing.T, gateway *Gateway, opts ...connect.HandlerOption) gatewayv1connect.NotificationsGatewayClient {
	t.Helper()
	mux := http.NewServeMux()
	handlerPath, handler := gatewayv1connect.NewNotificationsGatewayHandler(gateway, opts...)
	mux.Handle(handlerPath, handler)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return gatewayv1connect.NewNotificationsGatewayClient(server.Client(), server.URL)
}

func TestNotificationsGatewaySubscribeRequiresIdentity(t *testing.T) {
	client := &fakeNotificationsClient{}
	gatewayClient := newNotificationsGatewayClient(t, &Gateway{notifications: client})

	stream, err := gatewayClient.Subscribe(context.Background(), connect.NewRequest(&notificationsv1.SubscribeRequest{Rooms: []string{"room"}}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stream.Receive() {
		t.Fatalf("expected no messages")
	}
	if connect.CodeOf(stream.Err()) != connect.CodeUnauthenticated {
		t.Fatalf("expected CodeUnauthenticated, got %v", connect.CodeOf(stream.Err()))
	}
	if client.subscribeCalls != 0 {
		t.Fatalf("expected Subscribe not to be called, got %d", client.subscribeCalls)
	}
}

func TestNotificationsGatewaySubscribePassThrough(t *testing.T) {
	client := &fakeNotificationsClient{}
	client.subscribe = func(ctx context.Context, req *notificationsv1.SubscribeRequest) (grpc.ServerStreamingClient[notificationsv1.SubscribeResponse], error) {
		return &fakeNotificationsStream{ctx: ctx, responses: []*notificationsv1.SubscribeResponse{
			{Envelope: &notificationsv1.NotificationEnvelope{Event: "notification.created"}},
		}}, nil
	}
	resolved := identity.ResolvedIdentity{IdentityID: "identity-1", IdentityType: identity.IdentityTypeUser}
	options := connect.WithInterceptors(identityInterceptor{resolved: resolved})
	gatewayClient := newNotificationsGatewayClient(t, &Gateway{notifications: client}, options)

	request := connect.NewRequest(&notificationsv1.SubscribeRequest{Rooms: []string{"room"}})
	stream, err := gatewayClient.Subscribe(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stream.Receive() {
		t.Fatalf("expected message, got error: %v", stream.Err())
	}
	if stream.Msg().GetEnvelope().GetEvent() != "notification.created" {
		t.Fatalf("unexpected event: %s", stream.Msg().GetEnvelope().GetEvent())
	}
	if stream.Receive() {
		t.Fatalf("expected single message")
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}

	if client.subscribeCalls != 1 {
		t.Fatalf("expected Subscribe to be called once, got %d", client.subscribeCalls)
	}
	if client.subscribeReq == nil || len(client.subscribeReq.GetRooms()) != 1 || client.subscribeReq.GetRooms()[0] != "room" {
		t.Fatalf("expected rooms to be forwarded")
	}
}
