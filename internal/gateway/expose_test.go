package gateway

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	exposev1 "github.com/agynio/gateway/gen/agynio/api/expose/v1"
	"github.com/agynio/gateway/internal/clusteradminresolver"
	"github.com/agynio/gateway/internal/httpauth"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	agentID      = "agent-1"
	workloadID   = "workload-1"
	clusterToken = "cluster-token"
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

func agentContext() context.Context {
	return identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   agentID,
		IdentityType: identity.IdentityTypeAgent,
		WorkloadID:   workloadID,
	})
}

func userContext() context.Context {
	return identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   "user-1",
		IdentityType: identity.IdentityTypeUser,
		WorkloadID:   workloadID,
	})
}

func newClusterGateway(t *testing.T, client *fakeExposeClient) (*ExposeGateway, context.Context) {
	resolver, err := clusteradminresolver.NewResolver(clusterToken, "cluster-admin")
	if err != nil {
		t.Fatalf("unexpected resolver error: %v", err)
	}
	ctx := httpauth.WithBearerToken(context.Background(), clusterToken)
	return NewExposeGateway(client, resolver), ctx
}

func TestExposeGatewayAddExposureDefaultsFromAgentIdentity(t *testing.T) {
	client := &fakeExposeClient{
		addExposureResp: &exposev1.AddExposureResponse{
			Exposure: &exposev1.Exposure{Meta: &exposev1.EntityMeta{Id: "exposure-1"}, Port: 8080},
		},
	}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.AddExposureRequest{Port: 8080})
	resp, err := gateway.AddExposure(agentContext(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.addExposureCalls != 1 {
		t.Fatalf("expected add exposure to be called once, got %d", client.addExposureCalls)
	}
	if client.addExposureReq.GetWorkloadId() != workloadID {
		t.Fatalf("expected workload id %q, got %q", workloadID, client.addExposureReq.GetWorkloadId())
	}
	if client.addExposureReq.GetAgentId() != agentID {
		t.Fatalf("expected agent id %q, got %q", agentID, client.addExposureReq.GetAgentId())
	}
	if client.addExposureReq.GetPort() != 8080 {
		t.Fatalf("expected port 8080, got %d", client.addExposureReq.GetPort())
	}
	if resp.Msg != client.addExposureResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestExposeGatewayAddExposureRejectsWorkloadMismatch(t *testing.T) {
	client := &fakeExposeClient{}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.AddExposureRequest{WorkloadId: "workload-2", Port: 8080})
	resp, err := gateway.AddExposure(agentContext(), req)
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

func TestExposeGatewayAddExposureRejectsAgentMismatch(t *testing.T) {
	client := &fakeExposeClient{}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.AddExposureRequest{AgentId: "agent-2", Port: 8080})
	resp, err := gateway.AddExposure(agentContext(), req)
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

func TestExposeGatewayAddExposureRequiresIdentity(t *testing.T) {
	client := &fakeExposeClient{}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.AddExposureRequest{Port: 8080})
	resp, err := gateway.AddExposure(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Fatalf("expected CodeUnauthenticated, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.addExposureCalls != 0 {
		t.Fatalf("expected no add exposure calls, got %d", client.addExposureCalls)
	}
}

func TestExposeGatewayAddExposureRejectsNonAgentIdentity(t *testing.T) {
	client := &fakeExposeClient{}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.AddExposureRequest{Port: 8080})
	resp, err := gateway.AddExposure(userContext(), req)
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

func TestExposeGatewayAddExposureClusterAdminExplicitIDs(t *testing.T) {
	client := &fakeExposeClient{
		addExposureResp: &exposev1.AddExposureResponse{},
	}
	gateway, ctx := newClusterGateway(t, client)

	req := connect.NewRequest(&exposev1.AddExposureRequest{WorkloadId: "workload-9", AgentId: "agent-9", Port: 8080})
	resp, err := gateway.AddExposure(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.addExposureCalls != 1 {
		t.Fatalf("expected add exposure to be called once, got %d", client.addExposureCalls)
	}
	if client.addExposureReq.GetWorkloadId() != "workload-9" {
		t.Fatalf("expected workload id %q, got %q", "workload-9", client.addExposureReq.GetWorkloadId())
	}
	if client.addExposureReq.GetAgentId() != "agent-9" {
		t.Fatalf("expected agent id %q, got %q", "agent-9", client.addExposureReq.GetAgentId())
	}
}

func TestExposeGatewayAddExposureClusterAdminRequiresWorkloadID(t *testing.T) {
	client := &fakeExposeClient{}
	gateway, ctx := newClusterGateway(t, client)

	req := connect.NewRequest(&exposev1.AddExposureRequest{AgentId: "agent-9", Port: 8080})
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

func TestExposeGatewayAddExposureClusterAdminRequiresAgentID(t *testing.T) {
	client := &fakeExposeClient{}
	gateway, ctx := newClusterGateway(t, client)

	req := connect.NewRequest(&exposev1.AddExposureRequest{WorkloadId: "workload-9", Port: 8080})
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

func TestExposeGatewayAddExposureError(t *testing.T) {
	client := &fakeExposeClient{
		addExposureErr: status.Error(codes.InvalidArgument, "invalid"),
	}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.AddExposureRequest{Port: 8080})
	resp, err := gateway.AddExposure(agentContext(), req)
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
	if client.addExposureReq.GetWorkloadId() != workloadID {
		t.Fatalf("expected workload id %q, got %q", workloadID, client.addExposureReq.GetWorkloadId())
	}
	if client.addExposureReq.GetAgentId() != agentID {
		t.Fatalf("expected agent id %q, got %q", agentID, client.addExposureReq.GetAgentId())
	}
}

func TestExposeGatewayRemoveExposureDefaultsFromAgentIdentity(t *testing.T) {
	client := &fakeExposeClient{
		removeExposureResp: &exposev1.RemoveExposureResponse{},
	}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.RemoveExposureRequest{Port: 8080})
	resp, err := gateway.RemoveExposure(agentContext(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.removeExposureCalls != 1 {
		t.Fatalf("expected remove exposure to be called once, got %d", client.removeExposureCalls)
	}
	if client.removeExposureReq.GetWorkloadId() != workloadID {
		t.Fatalf("expected workload id %q, got %q", workloadID, client.removeExposureReq.GetWorkloadId())
	}
	if client.removeExposureReq.GetPort() != 8080 {
		t.Fatalf("expected port 8080, got %d", client.removeExposureReq.GetPort())
	}
	if resp.Msg != client.removeExposureResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestExposeGatewayRemoveExposureRejectsWorkloadMismatch(t *testing.T) {
	client := &fakeExposeClient{}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.RemoveExposureRequest{WorkloadId: "workload-9", Port: 8080})
	resp, err := gateway.RemoveExposure(agentContext(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected CodePermissionDenied, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.removeExposureCalls != 0 {
		t.Fatalf("expected no remove exposure calls, got %d", client.removeExposureCalls)
	}
}

func TestExposeGatewayRemoveExposureClusterAdminExplicitIDs(t *testing.T) {
	client := &fakeExposeClient{
		removeExposureResp: &exposev1.RemoveExposureResponse{},
	}
	gateway, ctx := newClusterGateway(t, client)

	req := connect.NewRequest(&exposev1.RemoveExposureRequest{WorkloadId: "workload-9", Port: 8080})
	resp, err := gateway.RemoveExposure(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.removeExposureCalls != 1 {
		t.Fatalf("expected remove exposure to be called once, got %d", client.removeExposureCalls)
	}
	if client.removeExposureReq.GetWorkloadId() != "workload-9" {
		t.Fatalf("expected workload id %q, got %q", "workload-9", client.removeExposureReq.GetWorkloadId())
	}
}

func TestExposeGatewayRemoveExposureClusterAdminRequiresWorkloadID(t *testing.T) {
	client := &fakeExposeClient{}
	gateway, ctx := newClusterGateway(t, client)

	req := connect.NewRequest(&exposev1.RemoveExposureRequest{Port: 8080})
	resp, err := gateway.RemoveExposure(ctx, req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected CodePermissionDenied, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.removeExposureCalls != 0 {
		t.Fatalf("expected no remove exposure calls, got %d", client.removeExposureCalls)
	}
}

func TestExposeGatewayRemoveExposureError(t *testing.T) {
	client := &fakeExposeClient{
		removeExposureErr: status.Error(codes.NotFound, "missing"),
	}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.RemoveExposureRequest{Port: 8080})
	resp, err := gateway.RemoveExposure(agentContext(), req)
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
}

func TestExposeGatewayListExposuresDefaultsFromAgentIdentity(t *testing.T) {
	client := &fakeExposeClient{
		listExposuresResp: &exposev1.ListExposuresResponse{
			Exposures:     []*exposev1.Exposure{{Meta: &exposev1.EntityMeta{Id: "exposure-1"}, Port: 8080}},
			NextPageToken: "next",
		},
	}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.ListExposuresRequest{PageSize: 10, PageToken: "token"})
	resp, err := gateway.ListExposures(agentContext(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.listExposuresCalls != 1 {
		t.Fatalf("expected list exposures to be called once, got %d", client.listExposuresCalls)
	}
	if client.listExposuresReq.GetWorkloadId() != workloadID {
		t.Fatalf("expected workload id %q, got %q", workloadID, client.listExposuresReq.GetWorkloadId())
	}
	if client.listExposuresReq.GetPageSize() != 10 {
		t.Fatalf("expected page size 10, got %d", client.listExposuresReq.GetPageSize())
	}
	if client.listExposuresReq.GetPageToken() != "token" {
		t.Fatalf("expected page token %q, got %q", "token", client.listExposuresReq.GetPageToken())
	}
	if resp.Msg != client.listExposuresResp {
		t.Fatalf("expected response to be forwarded")
	}
}

func TestExposeGatewayListExposuresRejectsWorkloadMismatch(t *testing.T) {
	client := &fakeExposeClient{}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.ListExposuresRequest{WorkloadId: "workload-9"})
	resp, err := gateway.ListExposures(agentContext(), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected CodePermissionDenied, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.listExposuresCalls != 0 {
		t.Fatalf("expected no list exposures calls, got %d", client.listExposuresCalls)
	}
}

func TestExposeGatewayListExposuresClusterAdminExplicitIDs(t *testing.T) {
	client := &fakeExposeClient{
		listExposuresResp: &exposev1.ListExposuresResponse{},
	}
	gateway, ctx := newClusterGateway(t, client)

	req := connect.NewRequest(&exposev1.ListExposuresRequest{WorkloadId: "workload-9"})
	resp, err := gateway.ListExposures(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if client.listExposuresCalls != 1 {
		t.Fatalf("expected list exposures to be called once, got %d", client.listExposuresCalls)
	}
	if client.listExposuresReq.GetWorkloadId() != "workload-9" {
		t.Fatalf("expected workload id %q, got %q", "workload-9", client.listExposuresReq.GetWorkloadId())
	}
}

func TestExposeGatewayListExposuresClusterAdminRequiresWorkloadID(t *testing.T) {
	client := &fakeExposeClient{}
	gateway, ctx := newClusterGateway(t, client)

	req := connect.NewRequest(&exposev1.ListExposuresRequest{PageSize: 10})
	resp, err := gateway.ListExposures(ctx, req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected CodePermissionDenied, got %v", connect.CodeOf(err))
	}
	if resp != nil {
		t.Fatalf("expected no response")
	}
	if client.listExposuresCalls != 0 {
		t.Fatalf("expected no list exposures calls, got %d", client.listExposuresCalls)
	}
}

func TestExposeGatewayListExposuresError(t *testing.T) {
	client := &fakeExposeClient{
		listExposuresErr: status.Error(codes.PermissionDenied, "denied"),
	}
	gateway := NewExposeGateway(client, nil)

	req := connect.NewRequest(&exposev1.ListExposuresRequest{PageSize: 10})
	resp, err := gateway.ListExposures(agentContext(), req)
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
}
