#!/usr/bin/env bash
set -e

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Path to the docker control script (relative to .agent/scripts/)
CONTROL_SCRIPT="$SCRIPT_DIR/../../docker/docker-control.sh"

# We use 'sh -c' to ensure variables are expanded inside the container
"$CONTROL_SCRIPT" raw exec -T webserver sh -c "echo \"\${FRONTEND_BASE_URL}\""
