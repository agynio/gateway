package gateway

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	appsv1 "github.com/agynio/gateway/gen/agynio/api/apps/v1"
	"github.com/agynio/gateway/internal/identity"
	"github.com/openziti/sdk-golang/ziti"
)

type AppProxyHandler struct {
	apps        appsv1.AppsServiceClient
	zitiContext ziti.Context
	transport   *http.Transport
	client      *http.Client
	cacheMu     sync.RWMutex
	cache       map[string]cachedApp
	cacheTTL    time.Duration
}

type cachedApp struct {
	serviceName string
	expiresAt   time.Time
}

func NewAppProxyHandler(apps appsv1.AppsServiceClient, zitiContext ziti.Context, cacheTTL time.Duration) *AppProxyHandler {
	if apps == nil {
		panic("apps client is required")
	}
	if zitiContext == nil {
		panic("ziti context is required")
	}
	if cacheTTL <= 0 {
		panic("cache ttl must be positive")
	}

	handler := &AppProxyHandler{
		apps:        apps,
		zitiContext: zitiContext,
		cache:       make(map[string]cachedApp),
		cacheTTL:    cacheTTL,
	}
	handler.transport = &http.Transport{
		DialContext: handler.dialContext,
	}
	handler.client = &http.Client{Transport: handler.transport}

	return handler
}

func (h *AppProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug, method, err := parseAppProxyPath(r.URL.Path)
	if err != nil {
		writeProblem(w, http.StatusBadRequest, err.Error())
		return
	}

	serviceName, err := h.resolveServiceName(r.Context(), slug)
	if err != nil {
		writeProblem(w, httpStatusFromError(err), err.Error())
		return
	}

	targetURL := url.URL{
		Scheme:   "http",
		Host:     serviceName,
		Path:     "/" + method,
		RawQuery: r.URL.RawQuery,
	}

	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL.String(), r.Body)
	if err != nil {
		writeProblem(w, http.StatusBadRequest, "invalid proxy request")
		return
	}
	proxyReq.Header = r.Header.Clone()
	proxyReq.ContentLength = r.ContentLength
	proxyReq.Host = serviceName

	if resolved, ok := identity.IdentityFromContext(r.Context()); ok {
		proxyReq.Header.Set(identity.MetadataKeyIdentityID, resolved.IdentityID)
		proxyReq.Header.Set(identity.MetadataKeyIdentityType, string(resolved.IdentityType))
	}

	resp, err := h.client.Do(proxyReq)
	if err != nil {
		writeProblem(w, http.StatusBadGateway, err.Error())
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func (h *AppProxyHandler) resolveServiceName(ctx context.Context, slug string) (string, error) {
	now := time.Now()
	h.cacheMu.RLock()
	if cached, ok := h.cache[slug]; ok && now.Before(cached.expiresAt) {
		h.cacheMu.RUnlock()
		return cached.serviceName, nil
	}
	h.cacheMu.RUnlock()

	response, err := h.apps.GetAppBySlug(ctx, &appsv1.GetAppBySlugRequest{Slug: slug})
	if err != nil {
		return "", err
	}

	app := response.GetApp()
	if app == nil {
		return "", fmt.Errorf("app response missing")
	}

	serviceID := strings.TrimSpace(app.GetZitiServiceId())
	if serviceID == "" {
		return "", fmt.Errorf("app service id missing")
	}
	// ZitiServiceId stores the dialable service identifier for the app.
	serviceName := serviceID

	h.cacheMu.Lock()
	h.cache[slug] = cachedApp{serviceName: serviceName, expiresAt: now.Add(h.cacheTTL)}
	h.cacheMu.Unlock()

	return serviceName, nil
}

func parseAppProxyPath(path string) (string, string, error) {
	trimmed := strings.TrimPrefix(path, "/")
	parts := strings.SplitN(trimmed, "/", 3)
	if len(parts) != 3 || parts[0] != "apps" {
		return "", "", fmt.Errorf("invalid app proxy path")
	}

	slug := strings.TrimSpace(parts[1])
	method := strings.TrimSpace(parts[2])
	if slug == "" || method == "" {
		return "", "", fmt.Errorf("invalid app proxy path")
	}

	return slug, method, nil
}

func (h *AppProxyHandler) dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	_ = network

	serviceName := strings.TrimSpace(addr)
	if host, _, err := net.SplitHostPort(addr); err == nil {
		serviceName = strings.TrimSpace(host)
	}
	if serviceName == "" {
		return nil, fmt.Errorf("service address missing")
	}

	return h.zitiContext.Dial(serviceName)
}
