#!/bin/bash
# Bash strict mode: http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail

generate_mocks () {
    # Switch to directory
    cd $1

    # Clear current mocks
    rm -rf mocks/

    # Regenerate mocks
    mockery --all --keeptree

    # Restore original directory
    cd -
}

generate_mocks services/alarm
generate_mocks services/physical
