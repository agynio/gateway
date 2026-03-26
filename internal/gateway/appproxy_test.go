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
	getAppBySlug func(ctx context.Context, req *appsv1.GetAppBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetAppBySlugResponse, error)
}

func (f *fakeAppsClient) RegisterApp(ctx context.Context, req *appsv1.RegisterAppRequest, opts ...grpc.CallOption) (*appsv1.RegisterAppResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) EnrollApp(ctx context.Context, req *appsv1.EnrollAppRequest, opts ...grpc.CallOption) (*appsv1.EnrollAppResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) GetApp(ctx context.Context, req *appsv1.GetAppRequest, opts ...grpc.CallOption) (*appsv1.GetAppResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeAppsClient) GetAppBySlug(ctx context.Context, req *appsv1.GetAppBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetAppBySlugResponse, error) {
	if f.getAppBySlug != nil {
		return f.getAppBySlug(ctx, req, opts...)
	}
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

func TestResolveServiceNameCache(t *testing.T) {
	var calls int
	appsClient := &fakeAppsClient{
		getAppBySlug: func(ctx context.Context, req *appsv1.GetAppBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetAppBySlugResponse, error) {
			calls++
			return &appsv1.GetAppBySlugResponse{App: &appsv1.App{ZitiServiceId: "service-1"}}, nil
		},
	}
	handler := &AppProxyHandler{
		apps:     appsClient,
		cache:    make(map[string]cachedApp),
		cacheTTL: time.Minute,
	}

	ctx := context.Background()
	serviceName, err := handler.resolveServiceName(ctx, "slug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if serviceName != "service-1" {
		t.Fatalf("expected service name %q, got %q", "service-1", serviceName)
	}

	serviceName, err = handler.resolveServiceName(ctx, "slug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if serviceName != "service-1" {
		t.Fatalf("expected service name %q, got %q", "service-1", serviceName)
	}
	if calls != 1 {
		t.Fatalf("expected apps service to be called once, got %d", calls)
	}
}

func TestResolveServiceNameErrors(t *testing.T) {
	t.Run("app lookup error", func(t *testing.T) {
		handler := &AppProxyHandler{
			apps: &fakeAppsClient{
				getAppBySlug: func(ctx context.Context, req *appsv1.GetAppBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetAppBySlugResponse, error) {
					return nil, status.Error(codes.NotFound, "missing")
				},
			},
			cache:    make(map[string]cachedApp),
			cacheTTL: time.Minute,
		}

		_, err := handler.resolveServiceName(context.Background(), "slug")
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("missing app", func(t *testing.T) {
		handler := &AppProxyHandler{
			apps: &fakeAppsClient{
				getAppBySlug: func(ctx context.Context, req *appsv1.GetAppBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetAppBySlugResponse, error) {
					return &appsv1.GetAppBySlugResponse{}, nil
				},
			},
			cache:    make(map[string]cachedApp),
			cacheTTL: time.Minute,
		}

		_, err := handler.resolveServiceName(context.Background(), "slug")
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("missing service id", func(t *testing.T) {
		handler := &AppProxyHandler{
			apps: &fakeAppsClient{
				getAppBySlug: func(ctx context.Context, req *appsv1.GetAppBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetAppBySlugResponse, error) {
					return &appsv1.GetAppBySlugResponse{App: &appsv1.App{ZitiServiceId: " "}}, nil
				},
			},
			cache:    make(map[string]cachedApp),
			cacheTTL: time.Minute,
		}

		_, err := handler.resolveServiceName(context.Background(), "slug")
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestAppProxyHandlerServeHTTP(t *testing.T) {
	var gotIdentityID string
	var gotIdentityType string
	var gotPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotIdentityID = r.Header.Get(identity.MetadataKeyIdentityID)
		gotIdentityType = r.Header.Get(identity.MetadataKeyIdentityType)
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
		getAppBySlug: func(ctx context.Context, req *appsv1.GetAppBySlugRequest, opts ...grpc.CallOption) (*appsv1.GetAppBySlugResponse, error) {
			return &appsv1.GetAppBySlugResponse{App: &appsv1.App{ZitiServiceId: "service"}}, nil
		},
	}
	handler := &AppProxyHandler{
		apps:     appsClient,
		client:   &http.Client{Transport: transport},
		cache:    make(map[string]cachedApp),
		cacheTTL: time.Minute,
	}

	resolved := identity.ResolvedIdentity{IdentityID: "identity-1", IdentityType: identity.IdentityTypeUser}
	req := httptest.NewRequest(http.MethodGet, "/apps/app-slug/health", nil)
	req = req.WithContext(identity.WithIdentity(req.Context(), resolved))
	resp := httptest.NewRecorder()

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
	if gotPath != "/health" {
		t.Fatalf("expected proxied path %q, got %q", "/health", gotPath)
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
