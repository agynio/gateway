package gen

// Problem represents the RFC 7807 error schema emitted by the gateway.
type Problem struct {
	Type     string               `json:"type,omitempty"`
	Title    string               `json:"title"`
	Status   int32                `json:"status"`
	Detail   *string              `json:"detail,omitempty"`
	Instance *string              `json:"instance,omitempty"`
	Errors   *map[string][]string `json:"errors,omitempty"`
}
