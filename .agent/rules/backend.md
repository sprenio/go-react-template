---
trigger: glob
description: General information about backend. Read it ALWAYS AND ONLY if the question or task is related to the backend
globs: backend/**
---

# Backend Development Standards (Go)

## 1. Project Location & Structure
- **Root Directory:** `backend`
- **Entrypoints (`cmd/`):**
  - `cmd/webserver/main.go`: REST API HTTP server.
  - `cmd/consumer/main.go`: RabbitMQ background consumer worker.
  - `cmd/cron/main.go`: Periodic scheduled jobs worker.

## 2. Layered Architecture (`internal/`)
Each domain module (e.g. `user`, `settings`, `session`, `confirmation`) follows a consistent layered architecture:
- **`http/`**: HTTP handlers, router registration (`routes.go`), request validation.
- **`dto/`**: Data transfer objects for API requests and responses.
- **`service/`**: Business logic, domain rules, transactions.
- **`repository/`**: Database queries (MariaDB) using raw SQL or prepared statements.

## 3. Shared Packages & Infrastructure
- **`config/`**: Configuration loading via Viper from YAML and environment variables.
- **`internal/contexthelper/`**: Type-safe helper functions to set/get database connections, redis client, user ID, client IP, and services from `context.Context`.
- **`internal/response/`**: Standardized HTTP response helpers (`SuccessResponse`, `SuccessDataResponse`, `TooManyRequestsErrorResponse`, `apiErrorResponse`).
- **`internal/middleware/`**: Chi HTTP middlewares (IP extraction, JWT authentication, rate limiting, logging, panic recovery).
- **`internal/ratelimit/`**: Sliding window rate limiting service via Redis ZSET/Lua script.
- **`pkg/logger/`**: Structured contextual logger (`logger.InfoCtx`, `logger.ErrorCtx`, `logger.WarnCtx`).
- **`pkg/apperror/`**: Application error wrapping with structured attributes.

## 4. Coding Conventions
- Always write idiomatic Go: handle errors explicitly, return early.
- Keep handlers thin: delegate business logic to services.
- Never hardcode credentials, URLs, or magic constants; use `config` and `apicodes`.
- Run tests with `go test ./...` before finalizing backend changes.
