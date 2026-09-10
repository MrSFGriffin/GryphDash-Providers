# GryphDash Providers

Core external providers for [GryphDash](https://github.com/MrSFGriffin/GryphDash):
Codex, Currency, and OpenRouter. Each command speaks GryphDash's versioned
line-delimited subprocess protocol and supports `-widgets`.

## Build

Build the providers into a directory used by a local GryphDash checkout:

```sh
mkdir -p ../GryphDash/bin
for provider in currency codex openrouter; do
  go build -o "../GryphDash/bin/gryphdash-provider-$provider" "./cmd/gryphdash-provider-$provider"
done
GRYPHDASH_PROVIDER_DIR="../GryphDash/bin" ../GryphDash/bin/gryphdash
```

This repository is the release source for managed provider artifacts. The
GryphDash repository still contains identical in-repository provider sources
and build scripts for offline development and fallback when no managed copy is
installed. A release can cross-compile each command for Linux amd64, Windows
amd64, macOS amd64, and macOS arm64, then publish checksummed artifacts in a
provider repository manifest.

Providers own their service API calls and credentials. Tests use local fixtures;
they never require a live Codex login or third-party account.

To create release artifacts for all supported platforms:

```sh
./build-release.sh
```

To build and publish the complete GitHub release in one command, authenticate
with `gh auth login` first and provide a semantic version:

```sh
./release.sh 1.0.0
```

The script builds all platform binaries, generates `manifest-currency.json`,
`manifest-codex.json`, and `manifest-openrouter.json` from the provider widget
responses and checksums, creates the `v1.0.0` release when it does not exist, or
replaces its assets when it does exist. Temporary build artifacts are removed
after upload.

Publish the resulting binaries and `SHA256SUMS` at the artifact URLs referenced
by the GryphDash provider repository manifests.

See the GryphDash [provider repository guide](https://github.com/MrSFGriffin/GryphDash/blob/main/PROVIDER_REPOSITORIES.md)
for the versioned manifest format, HTTPS/checksum trust model, repository setup,
installation lifecycle, and managed-provider fallback rules.
