# Stage 1: build
FROM golang:1.24 AS build

WORKDIR /app

# Kopiujemy pliki mod i pobieramy zależności
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

# Kopiujemy cały kod źródłowy
COPY ./backend/ ./

# Zmienna określająca, który "main" kompilujemy
ARG TARGET=webserver

# Budujemy binarkę w zależności od targetu
RUN go build -o /bin/app ./cmd/${TARGET}

# Stage 2: runtime
FROM debian:bookworm-slim AS runtime

WORKDIR /app

# Skopiuj binarkę z etapu build
COPY --from=build /bin/app /app/app

# (opcjonalnie) jeśli masz entrypoint
COPY ./docker/services/backend/entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["/app/app"]
