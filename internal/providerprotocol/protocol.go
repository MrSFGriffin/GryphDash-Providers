// Package providerprotocol defines the wire format used by external providers.
package providerprotocol

import "github.com/MrSFGriffin/GryphDash-Providers/internal/metrics"

const Version = 2

type Request struct {
	Version int    `json:"version"`
	Method  string `json:"method"`
}

type Response struct {
	Version     int                       `json:"version"`
	Provider    string                    `json:"provider"`
	Description *ProviderDescription      `json:"description,omitempty"`
	Results     map[string]metrics.Result `json:"results,omitempty"`
	Error       string                    `json:"error,omitempty"`
}
