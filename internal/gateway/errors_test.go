package gateway

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcCodeToConnectCode(t *testing.T) {
	tests := []struct {
		name string
		code codes.Code
		want connect.Code
	}{
		{name: "Canceled", code: codes.Canceled, want: connect.CodeCanceled},
		{name: "Unknown", code: codes.Unknown, want: connect.CodeUnknown},
		{name: "InvalidArgument", code: codes.InvalidArgument, want: connect.CodeInvalidArgument},
		{name: "DeadlineExceeded", code: codes.DeadlineExceeded, want: connect.CodeDeadlineExceeded},
		{name: "NotFound", code: codes.NotFound, want: connect.CodeNotFound},
		{name: "AlreadyExists", code: codes.AlreadyExists, want: connect.CodeAlreadyExists},
		{name: "PermissionDenied", code: codes.PermissionDenied, want: connect.CodePermissionDenied},
		{name: "ResourceExhausted", code: codes.ResourceExhausted, want: connect.CodeResourceExhausted},
		{name: "FailedPrecondition", code: codes.FailedPrecondition, want: connect.CodeFailedPrecondition},
		{name: "Aborted", code: codes.Aborted, want: connect.CodeAborted},
		{name: "OutOfRange", code: codes.OutOfRange, want: connect.CodeOutOfRange},
		{name: "Unimplemented", code: codes.Unimplemented, want: connect.CodeUnimplemented},
		{name: "Internal", code: codes.Internal, want: connect.CodeInternal},
		{name: "Unavailable", code: codes.Unavailable, want: connect.CodeUnavailable},
		{name: "DataLoss", code: codes.DataLoss, want: connect.CodeDataLoss},
		{name: "Unauthenticated", code: codes.Unauthenticated, want: connect.CodeUnauthenticated},
		{name: "OK", code: codes.OK, want: connect.CodeInternal},
		{name: "UnknownCode", code: codes.Code(99), want: connect.CodeInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := grpcCodeToConnectCode(tt.code); got != tt.want {
				t.Fatalf("grpcCodeToConnectCode(%v) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestToConnectError(t *testing.T) {
	t.Run("non-grpc error", func(t *testing.T) {
		connectErr := toConnectError(errors.New("boom"))
		if connectErr.Code() != connect.CodeInternal {
			t.Fatalf("expected CodeInternal, got %v", connectErr.Code())
		}
		if connectErr.Message() != "boom" {
			t.Fatalf("expected message boom, got %q", connectErr.Message())
		}
	})

	t.Run("grpc error uses message", func(t *testing.T) {
		connectErr := toConnectError(status.Error(codes.PermissionDenied, "nope"))
		if connectErr.Code() != connect.CodePermissionDenied {
			t.Fatalf("expected CodePermissionDenied, got %v", connectErr.Code())
		}
		if connectErr.Message() != "nope" {
			t.Fatalf("expected message nope, got %q", connectErr.Message())
		}
	})

	t.Run("grpc error defaults to code message", func(t *testing.T) {
		connectErr := toConnectError(status.Error(codes.NotFound, ""))
		if connectErr.Code() != connect.CodeNotFound {
			t.Fatalf("expected CodeNotFound, got %v", connectErr.Code())
		}
		if connectErr.Message() != codes.NotFound.String() {
			t.Fatalf("expected message %q, got %q", codes.NotFound.String(), connectErr.Message())
		}
	})
}
