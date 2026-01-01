#!/usr/bin/env sh

# another comment block

# Check if required commands are available
#
# Usage:
# require <cmd1> <cmd2>
#
# Example:
# require curl yq
#
# output:
# Missing curl
# Missing yq
require() {
    missing=0
    for cmd in "$@"; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            echo "Missing: $cmd" >&2
            missing=1
        fi
    done
    [ "$missing" -eq 0 ] || exit 1
}