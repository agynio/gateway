package gateway

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	egressv1 "github.com/agynio/gateway/gen/agynio/api/egress/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type fakeEgressRulesClient struct {
	createRuleReq                    *egressv1.CreateEgressRuleRequest
	createRuleResp                   *egressv1.CreateEgressRuleResponse
	createRuleErr                    error
	createRuleMetadata               metadata.MD
	getRuleReq                       *egressv1.GetEgressRuleRequest
	getRuleResp                      *egressv1.GetEgressRuleResponse
	getRuleErr                       error
	listRulesReq                     *egressv1.ListEgressRulesRequest
	listRulesResp                    *egressv1.ListEgressRulesResponse
	listRulesErr                     error
	updateRuleReq                    *egressv1.UpdateEgressRuleRequest
	updateRuleResp                   *egressv1.UpdateEgressRuleResponse
	updateRuleErr                    error
	deleteRuleReq                    *egressv1.DeleteEgressRuleRequest
	deleteRuleResp                   *egressv1.DeleteEgressRuleResponse
	deleteRuleErr                    error
	createAttachmentReq              *egressv1.CreateEgressRuleAttachmentRequest
	createAttachmentResp             *egressv1.CreateEgressRuleAttachmentResponse
	createAttachmentErr              error
	deleteAttachmentReq              *egressv1.DeleteEgressRuleAttachmentRequest
	deleteAttachmentResp             *egressv1.DeleteEgressRuleAttachmentResponse
	deleteAttachmentErr              error
	listAttachmentsReq               *egressv1.ListEgressRuleAttachmentsRequest
	listAttachmentsResp              *egressv1.ListEgressRuleAttachmentsResponse
	listAttachmentsErr               error
	listRulesByAgentReq              *egressv1.ListEgressRulesByAgentRequest
	listRulesByAgentResp             *egressv1.ListEgressRulesByAgentResponse
	listRulesByAgentErr              error
	countRulesReferencingSecretReq   *egressv1.CountRulesReferencingSecretRequest
	countRulesReferencingSecretResp  *egressv1.CountRulesReferencingSecretResponse
	countRulesReferencingSecretError error
}

func (f *fakeEgressRulesClient) CreateEgressRule(ctx context.Context, in *egressv1.CreateEgressRuleRequest, opts ...grpc.CallOption) (*egressv1.CreateEgressRuleResponse, error) {
	f.createRuleReq = in
	f.createRuleMetadata, _ = metadata.FromOutgoingContext(ctx)
	if f.createRuleErr != nil {
		return nil, f.createRuleErr
	}
	if f.createRuleResp == nil {
		f.createRuleResp = &egressv1.CreateEgressRuleResponse{}
	}
	return f.createRuleResp, nil
}

func (f *fakeEgressRulesClient) GetEgressRule(ctx context.Context, in *egressv1.GetEgressRuleRequest, opts ...grpc.CallOption) (*egressv1.GetEgressRuleResponse, error) {
	f.getRuleReq = in
	if f.getRuleErr != nil {
		return nil, f.getRuleErr
	}
	if f.getRuleResp == nil {
		f.getRuleResp = &egressv1.GetEgressRuleResponse{}
	}
	return f.getRuleResp, nil
}

func (f *fakeEgressRulesClient) ListEgressRules(ctx context.Context, in *egressv1.ListEgressRulesRequest, opts ...grpc.CallOption) (*egressv1.ListEgressRulesResponse, error) {
	f.listRulesReq = in
	if f.listRulesErr != nil {
		return nil, f.listRulesErr
	}
	if f.listRulesResp == nil {
		f.listRulesResp = &egressv1.ListEgressRulesResponse{}
	}
	return f.listRulesResp, nil
}

func (f *fakeEgressRulesClient) UpdateEgressRule(ctx context.Context, in *egressv1.UpdateEgressRuleRequest, opts ...grpc.CallOption) (*egressv1.UpdateEgressRuleResponse, error) {
	f.updateRuleReq = in
	if f.updateRuleErr != nil {
		return nil, f.updateRuleErr
	}
	if f.updateRuleResp == nil {
		f.updateRuleResp = &egressv1.UpdateEgressRuleResponse{}
	}
	return f.updateRuleResp, nil
}

func (f *fakeEgressRulesClient) DeleteEgressRule(ctx context.Context, in *egressv1.DeleteEgressRuleRequest, opts ...grpc.CallOption) (*egressv1.DeleteEgressRuleResponse, error) {
	f.deleteRuleReq = in
	if f.deleteRuleErr != nil {
		return nil, f.deleteRuleErr
	}
	if f.deleteRuleResp == nil {
		f.deleteRuleResp = &egressv1.DeleteEgressRuleResponse{}
	}
	return f.deleteRuleResp, nil
}

func (f *fakeEgressRulesClient) CreateEgressRuleAttachment(ctx context.Context, in *egressv1.CreateEgressRuleAttachmentRequest, opts ...grpc.CallOption) (*egressv1.CreateEgressRuleAttachmentResponse, error) {
	f.createAttachmentReq = in
	if f.createAttachmentErr != nil {
		return nil, f.createAttachmentErr
	}
	if f.createAttachmentResp == nil {
		f.createAttachmentResp = &egressv1.CreateEgressRuleAttachmentResponse{}
	}
	return f.createAttachmentResp, nil
}

func (f *fakeEgressRulesClient) DeleteEgressRuleAttachment(ctx context.Context, in *egressv1.DeleteEgressRuleAttachmentRequest, opts ...grpc.CallOption) (*egressv1.DeleteEgressRuleAttachmentResponse, error) {
	f.deleteAttachmentReq = in
	if f.deleteAttachmentErr != nil {
		return nil, f.deleteAttachmentErr
	}
	if f.deleteAttachmentResp == nil {
		f.deleteAttachmentResp = &egressv1.DeleteEgressRuleAttachmentResponse{}
	}
	return f.deleteAttachmentResp, nil
}

func (f *fakeEgressRulesClient) ListEgressRuleAttachments(ctx context.Context, in *egressv1.ListEgressRuleAttachmentsRequest, opts ...grpc.CallOption) (*egressv1.ListEgressRuleAttachmentsResponse, error) {
	f.listAttachmentsReq = in
	if f.listAttachmentsErr != nil {
		return nil, f.listAttachmentsErr
	}
	if f.listAttachmentsResp == nil {
		f.listAttachmentsResp = &egressv1.ListEgressRuleAttachmentsResponse{}
	}
	return f.listAttachmentsResp, nil
}

func (f *fakeEgressRulesClient) ListEgressRulesByAgent(ctx context.Context, in *egressv1.ListEgressRulesByAgentRequest, opts ...grpc.CallOption) (*egressv1.ListEgressRulesByAgentResponse, error) {
	f.listRulesByAgentReq = in
	if f.listRulesByAgentErr != nil {
		return nil, f.listRulesByAgentErr
	}
	if f.listRulesByAgentResp == nil {
		f.listRulesByAgentResp = &egressv1.ListEgressRulesByAgentResponse{}
	}
	return f.listRulesByAgentResp, nil
}

func (f *fakeEgressRulesClient) CountRulesReferencingSecret(ctx context.Context, in *egressv1.CountRulesReferencingSecretRequest, opts ...grpc.CallOption) (*egressv1.CountRulesReferencingSecretResponse, error) {
	f.countRulesReferencingSecretReq = in
	if f.countRulesReferencingSecretError != nil {
		return nil, f.countRulesReferencingSecretError
	}
	if f.countRulesReferencingSecretResp == nil {
		f.countRulesReferencingSecretResp = &egressv1.CountRulesReferencingSecretResponse{}
	}
	return f.countRulesReferencingSecretResp, nil
}

func TestEgressRulesGatewayCreateEgressRuleForwardsRequest(t *testing.T) {
	client := &fakeEgressRulesClient{createRuleResp: &egressv1.CreateEgressRuleResponse{}}
	gateway := NewEgressRulesGateway(client)

	req := connect.NewRequest(&egressv1.CreateEgressRuleRequest{
		OrganizationId: "org-1",
		Name:           "github",
		Matcher:        &egressv1.EgressRuleMatcher{DomainPattern: "*.github.com"},
	})
	resp, err := gateway.CreateEgressRule(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Msg != client.createRuleResp {
		t.Fatalf("expected response to be forwarded")
	}
	if client.createRuleReq.GetOrganizationId() != "org-1" {
		t.Fatalf("expected organization id to be forwarded")
	}
	if client.createRuleReq.GetMatcher().GetDomainPattern() != "*.github.com" {
		t.Fatalf("expected matcher to be forwarded")
	}
}

func TestEgressRulesGatewayCreateEgressRulePropagatesIdentityMetadata(t *testing.T) {
	client := &fakeEgressRulesClient{createRuleResp: &egressv1.CreateEgressRuleResponse{}}
	gateway := NewEgressRulesGateway(client)
	resolved := identity.ResolvedIdentity{IdentityID: "user-1", IdentityType: identity.IdentityTypeUser}
	ctx := identity.WithIdentity(context.Background(), resolved)

	_, err := gateway.CreateEgressRule(ctx, connect.NewRequest(&egressv1.CreateEgressRuleRequest{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertMetadataValue(t, client.createRuleMetadata, identity.MetadataKeyIdentityID, resolved.IdentityID)
	assertMetadataValue(t, client.createRuleMetadata, identity.MetadataKeyIdentityType, string(resolved.IdentityType))
}

func TestEgressRulesGatewayForwardsAllPublicMethods(t *testing.T) {
	client := &fakeEgressRulesClient{}
	gateway := NewEgressRulesGateway(client)

	if _, err := gateway.GetEgressRule(context.Background(), connect.NewRequest(&egressv1.GetEgressRuleRequest{Id: "rule-1"})); err != nil {
		t.Fatalf("get egress rule: %v", err)
	}
	if client.getRuleReq.GetId() != "rule-1" {
		t.Fatalf("expected get request to be forwarded")
	}

	if _, err := gateway.ListEgressRules(context.Background(), connect.NewRequest(&egressv1.ListEgressRulesRequest{OrganizationId: "org-1"})); err != nil {
		t.Fatalf("list egress rules: %v", err)
	}
	if client.listRulesReq.GetOrganizationId() != "org-1" {
		t.Fatalf("expected list request to be forwarded")
	}

	if _, err := gateway.UpdateEgressRule(context.Background(), connect.NewRequest(&egressv1.UpdateEgressRuleRequest{Id: "rule-2"})); err != nil {
		t.Fatalf("update egress rule: %v", err)
	}
	if client.updateRuleReq.GetId() != "rule-2" {
		t.Fatalf("expected update request to be forwarded")
	}

	if _, err := gateway.DeleteEgressRule(context.Background(), connect.NewRequest(&egressv1.DeleteEgressRuleRequest{Id: "rule-3"})); err != nil {
		t.Fatalf("delete egress rule: %v", err)
	}
	if client.deleteRuleReq.GetId() != "rule-3" {
		t.Fatalf("expected delete request to be forwarded")
	}

	if _, err := gateway.CreateEgressRuleAttachment(context.Background(), connect.NewRequest(&egressv1.CreateEgressRuleAttachmentRequest{RuleId: "rule-4", AgentId: "agent-1"})); err != nil {
		t.Fatalf("create egress rule attachment: %v", err)
	}
	if client.createAttachmentReq.GetRuleId() != "rule-4" || client.createAttachmentReq.GetAgentId() != "agent-1" {
		t.Fatalf("expected create attachment request to be forwarded")
	}

	if _, err := gateway.DeleteEgressRuleAttachment(context.Background(), connect.NewRequest(&egressv1.DeleteEgressRuleAttachmentRequest{Id: "attachment-1"})); err != nil {
		t.Fatalf("delete egress rule attachment: %v", err)
	}
	if client.deleteAttachmentReq.GetId() != "attachment-1" {
		t.Fatalf("expected delete attachment request to be forwarded")
	}

	if _, err := gateway.ListEgressRuleAttachments(context.Background(), connect.NewRequest(&egressv1.ListEgressRuleAttachmentsRequest{OrganizationId: "org-2"})); err != nil {
		t.Fatalf("list egress rule attachments: %v", err)
	}
	if client.listAttachmentsReq.GetOrganizationId() != "org-2" {
		t.Fatalf("expected list attachments request to be forwarded")
	}
}
