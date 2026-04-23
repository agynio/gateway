package gateway

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	gatewayv1connect "github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	runnerv1 "github.com/agynio/gateway/gen/agynio/api/runner/v1"
	runnersv1 "github.com/agynio/gateway/gen/agynio/api/runners/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakeRunnersClient struct {
	getWorkloadReq   *runnersv1.GetWorkloadRequest
	getWorkloadResp  *runnersv1.GetWorkloadResponse
	getWorkloadErr   error
	getWorkloadCalls int
	getRunnerReq     *runnersv1.GetRunnerRequest
	getRunnerResp    *runnersv1.GetRunnerResponse
	getRunnerErr     error
	getRunnerCalls   int
}

func (f *fakeRunnersClient) RegisterRunner(ctx context.Context, in *runnersv1.RegisterRunnerRequest, opts ...grpc.CallOption) (*runnersv1.RegisterRunnerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "RegisterRunner not implemented")
}

func (f *fakeRunnersClient) GetRunner(ctx context.Context, in *runnersv1.GetRunnerRequest, opts ...grpc.CallOption) (*runnersv1.GetRunnerResponse, error) {
	f.getRunnerCalls++
	f.getRunnerReq = in
	if f.getRunnerErr != nil {
		return nil, f.getRunnerErr
	}
	if f.getRunnerResp == nil {
		f.getRunnerResp = &runnersv1.GetRunnerResponse{}
	}
	return f.getRunnerResp, nil
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
	if f.getWorkloadErr != nil {
		return nil, f.getWorkloadErr
	}
	if f.getWorkloadResp == nil {
		f.getWorkloadResp = &runnersv1.GetWorkloadResponse{}
	}
	return f.getWorkloadResp, nil
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

type fakeRunnerServer struct {
	runnerv1.UnimplementedRunnerServiceServer
	streamWorkloadLogs func(*runnerv1.StreamWorkloadLogsRequest, grpc.ServerStreamingServer[runnerv1.StreamWorkloadLogsResponse]) error
}

func (f *fakeRunnerServer) StreamWorkloadLogs(req *runnerv1.StreamWorkloadLogsRequest, stream grpc.ServerStreamingServer[runnerv1.StreamWorkloadLogsResponse]) error {
	if f.streamWorkloadLogs == nil {
		return status.Error(codes.Unimplemented, "StreamWorkloadLogs not implemented")
	}
	return f.streamWorkloadLogs(req, stream)
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

func startRunnerServer(t *testing.T, streamWorkloadLogs func(*runnerv1.StreamWorkloadLogsRequest, grpc.ServerStreamingServer[runnerv1.StreamWorkloadLogsResponse]) error) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	runnerv1.RegisterRunnerServiceServer(server, &fakeRunnerServer{streamWorkloadLogs: streamWorkloadLogs})
	go func() {
		_ = server.Serve(listener)
	}()

	t.Cleanup(func() {
		server.Stop()
		_ = listener.Close()
	})

	return listener.Addr().String()
}

func TestStreamWorkloadLogs_Success(t *testing.T) {
	workloadID := "workload-1"
	instanceID := "instance-1"
	runnerID := "runner-1"
	serviceName := "runner-service"
	containerName := "container-1"
	sinceTime := timestamppb.New(time.Unix(1723123123, 0))
	requestChan := make(chan *runnerv1.StreamWorkloadLogsRequest, 1)

	addr := startRunnerServer(t, func(req *runnerv1.StreamWorkloadLogsRequest, stream grpc.ServerStreamingServer[runnerv1.StreamWorkloadLogsResponse]) error {
		requestChan <- req
		return stream.Send(&runnerv1.StreamWorkloadLogsResponse{
			Event: &runnerv1.StreamWorkloadLogsResponse_Chunk{
				Chunk: &runnerv1.LogChunk{Data: []byte("hello")},
			},
		})
	})

	client := &fakeRunnersClient{
		getWorkloadResp: &runnersv1.GetWorkloadResponse{
			Workload: &runnersv1.Workload{
				RunnerId:   runnerID,
				InstanceId: &instanceID,
			},
		},
		getRunnerResp: &runnersv1.GetRunnerResponse{
			Runner: &runnersv1.Runner{OpenzitiServiceName: serviceName},
		},
	}

	gateway := NewRunnersGateway(client)
	gateway.requiresZitiContext = false
	gateway.dialRunner = func(ctx context.Context, name string) (net.Conn, error) {
		if name != serviceName {
			return nil, status.Errorf(codes.InvalidArgument, "unexpected service name %q", name)
		}
		return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	}

	gatewayClient := newRunnersGatewayClient(t, gateway)
	stream, err := gatewayClient.StreamWorkloadLogs(context.Background(), connect.NewRequest(&runnerv1.StreamWorkloadLogsRequest{
		WorkloadId:    " " + workloadID + " ",
		ContainerName: " " + containerName + " ",
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

	select {
	case gotReq := <-requestChan:
		if gotReq.GetWorkloadId() != instanceID {
			t.Fatalf("expected instance id %q, got %q", instanceID, gotReq.GetWorkloadId())
		}
		if gotReq.GetContainerName() != containerName {
			t.Fatalf("expected container name %q, got %q", containerName, gotReq.GetContainerName())
		}
		if gotReq.GetTailLines() != 3 {
			t.Fatalf("expected tail lines 3, got %d", gotReq.GetTailLines())
		}
		if gotReq.GetFollow() != true {
			t.Fatalf("expected follow true")
		}
		if gotReq.GetSinceTime() == nil || gotReq.GetSinceTime().AsTime().Unix() != sinceTime.AsTime().Unix() {
			t.Fatalf("expected since time %v, got %v", sinceTime, gotReq.GetSinceTime())
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for runner request")
	}

	if client.getWorkloadCalls != 1 {
		t.Fatalf("expected GetWorkload to be called once, got %d", client.getWorkloadCalls)
	}
	if client.getWorkloadReq.GetId() != workloadID {
		t.Fatalf("expected GetWorkload id %q, got %q", workloadID, client.getWorkloadReq.GetId())
	}
	if client.getRunnerCalls != 1 {
		t.Fatalf("expected GetRunner to be called once, got %d", client.getRunnerCalls)
	}
	if client.getRunnerReq.GetId() != runnerID {
		t.Fatalf("expected GetRunner id %q, got %q", runnerID, client.getRunnerReq.GetId())
	}
}

func TestStreamWorkloadLogs_PropagatesGetWorkloadErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code connect.Code
	}{
		{
			name: "not found",
			err:  status.Error(codes.NotFound, "missing workload"),
			code: connect.CodeNotFound,
		},
		{
			name: "permission denied",
			err:  status.Error(codes.PermissionDenied, "denied"),
			code: connect.CodePermissionDenied,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &fakeRunnersClient{getWorkloadErr: test.err}
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
			if connect.CodeOf(streamErr) != test.code {
				t.Fatalf("expected code %v, got %v", test.code, connect.CodeOf(streamErr))
			}
			if client.getRunnerCalls != 0 {
				t.Fatalf("expected GetRunner not to be called")
			}
		})
	}
}

func TestStreamWorkloadLogs_MissingInstanceID(t *testing.T) {
	runnerID := "runner-1"
	client := &fakeRunnersClient{
		getWorkloadResp: &runnersv1.GetWorkloadResponse{
			Workload: &runnersv1.Workload{RunnerId: runnerID},
		},
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
	if connect.CodeOf(streamErr) != connect.CodeFailedPrecondition {
		t.Fatalf("expected failed precondition, got %v", connect.CodeOf(streamErr))
	}
	if client.getRunnerCalls != 0 {
		t.Fatalf("expected GetRunner not to be called")
	}
}

func TestStreamWorkloadLogs_ZitiUnavailable(t *testing.T) {
	instanceID := "instance-1"
	runnerID := "runner-1"
	client := &fakeRunnersClient{
		getWorkloadResp: &runnersv1.GetWorkloadResponse{
			Workload: &runnersv1.Workload{
				RunnerId:   runnerID,
				InstanceId: &instanceID,
			},
		},
		getRunnerResp: &runnersv1.GetRunnerResponse{
			Runner: &runnersv1.Runner{OpenzitiServiceName: "runner-service"},
		},
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
		t.Fatalf("expected unavailable, got %v", connect.CodeOf(streamErr))
	}
	if client.getWorkloadCalls != 1 {
		t.Fatalf("expected GetWorkload to be called once, got %d", client.getWorkloadCalls)
	}
	if client.getRunnerCalls != 1 {
		t.Fatalf("expected GetRunner to be called once, got %d", client.getRunnerCalls)
	}
}
