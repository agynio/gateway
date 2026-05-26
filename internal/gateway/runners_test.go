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

type fakeRunnersClient struct {
	getWorkload             func(ctx context.Context, req *runnersv1.GetWorkloadRequest) (*runnersv1.GetWorkloadResponse, error)
	getWorkloadReq          *runnersv1.GetWorkloadRequest
	getWorkloadMetadata     metadata.MD
	getWorkloadCalls        int
	streamWorkloadLogs      func(ctx context.Context, req *runnerv1.StreamWorkloadLogsRequest) (grpc.ServerStreamingClient[runnerv1.StreamWorkloadLogsResponse], error)
	streamWorkloadLogsReq   *runnerv1.StreamWorkloadLogsRequest
	streamWorkloadLogsCalls int
}

func (f *fakeRunnersClient) RegisterRunner(ctx context.Context, in *runnersv1.RegisterRunnerRequest, opts ...grpc.CallOption) (*runnersv1.RegisterRunnerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "RegisterRunner not implemented")
}

func (f *fakeRunnersClient) GetRunner(ctx context.Context, in *runnersv1.GetRunnerRequest, opts ...grpc.CallOption) (*runnersv1.GetRunnerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetRunner not implemented")
}

func (f *fakeRunnersClient) ListRunners(ctx context.Context, in *runnersv1.ListRunnersRequest, opts ...grpc.CallOption) (*runnersv1.ListRunnersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListRunners not implemented")
}

func (f *fakeRunnersClient) UpdateRunner(ctx context.Context, in *runnersv1.UpdateRunnerRequest, opts ...grpc.CallOption) (*runnersv1.UpdateRunnerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "UpdateRunner not implemented")
}

func (f *fakeRunnersClient) DeleteRunner(ctx context.Context, in *runnersv1.DeleteRunnerRequest, opts ...grpc.CallOption) (*runnersv1.DeleteRunnerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteRunner not implemented")
}

func (f *fakeRunnersClient) ValidateServiceToken(ctx context.Context, in *runnersv1.ValidateServiceTokenRequest, opts ...grpc.CallOption) (*runnersv1.ValidateServiceTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ValidateServiceToken not implemented")
}

func (f *fakeRunnersClient) EnrollRunner(ctx context.Context, in *runnersv1.EnrollRunnerRequest, opts ...grpc.CallOption) (*runnersv1.EnrollRunnerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "EnrollRunner not implemented")
}

func (f *fakeRunnersClient) CreateWorkload(ctx context.Context, in *runnersv1.CreateWorkloadRequest, opts ...grpc.CallOption) (*runnersv1.CreateWorkloadResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateWorkload not implemented")
}

func (f *fakeRunnersClient) UpdateWorkload(ctx context.Context, in *runnersv1.UpdateWorkloadRequest, opts ...grpc.CallOption) (*runnersv1.UpdateWorkloadResponse, error) {
	return nil, status.Error(codes.Unimplemented, "UpdateWorkload not implemented")
}

func (f *fakeRunnersClient) UpdateWorkloadStatus(ctx context.Context, in *runnersv1.UpdateWorkloadStatusRequest, opts ...grpc.CallOption) (*runnersv1.UpdateWorkloadStatusResponse, error) {
	return nil, status.Error(codes.Unimplemented, "UpdateWorkloadStatus not implemented")
}

func (f *fakeRunnersClient) TouchWorkload(ctx context.Context, in *runnersv1.TouchWorkloadRequest, opts ...grpc.CallOption) (*runnersv1.TouchWorkloadResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TouchWorkload not implemented")
}

func (f *fakeRunnersClient) DeleteWorkload(ctx context.Context, in *runnersv1.DeleteWorkloadRequest, opts ...grpc.CallOption) (*runnersv1.DeleteWorkloadResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteWorkload not implemented")
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

func (f *fakeRunnersClient) ListWorkloadsByThread(ctx context.Context, in *runnersv1.ListWorkloadsByThreadRequest, opts ...grpc.CallOption) (*runnersv1.ListWorkloadsByThreadResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListWorkloadsByThread not implemented")
}

func (f *fakeRunnersClient) ListWorkloads(ctx context.Context, in *runnersv1.ListWorkloadsRequest, opts ...grpc.CallOption) (*runnersv1.ListWorkloadsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListWorkloads not implemented")
}

func (f *fakeRunnersClient) BatchUpdateWorkloadSampledAt(ctx context.Context, in *runnersv1.BatchUpdateWorkloadSampledAtRequest, opts ...grpc.CallOption) (*runnersv1.BatchUpdateWorkloadSampledAtResponse, error) {
	return nil, status.Error(codes.Unimplemented, "BatchUpdateWorkloadSampledAt not implemented")
}

func (f *fakeRunnersClient) StreamWorkloadLogs(ctx context.Context, in *runnerv1.StreamWorkloadLogsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[runnerv1.StreamWorkloadLogsResponse], error) {
	f.streamWorkloadLogsCalls++
	f.streamWorkloadLogsReq = in
	if f.streamWorkloadLogs == nil {
		return nil, status.Error(codes.Unimplemented, "StreamWorkloadLogs not implemented")
	}
	return f.streamWorkloadLogs(ctx, in)
}

func (f *fakeRunnersClient) CreateVolume(ctx context.Context, in *runnersv1.CreateVolumeRequest, opts ...grpc.CallOption) (*runnersv1.CreateVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateVolume not implemented")
}

func (f *fakeRunnersClient) UpdateVolume(ctx context.Context, in *runnersv1.UpdateVolumeRequest, opts ...grpc.CallOption) (*runnersv1.UpdateVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "UpdateVolume not implemented")
}

func (f *fakeRunnersClient) GetVolume(ctx context.Context, in *runnersv1.GetVolumeRequest, opts ...grpc.CallOption) (*runnersv1.GetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetVolume not implemented")
}

func (f *fakeRunnersClient) ListVolumes(ctx context.Context, in *runnersv1.ListVolumesRequest, opts ...grpc.CallOption) (*runnersv1.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListVolumes not implemented")
}

func (f *fakeRunnersClient) ListVolumesByThread(ctx context.Context, in *runnersv1.ListVolumesByThreadRequest, opts ...grpc.CallOption) (*runnersv1.ListVolumesByThreadResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListVolumesByThread not implemented")
}

func (f *fakeRunnersClient) BatchUpdateVolumeSampledAt(ctx context.Context, in *runnersv1.BatchUpdateVolumeSampledAtRequest, opts ...grpc.CallOption) (*runnersv1.BatchUpdateVolumeSampledAtResponse, error) {
	return nil, status.Error(codes.Unimplemented, "BatchUpdateVolumeSampledAt not implemented")
}

func newRunnersGatewayClient(t *testing.T, gateway *RunnersGateway) gatewayv1connect.RunnersGatewayClient {
	t.Helper()
	mux := http.NewServeMux()
	handlerPath, handler := gatewayv1connect.NewRunnersGatewayHandler(gateway)
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
