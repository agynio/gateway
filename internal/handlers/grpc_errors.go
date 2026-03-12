package handlers

import (
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcErrorToProblem(err error) *ProblemError {
	if err == nil {
		return nil
	}

	grpcStatus, ok := status.FromError(err)
	if !ok {
		problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), err.Error())
		return NewProblemError(problem, err)
	}

	statusCode := grpcStatusToHTTP(grpcStatus.Code())
	message := strings.TrimSpace(grpcStatus.Message())
	if message == "" {
		message = http.StatusText(statusCode)
	}

	problem := NewProblem(statusCode, http.StatusText(statusCode), message)
	return NewProblemError(problem, err)
}

func grpcStatusToHTTP(code codes.Code) int {
	switch code {
	case codes.InvalidArgument, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		return http.StatusConflict
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.FailedPrecondition:
		// Map to 412 to reflect explicit precondition failures.
		return http.StatusPreconditionFailed
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		// Use 503 to signal a retryable upstream outage.
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusBadGateway
	}
}
