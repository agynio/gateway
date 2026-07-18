package gateway

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authorizationv1 "github.com/agynio/gateway/gen/agynio/api/authorization/v1"
	gatewayv1 "github.com/agynio/gateway/gen/agynio/api/gateway/v1"
	runnersv1 "github.com/agynio/gateway/gen/agynio/api/runners/v1"
	terminalproxyv1 "github.com/agynio/gateway/gen/agynio/api/terminal_proxy/v1"
	"github.com/agynio/gateway/internal/identity"
)

const (
	terminalListWorkloadsPageSize = 1000

	terminalIdentityObjectPrefix = "identity:"
	terminalSandboxObjectPrefix  = "sandbox:"
	terminalAgentObjectPrefix    = "agent:"

	terminalSandboxConnectRelation    = "can_connect"
	terminalAgentEditConfigRelation   = "can_edit_config"
	terminalWorkloadPaginationMaxPage = 1000
)

type terminalRunnersClient interface {
	ListWorkloads(context.Context, *runnersv1.ListWorkloadsRequest, ...grpc.CallOption) (*runnersv1.ListWorkloadsResponse, error)
}

type terminalAuthorizationClient interface {
	Check(context.Context, *authorizationv1.CheckRequest, ...grpc.CallOption) (*authorizationv1.CheckResponse, error)
}

type terminalProxyClient interface {
	IssueTicket(context.Context, *terminalproxyv1.IssueTicketRequest, ...grpc.CallOption) (*terminalproxyv1.IssueTicketResponse, error)
}

type TerminalGateway struct {
	runners       terminalRunnersClient
	authorization terminalAuthorizationClient
	terminalProxy terminalProxyClient
}

func NewTerminalGateway(runners terminalRunnersClient, authorization terminalAuthorizationClient, terminalProxy terminalProxyClient) *TerminalGateway {
	return &TerminalGateway{runners: runners, authorization: authorization, terminalProxy: terminalProxy}
}

func (g *TerminalGateway) CreateTerminalSession(ctx context.Context, req *connect.Request[gatewayv1.CreateTerminalSessionRequest]) (*connect.Response[gatewayv1.CreateTerminalSessionResponse], error) {
	caller, ok := identity.IdentityFromContext(ctx)
	if !ok || strings.TrimSpace(caller.IdentityID) == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	workloadID, err := validateTerminalWorkloadID(req.Msg.GetWorkloadId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	containerName, err := validateTerminalContainerName(req.Msg.GetContainerName())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	workload, err := g.resolveTerminalWorkload(ctx, workloadID)
	if err != nil {
		return nil, toConnectError(err)
	}
	authorizationTuple, err := terminalAuthorizationTuple(caller.IdentityID, workload)
	if err != nil {
		return nil, toConnectError(err)
	}

	allowed, err := g.authorization.Check(ctx, &authorizationv1.CheckRequest{TupleKey: authorizationTuple})
	if err != nil {
		return nil, toConnectError(err)
	}
	if !allowed.GetAllowed() {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("terminal access denied"))
	}

	ticket, err := g.terminalProxy.IssueTicket(ctx, &terminalproxyv1.IssueTicketRequest{
		IdentityId:    strings.TrimSpace(caller.IdentityID),
		WorkloadId:    workloadID,
		ContainerName: containerName,
		Command:       req.Msg.GetCommand(),
	})
	if err != nil {
		return nil, toConnectError(err)
	}

	return connect.NewResponse(&gatewayv1.CreateTerminalSessionResponse{
		Ticket:           ticket.GetTicket(),
		WebsocketUrl:     ticket.GetWebsocketUrl(),
		ExpiresInSeconds: ticket.GetExpiresInSeconds(),
	}), nil
}

func validateTerminalWorkloadID(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", errors.New("workload_id: value is empty")
	}
	if _, err := uuid.Parse(trimmed); err != nil {
		return "", errors.New("workload_id: invalid UUID")
	}
	return trimmed, nil
}

func validateTerminalContainerName(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", errors.New("container_name: value is empty")
	}
	return trimmed, nil
}

func (g *TerminalGateway) resolveTerminalWorkload(ctx context.Context, workloadID string) (*runnersv1.Workload, error) {
	pageToken := ""
	for page := 0; page < terminalWorkloadPaginationMaxPage; page++ {
		resp, err := g.runners.ListWorkloads(internalRunnersContext(ctx), &runnersv1.ListWorkloadsRequest{
			PageSize:  terminalListWorkloadsPageSize,
			PageToken: pageToken,
		})
		if err != nil {
			return nil, err
		}
		for _, workload := range resp.GetWorkloads() {
			if workload.GetMeta().GetId() == workloadID {
				return workload, nil
			}
		}
		pageToken = resp.GetNextPageToken()
		if pageToken == "" {
			return nil, status.Error(codes.NotFound, "workload not found")
		}
	}
	return nil, status.Error(codes.Internal, "workload lookup exceeded pagination limit")
}

func internalRunnersContext(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(identity.WithoutIdentity(ctx), metadata.MD{})
}

func terminalAuthorizationTuple(identityID string, workload *runnersv1.Workload) (*authorizationv1.TupleKey, error) {
	user := terminalIdentityObjectPrefix + strings.TrimSpace(identityID)
	switch workload.GetOwnerKind() {
	case runnersv1.RuntimeOwnerKind_RUNTIME_OWNER_KIND_SANDBOX:
		ownerID := strings.TrimSpace(workload.GetOwnerId())
		if ownerID == "" {
			return nil, status.Error(codes.FailedPrecondition, "sandbox workload missing owner_id")
		}
		return &authorizationv1.TupleKey{
			User:     user,
			Relation: terminalSandboxConnectRelation,
			Object:   terminalSandboxObjectPrefix + ownerID,
		}, nil
	case runnersv1.RuntimeOwnerKind_RUNTIME_OWNER_KIND_AGENT_INSTANCE:
		agentID := terminalAgentID(workload)
		if agentID == "" {
			return nil, status.Error(codes.FailedPrecondition, "agent-instance workload missing agent_id")
		}
		return &authorizationv1.TupleKey{
			User:     user,
			Relation: terminalAgentEditConfigRelation,
			Object:   terminalAgentObjectPrefix + agentID,
		}, nil
	case runnersv1.RuntimeOwnerKind_RUNTIME_OWNER_KIND_UNSPECIFIED:
		return nil, status.Error(codes.FailedPrecondition, "workload owner kind is unspecified")
	default:
		return nil, status.Error(codes.FailedPrecondition, "workload owner kind is unsupported")
	}
}

func terminalAgentID(workload *runnersv1.Workload) string {
	if agentClassID := strings.TrimSpace(workload.GetAgentClassId()); agentClassID != "" {
		return agentClassID
	}
	return strings.TrimSpace(workload.GetAgentId())
}
