#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/../.env"

if [[ ! -f "$ENV_FILE" ]]; then
    echo "Error: file .env doesn't exist in $ENV_FILE" >&2
    exit 1
fi

set -a
source "$ENV_FILE"
set +a

APP_ENV="${APP_ENV:-prod}"

COMPOSE_FILE="docker-compose.${APP_ENV}.yml"

if [ ! -f "${SCRIPT_DIR}/${COMPOSE_FILE}" ]; then
  echo "[ERROR] File \"$COMPOSE_FILE\" does not exist."
  echo "Please set APP_ENV correctly or create the corresponding compose file."
  exit 1
fi

usage() {
  echo "Usage: $0 {up|down|restart|softrestart|status}"
  echo
  echo "Commands:"
  echo "  up           Start docker-compose containers in detached mode"
  echo "  down         Stop containers and remove volumes and orphan containers"
  echo "  restart      Restart docker-compose (hard: down + up)"
  echo "  softrestart  Restart docker-compose (soft: docker-compose restart)"
  echo "  status       Show status of docker-compose containers"
  echo "  raw          Execute a raw docker-compose command (pass additional args)"
  exit 1
}


[ $# -eq 0 ] && usage
action="$1"
shift

cd "$SCRIPT_DIR"

docker_up(){
  echo "Starting docker-compose..."
  docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d "$@"
}

docker_down(){
  echo "Stopping docker-compose and cleaning up..."
  docker compose -f "$COMPOSE_FILE" down --volumes --remove-orphans "$@"
}



[ -z "$action" ] && usage

case "$action" in
  up)
    docker_up "$@"
    ;;
  down)
    docker_down "$@"
    ;;
  restart)
    echo "Restarting docker-compose (hard)..."
    docker_down "$@"
    docker_up "$@"
    ;;
  softrestart)
    echo "Restarting docker-compose (soft)..."
    docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" restart "$@"
    ;;
  status)
    echo "Docker containers status:"
    docker compose -f "$COMPOSE_FILE" ps "$@"
    ;;
  raw)
    echo "Executing raw docker-compose command..."
    docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" "$@"
    ;;
  *)
    usage
    ;;
esac
