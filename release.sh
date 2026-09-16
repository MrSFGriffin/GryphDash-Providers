#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
  printf 'Usage: %s VERSION\n' "$0" >&2
  printf 'Example: %s 1.0.0\n' "$0" >&2
  exit 2
fi

version="$1"
tag="$version"
if [[ "$tag" != v* ]]; then
  tag="v$tag"
fi
if [[ ! "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([-.][0-9A-Za-z.-]+)?$ ]]; then
  printf 'Version must look like 1.2.3 or v1.2.3\n' >&2
  exit 2
fi

repo_root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
release_dir=$(mktemp -d "${TMPDIR:-/tmp}/gryphdash-providers-${tag}-XXXXXX")
manifest_dir="$release_dir/manifests"
stable_manifest_dir="$repo_root/manifests"
trap 'rm -rf "$release_dir"' EXIT

if ! command -v gh >/dev/null 2>&1; then
  printf 'GitHub CLI (gh) is required. Install it and run gh auth login.\n' >&2
  exit 1
fi
if ! gh auth status >/dev/null 2>&1; then
  printf 'GitHub CLI is not authenticated. Run gh auth login first.\n' >&2
  exit 1
fi
if [[ -z "$(git -C "$repo_root" config user.name)" || -z "$(git -C "$repo_root" config user.email)" ]]; then
  printf 'Git author identity is required. Configure git user.name and user.email first.\n' >&2
  exit 1
fi

printf 'Building provider release %s…\n' "$tag"
(cd "$repo_root" && ./build-release.sh "$release_dir")

release_url="https://github.com/MrSFGriffin/GryphDash-Providers/releases/download/$tag"
(cd "$repo_root" && go run ./cmd/generate-manifests \
  -version "${tag#v}" \
  -release-url "$release_url" \
  -artifacts "$release_dir" \
  -output "$manifest_dir")
cp "$manifest_dir"/manifest-*.json "$manifest_dir"/repository.json "$release_dir"/

mapfile -t assets < <(find "$release_dir" -maxdepth 1 -type f \( -name 'gryphdash-provider-*' -o -name 'SHA256SUMS' -o -name 'manifest-*.json' -o -name 'repository.json' \) | sort)
if [[ ${#assets[@]} -ne 17 ]]; then
  printf 'Expected 17 release assets, found %d\n' "${#assets[@]}" >&2
  exit 1
fi

if gh release view "$tag" >/dev/null 2>&1; then
  printf 'Uploading assets to existing release %s…\n' "$tag"
  gh release upload "$tag" "${assets[@]}" --clobber
else
  printf 'Creating GitHub release %s…\n' "$tag"
  gh release create "$tag" "${assets[@]}" \
    --title "GryphDash Providers $tag" \
    --generate-notes
fi

printf 'Publishing stable manifests to the provider repository main branch…\n'
mkdir -p "$stable_manifest_dir"
cp "$manifest_dir"/manifest-*.json "$stable_manifest_dir"/
cp "$manifest_dir"/repository.json "$repo_root"/
(cd "$repo_root" && git add manifests repository.json && if git diff --cached --quiet -- manifests repository.json; then
  printf 'Stable manifests are already current.\n'
else
  git commit -m "Publish provider manifests for $tag"
  git push
fi)

printf 'Release ready: %s\n' "$(gh release view "$tag" --json url --jq .url)"
