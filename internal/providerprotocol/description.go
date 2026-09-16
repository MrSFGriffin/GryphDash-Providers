package providerprotocol

// DataSourceDescription describes one stable provider result namespace.
type DataSourceDescription struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Schema      map[string]any `json:"schema"`
	URL         string         `json:"url,omitempty"`
	Method      string         `json:"method,omitempty"`
}

// ProviderDescription is returned by a provider's describe method. Templates
// and definitions are intentionally optional until GryphDash owns rendering.
type ProviderDescription struct {
	ID              string                  `json:"id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	ProtocolVersion int                     `json:"protocolVersion"`
	Sources         []DataSourceDescription `json:"sources"`
	TemplatePacks   []any                   `json:"templatePacks,omitempty"`
	Definitions     []any                   `json:"definitions,omitempty"`
}

// ObjectSchema builds the small JSON Schema documents used by core providers.
func ObjectSchema(properties map[string]any, required ...string) map[string]any {
	schema := map[string]any{"type": "object", "properties": properties}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}
