package gateway

import (
	"context"
	"encoding/json"
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
	cacheMu     sync.RWMutex
	cache       map[string]cachedApp
	cacheTTL    time.Duration
}

type cachedApp struct {
	serviceName string
	expiresAt   time.Time
}

type proxyErrorResponse struct {
	Error string `json:"error"`
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

	return &AppProxyHandler{
		apps:        apps,
		zitiContext: zitiContext,
		cache:       make(map[string]cachedApp),
		cacheTTL:    cacheTTL,
	}
}

func (h *AppProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug, method, err := parseAppProxyPath(r.URL.Path)
	if err != nil {
		writeProxyError(w, http.StatusBadRequest, err.Error())
		return
	}

	serviceName, err := h.resolveServiceName(r.Context(), slug)
	if err != nil {
		writeProxyError(w, httpStatusFromError(err), err.Error())
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
		writeProxyError(w, http.StatusBadRequest, "invalid proxy request")
		return
	}
	proxyReq.Header = r.Header.Clone()
	proxyReq.ContentLength = r.ContentLength
	proxyReq.Host = serviceName

	if resolved, ok := identity.IdentityFromContext(r.Context()); ok {
		proxyReq.Header.Set(identity.MetadataKeyIdentityID, resolved.IdentityID)
		proxyReq.Header.Set(identity.MetadataKeyIdentityType, string(resolved.IdentityType))
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return h.zitiContext.Dial(serviceName)
		},
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Do(proxyReq)
	if err != nil {
		writeProxyError(w, http.StatusBadGateway, err.Error())
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	if resp.ContentLength >= 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	}

	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func (h *AppProxyHandler) resolveServiceName(ctx context.Context, slug string) (string, error) {
	trimmed := strings.TrimSpace(slug)
	if trimmed == "" {
		return "", fmt.Errorf("app slug is required")
	}

	now := time.Now()
	h.cacheMu.RLock()
	if cached, ok := h.cache[trimmed]; ok && now.Before(cached.expiresAt) {
		h.cacheMu.RUnlock()
		return cached.serviceName, nil
	}
	h.cacheMu.RUnlock()

	response, err := h.apps.GetAppBySlug(ctx, &appsv1.GetAppBySlugRequest{Slug: trimmed})
	if err != nil {
		return "", err
	}

	app := response.GetApp()
	if app == nil {
		return "", fmt.Errorf("app response missing")
	}

	serviceName := strings.TrimSpace(app.GetZitiServiceId())
	if serviceName == "" {
		return "", fmt.Errorf("app service name missing")
	}

	h.cacheMu.Lock()
	h.cache[trimmed] = cachedApp{serviceName: serviceName, expiresAt: now.Add(h.cacheTTL)}
	h.cacheMu.Unlock()

	return serviceName, nil
}

func parseAppProxyPath(path string) (string, string, error) {
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")
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

func writeProxyError(w http.ResponseWriter, status int, message string) {
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		trimmed = http.StatusText(status)
	}

	payload, err := json.Marshal(proxyErrorResponse{Error: trimmed})
	if err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}
