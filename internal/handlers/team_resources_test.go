package handlers

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func attachmentTestNodes() (string, string, map[string]graphNode) {
	agentID := "11111111-1111-1111-1111-111111111111"
	toolID := "22222222-2222-2222-2222-222222222222"
	nodes := map[string]graphNode{
		agentID: {ID: agentID, Template: agentTemplateName},
		toolID:  {ID: toolID, Template: "manageTool"},
	}
	return agentID, toolID, nodes
}

func attachmentTestEdge(id, source, target string) graphEdge {
	return graphEdge{
		ID:           ptr(id),
		Source:       source,
		SourceHandle: "tools",
		Target:       target,
		TargetHandle: "$self",
	}
}

func TestAttachmentIDFromEdge(t *testing.T) {
	t.Run("validUUID", func(t *testing.T) {
		raw := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
		expected := uuid.MustParse(raw)
		got := attachmentIDFromEdge(raw)
		if got != expected {
			t.Fatalf("expected %s, got %s", expected, got)
		}
	})

	t.Run("nonUUID", func(t *testing.T) {
		raw := "edge-non-uuid"
		expected := uuid.NewSHA1(uuid.NameSpaceURL, []byte(raw))
		first := attachmentIDFromEdge(raw)
		second := attachmentIDFromEdge(raw)
		if first != expected {
			t.Fatalf("expected %s, got %s", expected, first)
		}
		if second != expected {
			t.Fatalf("expected %s, got %s", expected, second)
		}
		if first != second {
			t.Fatalf("expected idempotent result, got %s and %s", first, second)
		}
	})
}

func TestAttachmentFromEdgeMissingTimestamps(t *testing.T) {
	agentID, toolID, nodes := attachmentTestNodes()
	edgeID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	edge := attachmentTestEdge(edgeID, agentID, toolID)

	attachment, ok, err := attachmentFromEdge(edge, nodes, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected attachment")
	}
	if !attachment.CreatedAt.IsZero() {
		t.Fatalf("expected zero createdAt, got %v", attachment.CreatedAt)
	}
	if attachment.UpdatedAt != nil {
		t.Fatalf("expected nil updatedAt, got %v", attachment.UpdatedAt)
	}
	if attachment.ID != uuid.MustParse(edgeID) {
		t.Fatalf("expected id %s, got %s", edgeID, attachment.ID)
	}
}

func TestAttachmentFromEdgeNonUUIDID(t *testing.T) {
	agentID, toolID, nodes := attachmentTestNodes()
	edgeID := "edge-non-uuid"
	edge := attachmentTestEdge(edgeID, agentID, toolID)
	vars := map[string]string{
		attachmentCreatedAtKey(edgeID): "2024-01-01T00:00:00Z",
	}

	attachment, ok, err := attachmentFromEdge(edge, nodes, vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected attachment")
	}

	expectedID := uuid.NewSHA1(uuid.NameSpaceURL, []byte(edgeID))
	if attachment.ID != expectedID {
		t.Fatalf("expected id %s, got %s", expectedID, attachment.ID)
	}
	if attachment.CreatedAt.IsZero() {
		t.Fatalf("expected createdAt")
	}
}

func TestAttachmentFromEdgeInvalidTimestamps(t *testing.T) {
	agentID, toolID, nodes := attachmentTestNodes()
	edgeID := "cccccccc-cccc-cccc-cccc-cccccccccccc"
	edge := attachmentTestEdge(edgeID, agentID, toolID)

	cases := []struct {
		name    string
		vars    map[string]string
		wantErr string
	}{
		{
			name:    "createdAt",
			vars:    map[string]string{attachmentCreatedAtKey(edgeID): "not-a-time"},
			wantErr: "invalid createdAt",
		},
		{
			name: "updatedAt",
			vars: map[string]string{
				attachmentCreatedAtKey(edgeID): "2024-01-01T00:00:00Z",
				attachmentUpdatedAtKey(edgeID): "not-a-time",
			},
			wantErr: "invalid updatedAt",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := attachmentFromEdge(edge, nodes, tc.vars)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tc.wantErr, err.Error())
			}
		})
	}
}
