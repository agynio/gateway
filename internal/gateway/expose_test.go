package gateway

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	exposev1 "github.com/agynio/gateway/gen/agynio/api/expose/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeExposeClient struct {
	addExposureReq      *exposev1.AddExposureRequest
	addExposureResp     *exposev1.AddExposureResponse
	addExposureErr      error
	addExposureCalls    int
	removeExposureReq   *exposev1.RemoveExposureRequest
	removeExposureResp  *exposev1.RemoveExposureResponse
	removeExposureErr   error
	removeExposureCalls int
	listExposuresReq    *exposev1.ListExposuresRequest
	listExposuresResp   *exposev1.ListExposuresResponse
	listExposuresErr    error
	listExposuresCalls  int
}

func (f *fakeExposeClient) AddExposure(ctx context.Context, in *exposev1.AddExposureRequest, opts ...grpc.CallOption) (*exposev1.AddExposureResponse, error) {
	f.addExposureCalls++
	f.addExposureReq = in
	if f.addExposureErr != nil {
		return nil, f.addExposureErr
	}
	if f.addExposureResp == nil {
		f.addExposureResp = &exposev1.AddExposureResponse{}
	}
	return f.addExposureResp, nil
}

func (f *fakeExposeClient) RemoveExposure(ctx context.Context, in *exposev1.RemoveExposureRequest, opts ...grpc.CallOption) (*exposev1.RemoveExposureResponse, error) {
	f.removeExposureCalls++
	f.removeExposureReq = in
	if f.removeExposureErr != nil {
		return nil, f.removeExposureErr
	}
	if f.removeExposureResp == nil {
		f.removeExposureResp = &exposev1.RemoveExposureResponse{}
	}
	return f.removeExposureResp, nil
}

func (f *fakeExposeClient) ListExposures(ctx context.Context, in *exposev1.ListExposuresRequest, opts ...grpc.CallOption) (*exposev1.ListExposuresResponse, error) {
	f.listExposuresCalls++
	f.listExposuresReq = in
	if f.listExposuresErr != nil {
		return nil, f.listExposuresErr
	}
	if f.listExposuresResp == nil {
		f.listExposuresResp = &exposev1.ListExposuresResponse{}
	}
	return f.listExposuresResp, nil
}

func TestExposeGatewayAddExposureSuccess(t *testing.T) {
	client := &fakeExposeClient{
		addExposureResp: &exposev1.AddExposureResponse{
			Exposure: &exposev1.Exposure{Meta: &exposev1.EntityMeta{Id: "exposure-1"}, Port: 8080},
		},
	}
	gateway := NewExposeGateway(client)

	req := connect.NewRequest(&exposev1.AddExposureRequest{WorkloadId: "workload-1", Port: 8080})
	resp, err := gateway.AddExposure(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.addExposureCalls != 1 {
		t.Fatalf("expected add exposure to be called once, got %d", client.addExposureCalls)
	}
	if client.addExposureReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.addExposureResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestExposeGatewayAddExposureError(t *testing.T) {
	client := &fakeExposeClient{
		addExposureErr: status.Error(codes.InvalidArgument, "invalid"),
	}
	gateway := NewExposeGateway(client)

	req := connect.NewRequest(&exposev1.AddExposureRequest{WorkloadId: "workload-1"})
	resp, err := gateway.AddExposure(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeInvalidArgument {
		t.Fatalf("expected CodeInvalidArgument, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.addExposureCalls != 1 {
		t.Fatalf("expected add exposure to be called once, got %d", client.addExposureCalls)
	}
	if client.addExposureReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
}

func TestExposeGatewayRemoveExposureSuccess(t *testing.T) {
	client := &fakeExposeClient{
		removeExposureResp: &exposev1.RemoveExposureResponse{},
	}
	gateway := NewExposeGateway(client)

	req := connect.NewRequest(&exposev1.RemoveExposureRequest{WorkloadId: "workload-1", Port: 8080})
	resp, err := gateway.RemoveExposure(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.removeExposureCalls != 1 {
		t.Fatalf("expected remove exposure to be called once, got %d", client.removeExposureCalls)
	}
	if client.removeExposureReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.removeExposureResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestExposeGatewayRemoveExposureError(t *testing.T) {
	client := &fakeExposeClient{
		removeExposureErr: status.Error(codes.NotFound, "missing"),
	}
	gateway := NewExposeGateway(client)

	req := connect.NewRequest(&exposev1.RemoveExposureRequest{WorkloadId: "workload-1"})
	resp, err := gateway.RemoveExposure(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeNotFound {
		t.Fatalf("expected CodeNotFound, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.removeExposureCalls != 1 {
		t.Fatalf("expected remove exposure to be called once, got %d", client.removeExposureCalls)
	}
	if client.removeExposureReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
}

func TestExposeGatewayListExposuresSuccess(t *testing.T) {
	client := &fakeExposeClient{
		listExposuresResp: &exposev1.ListExposuresResponse{
			Exposures: []*exposev1.Exposure{{Meta: &exposev1.EntityMeta{Id: "exposure-1"}, Port: 8080}},
		},
	}
	gateway := NewExposeGateway(client)

	req := connect.NewRequest(&exposev1.ListExposuresRequest{WorkloadId: "workload-1", PageSize: 10})
	resp, err := gateway.ListExposures(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.listExposuresCalls != 1 {
		t.Fatalf("expected list exposures to be called once, got %d", client.listExposuresCalls)
	}
	if client.listExposuresReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
	if resp.Msg != client.listExposuresResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestExposeGatewayListExposuresError(t *testing.T) {
	client := &fakeExposeClient{
		listExposuresErr: status.Error(codes.PermissionDenied, "denied"),
	}
	gateway := NewExposeGateway(client)

	req := connect.NewRequest(&exposev1.ListExposuresRequest{WorkloadId: "workload-1"})
	resp, err := gateway.ListExposures(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected CodePermissionDenied, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.listExposuresCalls != 1 {
		t.Fatalf("expected list exposures to be called once, got %d", client.listExposuresCalls)
	}
	if client.listExposuresReq != req.Msg {
		t.Fatalf("expected request to be forwarded")
	}
}

func TestExposeGatewayAddExposureInjectsWorkloadID(t *testing.T) {
	client := &fakeExposeClient{}
	gateway := NewExposeGateway(client)
	ctx := identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   "agent-1",
		IdentityType: identity.IdentityTypeAgent,
		WorkloadID:   "workload-1",
	})

	req := connect.NewRequest(&exposev1.AddExposureRequest{Port: 8080})
	_, err := gateway.AddExposure(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.addExposureCalls != 1 {
		t.Fatalf("expected add exposure to be called once, got %d", client.addExposureCalls)
	}
	if client.addExposureReq.GetWorkloadId() != "workload-1" {
		t.Fatalf("expected workload id to be injected, got %q", client.addExposureReq.GetWorkloadId())
	}
}

func TestExposeGatewayAddExposureWorkloadMismatch(t *testing.T) {
	client := &fakeExposeClient{}
	gateway := NewExposeGateway(client)
	ctx := identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   "agent-1",
		IdentityType: identity.IdentityTypeAgent,
		WorkloadID:   "workload-1",
	})

	req := connect.NewRequest(&exposev1.AddExposureRequest{WorkloadId: "workload-2"})
	resp, err := gateway.AddExposure(ctx, req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected CodePermissionDenied, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.addExposureCalls != 0 {
		t.Fatalf("expected no add exposure calls, got %d", client.addExposureCalls)
	}
}

func TestExposeGatewayRemoveExposureInjectsWorkloadID(t *testing.T) {
	client := &fakeExposeClient{}
	gateway := NewExposeGateway(client)
	ctx := identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   "agent-1",
		IdentityType: identity.IdentityTypeAgent,
		WorkloadID:   "workload-3",
	})

	req := connect.NewRequest(&exposev1.RemoveExposureRequest{Port: 8080})
	_, err := gateway.RemoveExposure(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.removeExposureCalls != 1 {
		t.Fatalf("expected remove exposure to be called once, got %d", client.removeExposureCalls)
	}
	if client.removeExposureReq.GetWorkloadId() != "workload-3" {
		t.Fatalf("expected workload id to be injected, got %q", client.removeExposureReq.GetWorkloadId())
	}
}

func TestExposeGatewayListExposuresInjectsWorkloadID(t *testing.T) {
	client := &fakeExposeClient{}
	gateway := NewExposeGateway(client)
	ctx := identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   "agent-1",
		IdentityType: identity.IdentityTypeAgent,
		WorkloadID:   "workload-4",
	})

	req := connect.NewRequest(&exposev1.ListExposuresRequest{PageSize: 10})
	_, err := gateway.ListExposures(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.listExposuresCalls != 1 {
		t.Fatalf("expected list exposures to be called once, got %d", client.listExposuresCalls)
	}
	if client.listExposuresReq.GetWorkloadId() != "workload-4" {
		t.Fatalf("expected workload id to be injected, got %q", client.listExposuresReq.GetWorkloadId())
	}
}
