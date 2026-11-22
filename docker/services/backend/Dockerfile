
# Stage 1: build
FROM golang:1.27 AS build

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

COPY ./backend/ ./

ARG TARGET=webserver

# Allow only known application targets
RUN case "$TARGET" in \
      webserver|consumer|cron) ;; \
      *) echo "Invalid TARGET: $TARGET" && exit 1 ;; \
    esac && \
    CGO_ENABLED=0 go build -o /bin/app ./cmd/${TARGET}


# Stage 2: runtime
FROM debian:bookworm-slim AS runtime

WORKDIR /app

ARG GIT_COMMIT
ARG GIT_TAG

ENV GIT_COMMIT=$GIT_COMMIT
ENV GIT_TAG=$GIT_TAG

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    groupadd --system app && \
    useradd --system --gid app --home-dir /app --no-create-home app

COPY --from=build --chown=app:app /bin/app /app/app

COPY ./docker/services/backend/entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN chmod +x /usr/local/bin/docker-entrypoint.sh

USER app:app

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["/app/app"]
