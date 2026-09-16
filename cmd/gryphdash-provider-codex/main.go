package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MrSFGriffin/GryphDash-Providers/internal/providerprotocol"
	"github.com/MrSFGriffin/GryphDash-Providers/providers/codex"
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
		respond(providerprotocol.Response{Version: providerprotocol.Version, Provider: "codex", Description: pointer(codex.Description())})
		return
	}
	adapter := codex.NewAdapter(os.Getenv("GRYPHDASH_CODEX_BIN"))
	if os.Getenv("GRYPHDASH_CODEX_BIN") == "" {
		adapter = codex.NewAdapter("codex")
	}
	results := adapter.Read(context.Background())
	respond(providerprotocol.Response{Version: providerprotocol.Version, Provider: "codex", Results: results})
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
	_ = json.NewEncoder(os.Stdout).Encode(providerprotocol.Response{Version: providerprotocol.Version, Provider: "codex", Error: err.Error()})
	os.Exit(1)
}
