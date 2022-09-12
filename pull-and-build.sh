#!/bin/bash
# Based on http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

# Based on https://stackoverflow.com/a/246128
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
BASE_DIR="${SCRIPT_DIR}/.."

# Pull latest version
# Ignoring failures to prevent blocking if not internet present
git pull --ff-only || true

# Compile binary
cd "${BASE_DIR}/src"
go build -o ../sunrise-alarm ./cmd/
