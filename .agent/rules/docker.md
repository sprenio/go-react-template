---
trigger: model_decision
description: This rule should be activated ONLY when the user asks to perform Docker operations (start, stop, etc), manage containers, or requires specific architectural information from the docker-compose.yml file (service names, ports, etc)
---

# Docker Infrastructure & Management Rules

## 1. Source of Truth
- **Main Compose File:** `docker/docker-compose.yml`
- **Environment Context:** Always assume the `.env` file is required but managed via the control script.

## 2. Service Inventory & Architecture
| Service Name | Technology | Role / Responsibility |
| :--- | :--- | :--- |
| `db` | MariaDB 13.0 | Primary transactional database. |
| `db-migrator` | Custom Script (PHP) | Runs DB migrations on startup and exits immediately. |
| `webserver` | Go 1.26 (API) | Main backend server handling REST/API requests. |
| `consumer` | Go 1.26 (Worker) | RabbitMQ consumer processing background queues. |
| `cron` | Go 1.26 (Worker) | Periodic background jobs and scheduled tasks runner. |
| `frontend` | Nginx 1.29 | Web server hosting frontend SPA. |
| `rabbitmq` | RabbitMQ 4.1 | Message broker for inter-service communication. |
| `redis` | Redis 8.8 (Alpine) | In-memory cache, distributed locks and rate limiting. |

## 3. Mandatory Execution Method
**DO NOT** run `docker compose` commands directly. All operations must be performed using the wrapper script to ensure proper `.env` loading and path resolution.

- **Control Script Path:** `docker/docker-control.sh`
- **Usage Pattern:** `./docker/docker-control.sh {up|down|build|restart|softrestart|status|raw} [args...]`
    Environment:
      Config       Using compose file: 'docker-compose.yml'

    Commands:
      up           Start containers (detached). Passes args to 'docker compose up'.
      down         Stop containers, remove volumes/orphans. Passes args to 'docker compose down'.
      build        Rebuild image(s) from scratch (--no-cache). Optional service name.
      restart      Hard restart (down then up). Passes args to commands.
      softrestart  Soft restart 'docker compose restart'. Passes args.
      status       Show status 'docker compose ps'. Passes args.
      raw          Run arbitrary docker compose command. Example: './docker/docker-control.sh raw logs -f'


## 4. Operational Guardrails
- **Migration Logic:** If the user asks to "reset" or "update" the database, ensure `db-migrator` is triggered/restarted.
- **Nginx Routing:** Remember that services are accessed via hostnames managed by Nginx. No need to suggest high-number ports for local dev access.
- **Safety:** Do not use `sudo` within the script calls unless the system explicitly returns a permission error.
