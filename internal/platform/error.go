package platform

import (
	"fmt"
	"net/http"
	"strings"
)

// Problem models the upstream RFC7807 error payload returned by the platform service.
type Problem struct {
	Type     string              `json:"type"`
	Title    string              `json:"title"`
	Status   int                 `json:"status"`
	Detail   *string             `json:"detail,omitempty"`
	Instance *string             `json:"instance,omitempty"`
	Errors   map[string][]string `json:"errors,omitempty"`
}

// Error represents an HTTP error response from the platform service.
type Error struct {
	Status  int
	Problem *Problem
	Body    []byte
	Err     error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	base := fmt.Sprintf("platform responded with status %d", e.Status)

	if e.Problem != nil {
		title := e.Problem.Title
		if strings.TrimSpace(title) == "" {
			title = http.StatusText(e.Status)
		}

		if detail := strings.TrimSpace(ptrValue(e.Problem.Detail)); detail != "" {
			return fmt.Sprintf("%s %s: %s", base, title, detail)
		}

		return fmt.Sprintf("%s %s", base, title)
	}

	if len(e.Body) > 0 {
		snippet := strings.TrimSpace(string(e.Body))
		if len(snippet) > 256 {
			snippet = snippet[:253] + "..."
		}
		if snippet != "" {
			return fmt.Sprintf("%s (%s)", base, snippet)
		}
	}

	if e.Err != nil {
		return fmt.Sprintf("%s (%v)", base, e.Err)
	}

	return base
}

// Unwrap exposes the underlying error, if any.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func ptrValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
