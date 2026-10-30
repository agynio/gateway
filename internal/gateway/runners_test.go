package gateway

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	gatewayv1connect "github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	runnerv1 "github.com/agynio/gateway/gen/agynio/api/runner/v1"
	runnersv1 "github.com/agynio/gateway/gen/agynio/api/runners/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakeStream struct {
	ctx       context.Context
	responses []*runnerv1.StreamWorkloadLogsResponse
	err       error
	idx       int
}

func (f *fakeStream) Recv() (*runnerv1.StreamWorkloadLogsResponse, error) {
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

func (f *fakeStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (f *fakeStream) Trailer() metadata.MD {
	return nil
}

func (f *fakeStream) CloseSend() error {
	return nil
}

func (f *fakeStream) Context() context.Context {
	if f.ctx != nil {
		return f.ctx
	}
	return context.Background()
}

func (f *fakeStream) SendMsg(any) error {
	return nil
}

func (f *fakeStream) RecvMsg(any) error {
	return nil
}

// Embedded for the same reason as fakeAgentsClient: this fake answers three
// RPCs and used to restate thirty.
type fakeRunnersClient struct {
	runnersv1.RunnersServiceClient
	getWorkload                func(ctx context.Context, req *runnersv1.GetWorkloadRequest) (*runnersv1.GetWorkloadResponse, error)
	getWorkloadReq             *runnersv1.GetWorkloadRequest
	getWorkloadMetadata        metadata.MD
	getWorkloadCalls           int
	reportWorkloadIdle         func(ctx context.Context, req *runnersv1.ReportWorkloadIdleRequest) (*runnersv1.ReportWorkloadIdleResponse, error)
	reportWorkloadIdleReq      *runnersv1.ReportWorkloadIdleRequest
	reportWorkloadIdleMetadata metadata.MD
	reportWorkloadIdleCalls    int
	streamWorkloadLogs         func(ctx context.Context, req *runnerv1.StreamWorkloadLogsRequest) (grpc.ServerStreamingClient[runnerv1.StreamWorkloadLogsResponse], error)
	streamWorkloadLogsReq      *runnerv1.StreamWorkloadLogsRequest
	streamWorkloadLogsCalls    int
}

func (f *fakeRunnersClient) GetWorkload(ctx context.Context, in *runnersv1.GetWorkloadRequest, opts ...grpc.CallOption) (*runnersv1.GetWorkloadResponse, error) {
	f.getWorkloadCalls++
	f.getWorkloadReq = in
	f.getWorkloadMetadata, _ = metadata.FromOutgoingContext(ctx)
	if f.getWorkload == nil {
		return nil, status.Error(codes.Unimplemented, "GetWorkload not implemented")
	}
	return f.getWorkload(ctx, in)
}

func (f *fakeRunnersClient) ReportWorkloadIdle(ctx context.Context, in *runnersv1.ReportWorkloadIdleRequest, opts ...grpc.CallOption) (*runnersv1.ReportWorkloadIdleResponse, error) {
	f.reportWorkloadIdleCalls++
	f.reportWorkloadIdleReq = in
	f.reportWorkloadIdleMetadata, _ = metadata.FromOutgoingContext(ctx)
	if f.reportWorkloadIdle == nil {
		return nil, status.Error(codes.Unimplemented, "ReportWorkloadIdle not implemented")
	}
	return f.reportWorkloadIdle(ctx, in)
}

func (f *fakeRunnersClient) StreamWorkloadLogs(ctx context.Context, in *runnerv1.StreamWorkloadLogsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[runnerv1.StreamWorkloadLogsResponse], error) {
	f.streamWorkloadLogsCalls++
	f.streamWorkloadLogsReq = in
	if f.streamWorkloadLogs == nil {
		return nil, status.Error(codes.Unimplemented, "StreamWorkloadLogs not implemented")
	}
	return f.streamWorkloadLogs(ctx, in)
}

func newRunnersGatewayClient(t *testing.T, gateway *RunnersGateway, opts ...connect.HandlerOption) gatewayv1connect.RunnersGatewayClient {
	t.Helper()
	mux := http.NewServeMux()
	handlerPath, handler := gatewayv1connect.NewRunnersGatewayHandler(gateway, opts...)
	mux.Handle(handlerPath, handler)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return gatewayv1connect.NewRunnersGatewayClient(server.Client(), server.URL)
}

func assertMetadataValue(t *testing.T, md metadata.MD, key string, expected string) {
	t.Helper()
	got := md.Get(key)
	if len(got) != 1 || got[0] != expected {
		t.Fatalf("expected %s %q, got %#v", key, expected, got)
	}
}

func TestRunnersGatewayGetWorkloadPropagatesIdentityMetadata(t *testing.T) {
	client := &fakeRunnersClient{}
	client.getWorkload = func(ctx context.Context, req *runnersv1.GetWorkloadRequest) (*runnersv1.GetWorkloadResponse, error) {
		return &runnersv1.GetWorkloadResponse{Workload: &runnersv1.Workload{}}, nil
	}
	gateway := NewRunnersGateway(client)
	resolved := identity.ResolvedIdentity{
		IdentityID:   "identity-1",
		IdentityType: identity.IdentityTypeUser,
	}
	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		identity.MetadataKeyIdentityID, "stale-identity",
		identity.MetadataKeyIdentityType, string(identity.IdentityTypeRunner),
	)
	ctx = identity.WithIdentity(ctx, resolved)

	_, err := gateway.GetWorkload(ctx, connect.NewRequest(&runnersv1.GetWorkloadRequest{Id: "workload-1"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.getWorkloadCalls != 1 {
		t.Fatalf("expected GetWorkload to be called once, got %d", client.getWorkloadCalls)
	}
	if client.getWorkloadReq.GetId() != "workload-1" {
		t.Fatalf("expected workload id to be forwarded")
	}
	assertMetadataValue(t, client.getWorkloadMetadata, identity.MetadataKeyIdentityID, resolved.IdentityID)
	assertMetadataValue(t, client.getWorkloadMetadata, identity.MetadataKeyIdentityType, string(resolved.IdentityType))
}

func TestRunnersGatewayReportWorkloadIdleForwardsCallerIdentity(t *testing.T) {
	client := &fakeRunnersClient{}
	client.reportWorkloadIdle = func(ctx context.Context, req *runnersv1.ReportWorkloadIdleRequest) (*runnersv1.ReportWorkloadIdleResponse, error) {
		return &runnersv1.ReportWorkloadIdleResponse{}, nil
	}
	resolved := identity.ResolvedIdentity{
		IdentityID:   "instance-1",
		IdentityType: identity.IdentityTypeAgentInstance,
		WorkloadID:   "workload-1",
	}
	gatewayClient := newRunnersGatewayClient(t, NewRunnersGateway(client), connect.WithInterceptors(identityInterceptor{resolved: resolved}))

	_, err := gatewayClient.ReportWorkloadIdle(context.Background(), connect.NewRequest(&runnersv1.ReportWorkloadIdleRequest{Id: "workload-1"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.reportWorkloadIdleCalls != 1 {
		t.Fatalf("expected ReportWorkloadIdle to be called once, got %d", client.reportWorkloadIdleCalls)
	}
	if client.reportWorkloadIdleReq.GetId() != "workload-1" {
		t.Fatalf("expected workload id to be forwarded")
	}
	assertMetadataValue(t, client.reportWorkloadIdleMetadata, identity.MetadataKeyIdentityID, resolved.IdentityID)
	assertMetadataValue(t, client.reportWorkloadIdleMetadata, identity.MetadataKeyIdentityType, string(resolved.IdentityType))
	assertMetadataValue(t, client.reportWorkloadIdleMetadata, identity.MetadataKeyWorkloadID, resolved.WorkloadID)
}

// The runners service decides whether the caller owns the workload; what the
// Gateway guarantees is that it decides on the identity the caller
// authenticated as, whatever the request itself claims.
func TestRunnersGatewayReportWorkloadIdleDeniesAnotherWorkload(t *testing.T) {
	owners := map[string]string{"workload-1": "instance-1", "workload-2": "instance-2"}
	client := &fakeRunnersClient{}
	client.reportWorkloadIdle = func(ctx context.Context, req *runnersv1.ReportWorkloadIdleRequest) (*runnersv1.ReportWorkloadIdleResponse, error) {
		md, _ := metadata.FromOutgoingContext(ctx)
		caller := md.Get(identity.MetadataKeyIdentityID)
		if len(caller) != 1 || caller[0] != owners[req.GetId()] {
			return nil, status.Error(codes.PermissionDenied, "workload is not owned by the caller")
		}
		return &runnersv1.ReportWorkloadIdleResponse{}, nil
	}
	resolved := identity.ResolvedIdentity{
		IdentityID:   "instance-1",
		IdentityType: identity.IdentityTypeAgentInstance,
		WorkloadID:   "workload-1",
	}
	gatewayClient := newRunnersGatewayClient(t, NewRunnersGateway(client), connect.WithInterceptors(identityInterceptor{resolved: resolved}))

	req := connect.NewRequest(&runnersv1.ReportWorkloadIdleRequest{Id: "workload-2"})
	req.Header().Set(identity.MetadataKeyIdentityID, "instance-2")
	_, err := gatewayClient.ReportWorkloadIdle(context.Background(), req)
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected code %v, got %v", connect.CodePermissionDenied, err)
	}
	assertMetadataValue(t, client.reportWorkloadIdleMetadata, identity.MetadataKeyIdentityID, resolved.IdentityID)
}

func TestRunnersGatewayReportWorkloadIdlePropagatesError(t *testing.T) {
	client := &fakeRunnersClient{}
	client.reportWorkloadIdle = func(ctx context.Context, req *runnersv1.ReportWorkloadIdleRequest) (*runnersv1.ReportWorkloadIdleResponse, error) {
		return nil, status.Error(codes.Unavailable, "runners unavailable")
	}
	gatewayClient := newRunnersGatewayClient(t, NewRunnersGateway(client))

	_, err := gatewayClient.ReportWorkloadIdle(context.Background(), connect.NewRequest(&runnersv1.ReportWorkloadIdleRequest{Id: "workload-1"}))
	if connect.CodeOf(err) != connect.CodeUnavailable {
		t.Fatalf("expected code %v, got %v", connect.CodeUnavailable, connect.CodeOf(err))
	}
}

func TestStreamWorkloadLogs_PassThrough(t *testing.T) {
	requestWorkloadID := " workload-1 "
	requestContainerName := " container-1 "
	sinceTime := timestamppb.New(time.Unix(1723123123, 0))
	client := &fakeRunnersClient{}
	client.streamWorkloadLogs = func(ctx context.Context, req *runnerv1.StreamWorkloadLogsRequest) (grpc.ServerStreamingClient[runnerv1.StreamWorkloadLogsResponse], error) {
		return &fakeStream{ctx: ctx, responses: []*runnerv1.StreamWorkloadLogsResponse{
			{
				Event: &runnerv1.StreamWorkloadLogsResponse_Chunk{
					Chunk: &runnerv1.LogChunk{Data: []byte("hello")},
				},
			},
		}}, nil
	}

	gatewayClient := newRunnersGatewayClient(t, NewRunnersGateway(client))
	stream, err := gatewayClient.StreamWorkloadLogs(context.Background(), connect.NewRequest(&runnerv1.StreamWorkloadLogsRequest{
		WorkloadId:    requestWorkloadID,
		ContainerName: requestContainerName,
		Follow:        true,
		TailLines:     3,
		SinceTime:     sinceTime,
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stream.Receive() {
		t.Fatalf("expected log chunk, got error: %v", stream.Err())
	}
	if string(stream.Msg().GetChunk().GetData()) != "hello" {
		t.Fatalf("unexpected log data: %q", string(stream.Msg().GetChunk().GetData()))
	}
	if stream.Receive() {
		t.Fatalf("expected single log chunk")
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}

	if client.streamWorkloadLogsCalls != 1 {
		t.Fatalf("expected StreamWorkloadLogs to be called once, got %d", client.streamWorkloadLogsCalls)
	}
	if client.streamWorkloadLogsReq == nil {
		t.Fatal("expected StreamWorkloadLogs request to be recorded")
	}
	if client.streamWorkloadLogsReq.GetWorkloadId() != requestWorkloadID {
		t.Fatalf("expected workload id %q, got %q", requestWorkloadID, client.streamWorkloadLogsReq.GetWorkloadId())
	}
	if client.streamWorkloadLogsReq.GetContainerName() != requestContainerName {
		t.Fatalf("expected container name %q, got %q", requestContainerName, client.streamWorkloadLogsReq.GetContainerName())
	}
	if client.streamWorkloadLogsReq.GetTailLines() != 3 {
		t.Fatalf("expected tail lines 3, got %d", client.streamWorkloadLogsReq.GetTailLines())
	}
	if client.streamWorkloadLogsReq.GetFollow() != true {
		t.Fatalf("expected follow true")
	}
	if client.streamWorkloadLogsReq.GetSinceTime() == nil || client.streamWorkloadLogsReq.GetSinceTime().AsTime().Unix() != sinceTime.AsTime().Unix() {
		t.Fatalf("expected since time %v, got %v", sinceTime, client.streamWorkloadLogsReq.GetSinceTime())
	}
}

func TestStreamWorkloadLogs_PropagatesClientError(t *testing.T) {
	client := &fakeRunnersClient{}
	client.streamWorkloadLogs = func(ctx context.Context, req *runnerv1.StreamWorkloadLogsRequest) (grpc.ServerStreamingClient[runnerv1.StreamWorkloadLogsResponse], error) {
		return nil, status.Error(codes.NotFound, "missing workload")
	}

	gatewayClient := newRunnersGatewayClient(t, NewRunnersGateway(client))
	stream, err := gatewayClient.StreamWorkloadLogs(context.Background(), connect.NewRequest(&runnerv1.StreamWorkloadLogsRequest{
		WorkloadId:    "workload-1",
		ContainerName: "container-1",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stream.Receive() {
		t.Fatalf("expected no messages")
	}
	streamErr := stream.Err()
	if streamErr == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(streamErr) != connect.CodeNotFound {
		t.Fatalf("expected code %v, got %v", connect.CodeNotFound, connect.CodeOf(streamErr))
	}
}

func TestStreamWorkloadLogs_PropagatesRecvError(t *testing.T) {
	client := &fakeRunnersClient{}
	client.streamWorkloadLogs = func(ctx context.Context, req *runnerv1.StreamWorkloadLogsRequest) (grpc.ServerStreamingClient[runnerv1.StreamWorkloadLogsResponse], error) {
		return &fakeStream{ctx: ctx, err: status.Error(codes.Unavailable, "runner unavailable")}, nil
	}

	gatewayClient := newRunnersGatewayClient(t, NewRunnersGateway(client))
	stream, err := gatewayClient.StreamWorkloadLogs(context.Background(), connect.NewRequest(&runnerv1.StreamWorkloadLogsRequest{
		WorkloadId:    "workload-1",
		ContainerName: "container-1",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stream.Receive() {
		t.Fatalf("expected no messages")
	}
	streamErr := stream.Err()
	if streamErr == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(streamErr) != connect.CodeUnavailable {
		t.Fatalf("expected code %v, got %v", connect.CodeUnavailable, connect.CodeOf(streamErr))
	}
}
