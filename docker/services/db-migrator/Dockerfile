
FROM php:8.5-cli-alpine3.24

RUN apk add --no-cache bash \
    && docker-php-ext-install mysqli pdo pdo_mysql \
    && docker-php-ext-enable pdo_mysql

WORKDIR /migrator

COPY ./docker/services/db-migrator/run-migrations.php ./run-migrations.php
COPY ./db_migrations ./migrations

ENTRYPOINT ["php", "run-migrations.php"]
