package gen

import (
	"encoding/json"
	"testing"
)

func TestWorkspaceLimitRoundTripStrings(t *testing.T) {
	input := `{"config":{"memory_limit":"2Gi","cpu_limit":"500m"}}`

	var payload PostWorkspaceConfigurationsJSONRequestBody
	if err := json.Unmarshal([]byte(input), &payload); err != nil {
		t.Fatalf("unmarshal input: %v", err)
	}

	output, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	config := decodeConfig(t, output)
	assertConfigValue(t, config, "memory_limit", "2Gi")
	assertConfigValue(t, config, "cpu_limit", "500m")
}

func TestWorkspaceLimitRoundTripNumbers(t *testing.T) {
	input := `{"config":{"memory_limit":4096,"cpu_limit":4}}`

	var payload PostWorkspaceConfigurationsJSONRequestBody
	if err := json.Unmarshal([]byte(input), &payload); err != nil {
		t.Fatalf("unmarshal input: %v", err)
	}

	output, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	config := decodeConfig(t, output)
	assertConfigValue(t, config, "memory_limit", float64(4096))
	assertConfigValue(t, config, "cpu_limit", float64(4))
}

func decodeConfig(t *testing.T, payload []byte) map[string]any {
	t.Helper()

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	configRaw, ok := decoded["config"]
	if !ok {
		t.Fatalf("missing config in output")
	}

	config, ok := configRaw.(map[string]any)
	if !ok {
		t.Fatalf("config has unexpected type %T", configRaw)
	}

	return config
}

func assertConfigValue(t *testing.T, config map[string]any, key string, expected any) {
	t.Helper()

	value, ok := config[key]
	if !ok {
		t.Fatalf("missing %s in output", key)
	}

	if value != expected {
		t.Fatalf("%s mismatch: got %v want %v", key, value, expected)
	}
}
