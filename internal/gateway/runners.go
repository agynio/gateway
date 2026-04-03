package gateway

import (
	"context"

	"connectrpc.com/connect"
	runnersv1 "github.com/agynio/gateway/gen/agynio/api/runners/v1"
)

type RunnersGateway struct {
	runners runnersv1.RunnersServiceClient
}

func NewRunnersGateway(runners runnersv1.RunnersServiceClient) *RunnersGateway {
	return &RunnersGateway{runners: runners}
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

func (g *RunnersGateway) GetWorkload(ctx context.Context, req *connect.Request[runnersv1.GetWorkloadRequest]) (*connect.Response[runnersv1.GetWorkloadResponse], error) {
	resp, err := g.runners.GetWorkload(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
