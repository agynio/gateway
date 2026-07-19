package gateway

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	appsv1 "github.com/agynio/gateway/gen/agynio/api/apps/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeAppsClient struct {
	getApp                func(ctx context.Context, req *appsv1.GetAppRequest, opts ...grpc.CallOption) (*appsv1.GetAppResponse, error)
	getInstallationBySlug func(ctx context.Context, req *appsv1.GetInstallationBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationBySlugResponse, error)
}

func (f *fakeAppsClient) CreateApp(ctx context.Context, req *appsv1.CreateAppRequest, opts ...grpc.CallOption) (*appsv1.CreateAppResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) UpdateApp(ctx context.Context, req *appsv1.UpdateAppRequest, opts ...grpc.CallOption) (*appsv1.UpdateAppResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) GetApp(ctx context.Context, req *appsv1.GetAppRequest, opts ...grpc.CallOption) (*appsv1.GetAppResponse, error) {
	if f.getApp != nil {
		return f.getApp(ctx, req, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) GetAppBySlug(ctx context.Context, req *appsv1.GetAppBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetAppBySlugResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) ListApps(ctx context.Context, req *appsv1.ListAppsRequest, opts ...grpc.CallOption) (*appsv1.ListAppsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) DeleteApp(ctx context.Context, req *appsv1.DeleteAppRequest, opts ...grpc.CallOption) (*appsv1.DeleteAppResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) GetAppProfile(ctx context.Context, req *appsv1.GetAppProfileRequest, opts ...grpc.CallOption) (*appsv1.GetAppProfileResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) ValidateServiceToken(ctx context.Context, req *appsv1.ValidateServiceTokenRequest, opts ...grpc.CallOption) (*appsv1.ValidateServiceTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) EnrollApp(ctx context.Context, req *appsv1.EnrollAppRequest, opts ...grpc.CallOption) (*appsv1.EnrollAppResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) InstallApp(ctx context.Context, req *appsv1.InstallAppRequest, opts ...grpc.CallOption) (*appsv1.InstallAppResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) GetInstallation(ctx context.Context, req *appsv1.GetInstallationRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) GetInstallationBySlug(ctx context.Context, req *appsv1.GetInstallationBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationBySlugResponse, error) {
	if f.getInstallationBySlug != nil {
		return f.getInstallationBySlug(ctx, req, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) GetInstallationByIdentityId(ctx context.Context, req *appsv1.GetInstallationByIdentityIdRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationByIdentityIdResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) ListInstallations(ctx context.Context, req *appsv1.ListInstallationsRequest, opts ...grpc.CallOption) (*appsv1.ListInstallationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) UpdateInstallation(ctx context.Context, req *appsv1.UpdateInstallationRequest, opts ...grpc.CallOption) (*appsv1.UpdateInstallationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) UninstallApp(ctx context.Context, req *appsv1.UninstallAppRequest, opts ...grpc.CallOption) (*appsv1.UninstallAppResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) GetInstallationConfiguration(ctx context.Context, req *appsv1.GetInstallationConfigurationRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationConfigurationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) ReportInstallationStatus(ctx context.Context, req *appsv1.ReportInstallationStatusRequest, opts ...grpc.CallOption) (*appsv1.ReportInstallationStatusResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) AppendInstallationAuditLogEntry(ctx context.Context, req *appsv1.AppendInstallationAuditLogEntryRequest, opts ...grpc.CallOption) (*appsv1.AppendInstallationAuditLogEntryResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) ListInstallationAuditLogEntries(ctx context.Context, req *appsv1.ListInstallationAuditLogEntriesRequest, opts ...grpc.CallOption) (*appsv1.ListInstallationAuditLogEntriesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

type fakeAgentsClient struct {
	resolveAgentIdentity func(ctx context.Context, req *agentsv1.ResolveAgentIdentityRequest, opts ...grpc.CallOption) (*agentsv1.ResolveAgentIdentityResponse, error)
}

func (f *fakeAgentsClient) CreateAgent(ctx context.Context, req *agentsv1.CreateAgentRequest, opts ...grpc.CallOption) (*agentsv1.CreateAgentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetAgent(ctx context.Context, req *agentsv1.GetAgentRequest, opts ...grpc.CallOption) (*agentsv1.GetAgentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ResolveAgentIdentity(ctx context.Context, req *agentsv1.ResolveAgentIdentityRequest, opts ...grpc.CallOption) (*agentsv1.ResolveAgentIdentityResponse, error) {
	if f.resolveAgentIdentity != nil {
		return f.resolveAgentIdentity(ctx, req, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) UpdateAgent(ctx context.Context, req *agentsv1.UpdateAgentRequest, opts ...grpc.CallOption) (*agentsv1.UpdateAgentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteAgent(ctx context.Context, req *agentsv1.DeleteAgentRequest, opts ...grpc.CallOption) (*agentsv1.DeleteAgentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListAgents(ctx context.Context, req *agentsv1.ListAgentsRequest, opts ...grpc.CallOption) (*agentsv1.ListAgentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) SetAgentRole(ctx context.Context, req *agentsv1.SetAgentRoleRequest, opts ...grpc.CallOption) (*agentsv1.SetAgentRoleResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) RemoveAgentRole(ctx context.Context, req *agentsv1.RemoveAgentRoleRequest, opts ...grpc.CallOption) (*agentsv1.RemoveAgentRoleResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListAgentRoles(ctx context.Context, req *agentsv1.ListAgentRolesRequest, opts ...grpc.CallOption) (*agentsv1.ListAgentRolesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListMyAgentRoles(ctx context.Context, req *agentsv1.ListMyAgentRolesRequest, opts ...grpc.CallOption) (*agentsv1.ListMyAgentRolesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateEnvironment(ctx context.Context, req *agentsv1.CreateEnvironmentRequest, opts ...grpc.CallOption) (*agentsv1.CreateEnvironmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetEnvironment(ctx context.Context, req *agentsv1.GetEnvironmentRequest, opts ...grpc.CallOption) (*agentsv1.GetEnvironmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) UpdateEnvironment(ctx context.Context, req *agentsv1.UpdateEnvironmentRequest, opts ...grpc.CallOption) (*agentsv1.UpdateEnvironmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteEnvironment(ctx context.Context, req *agentsv1.DeleteEnvironmentRequest, opts ...grpc.CallOption) (*agentsv1.DeleteEnvironmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListEnvironments(ctx context.Context, req *agentsv1.ListEnvironmentsRequest, opts ...grpc.CallOption) (*agentsv1.ListEnvironmentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateSandbox(ctx context.Context, req *agentsv1.CreateSandboxRequest, opts ...grpc.CallOption) (*agentsv1.CreateSandboxResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetSandbox(ctx context.Context, req *agentsv1.GetSandboxRequest, opts ...grpc.CallOption) (*agentsv1.GetSandboxResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListSandboxes(ctx context.Context, req *agentsv1.ListSandboxesRequest, opts ...grpc.CallOption) (*agentsv1.ListSandboxesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) StopSandbox(ctx context.Context, req *agentsv1.StopSandboxRequest, opts ...grpc.CallOption) (*agentsv1.StopSandboxResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteSandbox(ctx context.Context, req *agentsv1.DeleteSandboxRequest, opts ...grpc.CallOption) (*agentsv1.DeleteSandboxResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) EnsureSandboxRunning(ctx context.Context, req *agentsv1.EnsureSandboxRunningRequest, opts ...grpc.CallOption) (*agentsv1.EnsureSandboxRunningResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) UpdateSandboxLastSession(ctx context.Context, req *agentsv1.UpdateSandboxLastSessionRequest, opts ...grpc.CallOption) (*agentsv1.UpdateSandboxLastSessionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateInstance(ctx context.Context, req *agentsv1.CreateInstanceRequest, opts ...grpc.CallOption) (*agentsv1.CreateInstanceResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetInstance(ctx context.Context, req *agentsv1.GetInstanceRequest, opts ...grpc.CallOption) (*agentsv1.GetInstanceResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListInstances(ctx context.Context, req *agentsv1.ListInstancesRequest, opts ...grpc.CallOption) (*agentsv1.ListInstancesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) PauseInstance(ctx context.Context, req *agentsv1.PauseInstanceRequest, opts ...grpc.CallOption) (*agentsv1.PauseInstanceResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ResumeInstance(ctx context.Context, req *agentsv1.ResumeInstanceRequest, opts ...grpc.CallOption) (*agentsv1.ResumeInstanceResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteInstance(ctx context.Context, req *agentsv1.DeleteInstanceRequest, opts ...grpc.CallOption) (*agentsv1.DeleteInstanceResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) WriteInboxItem(ctx context.Context, req *agentsv1.WriteInboxItemRequest, opts ...grpc.CallOption) (*agentsv1.WriteInboxItemResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) FanoutInboxItem(ctx context.Context, req *agentsv1.FanoutInboxItemRequest, opts ...grpc.CallOption) (*agentsv1.FanoutInboxItemResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetUnackedInboxItems(ctx context.Context, req *agentsv1.GetUnackedInboxItemsRequest, opts ...grpc.CallOption) (*agentsv1.GetUnackedInboxItemsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) AckInboxItems(ctx context.Context, req *agentsv1.AckInboxItemsRequest, opts ...grpc.CallOption) (*agentsv1.AckInboxItemsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetUnackedInboxCount(ctx context.Context, req *agentsv1.GetUnackedInboxCountRequest, opts ...grpc.CallOption) (*agentsv1.GetUnackedInboxCountResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateVolume(ctx context.Context, req *agentsv1.CreateVolumeRequest, opts ...grpc.CallOption) (*agentsv1.CreateVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetVolume(ctx context.Context, req *agentsv1.GetVolumeRequest, opts ...grpc.CallOption) (*agentsv1.GetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) UpdateVolume(ctx context.Context, req *agentsv1.UpdateVolumeRequest, opts ...grpc.CallOption) (*agentsv1.UpdateVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteVolume(ctx context.Context, req *agentsv1.DeleteVolumeRequest, opts ...grpc.CallOption) (*agentsv1.DeleteVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListVolumes(ctx context.Context, req *agentsv1.ListVolumesRequest, opts ...grpc.CallOption) (*agentsv1.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateVolumeAttachment(ctx context.Context, req *agentsv1.CreateVolumeAttachmentRequest, opts ...grpc.CallOption) (*agentsv1.CreateVolumeAttachmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetVolumeAttachment(ctx context.Context, req *agentsv1.GetVolumeAttachmentRequest, opts ...grpc.CallOption) (*agentsv1.GetVolumeAttachmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteVolumeAttachment(ctx context.Context, req *agentsv1.DeleteVolumeAttachmentRequest, opts ...grpc.CallOption) (*agentsv1.DeleteVolumeAttachmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListVolumeAttachments(ctx context.Context, req *agentsv1.ListVolumeAttachmentsRequest, opts ...grpc.CallOption) (*agentsv1.ListVolumeAttachmentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateMcp(ctx context.Context, req *agentsv1.CreateMcpRequest, opts ...grpc.CallOption) (*agentsv1.CreateMcpResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetMcp(ctx context.Context, req *agentsv1.GetMcpRequest, opts ...grpc.CallOption) (*agentsv1.GetMcpResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) UpdateMcp(ctx context.Context, req *agentsv1.UpdateMcpRequest, opts ...grpc.CallOption) (*agentsv1.UpdateMcpResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteMcp(ctx context.Context, req *agentsv1.DeleteMcpRequest, opts ...grpc.CallOption) (*agentsv1.DeleteMcpResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListMcps(ctx context.Context, req *agentsv1.ListMcpsRequest, opts ...grpc.CallOption) (*agentsv1.ListMcpsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateSkill(ctx context.Context, req *agentsv1.CreateSkillRequest, opts ...grpc.CallOption) (*agentsv1.CreateSkillResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetSkill(ctx context.Context, req *agentsv1.GetSkillRequest, opts ...grpc.CallOption) (*agentsv1.GetSkillResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) UpdateSkill(ctx context.Context, req *agentsv1.UpdateSkillRequest, opts ...grpc.CallOption) (*agentsv1.UpdateSkillResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteSkill(ctx context.Context, req *agentsv1.DeleteSkillRequest, opts ...grpc.CallOption) (*agentsv1.DeleteSkillResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListSkills(ctx context.Context, req *agentsv1.ListSkillsRequest, opts ...grpc.CallOption) (*agentsv1.ListSkillsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateHook(ctx context.Context, req *agentsv1.CreateHookRequest, opts ...grpc.CallOption) (*agentsv1.CreateHookResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetHook(ctx context.Context, req *agentsv1.GetHookRequest, opts ...grpc.CallOption) (*agentsv1.GetHookResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) UpdateHook(ctx context.Context, req *agentsv1.UpdateHookRequest, opts ...grpc.CallOption) (*agentsv1.UpdateHookResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteHook(ctx context.Context, req *agentsv1.DeleteHookRequest, opts ...grpc.CallOption) (*agentsv1.DeleteHookResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListHooks(ctx context.Context, req *agentsv1.ListHooksRequest, opts ...grpc.CallOption) (*agentsv1.ListHooksResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateEnv(ctx context.Context, req *agentsv1.CreateEnvRequest, opts ...grpc.CallOption) (*agentsv1.CreateEnvResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetEnv(ctx context.Context, req *agentsv1.GetEnvRequest, opts ...grpc.CallOption) (*agentsv1.GetEnvResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) UpdateEnv(ctx context.Context, req *agentsv1.UpdateEnvRequest, opts ...grpc.CallOption) (*agentsv1.UpdateEnvResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteEnv(ctx context.Context, req *agentsv1.DeleteEnvRequest, opts ...grpc.CallOption) (*agentsv1.DeleteEnvResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListEnvs(ctx context.Context, req *agentsv1.ListEnvsRequest, opts ...grpc.CallOption) (*agentsv1.ListEnvsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateInitScript(ctx context.Context, req *agentsv1.CreateInitScriptRequest, opts ...grpc.CallOption) (*agentsv1.CreateInitScriptResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetInitScript(ctx context.Context, req *agentsv1.GetInitScriptRequest, opts ...grpc.CallOption) (*agentsv1.GetInitScriptResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) UpdateInitScript(ctx context.Context, req *agentsv1.UpdateInitScriptRequest, opts ...grpc.CallOption) (*agentsv1.UpdateInitScriptResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteInitScript(ctx context.Context, req *agentsv1.DeleteInitScriptRequest, opts ...grpc.CallOption) (*agentsv1.DeleteInitScriptResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListInitScripts(ctx context.Context, req *agentsv1.ListInitScriptsRequest, opts ...grpc.CallOption) (*agentsv1.ListInitScriptsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) CreateImagePullSecretAttachment(ctx context.Context, req *agentsv1.CreateImagePullSecretAttachmentRequest, opts ...grpc.CallOption) (*agentsv1.CreateImagePullSecretAttachmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) GetImagePullSecretAttachment(ctx context.Context, req *agentsv1.GetImagePullSecretAttachmentRequest, opts ...grpc.CallOption) (*agentsv1.GetImagePullSecretAttachmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) DeleteImagePullSecretAttachment(ctx context.Context, req *agentsv1.DeleteImagePullSecretAttachmentRequest, opts ...grpc.CallOption) (*agentsv1.DeleteImagePullSecretAttachmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAgentsClient) ListImagePullSecretAttachments(ctx context.Context, req *agentsv1.ListImagePullSecretAttachmentsRequest, opts ...grpc.CallOption) (*agentsv1.ListImagePullSecretAttachmentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func TestParseAppProxyPath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		expSlug   string
		expMethod string
		wantErr   bool
	}{
		{name: "simple", path: "/apps/my-app/run", expSlug: "my-app", expMethod: "run"},
		{name: "nested", path: "/apps/my-app/run/task", expSlug: "my-app", expMethod: "run/task"},
		{name: "invalid prefix", path: "/app/my-app/run", wantErr: true},
		{name: "missing method", path: "/apps/my-app", wantErr: true},
		{name: "missing slug", path: "/apps//run", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, method, err := parseAppProxyPath(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if slug != tt.expSlug {
				t.Fatalf("expected slug %q, got %q", tt.expSlug, slug)
			}
			if method != tt.expMethod {
				t.Fatalf("expected method %q, got %q", tt.expMethod, method)
			}
		})
	}
}

func TestResolveInstallationCache(t *testing.T) {
	var installationCalls int
	var appCalls int
	appsClient := &fakeAppsClient{
		getInstallationBySlug: func(ctx context.Context, req *appsv1.GetInstallationBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationBySlugResponse, error) {
			installationCalls++
			return &appsv1.GetInstallationBySlugResponse{Installation: &appsv1.Installation{
				Meta:  &appsv1.EntityMeta{Id: "inst-" + req.OrganizationId},
				AppId: "app-" + req.OrganizationId,
			}}, nil
		},
		getApp: func(ctx context.Context, req *appsv1.GetAppRequest, opts ...grpc.CallOption) (*appsv1.GetAppResponse, error) {
			appCalls++
			return &appsv1.GetAppResponse{App: &appsv1.App{Slug: "slug-" + req.Id}}, nil
		},
	}
	handler := &AppProxyHandler{
		apps:     appsClient,
		cache:    make(map[string]cachedInstallation),
		cacheTTL: time.Minute,
	}

	ctx := context.Background()
	resolved, err := handler.resolveInstallation(ctx, "org-1", "slug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.serviceName != "app-slug-app-org-1" {
		t.Fatalf("expected service name %q, got %q", "app-slug-app-org-1", resolved.serviceName)
	}
	if resolved.installationID != "inst-org-1" {
		t.Fatalf("expected installation id %q, got %q", "inst-org-1", resolved.installationID)
	}

	resolved, err = handler.resolveInstallation(ctx, "org-1", "slug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.serviceName != "app-slug-app-org-1" {
		t.Fatalf("expected service name %q, got %q", "app-slug-app-org-1", resolved.serviceName)
	}
	if installationCalls != 1 {
		t.Fatalf("expected installation lookup once, got %d", installationCalls)
	}
	if appCalls != 1 {
		t.Fatalf("expected app lookup once, got %d", appCalls)
	}

	t.Run("different organization", func(t *testing.T) {
		resolved, err = handler.resolveInstallation(ctx, "org-2", "slug")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resolved.serviceName != "app-slug-app-org-2" {
			t.Fatalf("expected service name %q, got %q", "app-slug-app-org-2", resolved.serviceName)
		}
		if resolved.installationID != "inst-org-2" {
			t.Fatalf("expected installation id %q, got %q", "inst-org-2", resolved.installationID)
		}
		if installationCalls != 2 {
			t.Fatalf("expected installation lookup twice, got %d", installationCalls)
		}
		if appCalls != 2 {
			t.Fatalf("expected app lookup twice, got %d", appCalls)
		}
	})
}

func TestResolveInstallationErrors(t *testing.T) {
	t.Run("installation lookup error", func(t *testing.T) {
		handler := &AppProxyHandler{
			apps: &fakeAppsClient{
				getInstallationBySlug: func(ctx context.Context, req *appsv1.GetInstallationBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationBySlugResponse, error) {
					return nil, status.Error(codes.NotFound, "missing")
				},
			},
			cache:    make(map[string]cachedInstallation),
			cacheTTL: time.Minute,
		}

		_, err := handler.resolveInstallation(context.Background(), "org", "slug")
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("missing installation", func(t *testing.T) {
		handler := &AppProxyHandler{
			apps: &fakeAppsClient{
				getInstallationBySlug: func(ctx context.Context, req *appsv1.GetInstallationBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationBySlugResponse, error) {
					return &appsv1.GetInstallationBySlugResponse{}, nil
				},
			},
			cache:    make(map[string]cachedInstallation),
			cacheTTL: time.Minute,
		}

		_, err := handler.resolveInstallation(context.Background(), "org", "slug")
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("get app error", func(t *testing.T) {
		handler := &AppProxyHandler{
			apps: &fakeAppsClient{
				getInstallationBySlug: func(ctx context.Context, req *appsv1.GetInstallationBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationBySlugResponse, error) {
					return &appsv1.GetInstallationBySlugResponse{Installation: &appsv1.Installation{
						Meta:  &appsv1.EntityMeta{Id: "inst"},
						AppId: "app",
					}}, nil
				},
				getApp: func(ctx context.Context, req *appsv1.GetAppRequest, opts ...grpc.CallOption) (*appsv1.GetAppResponse, error) {
					return nil, status.Error(codes.Internal, "boom")
				},
			},
			cache:    make(map[string]cachedInstallation),
			cacheTTL: time.Minute,
		}

		_, err := handler.resolveInstallation(context.Background(), "org", "slug")
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("missing slug", func(t *testing.T) {
		handler := &AppProxyHandler{
			apps: &fakeAppsClient{
				getInstallationBySlug: func(ctx context.Context, req *appsv1.GetInstallationBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationBySlugResponse, error) {
					return &appsv1.GetInstallationBySlugResponse{Installation: &appsv1.Installation{
						Meta:  &appsv1.EntityMeta{Id: "inst"},
						AppId: "app",
					}}, nil
				},
				getApp: func(ctx context.Context, req *appsv1.GetAppRequest, opts ...grpc.CallOption) (*appsv1.GetAppResponse, error) {
					return &appsv1.GetAppResponse{App: &appsv1.App{Slug: " "}}, nil
				},
			},
			cache:    make(map[string]cachedInstallation),
			cacheTTL: time.Minute,
		}

		_, err := handler.resolveInstallation(context.Background(), "org", "slug")
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestAppProxyHandlerServeHTTP(t *testing.T) {
	var gotIdentityID string
	var gotIdentityType string
	var gotInstallationID string
	var gotPath string
	var gotOrganizationID string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotIdentityID = r.Header.Get(identity.MetadataKeyIdentityID)
		gotIdentityType = r.Header.Get(identity.MetadataKeyIdentityType)
		gotInstallationID = r.Header.Get(identity.MetadataKeyAppInstallationID)
		gotOrganizationID = r.Header.Get(identity.MetadataKeyOrganizationID)
		gotPath = r.URL.Path
		w.Header().Set("X-Test", "ok")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("response"))
	}))
	defer server.Close()

	serverAddr := server.Listener.Addr().String()
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			_ = network
			_ = addr
			return (&net.Dialer{}).DialContext(ctx, "tcp", serverAddr)
		},
	}
	agentsClient := &fakeAgentsClient{
		resolveAgentIdentity: func(ctx context.Context, req *agentsv1.ResolveAgentIdentityRequest, opts ...grpc.CallOption) (*agentsv1.ResolveAgentIdentityResponse, error) {
			return &agentsv1.ResolveAgentIdentityResponse{OrganizationId: "org-1"}, nil
		},
	}
	appsClient := &fakeAppsClient{
		getInstallationBySlug: func(ctx context.Context, req *appsv1.GetInstallationBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationBySlugResponse, error) {
			return &appsv1.GetInstallationBySlugResponse{Installation: &appsv1.Installation{
				Meta:  &appsv1.EntityMeta{Id: "inst-1"},
				AppId: "app-1",
			}}, nil
		},
		getApp: func(ctx context.Context, req *appsv1.GetAppRequest, opts ...grpc.CallOption) (*appsv1.GetAppResponse, error) {
			return &appsv1.GetAppResponse{App: &appsv1.App{Slug: "app-slug"}}, nil
		},
	}
	handler := &AppProxyHandler{
		apps:     appsClient,
		agents:   agentsClient,
		client:   &http.Client{Transport: transport},
		cache:    make(map[string]cachedInstallation),
		cacheTTL: time.Minute,
	}

	resolved := identity.ResolvedIdentity{IdentityID: "identity-1", IdentityType: identity.IdentityTypeAgent}
	req := httptest.NewRequest(http.MethodGet, "/apps/app-slug/health", nil)
	req = req.WithContext(identity.WithIdentity(req.Context(), resolved))
	resp := httptest.NewRecorder()
	req.Header.Set(identity.MetadataKeyOrganizationID, "org-override")

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if resp.Header().Get("X-Test") != "ok" {
		t.Fatalf("expected response headers to be forwarded")
	}
	if resp.Body.String() != "response" {
		t.Fatalf("expected response body to be forwarded")
	}
	if gotIdentityID != resolved.IdentityID {
		t.Fatalf("expected identity header %q, got %q", resolved.IdentityID, gotIdentityID)
	}
	if gotIdentityType != string(resolved.IdentityType) {
		t.Fatalf("expected identity type header %q, got %q", resolved.IdentityType, gotIdentityType)
	}
	if gotInstallationID != "inst-1" {
		t.Fatalf("expected installation header %q, got %q", "inst-1", gotInstallationID)
	}
	if gotOrganizationID != "" {
		t.Fatalf("expected organization header stripped, got %q", gotOrganizationID)
	}
	if gotPath != "/health" {
		t.Fatalf("expected proxied path %q, got %q", "/health", gotPath)
	}
}

func TestAppProxyHandlerServeHTTPMissingOrgID(t *testing.T) {
	appsClient := &fakeAppsClient{}
	handler := &AppProxyHandler{apps: appsClient}
	req := httptest.NewRequest(http.MethodGet, "/apps/app-slug/health", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func TestAppProxyHandlerServeHTTPUserIdentityUnsupported(t *testing.T) {
	appsClient := &fakeAppsClient{}
	agentsClient := &fakeAgentsClient{}
	handler := &AppProxyHandler{apps: appsClient, agents: agentsClient}
	req := httptest.NewRequest(http.MethodGet, "/apps/app-slug/health", nil)
	req = req.WithContext(identity.WithIdentity(req.Context(), identity.ResolvedIdentity{IdentityID: "user-1", IdentityType: identity.IdentityTypeUser}))
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusPreconditionFailed {
		t.Fatalf("expected status %d, got %d", http.StatusPreconditionFailed, resp.Code)
	}
}

func TestAppProxyHandlerServeHTTPInvalidPath(t *testing.T) {
	handler := &AppProxyHandler{apps: &fakeAppsClient{}}
	req := httptest.NewRequest(http.MethodGet, "/apps/app-slug", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}
	if resp.Header().Get("Content-Type") != problemContentType {
		t.Fatalf("expected problem content type")
	}
}
