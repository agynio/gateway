package gateway

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	problemContentType = "application/problem+json"
	problemTypeDefault = "about:blank"
)

type MeResponse struct {
	IdentityID   string                `json:"identity_id"`
	IdentityType identity.IdentityType `json:"identity_type"`
	TenantID     string                `json:"tenant_id"`
	AuthMethod   identity.AuthMethod   `json:"auth_method"`
}

type problemResponse struct {
	Type   string `json:"type,omitempty"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	resolvedIdentity, ok := identity.IdentityFromContext(r.Context())
	if !ok {
		writeProblem(w, http.StatusUnauthorized, "identity not available")
		return
	}

	response := MeResponse{
		IdentityID:   resolvedIdentity.IdentityID,
		IdentityType: resolvedIdentity.IdentityType,
		TenantID:     resolvedIdentity.TenantID,
		AuthMethod:   resolvedIdentity.AuthMethod,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		writeProblem(w, http.StatusInternalServerError, "failed to encode response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

func writeProblem(w http.ResponseWriter, status int, detail string) {
	message := strings.TrimSpace(detail)
	if message == "" {
		message = http.StatusText(status)
	}

	problem := problemResponse{
		Type:   problemTypeDefault,
		Title:  http.StatusText(status),
		Status: status,
		Detail: message,
	}

	payload, err := json.Marshal(problem)
	if err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.Header().Set("Content-Type", problemContentType)
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}

// httpStatusFromError maps gRPC errors to HTTP status codes for the /me HTTP endpoint.
// ConnectRPC handlers use grpcCodeToConnectCode in errors.go, so both mappings exist
// to target different transports with different default fallbacks.
func httpStatusFromError(err error) int {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		return http.StatusBadGateway
	}

	return grpcStatusToHTTP(grpcStatus.Code())
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
		return http.StatusPreconditionFailed
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusBadGateway
	}
}
