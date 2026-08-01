package identity

import (
	"context"
	"testing"
)

func TestIdentityContextRoundTrip(t *testing.T) {
	ctx := context.Background()
	input := ResolvedIdentity{
		IdentityID:   "id-123",
		IdentityType: IdentityTypeUser,
		WorkloadID:   "workload-9",
	}

	ctx = WithIdentity(ctx, input)

	got, ok := IdentityFromContext(ctx)
	if !ok {
		t.Fatalf("expected identity in context")
	}
	if got != input {
		t.Fatalf("unexpected identity: %+v", got)
	}
}

func TestIdentityContextMissing(t *testing.T) {
	_, ok := IdentityFromContext(context.Background())
	if ok {
		t.Fatalf("expected no identity")
	}
}

func TestParseIdentityTypeSandbox(t *testing.T) {
	identityType, err := ParseIdentityType("sandbox")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identityType != IdentityTypeSandbox {
		t.Fatalf("unexpected identity type: %s", identityType)
	}
	if !identityType.IsWorkload() {
		t.Fatalf("expected a sandbox identity to be a workload identity")
	}
}

func TestSandboxIDIsTheIdentityItself(t *testing.T) {
	sandbox := ResolvedIdentity{
		IdentityID:   "sandbox-1",
		IdentityType: IdentityTypeSandbox,
		WorkloadID:   "workload-9",
	}
	if sandbox.SandboxID() != "sandbox-1" {
		t.Fatalf("unexpected sandbox id: %s", sandbox.SandboxID())
	}

	agent := ResolvedIdentity{
		IdentityID:   "agent-1",
		IdentityType: IdentityTypeAgent,
		WorkloadID:   "workload-9",
	}
	if agent.SandboxID() != "" {
		t.Fatalf("expected no sandbox id for an agent identity, got %s", agent.SandboxID())
	}
}

func TestParseIdentityTypeApp(t *testing.T) {
	identityType, err := ParseIdentityType("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identityType != IdentityTypeApp {
		t.Fatalf("unexpected identity type: %s", identityType)
	}
}

// An agent workload authenticates as its instance. The type has to parse and
// count as a workload, or the gateway rejects the identity outright and every
// call the workload makes fails at the door.
func TestParseIdentityTypeAcceptsAgentInstance(t *testing.T) {
	parsed, err := ParseIdentityType("agent_instance")
	if err != nil {
		t.Fatalf("parse agent_instance: %v", err)
	}
	if parsed != IdentityTypeAgentInstance {
		t.Fatalf("parsed = %q, want agent_instance", parsed)
	}
	if !parsed.IsWorkload() {
		t.Fatal("agent_instance must count as a workload identity")
	}
}
