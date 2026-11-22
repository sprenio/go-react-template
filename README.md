## go-react-template

Go + React + TypeScript + Tailwind CSS full‑stack boilerplate, ready for local development.

### Tech stack

- **Backend**
  - Go (module in `backend`)
  - REST API with JWT-based authentication
  - MariaDB (via Docker)
  - RabbitMQ for async jobs (emails, registration, password reset)
  - Redis for cache and rate limiting
  - Cron service for scheduled background jobs
  - YAML + environment-variable based configuration
- **Frontend**
  - React + TypeScript (Vite)
  - Tailwind CSS + CSS theme variables (light/dark/base)
  - Radix UI & Lucide icons
  - React Router
  - React Query for data fetching/caching
  - i18n support with JSON locales
- **Infrastructure**
  - Docker Compose for development environment
  - Production deployment managed via dedicated deployment repository (using images built via GitHub Actions Release workflow)
  - Nginx reverse proxy
  - Database migrations (SQL files in `db_migrations`)

### Project structure

- **backend**: Go REST API, config, services, middleware, queue consumers, cron worker.
- **frontend**: React + TS SPA with routes, hooks, providers, and Tailwind styling.
- **docker**: Docker Compose, development Nginx config, helper CLI script (`docker-control.sh`), release script (`release.sh`).
- **db_migrations**: SQL migrations for MariaDB schema.
- **storage**: Docker volumes for database, RabbitMQ, Redis, and Nginx logs.

### Prerequisites

- **Backend / Docker**
  - Docker and Docker Compose
  - Make sure you can create an external Docker network (see Docker section).
- **Frontend**
  - Node.js (LTS) and npm

### Configuration

- Create a `.env` file in the project root (you can copy `.env.sample`); it is used by `docker/docker-control.sh`.
- At minimum, define:
  - **APP_ENV**: `dev`
  - **database credentials**: `DB_ROOT_PASSWORD`, `DB_DATABASE`, `DB_USER`, `DB_PASSWORD`.
  - **RabbitMQ credentials**: `RABBITMQ_USER`, `RABBITMQ_PASSWORD`.
  - **Redis settings**: `REDIS_HOST`, `REDIS_PORT`.
  - **Vite dev port**: `VITE_PORT` (defaults to `5173`).
- Nginx configuration: copy `docker/services/nginx/default_server.conf.sample` to `docker/services/nginx/default_server.conf` (if missing, `docker-control.sh` will auto-create it for you on startup). This file is ignored by git so you can customize ports and local routing without dirtying your git tree.
- Backend-specific configuration is in `backend/config/config.yaml` and `backend/config/config_dev.yaml` and can be overridden with environment variables (see `backend/config/config.go`).

### Running the stack with Docker (recommended)

1. **Create external Docker networks** (once):

   ```bash
   docker network create reverse-proxy
   docker network create db-admin
   ```

2. **Ensure `.env` exists in repo root** and contains required dev settings.

3. **Start the stack** from the `docker` directory:

   ```bash
   cd docker
   ./docker-control.sh up
   ```

   The script:
   - Brings up MariaDB, runs migrations, and starts the backend (`webserver`), background workers (`consumer`, `cron`), RabbitMQ, Redis, and Nginx.

4. **Access the app**
   - Nginx listens on port 80 inside Docker; map that as appropriate in your host environment or reverse proxy.
   - All API routes are available under `/api/` (proxied to the Go backend).

### Local frontend development (with Docker backend)

1. Start Docker stack in **dev** mode (see previous section).
2. In a separate terminal, run the frontend:

   ```bash
   cd frontend
   npm install
   npm run dev
   ```

   The port is configured via `VITE_PORT` in your root `.env` (defaults to `5173` in `.env.sample`).

3. Nginx (`docker/services/nginx/default_server.conf`) is configured to proxy:
   - `/api/` to the backend container (`http://webserver:8080/`).
   - `/` to the Vite dev server (`http://host.docker.internal:5173` matching `VITE_PORT`), so hitting Nginx will show your frontend with hot reloading.

> [!NOTE]
> If you change `VITE_PORT` in `.env`, remember to update the corresponding `proxy_pass` port in `docker/services/nginx/default_server.conf`.

### Running backend without Docker (development only)

If you already have MariaDB, Redis, and RabbitMQ running locally, you can run the backend directly:

```bash
cd backend
go run ./cmd/webserver
```

Or using `task` (requires `go-task` installed):

```bash
cd backend
task run                  # TARGET defaults to webserver
TARGET=consumer task run
TARGET=cron task run
```

Ensure your local environment variables (`DB_HOST`, `DB_USER`, `DB_PASS`, `DB_NAME`, `RABBITMQ_HOST`, `REDIS_HOST`, etc.) are aligned with the configuration in `backend/config/config.go`.

### Tests and quality

- **Backend**
  - Run tests:

    ```bash
    cd backend
    task test
    ```

  - Lint/format:

    ```bash
    cd backend
    task lint
    ```

- **Frontend**
  - Run unit tests:

    ```bash
    cd frontend
    npm test
    ```

  - Lint:

    ```bash
    cd frontend
    npm run lint
    ```

  - Format:

    ```bash
    cd frontend
    npm run format
    ```

### Releases

Releases are triggered by pushing git tags (`vX.Y.Z`). You can use the helper script:

```bash
cd docker
./release.sh
```

The GitHub Actions workflow (`.github/workflows/release.yml`) builds and publishes the Docker images (`webserver`, `consumer`, `cron`, `frontend`, `db-migrator`) to GitHub Container Registry (GHCR).
