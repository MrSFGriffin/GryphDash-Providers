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
		ID              string              `json:"id"`
		Name            string              `json:"name"`
		Description     string              `json:"description"`
		Version         string              `json:"version"`
		ProtocolVersion int                 `json:"protocolVersion"`
		Widgets         interface{}         `json:"widgets"`
		Artifacts       map[string]artifact `json:"artifacts"`
	} `json:"provider"`
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
	for _, providerID := range []string{"currency", "codex", "openrouter"} {
		widgets, err := widgetsFor(filepath.Join(*artifactsDir, "gryphdash-provider-"+providerID+"-linux-amd64"))
		if err != nil {
			fatal("%s widgets: %v", providerID, err)
		}
		item := manifest{Version: 1}
		item.Repository.ID, item.Repository.Name, item.Repository.Description = "core", "GryphDash Core Providers", "Official GryphDash providers"
		item.Provider.ID, item.Provider.Name, item.Provider.Description = providerID, displayName(providerID), description(providerID)
		item.Provider.Version, item.Provider.ProtocolVersion = *version, providerprotocol.Version
		item.Provider.Widgets, item.Provider.Artifacts = widgets, map[string]artifact{}
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
	}
}

func widgetsFor(path string) (interface{}, error) {
	data, err := exec.Command(path, "-widgets").Output()
	if err != nil {
		return nil, err
	}
	var response providerprotocol.Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, fmt.Errorf("provider error: %s", response.Error)
	}
	return response.Widgets, nil
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
