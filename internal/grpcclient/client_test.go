package grpcclient

import (
	"context"
	"net"
	"testing"

	runnersv1 "github.com/agynio/gateway/gen/agynio/api/runners/v1"
	"github.com/agynio/gateway/internal/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

const bufConnSize = 1024 * 1024

type metadataCaptureRunnersServer struct {
	runnersv1.UnimplementedRunnersServiceServer
	metadata metadata.MD
}

func (s *metadataCaptureRunnersServer) GetWorkload(ctx context.Context, req *runnersv1.GetWorkloadRequest) (*runnersv1.GetWorkloadResponse, error) {
	s.metadata, _ = metadata.FromIncomingContext(ctx)
	return &runnersv1.GetWorkloadResponse{Workload: &runnersv1.Workload{}}, nil
}

func TestNewRunnersClientSendsSingleIdentityMetadataValue(t *testing.T) {
	listener := bufconn.Listen(bufConnSize)
	server := grpc.NewServer()
	capture := &metadataCaptureRunnersServer{}
	runnersv1.RegisterRunnersServiceServer(server, capture)
	t.Cleanup(server.Stop)
	go func() {
		_ = server.Serve(listener)
	}()

	client, err := New(
		"passthrough:///bufnet",
		runnersv1.NewRunnersServiceClient,
		WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
	)
	if err != nil {
		t.Fatalf("new runners client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
		identity.MetadataKeyIdentityID, "stale-identity",
		identity.MetadataKeyIdentityType, string(identity.IdentityTypeRunner),
		identity.MetadataKeyWorkloadID, "stale-workload",
	))
	ctx = identity.WithIdentity(ctx, identity.ResolvedIdentity{
		IdentityID:   "identity-1",
		IdentityType: identity.IdentityTypeUser,
	})

	_, err = client.Service().GetWorkload(ctx, &runnersv1.GetWorkloadRequest{Id: "workload-1"})
	if err != nil {
		t.Fatalf("get workload: %v", err)
	}

	assertMetadataValues(t, capture.metadata, identity.MetadataKeyIdentityID, []string{"identity-1"})
	assertMetadataValues(t, capture.metadata, identity.MetadataKeyIdentityType, []string{string(identity.IdentityTypeUser)})
	assertMetadataValues(t, capture.metadata, identity.MetadataKeyWorkloadID, nil)
}

func assertMetadataValues(t *testing.T, md metadata.MD, key string, expected []string) {
	t.Helper()
	values := md.Get(key)
	if len(values) != len(expected) {
		t.Fatalf("expected %s values %v, got %v", key, expected, values)
	}
	for i, expectedValue := range expected {
		if values[i] != expectedValue {
			t.Fatalf("expected %s values %v, got %v", key, expected, values)
		}
	}
}
