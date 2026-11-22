# Docker – Development Environment

This directory contains the Docker configuration and development tooling for local development of the application.

> **Production Deployment:** Production configuration has been moved to a dedicated deployment repository (e.g., using images built by `.github/workflows/release.yml`). Do not use this directory or `docker-control.sh` for production deployments.

---

## 1. Directory Structure

```text
docker/
├── docker-compose.yml              # Docker Compose configuration for dev services
├── docker-control.sh               # Helper CLI script for managing dev containers
├── release.sh                      # Tag creation and release script
└── services/
    ├── backend/                    # Dockerfiles and entrypoints for Go backend
    ├── db/                         # Init scripts for local MariaDB
    ├── db-migrator/                # Migrator runner Dockerfile
    ├── frontend/                   # Frontend Dockerfile (used by CI/CD)
    └── nginx/                      # Local Nginx proxy configuration
        ├── default_server.conf.sample  # Sample tracked Nginx dev config
        └── default_server.conf         # Local ignored config (auto-created if missing)
```

The application source code (`backend/`, `frontend/`, `db_migrations/`) is mounted directly into the development containers.

---

## 2. Managing the Development Environment

Use `docker-control.sh` to manage your local Docker environment. The script automatically sources variables from the root `.env` file and provides convenient shortcuts.

Always run commands from this directory:

```bash
cd docker
```

### Start the environment

Start all services in detached mode:

```bash
./docker-control.sh up
```

To build or rebuild images before starting:

```bash
./docker-control.sh up --build
```

### Check container status

```bash
./docker-control.sh status
```

### Rebuild images

Rebuild all images without using cache:

```bash
./docker-control.sh build
```

Rebuild a specific service:

```bash
./docker-control.sh build webserver
./docker-control.sh build consumer
./docker-control.sh build cron
```

### Restart containers

Restart running containers without recreating them (soft restart):

```bash
./docker-control.sh softrestart
```

Perform a full recreate/restart:

```bash
./docker-control.sh restart
```

### Stop the environment

Stop containers and remove project volumes:

```bash
./docker-control.sh down
```

> [!WARNING]
> The `down` command removes Docker volumes (`--volumes`). Be aware that this can reset your local database state.

### Run arbitrary Compose commands

For any standard Docker Compose operation not covered by the script shortcuts, use the `raw` sub-command:

```bash
./docker-control.sh raw logs -f
./docker-control.sh raw logs -f webserver
./docker-control.sh raw ps
./docker-control.sh raw exec db bash
```

---

## 3. Releases

To release a new version (creates and pushes a `vX.Y.Z` git tag, which triggers the GitHub Release workflow):

```bash
./release.sh
```
