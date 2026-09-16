package codex

import "github.com/MrSFGriffin/GryphDash-Providers/internal/providerprotocol"

func Description() providerprotocol.ProviderDescription {
	return providerprotocol.ProviderDescription{
		ID: "codex", Name: "Codex", Description: "Codex account and usage metrics", ProtocolVersion: providerprotocol.Version,
		Sources: []providerprotocol.DataSourceDescription{
			{ID: "codex/account", Name: "Codex account", Description: "Authenticated Codex account details", Schema: providerprotocol.ObjectSchema(map[string]any{"planType": map[string]any{"type": "string"}, "type": map[string]any{"type": "string"}})},
			{ID: "codex/limits", Name: "Codex limits", Description: "Account limit windows and reset credits", Schema: providerprotocol.ObjectSchema(map[string]any{"buckets": map[string]any{"type": "object"}, "rateLimitResetCredits": map[string]any{"type": "object"}})},
			{ID: "codex/usage", Name: "Codex usage", Description: "Codex token usage summaries and daily activity", Schema: providerprotocol.ObjectSchema(map[string]any{"summary": map[string]any{"type": "object"}, "dailyUsageBuckets": map[string]any{"type": "array"}})},
		},
	}
}
