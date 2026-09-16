package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MrSFGriffin/GryphDash-Providers/internal/providerprotocol"
	"github.com/MrSFGriffin/GryphDash-Providers/providers/currency"
)

func main() {
	describe := len(os.Args) > 1 && os.Args[1] == "-describe"
	var request providerprotocol.Request
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fail(fmt.Errorf("invalid provider request: %w", err))
	}
	if request.Version != providerprotocol.Version || (describe && request.Method != "describe") || (!describe && request.Method != "read") {
		fail(fmt.Errorf("unsupported provider request"))
	}
	pairs := currency.ParsePairs(os.Getenv("GRYPHDASH_CURRENCY_PAIRS"))
	if len(pairs) == 0 {
		pairs = []currency.Pair{{Base: "EUR", Quote: "USD"}, {Base: "EUR", Quote: "GBP"}, {Base: "GBP", Quote: "EUR"}, {Base: "EUR", Quote: "HUF"}}
	}
	if describe {
		respond(providerprotocol.Response{Version: providerprotocol.Version, Provider: "currency", Description: pointer(currency.Descriptions(pairs))})
		return
	}
	result := currency.Client{HTTPClient: &http.Client{}, BaseURL: os.Getenv("GRYPHDASH_CURRENCY_URL"), Pairs: pairs}.Read(context.Background())
	respond(providerprotocol.Response{Version: providerprotocol.Version, Provider: "currency", Results: result})
}

func pointer(value providerprotocol.ProviderDescription) *providerprotocol.ProviderDescription {
	return &value
}

func respond(response providerprotocol.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fail(err)
	}
}

func fail(err error) {
	_ = json.NewEncoder(os.Stdout).Encode(providerprotocol.Response{Version: providerprotocol.Version, Provider: "currency", Error: err.Error()})
	os.Exit(1)
}
