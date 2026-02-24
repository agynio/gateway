package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/agynio/gateway/internal/gen"
)

const (
	problemContentType = "application/problem+json"
	problemTypeDefault = "about:blank"
)

func NewProblem(status int, title, detail string, errs map[string][]string) gen.Problem {
	problem := gen.Problem{
		Status: int32(status),
		Title:  title,
		Type:   problemTypeDefault,
	}

	if detail != "" {
		problem.Detail = ptr(detail)
	}

	if len(errs) > 0 {
		copied := make(map[string][]string, len(errs))
		for field, messages := range errs {
			copied[field] = append([]string(nil), messages...)
		}
		problem.Errors = &copied
	}

	return problem
}

func WriteProblem(w http.ResponseWriter, problem gen.Problem) {
	w.Header().Set("Content-Type", problemContentType)
	w.WriteHeader(int(problem.Status))
	if err := json.NewEncoder(w).Encode(problem); err != nil {
		log.Printf("failed to encode problem response: %v", err)
	}
}

func RequestValidationError(w http.ResponseWriter, message string, statusCode int) {
	if statusCode == http.StatusNotFound {
		normalized := strings.ToLower(strings.TrimSpace(message))
		if strings.Contains(normalized, "method not allowed") {
			statusCode = http.StatusMethodNotAllowed
		}
	}

	title := http.StatusText(statusCode)
	if title == "" {
		title = "Request validation error"
	}
	problem := NewProblem(statusCode, title, message, nil)
	WriteProblem(w, problem)
}

func RequestValidationMultiError(me openapi3.MultiError) (int, error) {
	messages := make([]string, 0, len(me))
	for _, err := range me {
		messages = append(messages, err.Error())
	}

	detail := strings.Join(messages, "; ")
	if detail == "" {
		detail = "validation failed"
	}

	return http.StatusUnprocessableEntity, fmt.Errorf(detail)
}

func ptr[T any](v T) *T {
	return &v
}
