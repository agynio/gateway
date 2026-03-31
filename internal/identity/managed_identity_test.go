package identity

import "testing"

func TestParseManagedIdentityName(t *testing.T) {
	validUUID := "9e2b5e44-965a-4a6d-a1ea-2cf9b8b49ab2"
	cases := []struct {
		name   string
		input  string
		wantOK bool
		wantID string
	}{
		{
			name:   "valid",
			input:  "agent-" + validUUID + "-acde1234",
			wantOK: true,
			wantID: validUUID,
		},
		{
			name:   "hyphenated suffix",
			input:  "agent-" + validUUID + "-ab-cd",
			wantOK: true,
			wantID: validUUID,
		},
		{
			name:   "trim whitespace",
			input:  "  agent-" + validUUID + "-acde1234  ",
			wantOK: true,
			wantID: validUUID,
		},
		{
			name:   "empty input",
			input:  "",
			wantOK: false,
		},
		{
			name:   "non-agent prefix",
			input:  "svc-runner-acde1234",
			wantOK: false,
		},
		{
			name:   "invalid uuid",
			input:  "agent-not-a-uuid-acde1234",
			wantOK: false,
		},
		{
			name:   "prefix only",
			input:  "agent-",
			wantOK: false,
		},
		{
			name:   "no suffix",
			input:  "agent-" + validUUID,
			wantOK: false,
		},
		{
			name:   "empty suffix",
			input:  "agent-" + validUUID + "-",
			wantOK: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			identity, ok := ParseManagedIdentityName(tc.input)
			if ok != tc.wantOK {
				t.Fatalf("expected ok=%t, got %t", tc.wantOK, ok)
			}
			if !tc.wantOK {
				return
			}
			if identity.IdentityID != tc.wantID {
				t.Fatalf("unexpected identity id %q", identity.IdentityID)
			}
			if identity.IdentityType != IdentityTypeAgent {
				t.Fatalf("unexpected identity type %q", identity.IdentityType)
			}
		})
	}
}
