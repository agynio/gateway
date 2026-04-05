package gateway

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
			return &appsv1.GetAppResponse{App: &appsv1.App{ZitiServiceId: "service-" + req.Id}}, nil
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
	if resolved.serviceName != "service-app-org-1" {
		t.Fatalf("expected service name %q, got %q", "service-app-org-1", resolved.serviceName)
	}
	if resolved.installationID != "inst-org-1" {
		t.Fatalf("expected installation id %q, got %q", "inst-org-1", resolved.installationID)
	}

	resolved, err = handler.resolveInstallation(ctx, "org-1", "slug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.serviceName != "service-app-org-1" {
		t.Fatalf("expected service name %q, got %q", "service-app-org-1", resolved.serviceName)
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
		if resolved.serviceName != "service-app-org-2" {
			t.Fatalf("expected service name %q, got %q", "service-app-org-2", resolved.serviceName)
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

	t.Run("missing service id", func(t *testing.T) {
		handler := &AppProxyHandler{
			apps: &fakeAppsClient{
				getInstallationBySlug: func(ctx context.Context, req *appsv1.GetInstallationBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationBySlugResponse, error) {
					return &appsv1.GetInstallationBySlugResponse{Installation: &appsv1.Installation{
						Meta:  &appsv1.EntityMeta{Id: "inst"},
						AppId: "app",
					}}, nil
				},
				getApp: func(ctx context.Context, req *appsv1.GetAppRequest, opts ...grpc.CallOption) (*appsv1.GetAppResponse, error) {
					return &appsv1.GetAppResponse{App: &appsv1.App{ZitiServiceId: " "}}, nil
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotIdentityID = r.Header.Get(identity.MetadataKeyIdentityID)
		gotIdentityType = r.Header.Get(identity.MetadataKeyIdentityType)
		gotInstallationID = r.Header.Get(identity.MetadataKeyAppInstallationID)
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
	appsClient := &fakeAppsClient{
		getInstallationBySlug: func(ctx context.Context, req *appsv1.GetInstallationBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetInstallationBySlugResponse, error) {
			return &appsv1.GetInstallationBySlugResponse{Installation: &appsv1.Installation{
				Meta:  &appsv1.EntityMeta{Id: "inst-1"},
				AppId: "app-1",
			}}, nil
		},
		getApp: func(ctx context.Context, req *appsv1.GetAppRequest, opts ...grpc.CallOption) (*appsv1.GetAppResponse, error) {
			return &appsv1.GetAppResponse{App: &appsv1.App{ZitiServiceId: "service"}}, nil
		},
	}
	handler := &AppProxyHandler{
		apps:     appsClient,
		client:   &http.Client{Transport: transport},
		cache:    make(map[string]cachedInstallation),
		cacheTTL: time.Minute,
	}

	resolved := identity.ResolvedIdentity{IdentityID: "identity-1", IdentityType: identity.IdentityTypeUser}
	req := httptest.NewRequest(http.MethodGet, "/apps/app-slug/health", nil)
	req = req.WithContext(identity.WithIdentity(req.Context(), resolved))
	resp := httptest.NewRecorder()
	req.Header.Set(identity.MetadataKeyOrganizationID, "org-1")

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
	if gotPath != "/health" {
		t.Fatalf("expected proxied path %q, got %q", "/health", gotPath)
	}
}

func TestAppProxyHandlerServeHTTPMissingOrgID(t *testing.T) {
	handler := &AppProxyHandler{}
	req := httptest.NewRequest(http.MethodGet, "/apps/app-slug/health", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}
}

func TestAppProxyHandlerServeHTTPInvalidPath(t *testing.T) {
	handler := &AppProxyHandler{}
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
