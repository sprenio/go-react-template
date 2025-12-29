# Backend

This is a REST API service written in Go. It provides data and business logic for the frontend client.

## 🛠️ Stack

- Go 1.24+
- `net/http` or preferred web framework (e.g., Gin, Fiber)
- MySQL/MariaDB
- JWT authentication
- Docker

## 🔧 Environment variables

`.env` (used both locally and via Docker):

```
PORT=8080
DB_HOST=headless-db
DB_USER=myuser
DB_PASSWORD=mypass
DB_NAME=mydb
JWT_SECRET=your-secret-key
```

## 🚀 Local development

```bash
go run ./cmd/webserver
```

## 🧪 API structure

Basic structure:

- `POST /api/login` – authenticate and return JWT token
- `GET /api/some-data` – protected endpoint, requires JWT
- etc...

## ✅ TODO

- Swagger/OpenAPI docs
- Middleware for request logging/auth
- Unit tests
