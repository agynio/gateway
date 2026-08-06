package gateway

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/grpc"

	gatewayv1 "github.com/agynio/gateway/gen/agynio/api/gateway/v1"
	terminalproxyv1 "github.com/agynio/gateway/gen/agynio/api/terminal_proxy/v1"
	"github.com/agynio/gateway/internal/identity"
)

type terminalProxyClient interface {
	IssueTicket(context.Context, *terminalproxyv1.IssueTicketRequest, ...grpc.CallOption) (*terminalproxyv1.IssueTicketResponse, error)
}

type TerminalGateway struct {
	terminalProxy terminalProxyClient
}

func NewTerminalGateway(terminalProxy terminalProxyClient) *TerminalGateway {
	return &TerminalGateway{terminalProxy: terminalProxy}
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

	// The kind and its parameters are forwarded as given. What a kind means —
	// the command it resolves to, which parameters it accepts, whether the
	// caller may attach — belongs to the Terminal Proxy, which owns terminal
	// sessions. The Gateway routes and does not interpret.
	ticket, err := g.terminalProxy.IssueTicket(ctx, &terminalproxyv1.IssueTicketRequest{
		IdentityId:    strings.TrimSpace(caller.IdentityID),
		WorkloadId:    workloadID,
		ContainerName: containerName,
		Kind:          req.Msg.GetKind(),
		Command:       req.Msg.GetCommand(),
		SyncRoot:      req.Msg.GetSyncRoot(),
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
