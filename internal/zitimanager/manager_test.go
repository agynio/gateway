package zitimanager

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/openziti/sdk-golang/ziti"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	zitimgmtv1 "github.com/agynio/gateway/gen/agynio/api/ziti_management/v1"
	"github.com/agynio/gateway/internal/zitimgmtclient"
)

type fakeZitiManagementServer struct {
	zitimgmtv1.UnimplementedZitiManagementServiceServer
	requestServiceIdentity func(context.Context, *zitimgmtv1.RequestServiceIdentityRequest) (*zitimgmtv1.RequestServiceIdentityResponse, error)
	extendIdentityLease    func(context.Context, *zitimgmtv1.ExtendIdentityLeaseRequest) (*zitimgmtv1.ExtendIdentityLeaseResponse, error)
}

func (f *fakeZitiManagementServer) RequestServiceIdentity(ctx context.Context, req *zitimgmtv1.RequestServiceIdentityRequest) (*zitimgmtv1.RequestServiceIdentityResponse, error) {
	if f.requestServiceIdentity != nil {
		return f.requestServiceIdentity(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeZitiManagementServer) ExtendIdentityLease(ctx context.Context, req *zitimgmtv1.ExtendIdentityLeaseRequest) (*zitimgmtv1.ExtendIdentityLeaseResponse, error) {
	if f.extendIdentityLease != nil {
		return f.extendIdentityLease(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

type fakeZitiContext struct {
	ziti.Context
	mu     sync.Mutex
	closed bool
}

func (f *fakeZitiContext) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closed = true
}

func (f *fakeZitiContext) Closed() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.closed
}

func startZitiManagementServer(t *testing.T, server *fakeZitiManagementServer) (string, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	zitimgmtv1.RegisterZitiManagementServiceServer(grpcServer, server)
	go func() {
		_ = grpcServer.Serve(listener)
	}()

	return listener.Addr().String(), func() {
		grpcServer.Stop()
		_ = listener.Close()
	}
}

func TestManagerNewEnrollsAndCreatesListener(t *testing.T) {
	requestCh := make(chan zitimgmtv1.ServiceType, 1)
	server := &fakeZitiManagementServer{
		requestServiceIdentity: func(ctx context.Context, req *zitimgmtv1.RequestServiceIdentityRequest) (*zitimgmtv1.RequestServiceIdentityResponse, error) {
			requestCh <- req.ServiceType
			return &zitimgmtv1.RequestServiceIdentityResponse{
				ZitiIdentityId: "identity-1",
				IdentityJson:   []byte("{}"),
			}, nil
		},
	}

	addr, stop := startZitiManagementServer(t, server)
	t.Cleanup(stop)

	client, err := zitimgmtclient.NewClient(addr)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Close()
	})

	originalFactory := newZitiContext
	fakeContext := &fakeZitiContext{}
	newZitiContext = func(config *ziti.Config) (ziti.Context, error) {
		return fakeContext, nil
	}
	t.Cleanup(func() {
		newZitiContext = originalFactory
	})

	var receivedContext ziti.Context
	var listener net.Listener
	listenerFactory := func(zitiCtx ziti.Context) (net.Listener, error) {
		receivedContext = zitiCtx
		var err error
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, err
		}
		return listener, nil
	}

	onNewListenerCh := make(chan net.Listener, 1)
	onNewListener := func(newListener net.Listener) {
		onNewListenerCh <- newListener
	}

	mgr, err := New(
		context.Background(),
		client,
		zitimgmtv1.ServiceType_SERVICE_TYPE_GATEWAY,
		time.Second,
		time.Second,
		listenerFactory,
		onNewListener,
	)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	if receivedContext != fakeContext {
		t.Fatalf("expected listener factory to receive ziti context")
	}
	if mgr.ZitiContext() != fakeContext {
		t.Fatalf("expected ziti context to be stored")
	}
	if mgr.identityID != "identity-1" {
		t.Fatalf("expected identity id %q, got %q", "identity-1", mgr.identityID)
	}

	select {
	case got := <-requestCh:
		if got != zitimgmtv1.ServiceType_SERVICE_TYPE_GATEWAY {
			t.Fatalf("unexpected service type: %v", got)
		}
	case <-time.After(time.Second):
		t.Fatalf("expected request service identity call")
	}

	select {
	case got := <-onNewListenerCh:
		if got != listener {
			t.Fatalf("expected listener callback to receive listener")
		}
	case <-time.After(time.Second):
		t.Fatalf("expected on new listener callback")
	}

	if listener != nil {
		_ = listener.Close()
	}
}

func TestRunLeaseRenewalReenrollsOnNotFound(t *testing.T) {
	requestCh := make(chan string, 2)
	extendCh := make(chan string, 2)
	identityIDs := []string{"identity-1", "identity-2"}
	var requestMu sync.Mutex
	var requestIndex int
	server := &fakeZitiManagementServer{
		requestServiceIdentity: func(ctx context.Context, req *zitimgmtv1.RequestServiceIdentityRequest) (*zitimgmtv1.RequestServiceIdentityResponse, error) {
			requestMu.Lock()
			idx := requestIndex
			requestIndex++
			requestMu.Unlock()
			if idx >= len(identityIDs) {
				return nil, status.Error(codes.Internal, "too many requests")
			}
			identityID := identityIDs[idx]
			requestCh <- identityID
			return &zitimgmtv1.RequestServiceIdentityResponse{
				ZitiIdentityId: identityID,
				IdentityJson:   []byte("{}"),
			}, nil
		},
		extendIdentityLease: func(ctx context.Context, req *zitimgmtv1.ExtendIdentityLeaseRequest) (*zitimgmtv1.ExtendIdentityLeaseResponse, error) {
			extendCh <- req.ZitiIdentityId
			return nil, status.Error(codes.NotFound, "missing")
		},
	}

	addr, stop := startZitiManagementServer(t, server)
	t.Cleanup(stop)

	client, err := zitimgmtclient.NewClient(addr)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Close()
	})

	originalFactory := newZitiContext
	firstContext := &fakeZitiContext{}
	secondContext := &fakeZitiContext{}
	contexts := []ziti.Context{firstContext, secondContext}
	var contextMu sync.Mutex
	newZitiContext = func(config *ziti.Config) (ziti.Context, error) {
		contextMu.Lock()
		defer contextMu.Unlock()
		if len(contexts) == 0 {
			return nil, status.Error(codes.Internal, "no contexts available")
		}
		ctx := contexts[0]
		contexts = contexts[1:]
		return ctx, nil
	}
	t.Cleanup(func() {
		newZitiContext = originalFactory
	})

	listenerCh := make(chan net.Listener, 2)
	listenerFactory := func(zitiCtx ziti.Context) (net.Listener, error) {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, err
		}
		listenerCh <- listener
		return listener, nil
	}

	onNewListener := func(listener net.Listener) {}

	mgr, err := New(
		context.Background(),
		client,
		zitimgmtv1.ServiceType_SERVICE_TYPE_GATEWAY,
		time.Second,
		10*time.Millisecond,
		listenerFactory,
		onNewListener,
	)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	firstListener := <-listenerCh
	if firstListener != nil {
		defer func() {
			_ = firstListener.Close()
		}()
	}

	leaseCtx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go mgr.RunLeaseRenewal(leaseCtx)

	select {
	case got := <-extendCh:
		if got != "identity-1" {
			t.Fatalf("expected extend identity for %q, got %q", "identity-1", got)
		}
	case <-time.After(time.Second):
		t.Fatalf("expected extend identity lease call")
	}

	var secondListener net.Listener
	select {
	case secondListener = <-listenerCh:
	case <-time.After(2 * time.Second):
		t.Fatalf("expected re-enroll listener")
	}
	cancel()
	if secondListener != nil {
		defer func() {
			_ = secondListener.Close()
		}()
	}

	if mgr.ZitiContext() != secondContext {
		t.Fatalf("expected manager to store new ziti context")
	}
	if !firstContext.Closed() {
		t.Fatalf("expected old ziti context to be closed")
	}

	select {
	case got := <-requestCh:
		if got != "identity-1" {
			t.Fatalf("expected first identity %q, got %q", "identity-1", got)
		}
	case <-time.After(time.Second):
		t.Fatalf("expected initial enrollment")
	}
	select {
	case got := <-requestCh:
		if got != "identity-2" {
			t.Fatalf("expected second identity %q, got %q", "identity-2", got)
		}
	case <-time.After(time.Second):
		t.Fatalf("expected re-enrollment")
	}
}
