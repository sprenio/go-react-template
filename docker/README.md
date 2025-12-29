# Docker – MyApp

This directory contains the Docker configuration for running the entire stack locally.

## 📦 Structure

```
src/
├── db/                   # database schema and migrations
├── env/                  # shared environment variables used in Docker
├── docker-compose.yml    # defines all services (frontend, backend, database)
├── docker-control.bat    # script for windows to control Docker
├── docker-control.sh     # script for linux to control Docker
```
- `Dockerfile` – exists in `frontend/` and `backend/` for building app containers

## 🐳 Run everything

Use the appropriate script for your operating system:

- **Windows**
  ```bash
  .\docker-compose.bat up --build
  ```
- **Linux/macOS**
  ```bash
  ./docker-compose.sh up --build
  ```

## 🔁 Rebuild a single service

Use the appropriate script for your operating system:

- **Windows**
  ```bash
  .\docker-compose.bat up --build frontend
  ```
- **Linux/macOS**
  ```bash
  ./docker-compose.sh up --build frontend
  ```


## 📦 Used images

- Node.js + Vite for frontend
- Go for backend
- `mariadb:11.8.2-ubi` as the database

## ⚙️ Tips

- You can override environment values in `./env/.env*`
- Volumes are used to persist DB data
- Ports:
  - `3000` – frontend
  - `8080` – backend API
  - `3306` – MySQL/MariaDB
