package gateway

import (
	"context"
	"reflect"
	"testing"

	"connectrpc.com/connect"
	gatewayv1 "github.com/agynio/gateway/gen/agynio/api/gateway/v1"
	terminalproxyv1 "github.com/agynio/gateway/gen/agynio/api/terminal_proxy/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	terminalTestIdentityID = "00000000-0000-0000-0000-000000000001"
	terminalTestWorkloadID = "00000000-0000-0000-0000-000000000002"
)

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

func TestTerminalGatewayCreateTerminalSessionForwardsIssueTicket(t *testing.T) {
	command := &terminalproxyv1.TerminalCommand{Command: &terminalproxyv1.TerminalCommand_Shell{Shell: "bash"}}
	terminalProxy := &fakeTerminalProxyClient{issueTicketResp: &terminalproxyv1.IssueTicketResponse{
		Ticket:           "ticket-1",
		WebsocketUrl:     "wss://terminal.example/terminal?ticket=ticket-1",
		ExpiresInSeconds: 60,
	}}

	resp, err := NewTerminalGateway(terminalProxy).CreateTerminalSession(
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
	assertTerminalIssueTicket(t, terminalProxy.issueTicketReq, terminalTestIdentityID, terminalTestWorkloadID, "main", command)
}

func TestTerminalGatewayCreateTerminalSessionRequiresIdentity(t *testing.T) {
	terminalProxy := &fakeTerminalProxyClient{}

	_, err := NewTerminalGateway(terminalProxy).CreateTerminalSession(
		context.Background(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Fatalf("expected unauthenticated, got %v: %v", connect.CodeOf(err), err)
	}
	if terminalProxy.issueTicketCalls != 0 {
		t.Fatalf("expected no terminal proxy calls, got %d", terminalProxy.issueTicketCalls)
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
			terminalProxy := &fakeTerminalProxyClient{}
			_, err := NewTerminalGateway(terminalProxy).CreateTerminalSession(
				terminalIdentityContext(),
				connect.NewRequest(tt.req),
			)
			if connect.CodeOf(err) != connect.CodeInvalidArgument {
				t.Fatalf("expected invalid argument, got %v: %v", connect.CodeOf(err), err)
			}
			if terminalProxy.issueTicketCalls != 0 {
				t.Fatalf("expected no terminal proxy calls, got %d", terminalProxy.issueTicketCalls)
			}
		})
	}
}

func TestTerminalGatewayCreateTerminalSessionPropagatesTicketError(t *testing.T) {
	terminalProxy := &fakeTerminalProxyClient{issueTicketErr: status.Error(codes.PermissionDenied, "terminal access denied")}

	_, err := NewTerminalGateway(terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{WorkloadId: terminalTestWorkloadID, ContainerName: "main"}),
	)
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected permission denied, got %v: %v", connect.CodeOf(err), err)
	}
	if terminalProxy.issueTicketCalls != 1 {
		t.Fatalf("expected one terminal proxy call, got %d", terminalProxy.issueTicketCalls)
	}
}

func TestTerminalGatewayHasNoAuthorizationOrRunnersDependencies(t *testing.T) {
	gatewayType := reflect.TypeOf(TerminalGateway{})
	for _, fieldName := range []string{"authorization", "runners"} {
		if _, ok := gatewayType.FieldByName(fieldName); ok {
			t.Fatalf("terminal gateway must not keep %s dependency", fieldName)
		}
	}
}

func terminalIdentityContext() context.Context {
	return identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   terminalTestIdentityID,
		IdentityType: identity.IdentityTypeUser,
	})
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

// The Gateway forwards the kind and its parameters untouched. Interpreting them
// here would put a service's domain rules in its router.
// A shell a caller names has to survive the hop. The Gateway dropping it is
// invisible -- the proxy issues a working ticket for an unnamed shell, and the
// only symptom is that persistence silently does not happen.
func TestTerminalGatewayForwardsShellParameters(t *testing.T) {
	terminalProxy := &fakeTerminalProxyClient{issueTicketResp: &terminalproxyv1.IssueTicketResponse{
		Ticket: "t", WebsocketUrl: "wss://proxy/terminal", ExpiresInSeconds: 30,
	}}
	_, err := NewTerminalGateway(terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{
			WorkloadId:    terminalTestWorkloadID,
			ContainerName: "main",
			Kind:          terminalproxyv1.SessionKind_SESSION_KIND_SHELL,
			ShellId:       "shell-a",
			ShellCwd:      "/workspace",
		}),
	)
	if err != nil {
		t.Fatalf("CreateTerminalSession: %v", err)
	}
	if got := terminalProxy.issueTicketReq.GetShellId(); got != "shell-a" {
		t.Fatalf("shell_id not forwarded, got %q", got)
	}
	if got := terminalProxy.issueTicketReq.GetShellCwd(); got != "/workspace" {
		t.Fatalf("shell_cwd not forwarded, got %q", got)
	}
}

func TestTerminalGatewayForwardsKindAndSyncRoot(t *testing.T) {
	terminalProxy := &fakeTerminalProxyClient{issueTicketResp: &terminalproxyv1.IssueTicketResponse{
		Ticket: "t", WebsocketUrl: "wss://proxy/terminal", ExpiresInSeconds: 30,
	}}
	_, err := NewTerminalGateway(terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{
			WorkloadId:    terminalTestWorkloadID,
			ContainerName: "main",
			Kind:          terminalproxyv1.SessionKind_SESSION_KIND_SYNC,
			SyncRoot:      "/workspace/project",
		}),
	)
	if err != nil {
		t.Fatalf("create terminal session: %v", err)
	}
	if terminalProxy.issueTicketReq.GetKind() != terminalproxyv1.SessionKind_SESSION_KIND_SYNC {
		t.Fatalf("kind not forwarded, got %s", terminalProxy.issueTicketReq.GetKind())
	}
	if terminalProxy.issueTicketReq.GetSyncRoot() != "/workspace/project" {
		t.Fatalf("sync_root not forwarded, got %q", terminalProxy.issueTicketReq.GetSyncRoot())
	}
}

// An unspecified kind is the proxy's to reject, not the Gateway's to default.
func TestTerminalGatewayDoesNotDefaultTheKind(t *testing.T) {
	terminalProxy := &fakeTerminalProxyClient{issueTicketResp: &terminalproxyv1.IssueTicketResponse{Ticket: "t"}}
	_, _ = NewTerminalGateway(terminalProxy).CreateTerminalSession(
		terminalIdentityContext(),
		connect.NewRequest(&gatewayv1.CreateTerminalSessionRequest{
			WorkloadId:    terminalTestWorkloadID,
			ContainerName: "main",
		}),
	)
	if terminalProxy.issueTicketReq.GetKind() != terminalproxyv1.SessionKind_SESSION_KIND_UNSPECIFIED {
		t.Fatalf("the Gateway substituted a kind: %s", terminalProxy.issueTicketReq.GetKind())
	}
}
