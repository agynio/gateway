package handlers

import (
	"encoding/json"
	"errors"
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

type ProblemError struct {
	Problem gen.Problem
	Err     error
}

func NewProblemError(problem gen.Problem, err error) *ProblemError {
	return &ProblemError{Problem: problem, Err: err}
}

func (e *ProblemError) Error() string {
	ititle := strings.TrimSpace(e.Problem.Title)
	if ititle == "" {
		ititle = http.StatusText(int(e.Problem.Status))
	}
	detail := ""
	if e.Problem.Detail != nil {
		detail = strings.TrimSpace(*e.Problem.Detail)
	}
	if detail != "" {
		return fmt.Sprintf("%s: %s", ititle, detail)
	}
	return ititle
}

func (e *ProblemError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
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

func StrictErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var problemErr *ProblemError
	if errors.As(err, &problemErr) {
		if problemErr.Err != nil {
			log.Printf("handler error: %v", problemErr.Err)
		} else {
			log.Printf("handler error: %s %s", r.Method, r.URL.Path)
		}
		WriteProblem(w, problemErr.Problem)
		return
	}

	log.Printf("handler unexpected error: %v", err)
	problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), "upstream request failed", nil)
	WriteProblem(w, problem)
}
