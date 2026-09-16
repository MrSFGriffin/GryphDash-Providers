package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MrSFGriffin/GryphDash-Providers/internal/providerprotocol"
	"github.com/MrSFGriffin/GryphDash-Providers/providers/openrouter"
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
	if describe {
		respond(providerprotocol.Response{Version: providerprotocol.Version, Provider: "openrouter", Description: pointer(openrouter.Description())})
		return
	}
	results := openrouter.NewAdapter(os.Getenv("OPENROUTER_API_KEY"), &http.Client{}).Read(context.Background())
	respond(providerprotocol.Response{Version: providerprotocol.Version, Provider: "openrouter", Results: results})
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
	_ = json.NewEncoder(os.Stdout).Encode(providerprotocol.Response{Version: providerprotocol.Version, Provider: "openrouter", Error: err.Error()})
	os.Exit(1)
}
