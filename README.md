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

Publish the resulting binaries and `SHA256SUMS` at the artifact URLs referenced
by the GryphDash provider repository manifests.
