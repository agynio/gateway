package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ProxyTarget interface {
	BaseURL() *url.URL
	DefaultHeaders() http.Header
}

func NewUpstreamProxy(target ProxyTarget) http.Handler {
	if target == nil {
		panic("proxy target is required")
	}

	upstream := target.BaseURL()
	if upstream == nil {
		panic("proxy target base URL is required")
	}

	headers := target.DefaultHeaders()
	director := func(req *http.Request) {
		req.URL.Scheme = upstream.Scheme
		req.URL.Host = upstream.Host
		req.Host = upstream.Host
		req.URL.Path = joinURLPath(upstream.Path, req.URL.Path)
		if req.URL.RawPath != "" || upstream.RawPath != "" {
			req.URL.RawPath = joinURLPath(upstream.RawPath, req.URL.RawPath)
		}

		req.URL.RawQuery = joinQueries(upstream.RawQuery, req.URL.RawQuery)

		for key := range headers {
			req.Header.Del(key)
		}
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	proxy := &httputil.ReverseProxy{Director: director}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}

	return proxy
}

func joinURLPath(basePath, requestPath string) string {
	switch {
	case basePath == "":
		if requestPath == "" {
			return "/"
		}
		return requestPath
	case requestPath == "":
		return basePath
	case strings.HasSuffix(basePath, "/") && strings.HasPrefix(requestPath, "/"):
		return basePath + requestPath[1:]
	case !strings.HasSuffix(basePath, "/") && !strings.HasPrefix(requestPath, "/"):
		return basePath + "/" + requestPath
	default:
		return basePath + requestPath
	}
}

func joinQueries(baseQuery, requestQuery string) string {
	switch {
	case baseQuery == "":
		return requestQuery
	case requestQuery == "":
		return baseQuery
	default:
		return baseQuery + "&" + requestQuery
	}
}
