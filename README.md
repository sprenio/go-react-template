## go-react-template

Go + React + TypeScript + Tailwind CSS full‑stack boilerplate, ready for local development and containerized deployments.

### Tech stack

- **Backend**
  - Go (module in `backend`)
  - REST API with JWT-based authentication
  - MariaDB (via Docker)
  - RabbitMQ for async jobs (emails, reports)
  - YAML + environment-variable based configuration
- **Frontend**
  - React + TypeScript (Vite)
  - Tailwind CSS
  - React Router
  - React Query for data fetching/caching
  - i18n support with JSON locales
- **Infrastructure**
  - Docker Compose for dev/prod
  - Nginx reverse proxy
  - Database migrations (SQL files in `db_migrations`)

### Project structure

- **backend**: Go REST API, config, services, middleware, queue consumers.
- **frontend**: React + TS SPA with routes, hooks, providers, and Tailwind styling.
- **docker**: Dockerfiles, compose files, Nginx config, helper script.
- **db_migrations**: SQL migrations for MariaDB schema.
- **storage**: Docker volumes for database, RabbitMQ, and Nginx logs.

### Prerequisites

- **Backend / Docker**
  - Docker and Docker Compose
  - Make sure you can create an external Docker network (see Docker section).
- **Frontend**
  - Node.js (LTS) and npm

### Configuration

- Create a `.env` file in the project root; it is used by `docker/docker-control.sh`.
- At minimum, define:
  - **APP_ENV**: `dev` or `prod` (selects `docker-compose.dev.yml` or `docker-compose.prod.yml`).
  - **database credentials**: `DB_ROOT_PASSWORD`, `DB_DATABASE`, `DB_USER`, `DB_PASSWORD` (for dev compose) or corresponding variables for prod.
  - **RabbitMQ credentials** where applicable: `RABBITMQ_USER`, `RABBITMQ_PASSWORD`.
- Backend-specific configuration is in `backend/config/config.yaml` and `backend/config/config_dev.yaml` and can be overridden with environment variables (see `backend/config/config.go`).

### Running the stack with Docker (recommended)

1. **Create external Docker network** (once):

   ```bash
   docker network create reverse-proxy
   ```

2. **Ensure `.env` exists in repo root** and contains at least `APP_ENV` and DB/RabbitMQ settings.

3. **Start the stack** from the `docker` directory:

   ```bash
   cd docker
   ./docker-control.sh up
   ```

   The script:
   - Selects the compose file based on `APP_ENV` (for example `docker-compose.dev.yml`).
   - Brings up MariaDB, migrations, backend (`gr-webserver` / `headless-webserver`), consumer, RabbitMQ, and Nginx.

4. **Access the app**
   - Nginx listens on port 80 inside Docker; map that as appropriate in your host environment or reverse proxy.
   - All API routes are available under `/api/` (proxied to the Go backend).

### Local frontend development (with Docker backend)

1. Start Docker stack in **dev** mode (see previous section).
2. In a separate terminal, run the frontend:

   ```bash
   cd frontend
   npm install
   npm run dev -- --port 5174
   ```

3. Nginx (in `docker-compose.dev.yml`) is configured to proxy:
   - `/api/` to the backend container.
   - `/` to `host.docker.internal:5174`, so hitting Nginx will show the Vite dev server.

You can also access the Vite dev server directly on `http://localhost:5174` if you configure the frontend API base URL accordingly.

### Running backend without Docker (development only)

If you already have MariaDB and RabbitMQ running locally, you can run the backend directly:

```bash
cd backend
go run ./cmd/webserver
```

Or using `task` (requires `go-task` installed):

```bash
cd backend
task run          # TARGET defaults to webserver
TARGET=consumer task run
```

Ensure your local environment variables (`DB_HOST`, `DB_USER`, `DB_PASS`, `DB_NAME`, `RABBITMQ_HOST`, etc.) are aligned with the configuration in `backend/config/config.go`.

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

### Where to look next

- **Backend details**: see `backend/README.md` and the `internal` package (handlers, services, repositories, middleware).
- **Frontend details**: see `frontend/README.md` and the `src` directory for routes, hooks, and providers.
- **Docker details**: see `docker/README.md` and `docker/docker-control.sh` for controlling the stack in different environments.
