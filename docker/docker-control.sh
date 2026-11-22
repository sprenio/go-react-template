#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/../.env"
NGINX_CONF="$SCRIPT_DIR/services/nginx/default_server.conf"
NGINX_SAMPLE="$SCRIPT_DIR/services/nginx/default_server.conf.sample"

if [[ ! -f "$ENV_FILE" ]]; then
    echo "Error: file .env doesn't exist in $ENV_FILE" >&2
    exit 1
fi

set -a
source "$ENV_FILE"
set +a

if [[ ! -f "$NGINX_CONF" && -f "$NGINX_SAMPLE" ]]; then
    target_port="${VITE_PORT:-5173}"
    echo "Notice: $NGINX_CONF not found. Generating from $NGINX_SAMPLE with VITE_PORT=$target_port..."
    sed "s/{VITE_PORT}/$target_port/g" "$NGINX_SAMPLE" > "$NGINX_CONF"
fi

usage() {
  echo "Usage: $0 {up|down|build|restart|softrestart|status|raw} [args...]"
  echo
  echo "Commands:"
  echo "  up           Start containers (detached). Passes args to 'docker compose up'."
  echo "  down         Stop containers, remove volumes/orphans. Passes args to 'docker compose down'."
  echo "  build        Rebuild image(s) from scratch (--no-cache). Optional service name."
  echo "  restart      Hard restart (down then up). Passes args to commands."
  echo "  softrestart  Soft restart 'docker compose restart'. Passes args."
  echo "  status       Show status 'docker compose ps'. Passes args."
  echo "  raw          Run arbitrary docker compose command. Example: '$0 raw logs -f'"
  exit 1
}

[ $# -eq 0 ] && usage
action="$1"
shift

cd "$SCRIPT_DIR"

# Ensure database data and logs directories exist and have proper permissions for MariaDB/MySQL (UID:GID 999:999)
ensure_db_permissions() {
  local db_dirs=()

  # Skip changing ownership on macOS (Darwin) – Docker mount layer handles permission mapping automatically
  if [[ "$(uname -s)" == "Darwin" ]]; then
    echo "macOS detected – skipping 999:999 ownership change for DB directories."
    for dir in "${DB_DATA_DIR:-}" "${DB_LOGS_DIR:-}"; do
      [[ -n "$dir" && ! -d "$dir" ]] && mkdir -p "$dir"
    done
    return 0
  fi

  # Linux logic (e.g. Ubuntu production server)
  [[ -n "${DB_DATA_DIR:-}" ]] && db_dirs+=("$DB_DATA_DIR")
  [[ -n "${DB_LOGS_DIR:-}" ]] && db_dirs+=("$DB_LOGS_DIR")

  if [[ ${#db_dirs[@]} -gt 0 ]]; then
    echo "Ensuring DB directories exist and have correct permissions (999:999)..."
    for dir in "${db_dirs[@]}"; do
      if [[ ! -d "$dir" ]]; then
        mkdir -p "$dir"
      fi

      # Official MariaDB/MySQL images run as UID/GID 999
      if [[ $EUID -eq 0 ]]; then
        chown -R 999:999 "$dir"
        chmod -R 755 "$dir"
      else
        sudo chown -R 999:999 "$dir"
        sudo chmod -R 755 "$dir"
      fi
    done
  fi
}

docker_up(){
  echo "Checking DB permissions..."
  ensure_db_permissions

  echo "Starting docker-compose..."
  docker compose --env-file "$ENV_FILE" up -d "$@"
}

docker_down(){
  echo "Stopping docker-compose and cleaning up..."
  docker compose down --volumes --remove-orphans "$@"
}

docker_build(){
  echo "Rebuilding docker-compose image(s) from scratch (--no-cache)..."
  docker compose --env-file "$ENV_FILE" build --no-cache "$@"
}

[ -z "$action" ] && usage

case "$action" in
  up)
    docker_up "$@"
    ;;
  down)
    docker_down "$@"
    ;;
  build)
    docker_build "$@"
    ;;
  restart)
    echo "Restarting docker-compose (hard)..."
    docker_down "$@"
    docker_up "$@"
    ;;
  softrestart)
    echo "Restarting docker-compose (soft)..."
    docker compose --env-file "$ENV_FILE" restart "$@"
    ;;
  status)
    echo "Docker containers status:"
    docker compose ps "$@"
    ;;
  raw)
    echo "Executing raw docker-compose command..."
    docker compose --env-file "$ENV_FILE" "$@"
    ;;
  *)
    usage
    ;;
esac