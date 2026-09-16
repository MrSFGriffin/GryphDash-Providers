package currency

import (
	"fmt"

	"github.com/MrSFGriffin/GryphDash-Providers/internal/providerprotocol"
)

// Descriptions returns one concrete source for every configured currency pair.
func Descriptions(pairs []Pair) providerprotocol.ProviderDescription {
	sources := make([]providerprotocol.DataSourceDescription, 0, len(pairs))
	for _, pair := range pairs {
		id := pair.ID()
		sources = append(sources, providerprotocol.DataSourceDescription{
			ID: id, Name: fmt.Sprintf("%s/%s exchange rate", pair.Base, pair.Quote),
			Description: fmt.Sprintf("%s per %s · Frankfurter reference rate", pair.Quote, pair.Base),
			Schema:      providerprotocol.ObjectSchema(map[string]any{"date": map[string]any{"type": "string"}, "base": map[string]any{"type": "string"}, "quote": map[string]any{"type": "string"}, "rate": map[string]any{"type": "number"}}, "date", "base", "quote", "rate"),
			URL:         "https://api.frankfurter.dev/v2/rate/" + pair.Base + "/" + pair.Quote, Method: "GET",
		})
	}
	return providerprotocol.ProviderDescription{ID: "currency", Name: "Currency", Description: "Reference exchange rates", ProtocolVersion: providerprotocol.Version, Sources: sources}
}
