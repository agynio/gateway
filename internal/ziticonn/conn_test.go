package ziticonn

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestSourceIdentityFromConnDialerID(t *testing.T) {
	conn := dialerConn{dialerID: "  dialer-id  "}

	identity, ok := SourceIdentityFromConn(conn)
	if !ok {
		t.Fatalf("expected identity")
	}
	if identity != "dialer-id" {
		t.Fatalf("expected dialer id, got %q", identity)
	}
}

func TestSourceIdentityFromConnMissingIdentity(t *testing.T) {
	_, ok := SourceIdentityFromConn(stubConn{})
	if ok {
		t.Fatalf("expected no identity")
	}
}

func TestSourceIdentityFromConnDialerOnlyEmptyFallsThrough(t *testing.T) {
	_, ok := SourceIdentityFromConn(dialerConn{dialerID: "   "})
	if ok {
		t.Fatalf("expected no identity")
	}
}

type stubConn struct{}

func (stubConn) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (stubConn) Write(data []byte) (int, error) {
	return len(data), nil
}

func (stubConn) Close() error {
	return nil
}

func (stubConn) LocalAddr() net.Addr {
	return stubAddr{}
}

func (stubConn) RemoteAddr() net.Addr {
	return stubAddr{}
}

func (stubConn) SetDeadline(time.Time) error {
	return nil
}

func (stubConn) SetReadDeadline(time.Time) error {
	return nil
}

func (stubConn) SetWriteDeadline(time.Time) error {
	return nil
}

type stubAddr struct{}

func (stubAddr) Network() string {
	return "stub"
}

func (stubAddr) String() string {
	return "stub"
}

type dialerConn struct {
	stubConn
	dialerID string
}

func (conn dialerConn) GetDialerIdentityId() string {
	return conn.dialerID
}
