#!/usr/bin/env bash
set -e

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Path to the docker control script (relative to .agent/scripts/)
CONTROL_SCRIPT="$SCRIPT_DIR/../../docker/docker-control.sh"

# Check if query argument is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <SQL_QUERY>"
    echo "Example: $0 \"SELECT * FROM users LIMIT 1;\""
    exit 1
fi

QUERY="$1"

# Execute the query inside the container using environment variables for credentials
# We use 'sh -c' to ensure variables are expanded inside the container
"$CONTROL_SCRIPT" raw exec -T db sh -c "mariadb -u\"\${MYSQL_USER}\" -p\"\${MYSQL_PASSWORD}\" \"\${MYSQL_DATABASE}\" -e \"$QUERY\""
