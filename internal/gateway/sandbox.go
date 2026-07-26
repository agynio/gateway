package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/grpc"

	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	"github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	"github.com/agynio/gateway/internal/identity"
)

// sandboxProcedures is the complete set of Gateway operations a sandbox
// workload identity may call. A sandbox is not an organization member and holds
// no organization tuples, so nothing outside this set is reachable for it: port
// exposure, so `agyn expose` works from inside the sandbox, and file upload and
// download, so agent tooling can move files. LLM calls are authorized by the LLM
// Proxy, which the sandbox reaches directly over the overlay rather than through
// the Gateway.
var sandboxProcedures = map[string]struct{}{
	gatewayv1connect.ExposeGatewayAddExposureProcedure:    {},
	gatewayv1connect.ExposeGatewayRemoveExposureProcedure: {},
	gatewayv1connect.ExposeGatewayListExposuresProcedure:  {},
	gatewayv1connect.FilesGatewayUploadFileProcedure:      {},
	gatewayv1connect.FilesGatewayGetFileMetadataProcedure: {},
	gatewayv1connect.FilesGatewayGetDownloadUrlProcedure:  {},
	gatewayv1connect.FilesGatewayGetFileContentProcedure:  {},
}

// SandboxRecordResolver reads the sandbox record behind a sandbox workload
// identity. It is the Agents service.
type SandboxRecordResolver interface {
	GetSandbox(ctx context.Context, in *agentsv1.GetSandboxRequest, opts ...grpc.CallOption) (*agentsv1.GetSandboxResponse, error)
}

// SandboxAuthorizer decides what a sandbox workload identity may do. Identities
// of every other type pass through it untouched — they are authorized by the
// services behind the Gateway against the tuples they hold.
type SandboxAuthorizer struct {
	agents SandboxRecordResolver
}

func NewSandboxAuthorizer(agents SandboxRecordResolver) *SandboxAuthorizer {
	if agents == nil {
		panic("sandbox record resolver is required")
	}
	return &SandboxAuthorizer{agents: agents}
}

// Authorize admits a sandbox identity to the named procedure. Reachability over
// the overlay is not authorization: a sandbox dials the Gateway the way an agent
// workload does, and this is where that reach is narrowed to the operations a
// sandbox is meant to have. The sandbox record is read on every call rather than
// cached — it is one in-cluster hop, and a cache would keep answering for a
// sandbox whose record has since been terminated.
func (a *SandboxAuthorizer) Authorize(ctx context.Context, procedure string) error {
	caller, ok := identity.IdentityFromContext(ctx)
	if !ok || caller.IdentityType != identity.IdentityTypeSandbox {
		return nil
	}

	if _, allowed := sandboxProcedures[procedure]; !allowed {
		return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("sandbox identities cannot call %s", procedure))
	}

	sandboxID := caller.SandboxID()
	if sandboxID == "" {
		return connect.NewError(connect.CodeUnauthenticated, errors.New("sandbox id missing for sandbox identity"))
	}

	response, err := a.agents.GetSandbox(downstreamContext(ctx), &agentsv1.GetSandboxRequest{
		Ref: &agentsv1.GetSandboxRequest_Id{Id: sandboxID},
	})
	if err != nil {
		return toConnectError(err)
	}

	sandbox := response.GetSandbox()
	if strings.TrimSpace(sandbox.GetOrganizationId()) == "" || strings.TrimSpace(sandbox.GetOwnerId()) == "" {
		return connect.NewError(connect.CodePermissionDenied, errors.New("sandbox has no organization or owner"))
	}
	if sandbox.GetStatus() == agentsv1.SandboxStatus_SANDBOX_STATUS_TERMINATED {
		return connect.NewError(connect.CodePermissionDenied, errors.New("sandbox is terminated"))
	}

	return nil
}

// NewSandboxInterceptor enforces the sandbox operation set on the ConnectRPC
// surface. It must run after the auth interceptor, which is what puts the
// resolved identity on the context.
func NewSandboxInterceptor(authorizer *SandboxAuthorizer) connect.Interceptor {
	if authorizer == nil {
		panic("sandbox authorizer is required")
	}
	return sandboxInterceptor{authorizer: authorizer}
}

type sandboxInterceptor struct {
	authorizer *SandboxAuthorizer
}

func (s sandboxInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if err := s.authorizer.Authorize(ctx, req.Spec().Procedure); err != nil {
			return nil, err
		}
		return next(ctx, req)
	}
}

func (s sandboxInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (s sandboxInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		if err := s.authorizer.Authorize(ctx, conn.Spec().Procedure); err != nil {
			return err
		}
		return next(ctx, conn)
	}
}

// NewSandboxMiddleware refuses sandbox identities on the Gateway's plain HTTP
// surface — `/me` and the app proxy. Neither is in the sandbox operation set, so
// there is no procedure to match: a sandbox has no user profile to read and no
// business calling an installed app.
func NewSandboxMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			caller, ok := identity.IdentityFromContext(r.Context())
			if ok && caller.IdentityType == identity.IdentityTypeSandbox {
				writeProblem(w, http.StatusForbidden, "sandbox identities cannot call "+r.URL.Path)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
