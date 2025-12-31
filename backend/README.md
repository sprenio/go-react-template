## Backend

Go REST API that provides authentication, user management, and settings for the React frontend.

### Stack

- **Language**: Go (module located in `backend`)
- **Database**: MariaDB/MySQL
- **Messaging**: RabbitMQ (email + report queues)
- **Config**: YAML files with environment variable overrides (`config.yaml`, `config_dev.yaml`, `config.go`)
- **Processes**:
  - `webserver` – HTTP API
  - `consumer` – background queue consumer

### Configuration

The backend configuration is a combination of **embedded YAML** files and **environment variables** (`backend/config/config.go`).

- Base files:
  - `config.yaml`
  - `config_dev.yaml` (loaded when `APP_ENV=dev`)
- Important sections:
  - `web_server`: host and `http_port`
  - `database`: `user`, `password`, `host`, `port`, `dbname`
  - `rabbitmq`: `user`, `password`, `host`, `port`
  - `frontend`: URLs used in email links (`base_url`, `confirmation_endpoint`)
  - `register`, `reset_password`, `email_change`: feature flags and expiration settings
- Environment variables can override config; examples:
  - `APP_NAME`, `LOG_LEVEL`
  - `BACKEND_HOST`, `WEBSERVER_PORT`
  - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`
  - `RABBITMQ_HOST`, `RABBITMQ_PORT`, `RABBITMQ_USER`, `RABBITMQ_PASS`
  - `FRONTEND_BASE_URL`, `CONFIRMATION_ENDPOINT`
  - `REGISTER_ENABLED`, `REGISTER_CONFIRMATION_ENDPOINT`, `REGISTER_EXPIRATION_DAYS`
  - `RESET_PASSWORD_ENABLED`, `RESET_PASSWORD_EXPIRATION_DAYS`
  - `EMAIL_CHANGE_EXPIRATION_DAYS`
  - `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM`

When running under Docker, most of these values are provided by `docker-compose.dev.yml` / `docker-compose.prod.yml` and the root `.env` file.

### Local development (without Docker)

Ensure you have:

- A running MariaDB/MySQL instance.
- Optionally a running RabbitMQ instance (for queues).

Then:

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

# Hot reload webserver (uses air)
task watch
```

Make sure your environment variables match your local services (see “Configuration”).

### Running via Docker

The backend is typically run as part of the full stack using Docker Compose:

- **Dev**: `gr-webserver` and `gr-consumer` services in `docker/docker-compose.dev.yml`.
- **Prod**: `headless-webserver` and `headless-consumer` services in `docker/docker-compose.prod.yml`.

From the root `docker` directory:

```bash
cd docker
./docker-control.sh up
```

This will:

- Start MariaDB + migrations.
- Start RabbitMQ.
- Build and run the backend webserver and consumer with correct env vars.

### API notes

- All endpoints are exposed under `/api/` behind Nginx.
- JWT-based authentication and user flows are implemented (register, confirm, login, settings, password reset/change, email change).
- Handlers, services, and repositories live under `backend/internal`.

