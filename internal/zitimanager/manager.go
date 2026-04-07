package zitimanager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/openziti/sdk-golang/ziti"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	zitimgmtv1 "github.com/agynio/gateway/gen/agynio/api/ziti_management/v1"
	"github.com/agynio/gateway/internal/zitimgmtclient"
)

const (
	retryInitialBackoff = 1 * time.Second
	retryMaxBackoff     = 15 * time.Second
)

var (
	leaseRetryBackoffs = []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	newZitiContext     = ziti.NewContext
)

type ListenerFactory func(zitiCtx ziti.Context) (net.Listener, error)
type OnNewListener func(listener net.Listener)

type Manager struct {
	mu              sync.RWMutex
	zitiCtx         ziti.Context
	identityID      string
	mgmtClient      *zitimgmtclient.Client
	serviceType     zitimgmtv1.ServiceType
	renewalInterval time.Duration
	enrollTimeout   time.Duration
	appCtx          context.Context

	listenerFactory ListenerFactory
	onNewListener   OnNewListener
}

func New(
	appCtx context.Context,
	client *zitimgmtclient.Client,
	serviceType zitimgmtv1.ServiceType,
	enrollTimeout time.Duration,
	renewalInterval time.Duration,
	listenerFactory ListenerFactory,
	onNewListener OnNewListener,
) (*Manager, error) {
	if appCtx == nil {
		return nil, fmt.Errorf("app context is required")
	}
	if client == nil {
		return nil, fmt.Errorf("ziti management client is required")
	}
	if serviceType == zitimgmtv1.ServiceType_SERVICE_TYPE_UNSPECIFIED {
		return nil, fmt.Errorf("service type is required")
	}
	if enrollTimeout <= 0 {
		return nil, fmt.Errorf("enroll timeout must be positive")
	}
	if renewalInterval <= 0 {
		return nil, fmt.Errorf("renewal interval must be positive")
	}
	if listenerFactory == nil {
		return nil, fmt.Errorf("listener factory is required")
	}
	if onNewListener == nil {
		return nil, fmt.Errorf("on new listener callback is required")
	}

	mgr := &Manager{
		mgmtClient:      client,
		serviceType:     serviceType,
		renewalInterval: renewalInterval,
		enrollTimeout:   enrollTimeout,
		appCtx:          appCtx,
		listenerFactory: listenerFactory,
		onNewListener:   onNewListener,
	}

	if err := mgr.reEnroll(appCtx); err != nil {
		return nil, err
	}

	return mgr, nil
}

func (m *Manager) ZitiContext() ziti.Context {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.zitiCtx
}

func (m *Manager) RunLeaseRenewal(ctx context.Context) {
	ticker := time.NewTicker(m.renewalInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ctx.Err() != nil {
				return
			}
			err := m.extendLeaseWithRetry(ctx)
			if err == nil {
				continue
			}
			if status.Code(err) == codes.NotFound {
				if reEnrollErr := m.reEnroll(ctx); reEnrollErr != nil {
					log.Printf("failed to re-enroll ziti identity: %v", reEnrollErr)
				}
				continue
			}
			log.Printf("failed to extend ziti lease: %v", err)
		}
	}
}

func (m *Manager) reEnroll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.zitiCtx != nil {
		m.zitiCtx.Close()
		m.zitiCtx = nil
	}

	enrollmentCtx, cancel := context.WithTimeout(ctx, m.enrollTimeout)
	defer cancel()

	var identityID string
	var identityJSON []byte
	if err := retryWithBackoff(enrollmentCtx, "ziti enrollment", func(attemptCtx context.Context) error {
		var requestErr error
		identityID, identityJSON, requestErr = m.mgmtClient.RequestServiceIdentity(attemptCtx, m.serviceType)
		return requestErr
	}); err != nil {
		return err
	}

	zitiConfig := &ziti.Config{}
	if err := json.Unmarshal(identityJSON, zitiConfig); err != nil {
		return fmt.Errorf("failed to parse ziti identity: %w", err)
	}

	zitiCtx, err := newZitiContext(zitiConfig)
	if err != nil {
		return fmt.Errorf("failed to create ziti context: %w", err)
	}

	listener, err := m.listenerFactory(zitiCtx)
	if err != nil {
		zitiCtx.Close()
		return err
	}

	m.zitiCtx = zitiCtx
	m.identityID = identityID
	m.onNewListener(listener)

	return nil
}

func (m *Manager) extendLeaseWithRetry(ctx context.Context) error {
	var lastErr error
	for attempt := 0; attempt <= len(leaseRetryBackoffs); attempt++ {
		if attempt > 0 {
			if err := sleepWithContext(ctx, leaseRetryBackoffs[attempt-1]); err != nil {
				return err
			}
		}
		identityID := m.identity()
		if identityID == "" {
			return fmt.Errorf("ziti identity id missing")
		}
		lastErr = m.mgmtClient.ExtendIdentityLease(ctx, identityID)
		if lastErr == nil {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if status.Code(lastErr) == codes.NotFound {
			return lastErr
		}
	}

	return lastErr
}

func (m *Manager) identity() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.identityID
}

func retryWithBackoff(ctx context.Context, operationName string, fn func(context.Context) error) error {
	backoff := retryInitialBackoff
	attempt := 1
	for {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !isRetryableGrpcError(err) {
			return err
		}

		delay := backoff
		if delay > retryMaxBackoff {
			delay = retryMaxBackoff
		}

		log.Printf("%s failed (attempt %d), retrying in %s: %v", operationName, attempt, delay, err)

		if err := sleepWithContext(ctx, delay); err != nil {
			return err
		}

		backoff *= 2
		if backoff > retryMaxBackoff {
			backoff = retryMaxBackoff
		}
		attempt++
	}
}

func isRetryableGrpcError(err error) bool {
	statusErr, ok := status.FromError(err)
	if !ok {
		return false
	}
	return statusErr.Code() == codes.Unavailable || statusErr.Code() == codes.Unknown
}

func sleepWithContext(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
