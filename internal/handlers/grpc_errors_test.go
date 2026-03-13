package handlers

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestGrpcStatusToHTTP(t *testing.T) {
	tests := []struct {
		name string
		code codes.Code
		want int
	}{
		{name: "not found", code: codes.NotFound, want: http.StatusNotFound},
		{name: "invalid argument", code: codes.InvalidArgument, want: http.StatusBadRequest},
		{name: "already exists", code: codes.AlreadyExists, want: http.StatusConflict},
		{name: "failed precondition", code: codes.FailedPrecondition, want: http.StatusPreconditionFailed},
		{name: "unavailable", code: codes.Unavailable, want: http.StatusServiceUnavailable},
		{name: "permission denied", code: codes.PermissionDenied, want: http.StatusForbidden},
		{name: "unauthenticated", code: codes.Unauthenticated, want: http.StatusUnauthorized},
		{name: "default", code: codes.Internal, want: http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := grpcStatusToHTTP(tt.code); got != tt.want {
				t.Fatalf("unexpected status: got %d want %d", got, tt.want)
			}
		})
	}
}
