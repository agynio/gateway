package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultRetryWait = 100 * time.Millisecond

// HTTPClient defines the subset of the http.Client API used by the platform client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client provides typed access to the platform HTTP API.
type Client struct {
	baseURL        *url.URL
	httpClient     HTTPClient
	retries        int
	defaultHeaders http.Header
	retryWait      time.Duration
}

// NewClient constructs a Client from the provided configuration.
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	baseURL := cloneURL(cfg.BaseURL)
	headers := make(http.Header, len(cfg.Headers))
	for key, values := range cfg.Headers {
		cloned := append([]string(nil), values...)
		headers[key] = cloned
	}

	if token := strings.TrimSpace(cfg.AuthToken); token != "" {
		if _, exists := headers["Authorization"]; !exists {
			headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		}
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	retries := cfg.Retries
	if retries < 0 {
		retries = defaultRetries
	}

	return &Client{
		baseURL:        baseURL,
		httpClient:     &http.Client{Timeout: timeout},
		retries:        retries,
		defaultHeaders: headers,
		retryWait:      defaultRetryWait,
	}, nil
}

// BaseURL returns a defensive copy of the configured platform base URL.
func (c *Client) BaseURL() *url.URL {
	if c == nil {
		return nil
	}
	return cloneURL(c.baseURL)
}

// DefaultHeaders returns a defensive copy of the default headers applied to each request.
func (c *Client) DefaultHeaders() http.Header {
	if c == nil {
		return nil
	}

	cloned := make(http.Header, len(c.defaultHeaders))
	for key, values := range c.defaultHeaders {
		cloned[key] = append([]string(nil), values...)
	}
	return cloned
}

// HTTPClient exposes the underlying HTTP client used for upstream calls.
func (c *Client) HTTPClient() HTTPClient {
	if c == nil {
		return nil
	}
	return c.httpClient
}

// Retries returns the configured retry count for upstream requests.
func (c *Client) Retries() int {
	if c == nil {
		return 0
	}
	return c.retries
}

// Do performs an HTTP request against the platform API, decoding a successful JSON response into out.
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, body any, out any) (int, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	payload, err := encodeBody(body)
	if err != nil {
		return 0, fmt.Errorf("encode request body: %w", err)
	}

	attempts := c.retries + 1

	for attempt := 0; attempt < attempts; attempt++ {
		req, err := c.newRequest(ctx, method, path, query, payload)
		if err != nil {
			return 0, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if attempt < c.retries && shouldRetryError(err, method) {
				c.sleep(attempt)
				continue
			}
			return 0, err
		}

		status := resp.StatusCode
		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			if attempt < c.retries && shouldRetryStatus(status, method) {
				c.sleep(attempt)
				continue
			}
			return status, fmt.Errorf("read response body: %w", readErr)
		}

		if status >= 200 && status < 300 {
			if out != nil && len(bodyBytes) > 0 {
				if err := json.Unmarshal(bodyBytes, out); err != nil {
					return status, fmt.Errorf("decode response body: %w", err)
				}
			}
			return status, nil
		}

		pErr := &Error{Status: status, Body: bodyBytes}
		if len(bodyBytes) > 0 {
			var problem Problem
			if err := json.Unmarshal(bodyBytes, &problem); err == nil {
				if problem.Status == 0 {
					problem.Status = status
				}
				pErr.Problem = &problem
			} else {
				pErr.Err = err
			}
		}

		if attempt < c.retries && shouldRetryStatus(status, method) {
			c.sleep(attempt)
			continue
		}

		return status, pErr
	}

	return 0, fmt.Errorf("exhausted retries")
}

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, payload []byte) (*http.Request, error) {
	if method == "" {
		return nil, fmt.Errorf("http method is required")
	}

	normalizedPath := path
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	rel := &url.URL{Path: normalizedPath}
	if query != nil {
		rel.RawQuery = query.Encode()
	}

	target := c.baseURL.ResolveReference(rel)

	var bodyReader io.Reader
	if payload != nil {
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, target.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	for key, values := range c.defaultHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) sleep(attempt int) {
	if c.retryWait <= 0 {
		return
	}
	time.Sleep(time.Duration(attempt+1) * c.retryWait)
}

func encodeBody(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 || string(bytes) == "null" {
		return nil, nil
	}
	return bytes, nil
}

func shouldRetryError(err error, method string) bool {
	if err == nil || !isIdempotent(method) {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}

	return true
}

func shouldRetryStatus(status int, method string) bool {
	if !isIdempotent(method) {
		return false
	}
	return status >= 500 && status != http.StatusNotImplemented
}

func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodDelete, http.MethodOptions:
		return true
	default:
		return false
	}
}

func cloneURL(u *url.URL) *url.URL {
	if u == nil {
		return nil
	}
	cloned := *u
	return &cloned
}
