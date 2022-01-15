#!/bin/bash
# Bash strict mode: http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail

# Clear current mocks
rm -rf mocks/

# Regenerate mocks
mockery --all --keeptree