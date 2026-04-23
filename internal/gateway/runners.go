package gateway

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync"

	"connectrpc.com/connect"
	runnerv1 "github.com/agynio/gateway/gen/agynio/api/runner/v1"
	runnersv1 "github.com/agynio/gateway/gen/agynio/api/runners/v1"
	"github.com/agynio/gateway/internal/grpcclient"
	"github.com/openziti/sdk-golang/ziti"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var errZitiContextUnavailable = errors.New("ziti context unavailable")

type runnerDialer func(ctx context.Context, serviceName string) (net.Conn, error)

type RunnersGateway struct {
	runners             runnersv1.RunnersServiceClient
	zitiMu              sync.RWMutex
	zitiContextProvider ZitiContextProvider
	dialRunner          runnerDialer
	requiresZitiContext bool
}

func NewRunnersGateway(runners runnersv1.RunnersServiceClient) *RunnersGateway {
	gateway := &RunnersGateway{
		runners:             runners,
		requiresZitiContext: true,
	}
	gateway.dialRunner = gateway.dialWithZiti
	return gateway
}

func (g *RunnersGateway) SetZitiContextProvider(provider ZitiContextProvider) {
	g.zitiMu.Lock()
	g.zitiContextProvider = provider
	g.zitiMu.Unlock()
}

func (g *RunnersGateway) zitiContext() ziti.Context {
	g.zitiMu.RLock()
	provider := g.zitiContextProvider
	g.zitiMu.RUnlock()
	if provider == nil {
		return nil
	}
	return provider.ZitiContext()
}

func (g *RunnersGateway) dialWithZiti(ctx context.Context, serviceName string) (net.Conn, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	zitiCtx := g.zitiContext()
	if zitiCtx == nil {
		return nil, errZitiContextUnavailable
	}

	return zitiCtx.DialContext(ctx, serviceName)
}

func (g *RunnersGateway) dialRunnerConn(ctx context.Context, serviceName string) (*grpc.ClientConn, error) {
	if g.dialRunner == nil {
		return nil, errZitiContextUnavailable
	}

	target := "passthrough:///" + serviceName
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return g.dialRunner(ctx, serviceName)
		}),
		grpc.WithChainUnaryInterceptor(grpcclient.IdentityUnaryInterceptor),
		grpc.WithChainStreamInterceptor(grpcclient.IdentityStreamInterceptor),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (g *RunnersGateway) RegisterRunner(ctx context.Context, req *connect.Request[runnersv1.RegisterRunnerRequest]) (*connect.Response[runnersv1.RegisterRunnerResponse], error) {
	resp, err := g.runners.RegisterRunner(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) EnrollRunner(ctx context.Context, req *connect.Request[runnersv1.EnrollRunnerRequest]) (*connect.Response[runnersv1.EnrollRunnerResponse], error) {
	resp, err := g.runners.EnrollRunner(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) GetRunner(ctx context.Context, req *connect.Request[runnersv1.GetRunnerRequest]) (*connect.Response[runnersv1.GetRunnerResponse], error) {
	resp, err := g.runners.GetRunner(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) ListRunners(ctx context.Context, req *connect.Request[runnersv1.ListRunnersRequest]) (*connect.Response[runnersv1.ListRunnersResponse], error) {
	resp, err := g.runners.ListRunners(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) UpdateRunner(ctx context.Context, req *connect.Request[runnersv1.UpdateRunnerRequest]) (*connect.Response[runnersv1.UpdateRunnerResponse], error) {
	resp, err := g.runners.UpdateRunner(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) DeleteRunner(ctx context.Context, req *connect.Request[runnersv1.DeleteRunnerRequest]) (*connect.Response[runnersv1.DeleteRunnerResponse], error) {
	resp, err := g.runners.DeleteRunner(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) ListWorkloadsByThread(ctx context.Context, req *connect.Request[runnersv1.ListWorkloadsByThreadRequest]) (*connect.Response[runnersv1.ListWorkloadsByThreadResponse], error) {
	resp, err := g.runners.ListWorkloadsByThread(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) ListWorkloads(ctx context.Context, req *connect.Request[runnersv1.ListWorkloadsRequest]) (*connect.Response[runnersv1.ListWorkloadsResponse], error) {
	resp, err := g.runners.ListWorkloads(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) GetWorkload(ctx context.Context, req *connect.Request[runnersv1.GetWorkloadRequest]) (*connect.Response[runnersv1.GetWorkloadResponse], error) {
	resp, err := g.runners.GetWorkload(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) TouchWorkload(ctx context.Context, req *connect.Request[runnersv1.TouchWorkloadRequest]) (*connect.Response[runnersv1.TouchWorkloadResponse], error) {
	resp, err := g.runners.TouchWorkload(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) StreamWorkloadLogs(ctx context.Context, req *connect.Request[runnerv1.StreamWorkloadLogsRequest], stream *connect.ServerStream[runnerv1.StreamWorkloadLogsResponse]) error {
	workloadID := strings.TrimSpace(req.Msg.GetWorkloadId())
	if workloadID == "" {
		return toConnectError(status.Error(codes.InvalidArgument, "workload_id is required"))
	}

	containerName := strings.TrimSpace(req.Msg.GetContainerName())
	if containerName == "" {
		return toConnectError(status.Error(codes.InvalidArgument, "container_name is required"))
	}

	workloadResp, err := g.runners.GetWorkload(ctx, &runnersv1.GetWorkloadRequest{Id: workloadID})
	if err != nil {
		return toConnectError(err)
	}
	workload := workloadResp.GetWorkload()
	if workload == nil {
		return toConnectError(status.Error(codes.Internal, "workload missing from response"))
	}

	instanceID := strings.TrimSpace(workload.GetInstanceId())
	if instanceID == "" {
		return toConnectError(status.Error(codes.FailedPrecondition, "workload instance_id is required"))
	}

	runnerID := strings.TrimSpace(workload.GetRunnerId())
	if runnerID == "" {
		return toConnectError(status.Error(codes.Internal, "workload runner_id is required"))
	}

	runnerResp, err := g.runners.GetRunner(ctx, &runnersv1.GetRunnerRequest{Id: runnerID})
	if err != nil {
		return toConnectError(err)
	}
	runner := runnerResp.GetRunner()
	if runner == nil {
		return toConnectError(status.Error(codes.Internal, "runner missing from response"))
	}
	serviceName := strings.TrimSpace(runner.GetOpenzitiServiceName())
	if serviceName == "" {
		return toConnectError(status.Error(codes.Internal, "runner ziti service name is required"))
	}

	if g.requiresZitiContext && g.zitiContext() == nil {
		return toConnectError(status.Error(codes.Unavailable, errZitiContextUnavailable.Error()))
	}

	conn, err := g.dialRunnerConn(ctx, serviceName)
	if err != nil {
		if errors.Is(err, errZitiContextUnavailable) {
			return toConnectError(status.Error(codes.Unavailable, err.Error()))
		}
		if errors.Is(err, context.Canceled) {
			return toConnectError(status.Error(codes.Canceled, err.Error()))
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return toConnectError(status.Error(codes.DeadlineExceeded, err.Error()))
		}
		return toConnectError(status.Errorf(codes.Unavailable, "dial runner: %v", err))
	}
	defer conn.Close()

	runnerReq := &runnerv1.StreamWorkloadLogsRequest{
		WorkloadId:    instanceID,
		Follow:        req.Msg.GetFollow(),
		Since:         req.Msg.GetSince(),
		Tail:          req.Msg.GetTail(),
		Stdout:        req.Msg.GetStdout(),
		Stderr:        req.Msg.GetStderr(),
		Timestamps:    req.Msg.GetTimestamps(),
		ContainerName: containerName,
		TailLines:     req.Msg.GetTailLines(),
		SinceTime:     req.Msg.GetSinceTime(),
	}

	runnerClient := runnerv1.NewRunnerServiceClient(conn)
	grpcStream, err := runnerClient.StreamWorkloadLogs(ctx, runnerReq)
	if err != nil {
		return toConnectError(err)
	}

	for {
		msg, err := grpcStream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return toConnectError(err)
		}
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
}

func (g *RunnersGateway) GetVolume(ctx context.Context, req *connect.Request[runnersv1.GetVolumeRequest]) (*connect.Response[runnersv1.GetVolumeResponse], error) {
	resp, err := g.runners.GetVolume(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) ListVolumes(ctx context.Context, req *connect.Request[runnersv1.ListVolumesRequest]) (*connect.Response[runnersv1.ListVolumesResponse], error) {
	resp, err := g.runners.ListVolumes(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *RunnersGateway) ListVolumesByThread(ctx context.Context, req *connect.Request[runnersv1.ListVolumesByThreadRequest]) (*connect.Response[runnersv1.ListVolumesByThreadResponse], error) {
	resp, err := g.runners.ListVolumesByThread(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
