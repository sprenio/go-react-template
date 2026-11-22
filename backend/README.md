## Backend

Go REST API that provides authentication, user management, and settings for the React frontend.

### Stack

- **Language**: Go (module located in `backend`)
- **Database**: MariaDB/MySQL
- **Cache & Rate Limiting**: Redis
- **Messaging**: RabbitMQ (emails, registration, password resets)
- **Background Cron**: Scheduled worker (`cmd/cron`)
- **Config**: YAML files with environment variable overrides (`config.yaml`, `config.go`)
- **Processes**:
  - `webserver` – HTTP API
  - `consumer` – background queue consumer
  - `cron` – scheduled background tasks

### Configuration

The backend configuration is a combination of **embedded YAML** files and **environment variables** (`backend/config/config.go`).

- Base file: `config.yaml` (overridable with `config_dev.yaml` when `APP_ENV=dev`)
- Important sections:
  - `web_server`: host and `http_port`
  - `database`: `user`, `password`, `host`, `port`, `dbname`
  - `rabbitmq`: `user`, `password`, `host`, `port`
  - `redis`: `host`, `port`, `db`
  - `rate_limit`: endpoint and global limit rules
  - `frontend`: URLs used in email links (`base_url`, `confirmation_endpoint`)
  - `register`, `reset_password`, `email_change`: feature flags and expiration settings
- Environment variables can override config; examples:
  - `APP_NAME`, `LOG_LEVEL`
  - `BACKEND_HOST`, `WEBSERVER_PORT`
  - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`
  - `RABBITMQ_HOST`, `RABBITMQ_PORT`, `RABBITMQ_USER`, `RABBITMQ_PASS`
  - `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`
  - `FRONTEND_BASE_URL`, `CONFIRMATION_ENDPOINT`
  - `REGISTER_ENABLED`, `RESET_PASSWORD_ENABLED`
  - `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM`

When running under Docker, values are provided by `docker-compose.yml` and the root `.env` file.

### Local development (without Docker)

Ensure you have running instances of MariaDB/MySQL, Redis, and RabbitMQ.

```bash
cd backend
go run ./cmd/webserver
```

or using the `Taskfile`:

```bash
cd backend
# Run webserver (default TARGET)
task run

# Run consumer
TARGET=consumer task run

# Run cron worker
TARGET=cron task run

# Hot reload (uses air)
task watch
```

Make sure your environment variables match your local services (see “Configuration”).

### Running via Docker

The backend is run as part of the development stack using Docker Compose:

From the root `docker` directory:

```bash
cd docker
./docker-control.sh up
```

This will:
- Start MariaDB + run migrations via `db-migrator`.
- Start RabbitMQ and Redis.
- Build and run the backend services (`webserver`, `consumer`, `cron`) with hot reload.
- Start Nginx proxying `/api/` to `webserver` and `/` to Vite dev server.

### API notes

- All endpoints are exposed under `/api/` behind Nginx.
- JWT-based authentication and user flows are implemented (register, confirm, login, settings, password reset/change, email change).
- Modular architecture lives under `backend/internal` (e.g. `user`, `session`, `settings`, `confirmation`, `ratelimit`, `queue`, `cron`).
