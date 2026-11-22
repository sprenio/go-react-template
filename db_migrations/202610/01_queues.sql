ALTER TABLE `confirmation_tokens`
    CHANGE `status` `status` ENUM('NEW','WAITING','PROCESSING','CONSUMED','FAILED','EXPIRED','CANCELED') NOT NULL DEFAULT 'NEW';

CREATE TABLE `queues` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `task` ENUM('register_user','reset_password','email') NOT NULL,
    `payload` JSON NOT NULL,
    `status` ENUM('PENDING', 'QUEUED', 'PROCESSING', 'DONE', 'FAILED', 'RETRY_PENDING') NOT NULL DEFAULT 'PENDING',
    `retries_count` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `max_retries` TINYINT UNSIGNED NOT NULL DEFAULT 3,
    `generation` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `last_error` TEXT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `next_retry_at` DATETIME NULL,
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_task_status` (`task`, `status`)
) ENGINE = InnoDB  DEFAULT CHARSET=utf8mb4;

CREATE TABLE `email_log` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `queue_id` BIGINT UNSIGNED NOT NULL,
    `generation` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `retry_number` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `status` ENUM('STARTED', 'SENT', 'FAILED') NOT NULL DEFAULT 'STARTED',
    `log` text NULL DEFAULT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_queue_id_generation` (`queue_id`, `generation`, `retry_number`),
    KEY `idx_status` (`status`)
) ENGINE = InnoDB  DEFAULT CHARSET=utf8mb4;
