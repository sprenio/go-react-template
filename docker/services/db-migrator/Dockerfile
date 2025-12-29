FROM php:8.3-rc-cli-alpine3.22
RUN apk add --no-cache bash
RUN docker-php-ext-install mysqli pdo pdo_mysql && docker-php-ext-enable pdo_mysql

WORKDIR /migrator

COPY ./docker/services/db-migrator/run-migrations.php .

ENTRYPOINT ["php", "run-migrations.php"]
