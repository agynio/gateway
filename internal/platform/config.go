package platform

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
	defaultRetries = 2
)

// Config holds the runtime configuration for communicating with the platform service.
type Config struct {
	BaseURL           *url.URL
	FilesBaseURL      *url.URL
	Timeout           time.Duration
	Retries           int
	RetriesConfigured bool
	AuthToken         string
	Headers           http.Header
}

// LoadConfigFromEnv constructs a Config instance from environment variables.
func LoadConfigFromEnv() (*Config, error) {
	rawBaseURL := strings.TrimSpace(os.Getenv("PLATFORM_BASE_URL"))
	if rawBaseURL == "" {
		return nil, fmt.Errorf("PLATFORM_BASE_URL is required")
	}

	parsedURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse PLATFORM_BASE_URL: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("PLATFORM_BASE_URL must include scheme and host")
	}

	rawFilesBaseURL := strings.TrimSpace(os.Getenv("FILES_BASE_URL"))
	var filesBaseURL *url.URL
	if rawFilesBaseURL != "" {
		parsedFilesURL, err := url.Parse(rawFilesBaseURL)
		if err != nil {
			return nil, fmt.Errorf("parse FILES_BASE_URL: %w", err)
		}
		if parsedFilesURL.Scheme == "" || parsedFilesURL.Host == "" {
			return nil, fmt.Errorf("FILES_BASE_URL must include scheme and host")
		}
		filesBaseURL = parsedFilesURL
	}

	timeout := defaultTimeout
	if rawTimeout := strings.TrimSpace(os.Getenv("PLATFORM_TIMEOUT_MS")); rawTimeout != "" {
		ms, err := strconv.Atoi(rawTimeout)
		if err != nil {
			return nil, fmt.Errorf("parse PLATFORM_TIMEOUT_MS: %w", err)
		}
		if ms > 0 {
			timeout = time.Duration(ms) * time.Millisecond
		}
	}

	retries := defaultRetries
	retriesConfigured := false
	if rawRetries := strings.TrimSpace(os.Getenv("PLATFORM_RETRIES")); rawRetries != "" {
		value, err := strconv.Atoi(rawRetries)
		if err != nil {
			return nil, fmt.Errorf("parse PLATFORM_RETRIES: %w", err)
		}
		if value < 0 {
			return nil, fmt.Errorf("PLATFORM_RETRIES must be >= 0")
		}
		retries = value
		retriesConfigured = true
	}

	headers := make(http.Header)
	if rawHeaders := strings.TrimSpace(os.Getenv("PLATFORM_REQUEST_HEADERS_JSON")); rawHeaders != "" {
		var input map[string]string
		if err := json.Unmarshal([]byte(rawHeaders), &input); err != nil {
			return nil, fmt.Errorf("parse PLATFORM_REQUEST_HEADERS_JSON: %w", err)
		}
		for key, value := range input {
			trimmedKey := strings.TrimSpace(key)
			if trimmedKey == "" {
				continue
			}
			headers.Set(trimmedKey, value)
		}
	}

	authToken := strings.TrimSpace(os.Getenv("PLATFORM_AUTH_TOKEN"))

	return &Config{
		BaseURL:           parsedURL,
		FilesBaseURL:      filesBaseURL,
		Timeout:           timeout,
		Retries:           retries,
		RetriesConfigured: retriesConfigured,
		AuthToken:         authToken,
		Headers:           headers,
	}, nil
}
