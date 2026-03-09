package gen

import "encoding/json"

func marshalWorkspaceLimit(raw json.RawMessage) ([]byte, error) {
	if len(raw) == 0 {
		return []byte("null"), nil
	}
	return raw, nil
}

func unmarshalWorkspaceLimit(raw *json.RawMessage, data []byte) error {
	if data == nil {
		*raw = nil
		return nil
	}
	*raw = append((*raw)[:0], data...)
	return nil
}

// Request union types.
func (t *PostWorkspaceConfigurationsJSONBody_Config_CpuLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t PostWorkspaceConfigurationsJSONBody_Config_CpuLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

func (t *PostWorkspaceConfigurationsJSONBody_Config_MemoryLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t PostWorkspaceConfigurationsJSONBody_Config_MemoryLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

func (t *PatchWorkspaceConfigurationsIdJSONBody_Config_CpuLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t PatchWorkspaceConfigurationsIdJSONBody_Config_CpuLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

func (t *PatchWorkspaceConfigurationsIdJSONBody_Config_MemoryLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t PatchWorkspaceConfigurationsIdJSONBody_Config_MemoryLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

// Response union types.
type GetWorkspaceConfigurations200JSONResponse_Items_Config_CpuLimit struct {
	union json.RawMessage
}

type GetWorkspaceConfigurations200JSONResponse_Items_Config_MemoryLimit struct {
	union json.RawMessage
}

type PostWorkspaceConfigurations201JSONResponse_Config_CpuLimit struct {
	union json.RawMessage
}

type PostWorkspaceConfigurations201JSONResponse_Config_MemoryLimit struct {
	union json.RawMessage
}

type GetWorkspaceConfigurationsId200JSONResponse_Config_CpuLimit struct {
	union json.RawMessage
}

type GetWorkspaceConfigurationsId200JSONResponse_Config_MemoryLimit struct {
	union json.RawMessage
}

type PatchWorkspaceConfigurationsId200JSONResponse_Config_CpuLimit struct {
	union json.RawMessage
}

type PatchWorkspaceConfigurationsId200JSONResponse_Config_MemoryLimit struct {
	union json.RawMessage
}

func (t *GetWorkspaceConfigurations200JSONResponse_Items_Config_CpuLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t GetWorkspaceConfigurations200JSONResponse_Items_Config_CpuLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

func (t *GetWorkspaceConfigurations200JSONResponse_Items_Config_MemoryLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t GetWorkspaceConfigurations200JSONResponse_Items_Config_MemoryLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

func (t *PostWorkspaceConfigurations201JSONResponse_Config_CpuLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t PostWorkspaceConfigurations201JSONResponse_Config_CpuLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

func (t *PostWorkspaceConfigurations201JSONResponse_Config_MemoryLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t PostWorkspaceConfigurations201JSONResponse_Config_MemoryLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

func (t *GetWorkspaceConfigurationsId200JSONResponse_Config_CpuLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t GetWorkspaceConfigurationsId200JSONResponse_Config_CpuLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

func (t *GetWorkspaceConfigurationsId200JSONResponse_Config_MemoryLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t GetWorkspaceConfigurationsId200JSONResponse_Config_MemoryLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

func (t *PatchWorkspaceConfigurationsId200JSONResponse_Config_CpuLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t PatchWorkspaceConfigurationsId200JSONResponse_Config_CpuLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}

func (t *PatchWorkspaceConfigurationsId200JSONResponse_Config_MemoryLimit) UnmarshalJSON(b []byte) error {
	return unmarshalWorkspaceLimit(&t.union, b)
}

func (t PatchWorkspaceConfigurationsId200JSONResponse_Config_MemoryLimit) MarshalJSON() ([]byte, error) {
	return marshalWorkspaceLimit(t.union)
}
