package gateway

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	authorizationv1 "github.com/agynio/gateway/gen/agynio/api/authorization/v1"
	gatewayv1 "github.com/agynio/gateway/gen/agynio/api/gateway/v1"
	runnersv1 "github.com/agynio/gateway/gen/agynio/api/runners/v1"
	terminalproxyv1 "github.com/agynio/gateway/gen/agynio/api/terminal_proxy/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	terminalTestIdentityID    = "00000000-0000-0000-0000-000000000001"
	terminalTestWorkloadID    = "00000000-0000-0000-0000-000000000002"
	terminalTestSandboxID     = "00000000-0000-0000-0000-000000000003"
	terminalTestAgentID       = "00000000-0000-0000-0000-000000000004"
	terminalTestLegacyAgentID = "00000000-0000-0000-0000-000000000005"
)

type fakeTerminalRunnersClient struct {
	listWorkloadsCalls    int
	listWorkloadsRequests []*runnersv1.ListWorkloadsRequest
	listWorkloadsMetadata []metadata.MD
	listWorkloadsPages    []*runnersv1.ListWorkloadsResponse
	listWorkloadsErr      error
}

func (f *fakeTerminalRunnersClient) ListWorkloads(ctx context.Context, in *runnersv1.ListWorkloadsRequest, opts ...grpc.CallOption) (*runnersv1.ListWorkloadsResponse, error) {
	f.listWorkloadsCalls++
	f.listWorkloadsRequests = append(f.listWorkloadsRequests, in)
	ctx = identity.AppendToOutgoingContext(ctx)
	md, _ := metadata.FromOutgoingContext(ctx)
	f.listWorkloadsMetadata = append(f.listWorkloadsMetadata, md)
	if f.listWorkloadsErr != nil {
		return nil, f.listWorkloadsErr
	}
	if len(f.listWorkloadsPages) == 0 {
		return &runnersv1.ListWorkloadsResponse{}, nil
	}
	index := f.listWorkloadsCalls - 1
	if index >= len(f.listWorkloadsPages) {
		index = len(f.listWorkloadsPages) - 1
	}
	return f.listWorkloadsPages[index], nil
}

type fakeTerminalAuthorizationClient struct {
	checkCalls int
	checkReq   *authorizationv1.CheckRequest
	checkResp  *authorizationv1.CheckResponse
	checkErr   error
}

func (f *fakeTerminalAuthorizationClient) Check(ctx context.Context, in *authorizationv1.CheckRequest, opts ...grpc.CallOption) (*authorizationv1.CheckResponse, error) {
	f.checkCalls++
	f.checkReq = in
	if f.checkErr != nil {
		return nil, f.checkErr
	}
	if f.checkResp == nil {
		return &authorizationv1.CheckResponse{}, nil
	}
	return f.checkResp, nil
}

type fakeTerminalProxyClient struct {
	issueTicketCalls int
	issueTicketReq   *terminalproxyv1.IssueTicketRequest
	issueTicketResp  *terminalproxyv1.IssueTicketResponse
	issueTicketErr   error
}

func (f *fakeTerminalProxyClient) IssueTicket(ctx context.Context, in *terminalproxyv1.IssueTicketRequest, opts ...grpc.CallOption) (*terminalproxyv1.IssueTicketResponse, error) {
	f.issueTicketCalls++
	f.issueTicketReq = in
	if f.issueTicketErr != nil {
		return nil, f.issueTicketErr
	}
	if f.issueTicketResp == nil {
		return &terminalproxyv1.IssueTicketResponse{}, nil
	}
	return f.issueTicketResp, nil
}

func TestTerminalGatewayCreateTerminalSessionSandboxAllowed(t *testing.T) {
	command := &terminalproxyv1.TerminalCommand{Command: &terminalproxyv1.TerminalCommand_Shell{Shell: "bash"}}
	runners := &fakeTerminalRunnersClient{listWorkloadsPages: []*runnersv1.ListWorkloadsResponse{
		{Workloads: []*runnersv1.Workload{terminalSandboxWorkload(terminalTestWorkloadID, terminalTestSandboxID)}},
	}}
	authorization := &fakeTerminalAuthorizationClient{checkResp: &authorizationv1.CheckResponse{Allowed: true}}
	terminalProxy := &fakeTerminalProxyClient{issueTicketResp: &terminalproxyv1.IssueTicketResponse{
		Ticket:           "ticket-1",
		WebsocketUrl:     "wss://terminal.example/terminal?ticket=ticket-1",
		ExpiresInSeconds: 60,
	}}

	resp, err := NewTerminalGateway(runners, authorization, terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{
			WorkloadId:    " " + terminalTestWorkloadID + " ",
			ContainerName: " main ",
			Command:       command,
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Msg.GetTicket() != "ticket-1" {
		t.Fatalf("unexpected ticket: %s", resp.Msg.GetTicket())
	}
	if resp.Msg.GetWebsocketUrl() != "wss://terminal.example/terminal?ticket=ticket-1" {
		t.Fatalf("unexpected websocket url: %s", resp.Msg.GetWebsocketUrl())
	}
	if resp.Msg.GetExpiresInSeconds() != 60 {
		t.Fatalf("unexpected expiry: %d", resp.Msg.GetExpiresInSeconds())
	}
	assertTerminalTuple(t, authorization.checkReq.GetTupleKey(), "identity:"+terminalTestIdentityID, "can_connect", "sandbox:"+terminalTestSandboxID)
	assertTerminalIssueTicket(t, terminalProxy.issueTicketReq, terminalTestIdentityID, terminalTestWorkloadID, "main", command)
	if runners.listWorkloadsCalls != 1 {
		t.Fatalf("expected one workload list call, got %d", runners.listWorkloadsCalls)
	}
	if runners.listWorkloadsRequests[0].GetPageSize() != terminalListWorkloadsPageSize {
		t.Fatalf("unexpected page size: %d", runners.listWorkloadsRequests[0].GetPageSize())
	}
	assertNoTerminalIdentityMetadata(t, runners.listWorkloadsMetadata[0])
}

func TestTerminalGatewayCreateTerminalSessionSandboxDenied(t *testing.T) {
	runners := &fakeTerminalRunnersClient{listWorkloadsPages: []*runnersv1.ListWorkloadsResponse{
		{Workloads: []*runnersv1.Workload{terminalSandboxWorkload(terminalTestWorkloadID, terminalTestSandboxID)}},
	}}
	authorization := &fakeTerminalAuthorizationClient{checkResp: &authorizationv1.CheckResponse{Allowed: false}}
	terminalProxy := &fakeTerminalProxyClient{}

	_, err := NewTerminalGateway(runners, authorization, terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected permission denied, got %v: %v", connect.CodeOf(err), err)
	}
	if terminalProxy.issueTicketCalls != 0 {
		t.Fatalf("expected no ticket call, got %d", terminalProxy.issueTicketCalls)
	}
}

func TestTerminalGatewayCreateTerminalSessionAgentInstanceAllowed(t *testing.T) {
	agentClassID := terminalTestAgentID
	runners := &fakeTerminalRunnersClient{listWorkloadsPages: []*runnersv1.ListWorkloadsResponse{
		{Workloads: []*runnersv1.Workload{terminalAgentWorkload(terminalTestWorkloadID, &agentClassID, terminalTestLegacyAgentID)}},
	}}
	authorization := &fakeTerminalAuthorizationClient{checkResp: &authorizationv1.CheckResponse{Allowed: true}}
	terminalProxy := &fakeTerminalProxyClient{issueTicketResp: &terminalproxyv1.IssueTicketResponse{Ticket: "ticket-1"}}

	_, err := NewTerminalGateway(runners, authorization, terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTerminalTuple(t, authorization.checkReq.GetTupleKey(), "identity:"+terminalTestIdentityID, "can_edit_config", "agent:"+terminalTestAgentID)
	if terminalProxy.issueTicketCalls != 1 {
		t.Fatalf("expected ticket call, got %d", terminalProxy.issueTicketCalls)
	}
}

func TestTerminalGatewayCreateTerminalSessionAgentInstanceUsesLegacyAgentID(t *testing.T) {
	runners := &fakeTerminalRunnersClient{listWorkloadsPages: []*runnersv1.ListWorkloadsResponse{
		{Workloads: []*runnersv1.Workload{terminalAgentWorkload(terminalTestWorkloadID, nil, terminalTestLegacyAgentID)}},
	}}
	authorization := &fakeTerminalAuthorizationClient{checkResp: &authorizationv1.CheckResponse{Allowed: true}}
	terminalProxy := &fakeTerminalProxyClient{issueTicketResp: &terminalproxyv1.IssueTicketResponse{Ticket: "ticket-1"}}

	_, err := NewTerminalGateway(runners, authorization, terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTerminalTuple(t, authorization.checkReq.GetTupleKey(), "identity:"+terminalTestIdentityID, "can_edit_config", "agent:"+terminalTestLegacyAgentID)
}

func TestTerminalGatewayCreateTerminalSessionAgentInstanceDenied(t *testing.T) {
	agentClassID := terminalTestAgentID
	runners := &fakeTerminalRunnersClient{listWorkloadsPages: []*runnersv1.ListWorkloadsResponse{
		{Workloads: []*runnersv1.Workload{terminalAgentWorkload(terminalTestWorkloadID, &agentClassID, "")}},
	}}
	authorization := &fakeTerminalAuthorizationClient{checkResp: &authorizationv1.CheckResponse{Allowed: false}}
	terminalProxy := &fakeTerminalProxyClient{}

	_, err := NewTerminalGateway(runners, authorization, terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected permission denied, got %v: %v", connect.CodeOf(err), err)
	}
	if terminalProxy.issueTicketCalls != 0 {
		t.Fatalf("expected no ticket call, got %d", terminalProxy.issueTicketCalls)
	}
}

func TestTerminalGatewayCreateTerminalSessionPaginatesWorkloadLookup(t *testing.T) {
	runners := &fakeTerminalRunnersClient{listWorkloadsPages: []*runnersv1.ListWorkloadsResponse{
		{NextPageToken: "page-2"},
		{Workloads: []*runnersv1.Workload{terminalSandboxWorkload(terminalTestWorkloadID, terminalTestSandboxID)}},
	}}
	authorization := &fakeTerminalAuthorizationClient{checkResp: &authorizationv1.CheckResponse{Allowed: true}}
	terminalProxy := &fakeTerminalProxyClient{issueTicketResp: &terminalproxyv1.IssueTicketResponse{Ticket: "ticket-1"}}

	_, err := NewTerminalGateway(runners, authorization, terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if runners.listWorkloadsCalls != 2 {
		t.Fatalf("expected two workload list calls, got %d", runners.listWorkloadsCalls)
	}
	if runners.listWorkloadsRequests[1].GetPageToken() != "page-2" {
		t.Fatalf("unexpected second page token: %q", runners.listWorkloadsRequests[1].GetPageToken())
	}
}

func TestTerminalGatewayCreateTerminalSessionRequiresIdentity(t *testing.T) {
	_, err := NewTerminalGateway(&fakeTerminalRunnersClient{}, &fakeTerminalAuthorizationClient{}, &fakeTerminalProxyClient{}).CreateTerminalSession(
		context.Background(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Fatalf("expected unauthenticated, got %v: %v", connect.CodeOf(err), err)
	}
}

func TestTerminalGatewayCreateTerminalSessionValidatesRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *gatewayv1.CreateTerminalSessionRequest
	}{
		{name: "missing workload", req: &gatewayv1.CreateTerminalSessionRequest{ContainerName: "main"}},
		{name: "invalid workload", req: &gatewayv1.CreateTerminalSessionRequest{WorkloadId: "not-a-uuid", ContainerName: "main"}},
		{name: "missing container", req: &gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTerminalGateway(&fakeTerminalRunnersClient{}, &fakeTerminalAuthorizationClient{}, &fakeTerminalProxyClient{}).CreateTerminalSession(
				terminalIdentityContext(),
				connect.NewRequest(tt.req),
			)
			if connect.CodeOf(err) != connect.CodeInvalidArgument {
				t.Fatalf("expected invalid argument, got %v: %v", connect.CodeOf(err), err)
			}
		})
	}
}

func TestTerminalGatewayCreateTerminalSessionWorkloadNotFound(t *testing.T) {
	_, err := NewTerminalGateway(
		&fakeTerminalRunnersClient{listWorkloadsPages: []*runnersv1.ListWorkloadsResponse{{}}},
		&fakeTerminalAuthorizationClient{},
		&fakeTerminalProxyClient{},
	).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if connect.CodeOf(err) != connect.CodeNotFound {
		t.Fatalf("expected not found, got %v: %v", connect.CodeOf(err), err)
	}
}

func TestTerminalGatewayCreateTerminalSessionPropagatesRunnersError(t *testing.T) {
	_, err := NewTerminalGateway(
		&fakeTerminalRunnersClient{listWorkloadsErr: status.Error(codes.Unavailable, "runners unavailable")},
		&fakeTerminalAuthorizationClient{},
		&fakeTerminalProxyClient{},
	).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if connect.CodeOf(err) != connect.CodeUnavailable {
		t.Fatalf("expected unavailable, got %v: %v", connect.CodeOf(err), err)
	}
}

func TestTerminalGatewayCreateTerminalSessionUnsupportedOwnerKind(t *testing.T) {
	runners := &fakeTerminalRunnersClient{listWorkloadsPages: []*runnersv1.ListWorkloadsResponse{
		{Workloads: []*runnersv1.Workload{{Meta: &runnersv1.EntityMeta{Id: terminalTestWorkloadID}}}},
	}}

	_, err := NewTerminalGateway(runners, &fakeTerminalAuthorizationClient{}, &fakeTerminalProxyClient{}).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if connect.CodeOf(err) != connect.CodeFailedPrecondition {
		t.Fatalf("expected failed precondition, got %v: %v", connect.CodeOf(err), err)
	}
}

func TestTerminalGatewayCreateTerminalSessionPropagatesAuthorizationError(t *testing.T) {
	runners := &fakeTerminalRunnersClient{listWorkloadsPages: []*runnersv1.ListWorkloadsResponse{
		{Workloads: []*runnersv1.Workload{terminalSandboxWorkload(terminalTestWorkloadID, terminalTestSandboxID)}},
	}}
	authorization := &fakeTerminalAuthorizationClient{checkErr: status.Error(codes.Unavailable, "authorization unavailable")}

	_, err := NewTerminalGateway(runners, authorization, &fakeTerminalProxyClient{}).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if connect.CodeOf(err) != connect.CodeUnavailable {
		t.Fatalf("expected unavailable, got %v: %v", connect.CodeOf(err), err)
	}
}

func TestTerminalGatewayCreateTerminalSessionPropagatesTicketError(t *testing.T) {
	runners := &fakeTerminalRunnersClient{listWorkloadsPages: []*runnersv1.ListWorkloadsResponse{
		{Workloads: []*runnersv1.Workload{terminalSandboxWorkload(terminalTestWorkloadID, terminalTestSandboxID)}},
	}}
	authorization := &fakeTerminalAuthorizationClient{checkResp: &authorizationv1.CheckResponse{Allowed: true}}
	terminalProxy := &fakeTerminalProxyClient{issueTicketErr: status.Error(codes.Unavailable, "terminal proxy unavailable")}

	_, err := NewTerminalGateway(runners, authorization, terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if connect.CodeOf(err) != connect.CodeUnavailable {
		t.Fatalf("expected unavailable, got %v: %v", connect.CodeOf(err), err)
	}
}

func terminalIdentityContext() context.Context {
	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		identity.MetadataKeyIdentityID, "stale-identity",
		identity.MetadataKeyIdentityType, string(identity.IdentityTypeRunner),
	)
	return identity.WithIdentity(ctx, identity.ResolvedIdentity{
		IdentityID:   terminalTestIdentityID,
		IdentityType: identity.IdentityTypeUser,
	})
}

func terminalSandboxWorkload(workloadID, sandboxID string) *runnersv1.Workload {
	return &runnersv1.Workload{
		Meta:      &runnersv1.EntityMeta{Id: workloadID},
		OwnerKind: runnersv1.RuntimeOwnerKind_RUNTIME_OWNER_KIND_SANDBOX,
		OwnerId:   sandboxID,
	}
}

func terminalAgentWorkload(workloadID string, agentClassID *string, legacyAgentID string) *runnersv1.Workload {
	return &runnersv1.Workload{
		Meta:         &runnersv1.EntityMeta{Id: workloadID},
		OwnerKind:    runnersv1.RuntimeOwnerKind_RUNTIME_OWNER_KIND_AGENT_INSTANCE,
		AgentClassId: agentClassID,
		AgentId:      legacyAgentID,
	}
}

func assertTerminalTuple(t *testing.T, tuple *authorizationv1.TupleKey, user, relation, object string) {
	t.Helper()
	if tuple == nil {
		t.Fatalf("expected authorization tuple")
	}
	if tuple.GetUser() != user {
		t.Fatalf("expected tuple user %q, got %q", user, tuple.GetUser())
	}
	if tuple.GetRelation() != relation {
		t.Fatalf("expected tuple relation %q, got %q", relation, tuple.GetRelation())
	}
	if tuple.GetObject() != object {
		t.Fatalf("expected tuple object %q, got %q", object, tuple.GetObject())
	}
}

func assertTerminalIssueTicket(t *testing.T, req *terminalproxyv1.IssueTicketRequest, identityID, workloadID, containerName string, command *terminalproxyv1.TerminalCommand) {
	t.Helper()
	if req == nil {
		t.Fatalf("expected issue ticket request")
	}
	if req.GetIdentityId() != identityID {
		t.Fatalf("expected identity id %q, got %q", identityID, req.GetIdentityId())
	}
	if req.GetWorkloadId() != workloadID {
		t.Fatalf("expected workload id %q, got %q", workloadID, req.GetWorkloadId())
	}
	if req.GetContainerName() != containerName {
		t.Fatalf("expected container name %q, got %q", containerName, req.GetContainerName())
	}
	if req.GetCommand() != command {
		t.Fatalf("expected command to be forwarded unchanged")
	}
}

func assertNoTerminalIdentityMetadata(t *testing.T, md metadata.MD) {
	t.Helper()
	if values := md.Get(identity.MetadataKeyIdentityID); len(values) != 0 {
		t.Fatalf("expected no identity metadata, got %v", values)
	}
	if values := md.Get(identity.MetadataKeyIdentityType); len(values) != 0 {
		t.Fatalf("expected no identity type metadata, got %v", values)
	}
}
