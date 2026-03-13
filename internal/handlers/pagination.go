package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

func parsePagination(r *http.Request) (int, string, error) {
	pageSize := defaultPageSize
	pageToken := strings.TrimSpace(r.URL.Query().Get("pageToken"))

	if raw := strings.TrimSpace(r.URL.Query().Get("pageSize")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			return 0, "", fmt.Errorf("pageSize must be a positive integer")
		}
		pageSize = parsed
	}

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return pageSize, pageToken, nil
}
