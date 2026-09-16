package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MrSFGriffin/GryphDash-Providers/internal/providerprotocol"
)

type artifact struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

type manifest struct {
	Version    int `json:"version"`
	Repository struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"repository"`
	Provider struct {
		ID              string                                   `json:"id"`
		Name            string                                   `json:"name"`
		Description     string                                   `json:"description"`
		Version         string                                   `json:"version"`
		ProtocolVersion int                                      `json:"protocolVersion"`
		Sources         []providerprotocol.DataSourceDescription `json:"sources"`
		TemplatePacks   []any                                    `json:"templatePacks,omitempty"`
		Definitions     []any                                    `json:"definitions,omitempty"`
		Artifacts       map[string]artifact                      `json:"artifacts"`
	} `json:"provider"`
}

type repositoryIndex struct {
	Version    int `json:"version"`
	Repository struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"repository"`
	Providers []manifestProvider `json:"providers"`
}

type manifestProvider = struct {
	ID              string                                   `json:"id"`
	Name            string                                   `json:"name"`
	Description     string                                   `json:"description"`
	Version         string                                   `json:"version"`
	ProtocolVersion int                                      `json:"protocolVersion"`
	Sources         []providerprotocol.DataSourceDescription `json:"sources"`
	TemplatePacks   []any                                    `json:"templatePacks,omitempty"`
	Definitions     []any                                    `json:"definitions,omitempty"`
	Artifacts       map[string]artifact                      `json:"artifacts"`
}

func main() {
	version := flag.String("version", "", "provider version")
	releaseURL := flag.String("release-url", "", "base GitHub release URL")
	artifactsDir := flag.String("artifacts", "", "directory containing release artifacts")
	outputDir := flag.String("output", "", "directory for generated manifests")
	flag.Parse()
	if *version == "" || *releaseURL == "" || *artifactsDir == "" || *outputDir == "" {
		fatal("-version, -release-url, -artifacts, and -output are required")
	}
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fatal("create output directory: %v", err)
	}
	index := repositoryIndex{Version: 1}
	index.Repository.ID, index.Repository.Name, index.Repository.Description = "core", "GryphDash Core Providers", "Official GryphDash providers"
	for _, providerID := range []string{"currency", "codex", "openrouter"} {
		desc, err := descriptionFor(filepath.Join(*artifactsDir, "gryphdash-provider-"+providerID+"-linux-amd64"), providerID)
		if err != nil {
			fatal("%s description: %v", providerID, err)
		}
		item := manifest{Version: 1}
		item.Repository.ID, item.Repository.Name, item.Repository.Description = "core", "GryphDash Core Providers", "Official GryphDash providers"
		item.Provider.ID, item.Provider.Name, item.Provider.Description = providerID, displayName(providerID), description(providerID)
		item.Provider.Version, item.Provider.ProtocolVersion = *version, providerprotocol.Version
		item.Provider.Sources = desc.Sources
		item.Provider.TemplatePacks, item.Provider.Definitions = desc.TemplatePacks, desc.Definitions
		item.Provider.Artifacts = map[string]artifact{}
		for _, platform := range []string{"linux-amd64", "windows-amd64", "darwin-amd64", "darwin-arm64"} {
			name := "gryphdash-provider-" + providerID + "-" + platform
			if strings.HasPrefix(platform, "windows-") {
				name += ".exe"
			}
			path := filepath.Join(*artifactsDir, name)
			checksum, err := fileChecksum(path)
			if err != nil {
				fatal("%s %s checksum: %v", providerID, platform, err)
			}
			item.Provider.Artifacts[platform] = artifact{URL: strings.TrimRight(*releaseURL, "/") + "/" + name, SHA256: checksum}
		}
		data, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			fatal("encode %s manifest: %v", providerID, err)
		}
		if err := os.WriteFile(filepath.Join(*outputDir, "manifest-"+providerID+".json"), append(data, '\n'), 0644); err != nil {
			fatal("write %s manifest: %v", providerID, err)
		}
		index.Providers = append(index.Providers, item.Provider)
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		fatal("encode repository index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(*outputDir, "repository.json"), append(data, '\n'), 0644); err != nil {
		fatal("write repository index: %v", err)
	}
}

func descriptionFor(path, providerID string) (providerprotocol.ProviderDescription, error) {
	request, err := json.Marshal(providerprotocol.Request{Version: providerprotocol.Version, Method: "describe"})
	if err != nil {
		return providerprotocol.ProviderDescription{}, err
	}
	command := exec.Command(path, "-describe")
	command.Stdin = strings.NewReader(string(request) + "\n")
	data, err := command.Output()
	if err != nil {
		return providerprotocol.ProviderDescription{}, err
	}
	var response providerprotocol.Response
	if err := json.Unmarshal(data, &response); err != nil {
		return providerprotocol.ProviderDescription{}, err
	}
	if response.Error != "" {
		return providerprotocol.ProviderDescription{}, fmt.Errorf("provider error: %s", response.Error)
	}
	if response.Description == nil || response.Description.ID != providerID {
		return providerprotocol.ProviderDescription{}, fmt.Errorf("provider returned invalid description")
	}
	return *response.Description, nil
}

func fileChecksum(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func displayName(id string) string {
	return map[string]string{"currency": "Currency", "codex": "Codex", "openrouter": "OpenRouter"}[id]
}
func description(id string) string {
	return map[string]string{"currency": "Reference exchange rates", "codex": "Codex account and usage metrics", "openrouter": "OpenRouter API key usage and credits"}[id]
}
func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
