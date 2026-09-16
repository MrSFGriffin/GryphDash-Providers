package openrouter

import "github.com/MrSFGriffin/GryphDash-Providers/internal/providerprotocol"

func Description() providerprotocol.ProviderDescription {
	return providerprotocol.ProviderDescription{
		ID: "openrouter", Name: "OpenRouter", Description: "OpenRouter API key usage and credits", ProtocolVersion: providerprotocol.Version,
		Sources: []providerprotocol.DataSourceDescription{
			{ID: "openrouter/key", Name: "OpenRouter key", Description: "Usage and limits for the configured API key", Schema: providerprotocol.ObjectSchema(map[string]any{"usage_monthly": map[string]any{"type": "number"}, "usage_daily": map[string]any{"type": "number"}, "usage_weekly": map[string]any{"type": "number"}, "limit": map[string]any{"type": "number"}, "limit_remaining": map[string]any{"type": "number"}, "limit_reset": map[string]any{"type": "string"}, "expires_at": map[string]any{"type": "string"}}), URL: "https://openrouter.ai/api/v1/key", Method: "GET"},
			{ID: "openrouter/credits", Name: "OpenRouter credits", Description: "Purchased and consumed OpenRouter credits", Schema: providerprotocol.ObjectSchema(map[string]any{"total_credits": map[string]any{"type": "number"}, "total_usage": map[string]any{"type": "number"}}), URL: "https://openrouter.ai/api/v1/credits", Method: "GET"},
		},
	}
}
