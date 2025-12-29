#!/bin/sh
set -e

# Run all executable scripts in /entrypoint.d/
if [ -d "/entrypoint.d" ]; then
  for f in $(find /entrypoint.d -type f -executable | sort); do
    echo "[entrypoint] running $f"
    "$f"
  done
fi

echo "[entrypoint] starting main process: $@"

# If CMD/command is provided, run it
if [ $# -gt 0 ]; then
  exec "$@"
fi