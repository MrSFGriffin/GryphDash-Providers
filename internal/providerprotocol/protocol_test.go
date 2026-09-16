package providerprotocol

import (
	"encoding/json"
	"testing"
)

func TestDescriptionResponseShape(t *testing.T) {
	response := Response{Version: Version, Provider: "test", Description: &ProviderDescription{
		ID: "test", Name: "Test", Description: "Test provider", ProtocolVersion: Version,
		Sources: []DataSourceDescription{{ID: "test/value", Name: "Value", Description: "A value", Schema: map[string]any{"type": "object"}}},
	}}
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["description"] == nil || decoded["widgets"] != nil {
		t.Fatalf("response = %s", data)
	}
}
