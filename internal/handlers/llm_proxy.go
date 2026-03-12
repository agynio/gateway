package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewLLMResponseProxy(llmHTTPBaseURL *url.URL) http.Handler {
	if llmHTTPBaseURL == nil {
		panic("llm http base url is required")
	}

	director := func(req *http.Request) {
		req.URL.Scheme = llmHTTPBaseURL.Scheme
		req.URL.Host = llmHTTPBaseURL.Host
		req.Host = llmHTTPBaseURL.Host
		req.URL.Path = "/v1/responses"
		req.URL.RawPath = "/v1/responses"
	}

	proxy := &httputil.ReverseProxy{
		Director:      director,
		FlushInterval: -1,
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("llm response proxy error: %v", err)
		problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), "upstream request failed")
		WriteProblem(w, problem)
	}

	return proxy
}
