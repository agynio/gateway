package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	"github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	"github.com/agynio/gateway/internal/identity"
)

const (
	testSandboxID         = "11111111-1111-1111-1111-111111111111"
	testSandboxOwnerID    = "22222222-2222-2222-2222-222222222222"
	testSandboxOrgID      = "33333333-3333-3333-3333-333333333333"
	testSandboxWorkloadID = "44444444-4444-4444-4444-444444444444"
)

type fakeSandboxResolver struct {
	sandbox      *agentsv1.Sandbox
	err          error
	calls        int
	lastID       string
	lastMetadata metadata.MD
}

func (f *fakeSandboxResolver) GetSandbox(ctx context.Context, in *agentsv1.GetSandboxRequest, _ ...grpc.CallOption) (*agentsv1.GetSandboxResponse, error) {
	f.calls++
	f.lastID = in.GetId()
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		f.lastMetadata = md
	}
	if f.err != nil {
		return nil, f.err
	}
	return &agentsv1.GetSandboxResponse{Sandbox: f.sandbox}, nil
}

func runningSandboxResolver() *fakeSandboxResolver {
	return &fakeSandboxResolver{
		sandbox: &agentsv1.Sandbox{
			Meta:           &agentsv1.EntityMeta{Id: testSandboxID},
			OrganizationId: testSandboxOrgID,
			OwnerId:        testSandboxOwnerID,
			Status:         agentsv1.SandboxStatus_SANDBOX_STATUS_RUNNING,
		},
	}
}

func sandboxContext() context.Context {
	return identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   testSandboxID,
		IdentityType: identity.IdentityTypeSandbox,
		WorkloadID:   testSandboxWorkloadID,
	})
}

// A sandbox holds no organization tuples and is not a member of the
// organization it runs in. Everything an organization member can do through the
// Gateway has to be refused for it, whatever the services behind the Gateway
// would answer — this is the check that keeps overlay reach from turning into
// organization access.
func TestSandboxAuthorizerRefusesOrganizationMemberOperations(t *testing.T) {
	denied := []string{
		gatewayv1connect.AgentsGatewayListAgentsProcedure,
		gatewayv1connect.AgentsGatewayGetAgentProcedure,
		gatewayv1connect.AgentsGatewayCreateSandboxProcedure,
		gatewayv1connect.AgentsGatewayListSandboxesProcedure,
		gatewayv1connect.AgentsGatewayDeleteSandboxProcedure,
		gatewayv1connect.AgentsGatewayCreateInstanceProcedure,
		gatewayv1connect.AgentsGatewayWriteInboxItemProcedure,
		gatewayv1connect.AgentsGatewayListEnvironmentsProcedure,
		gatewayv1connect.ThreadsGatewayCreateThreadProcedure,
		gatewayv1connect.SecretsGatewayListSecretsProcedure,
		gatewayv1connect.LLMGatewayListModelsProcedure,
		gatewayv1connect.OrganizationsGatewayListOrganizationsProcedure,
		gatewayv1connect.UsersGatewayListUsersProcedure,
		gatewayv1connect.TerminalGatewayCreateTerminalSessionProcedure,
		gatewayv1connect.RunnersGatewayListWorkloadsProcedure,
		gatewayv1connect.TokenCountingGatewayCountTokensProcedure,
	}

	for _, procedure := range denied {
		resolver := runningSandboxResolver()
		err := NewSandboxAuthorizer(resolver).Authorize(sandboxContext(), procedure)
		if err == nil {
			t.Fatalf("%s: expected refusal", procedure)
		}
		if connect.CodeOf(err) != connect.CodePermissionDenied {
			t.Fatalf("%s: expected permission denied, got %v", procedure, connect.CodeOf(err))
		}
		if resolver.calls != 0 {
			t.Fatalf("%s: a refused procedure must not read the sandbox record", procedure)
		}
	}
}

func TestSandboxAuthorizerAllowsExposureAndFiles(t *testing.T) {
	allowed := []string{
		gatewayv1connect.ExposeGatewayAddExposureProcedure,
		gatewayv1connect.ExposeGatewayRemoveExposureProcedure,
		gatewayv1connect.ExposeGatewayListExposuresProcedure,
		gatewayv1connect.FilesGatewayUploadFileProcedure,
		gatewayv1connect.FilesGatewayGetFileMetadataProcedure,
		gatewayv1connect.FilesGatewayGetDownloadUrlProcedure,
		gatewayv1connect.FilesGatewayGetFileContentProcedure,
	}

	for _, procedure := range allowed {
		resolver := runningSandboxResolver()
		if err := NewSandboxAuthorizer(resolver).Authorize(sandboxContext(), procedure); err != nil {
			t.Fatalf("%s: unexpected error: %v", procedure, err)
		}
		if resolver.calls != 1 {
			t.Fatalf("%s: expected one sandbox lookup, got %d", procedure, resolver.calls)
		}
		if resolver.lastID != testSandboxID {
			t.Fatalf("%s: resolved %s", procedure, resolver.lastID)
		}
	}
}

func TestSandboxAuthorizerRefusesUnknownSandbox(t *testing.T) {
	resolver := &fakeSandboxResolver{err: status.Error(codes.NotFound, "sandbox not found")}

	err := NewSandboxAuthorizer(resolver).Authorize(sandboxContext(), gatewayv1connect.ExposeGatewayAddExposureProcedure)
	if err == nil {
		t.Fatalf("expected refusal")
	}
	if connect.CodeOf(err) != connect.CodeNotFound {
		t.Fatalf("expected not found, got %v", connect.CodeOf(err))
	}
}

func TestSandboxAuthorizerRefusesTerminatedSandbox(t *testing.T) {
	resolver := runningSandboxResolver()
	resolver.sandbox.Status = agentsv1.SandboxStatus_SANDBOX_STATUS_TERMINATED

	err := NewSandboxAuthorizer(resolver).Authorize(sandboxContext(), gatewayv1connect.ExposeGatewayAddExposureProcedure)
	if err == nil {
		t.Fatalf("expected refusal")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected permission denied, got %v", connect.CodeOf(err))
	}
}

func TestSandboxAuthorizerRefusesOwnerlessSandbox(t *testing.T) {
	resolver := runningSandboxResolver()
	resolver.sandbox.OwnerId = ""

	err := NewSandboxAuthorizer(resolver).Authorize(sandboxContext(), gatewayv1connect.ExposeGatewayAddExposureProcedure)
	if err == nil {
		t.Fatalf("expected refusal")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected permission denied, got %v", connect.CodeOf(err))
	}
}

func TestSandboxAuthorizerResolvesAsTheSandboxItself(t *testing.T) {
	resolver := runningSandboxResolver()

	if err := NewSandboxAuthorizer(resolver).Authorize(sandboxContext(), gatewayv1connect.FilesGatewayUploadFileProcedure); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := resolver.lastMetadata.Get(identity.MetadataKeyIdentityID); len(got) != 1 || got[0] != testSandboxID {
		t.Fatalf("unexpected identity metadata: %v", got)
	}
	if got := resolver.lastMetadata.Get(identity.MetadataKeyIdentityType); len(got) != 1 || got[0] != string(identity.IdentityTypeSandbox) {
		t.Fatalf("unexpected identity type metadata: %v", got)
	}
}

func TestSandboxAuthorizerIgnoresOtherIdentityTypes(t *testing.T) {
	resolver := runningSandboxResolver()
	ctx := identity.WithIdentity(context.Background(), identity.ResolvedIdentity{
		IdentityID:   testSandboxOwnerID,
		IdentityType: identity.IdentityTypeUser,
	})

	if err := NewSandboxAuthorizer(resolver).Authorize(ctx, gatewayv1connect.AgentsGatewayListAgentsProcedure); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolver.calls != 0 {
		t.Fatalf("expected no sandbox lookup, got %d", resolver.calls)
	}
}

// The interceptor has to read the procedure the same way ConnectRPC dispatches
// it, so the refusal is exercised over a real handler rather than a synthesized
// spec. The Gateway is constructed without downstream clients: a refused call
// must never reach one.
func TestSandboxInterceptorRefusesListAgentsOverConnect(t *testing.T) {
	resolver := runningSandboxResolver()
	handlerPath, handler := gatewayv1connect.NewAgentsGatewayHandler(
		New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil),
		connect.WithInterceptors(
			staticIdentityInterceptor{},
			NewSandboxInterceptor(NewSandboxAuthorizer(resolver)),
		),
	)

	mux := http.NewServeMux()
	mux.Handle(handlerPath, handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	client := gatewayv1connect.NewAgentsGatewayClient(server.Client(), server.URL)
	_, err := client.ListAgents(context.Background(), connect.NewRequest(&agentsv1.ListAgentsRequest{
		OrganizationId: testSandboxOrgID,
	}))
	if err == nil {
		t.Fatalf("expected refusal")
	}
	if connect.CodeOf(err) != connect.CodePermissionDenied {
		t.Fatalf("expected permission denied, got %v", connect.CodeOf(err))
	}
	if resolver.calls != 0 {
		t.Fatalf("a refused procedure must not read the sandbox record")
	}
}

type staticIdentityInterceptor struct{}

func (staticIdentityInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return next(sandboxIdentityOn(ctx), req)
	}
}

func (staticIdentityInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (staticIdentityInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		return next(sandboxIdentityOn(ctx), conn)
	}
}

func sandboxIdentityOn(ctx context.Context) context.Context {
	return identity.WithIdentity(ctx, identity.ResolvedIdentity{
		IdentityID:   testSandboxID,
		IdentityType: identity.IdentityTypeSandbox,
		WorkloadID:   testSandboxWorkloadID,
	})
}

func TestSandboxMiddlewareRefusesSandboxIdentity(t *testing.T) {
	called := false
	handler := NewSandboxMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/me", nil).WithContext(sandboxContext())
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", recorder.Code)
	}
	if called {
		t.Fatalf("handler reached")
	}
}

func TestSandboxMiddlewarePassesOtherIdentityTypes(t *testing.T) {
	called := false
	handler := NewSandboxMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/me", nil))

	if !called {
		t.Fatalf("handler not reached")
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}
