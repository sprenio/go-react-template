CREATE TABLE `confirmation_tokens`
(
    `id`          INT UNSIGNED                                       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `token`       char(32)                                           NOT NULL UNIQUE,
    `type`        enum ('register','email_change','password_reset')  NOT NULL,
    `payload`     longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`payload`)),
    `expires_at`  TIMESTAMP                                          NOT NULL,
    `consumed_at` TIMESTAMP                                          NULL     DEFAULT NULL,
    `created_at`  TIMESTAMP                                          NOT NULL DEFAULT current_timestamp()
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE users
(
    id            INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password      VARCHAR(255) NOT NULL,
    registered_at TIMESTAMP,
    confirmed_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE user_sessions
(
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id      INT UNSIGNED  NOT NULL,
    token_hash   CHAR(64)  NOT NULL UNIQUE,
    created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_date TIMESTAMP NOT NULL,
    revoked      BOOLEAN            DEFAULT FALSE,
    INDEX (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `user_sessions` ADD FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT ON UPDATE RESTRICT;