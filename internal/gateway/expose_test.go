package gateway

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	exposev1 "github.com/agynio/gateway/gen/agynio/api/expose/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type fakeExposeClient struct {
	addExposureReq      *exposev1.AddExposureRequest
	addExposureResp     *exposev1.AddExposureResponse
	addExposureErr      error
	addExposureCalls    int
	addExposureMetadata metadata.MD
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
	f.addExposureMetadata, _ = metadata.FromOutgoingContext(ctx)
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

func TestExposeGatewayAddExposureForwardsRequest(t *testing.T) {
	client := &fakeExposeClient{addExposureResp: &exposev1.AddExposureResponse{}}
	gateway := NewExposeGateway(client)

	req := connect.NewRequest(&exposev1.AddExposureRequest{WorkloadId: "workload-1", AgentId: "agent-1", Port: 8080})
	resp, err := gateway.AddExposure(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.addExposureCalls != 1 {
		t.Fatalf("expected add exposure call, got %d", client.addExposureCalls)
	}
	if client.addExposureReq.GetWorkloadId() != "workload-1" {
		t.Fatalf("expected workload id to be forwarded")
	}
	if client.addExposureReq.GetAgentId() != "agent-1" {
		t.Fatalf("expected agent id to be forwarded")
	}
	if client.addExposureReq.GetPort() != 8080 {
		t.Fatalf("expected port to be forwarded")
	}
	if resp.Msg != client.addExposureResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestExposeGatewayAddExposurePropagatesIdentityMetadata(t *testing.T) {
	client := &fakeExposeClient{addExposureResp: &exposev1.AddExposureResponse{}}
	gateway := NewExposeGateway(client)
	resolved := identity.ResolvedIdentity{
		IdentityID:   "agent-1",
		IdentityType: identity.IdentityTypeAgent,
		WorkloadID:   "workload-1",
	}
	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		identity.MetadataKeyIdentityID, "stale-identity",
		identity.MetadataKeyIdentityType, string(identity.IdentityTypeUser),
		identity.MetadataKeyWorkloadID, "stale-workload",
	)
	ctx = identity.WithIdentity(ctx, resolved)

	_, err := gateway.AddExposure(ctx, connect.NewRequest(&exposev1.AddExposureRequest{Port: 8080}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertMetadataValue(t, client.addExposureMetadata, identity.MetadataKeyIdentityID, resolved.IdentityID)
	assertMetadataValue(t, client.addExposureMetadata, identity.MetadataKeyIdentityType, string(resolved.IdentityType))
	assertMetadataValue(t, client.addExposureMetadata, identity.MetadataKeyWorkloadID, resolved.WorkloadID)
}

func TestExposeGatewayRemoveExposureForwardsRequest(t *testing.T) {
	client := &fakeExposeClient{removeExposureResp: &exposev1.RemoveExposureResponse{}}
	gateway := NewExposeGateway(client)

	req := connect.NewRequest(&exposev1.RemoveExposureRequest{WorkloadId: "workload-2", Port: 9090})
	resp, err := gateway.RemoveExposure(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.removeExposureCalls != 1 {
		t.Fatalf("expected remove exposure call, got %d", client.removeExposureCalls)
	}
	if client.removeExposureReq.GetWorkloadId() != "workload-2" {
		t.Fatalf("expected workload id to be forwarded")
	}
	if client.removeExposureReq.GetPort() != 9090 {
		t.Fatalf("expected port to be forwarded")
	}
	if resp.Msg != client.removeExposureResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestExposeGatewayListExposuresForwardsRequest(t *testing.T) {
	client := &fakeExposeClient{listExposuresResp: &exposev1.ListExposuresResponse{}}
	gateway := NewExposeGateway(client)

	req := connect.NewRequest(&exposev1.ListExposuresRequest{WorkloadId: "workload-3", PageSize: 10, PageToken: "next"})
	resp, err := gateway.ListExposures(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.listExposuresCalls != 1 {
		t.Fatalf("expected list exposures call, got %d", client.listExposuresCalls)
	}
	if client.listExposuresReq.GetWorkloadId() != "workload-3" {
		t.Fatalf("expected workload id to be forwarded")
	}
	if client.listExposuresReq.GetPageSize() != 10 {
		t.Fatalf("expected page size to be forwarded")
	}
	if client.listExposuresReq.GetPageToken() != "next" {
		t.Fatalf("expected page token to be forwarded")
	}
	if resp.Msg != client.listExposuresResp {
		t.Fatalf("expected response to be forwarded")
	}
}
