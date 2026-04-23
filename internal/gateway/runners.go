package gateway

import (
	"context"
	"errors"
	"io"

	"connectrpc.com/connect"
	runnerv1 "github.com/agynio/gateway/gen/agynio/api/runner/v1"
	runnersv1 "github.com/agynio/gateway/gen/agynio/api/runners/v1"
)

type RunnersGateway struct {
	runners runnersv1.RunnersServiceClient
	runner  runnerv1.RunnerServiceClient
}

func NewRunnersGateway(runners runnersv1.RunnersServiceClient, runner runnerv1.RunnerServiceClient) *RunnersGateway {
	return &RunnersGateway{runners: runners, runner: runner}
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
	grpcStream, err := g.runners.StreamWorkloadLogs(ctx, req.Msg)
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
