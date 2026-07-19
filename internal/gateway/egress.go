package gateway

import (
	"context"

	"connectrpc.com/connect"
	egressv1 "github.com/agynio/gateway/gen/agynio/api/egress/v1"
	"google.golang.org/grpc"
)

type EgressRulesGateway struct {
	egressRules egressRulesClient
}

type egressRulesClient interface {
	CreateEgressRule(context.Context, *egressv1.CreateEgressRuleRequest, ...grpc.CallOption) (*egressv1.CreateEgressRuleResponse, error)
	GetEgressRule(context.Context, *egressv1.GetEgressRuleRequest, ...grpc.CallOption) (*egressv1.GetEgressRuleResponse, error)
	ListEgressRules(context.Context, *egressv1.ListEgressRulesRequest, ...grpc.CallOption) (*egressv1.ListEgressRulesResponse, error)
	UpdateEgressRule(context.Context, *egressv1.UpdateEgressRuleRequest, ...grpc.CallOption) (*egressv1.UpdateEgressRuleResponse, error)
	DeleteEgressRule(context.Context, *egressv1.DeleteEgressRuleRequest, ...grpc.CallOption) (*egressv1.DeleteEgressRuleResponse, error)
	CreateEgressRuleAttachment(context.Context, *egressv1.CreateEgressRuleAttachmentRequest, ...grpc.CallOption) (*egressv1.CreateEgressRuleAttachmentResponse, error)
	DeleteEgressRuleAttachment(context.Context, *egressv1.DeleteEgressRuleAttachmentRequest, ...grpc.CallOption) (*egressv1.DeleteEgressRuleAttachmentResponse, error)
	ListEgressRuleAttachments(context.Context, *egressv1.ListEgressRuleAttachmentsRequest, ...grpc.CallOption) (*egressv1.ListEgressRuleAttachmentsResponse, error)
}

func NewEgressRulesGateway(egressRules egressRulesClient) *EgressRulesGateway {
	return &EgressRulesGateway{egressRules: egressRules}
}

func (g *EgressRulesGateway) CreateEgressRule(ctx context.Context, req *connect.Request[egressv1.CreateEgressRuleRequest]) (*connect.Response[egressv1.CreateEgressRuleResponse], error) {
	resp, err := g.egressRules.CreateEgressRule(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *EgressRulesGateway) GetEgressRule(ctx context.Context, req *connect.Request[egressv1.GetEgressRuleRequest]) (*connect.Response[egressv1.GetEgressRuleResponse], error) {
	resp, err := g.egressRules.GetEgressRule(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *EgressRulesGateway) ListEgressRules(ctx context.Context, req *connect.Request[egressv1.ListEgressRulesRequest]) (*connect.Response[egressv1.ListEgressRulesResponse], error) {
	resp, err := g.egressRules.ListEgressRules(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *EgressRulesGateway) UpdateEgressRule(ctx context.Context, req *connect.Request[egressv1.UpdateEgressRuleRequest]) (*connect.Response[egressv1.UpdateEgressRuleResponse], error) {
	resp, err := g.egressRules.UpdateEgressRule(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *EgressRulesGateway) DeleteEgressRule(ctx context.Context, req *connect.Request[egressv1.DeleteEgressRuleRequest]) (*connect.Response[egressv1.DeleteEgressRuleResponse], error) {
	resp, err := g.egressRules.DeleteEgressRule(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *EgressRulesGateway) CreateEgressRuleAttachment(ctx context.Context, req *connect.Request[egressv1.CreateEgressRuleAttachmentRequest]) (*connect.Response[egressv1.CreateEgressRuleAttachmentResponse], error) {
	resp, err := g.egressRules.CreateEgressRuleAttachment(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *EgressRulesGateway) DeleteEgressRuleAttachment(ctx context.Context, req *connect.Request[egressv1.DeleteEgressRuleAttachmentRequest]) (*connect.Response[egressv1.DeleteEgressRuleAttachmentResponse], error) {
	resp, err := g.egressRules.DeleteEgressRuleAttachment(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *EgressRulesGateway) ListEgressRuleAttachments(ctx context.Context, req *connect.Request[egressv1.ListEgressRuleAttachmentsRequest]) (*connect.Response[egressv1.ListEgressRuleAttachmentsResponse], error) {
	resp, err := g.egressRules.ListEgressRuleAttachments(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
