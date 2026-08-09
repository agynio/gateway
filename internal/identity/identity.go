package identity

import (
	"context"
	"fmt"
	"strings"
)

type IdentityType string

const (
	IdentityTypeUser IdentityType = "user"
	// IdentityTypeAgent names an agent class. It predates instances and is kept
	// because identities minted before the migration still present it.
	IdentityTypeAgent IdentityType = "agent"
	// IdentityTypeAgentInstance names one running instance of an agent class.
	// Workloads authenticate as this: an instance owns its inbox, its volumes
	// and its runner pinning, none of which the class can stand in for.
	IdentityTypeAgentInstance IdentityType = "agent_instance"
	IdentityTypeApp           IdentityType = "app"
	IdentityTypeRunner        IdentityType = "runner"
	IdentityTypeSandbox       IdentityType = "sandbox"
	// IdentityTypePlatform names the platform admin identity, which provisions
	// what a release ships. It has no account in Users and never becomes one.
	IdentityTypePlatform IdentityType = "platform"
)

// IsWorkload reports whether the type belongs to a workload the platform starts
// and stops — one whose identity is only meaningful together with the workload
// it runs in.
func (t IdentityType) IsWorkload() bool {
	return t == IdentityTypeAgent || t == IdentityTypeAgentInstance || t == IdentityTypeSandbox
}

type ResolvedIdentity struct {
	IdentityID   string
	IdentityType IdentityType
	WorkloadID   string
}

// SandboxID is the sandbox record a sandbox workload identity belongs to, the
// way WorkloadID is the workload an agent identity runs in. A sandbox workload
// authenticates as its sandbox — Ziti Management registers the managed identity
// with the sandbox id as its identity id — so the two values are one and the
// same, and callers that need the record read it from here rather than assuming
// an identity id doubles as one.
func (r ResolvedIdentity) SandboxID() string {
	if r.IdentityType != IdentityTypeSandbox {
		return ""
	}
	return strings.TrimSpace(r.IdentityID)
}

func ParseIdentityType(value string) (IdentityType, error) {
	trimmed := strings.TrimSpace(value)
	switch trimmed {
	case string(IdentityTypeUser):
		return IdentityTypeUser, nil
	case string(IdentityTypeAgent):
		return IdentityTypeAgent, nil
	case string(IdentityTypeAgentInstance):
		return IdentityTypeAgentInstance, nil
	case string(IdentityTypeApp):
		return IdentityTypeApp, nil
	case string(IdentityTypeRunner):
		return IdentityTypeRunner, nil
	case string(IdentityTypeSandbox):
		return IdentityTypeSandbox, nil
	case string(IdentityTypePlatform):
		return IdentityTypePlatform, nil
	default:
		return "", fmt.Errorf("unsupported identity type: %q", value)
	}
}

type contextKey struct{}

func WithIdentity(ctx context.Context, identity ResolvedIdentity) context.Context {
	return context.WithValue(ctx, contextKey{}, identity)
}

func IdentityFromContext(ctx context.Context) (ResolvedIdentity, bool) {
	identity, ok := ctx.Value(contextKey{}).(ResolvedIdentity)
	return identity, ok
}
