#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
output_dir="${1:-$repo_root/dist}"
mkdir -p "$output_dir"

for platform in linux-amd64 windows-amd64 darwin-amd64 darwin-arm64; do
  IFS=- read -r goos goarch <<< "$platform"
  for provider in currency codex openrouter; do
    suffix=""
    if [[ "$goos" == windows ]]; then suffix=".exe"; fi
    output_file="$output_dir/gryphdash-provider-$provider-$platform$suffix"
    (cd "$repo_root" && GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -trimpath -o "$output_file" "./cmd/gryphdash-provider-$provider")
    printf 'Built %s\n' "$output_file"
  done
done

sha256sum "$output_dir"/gryphdash-provider-* > "$output_dir/SHA256SUMS"
