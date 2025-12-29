CREATE TABLE `user_settings`
(
    `id`         INT UNSIGNED               NOT NULL AUTO_INCREMENT,
    `user_id`    INT UNSIGNED               NOT NULL UNIQUE,
    `lang_id`    TINYINT UNSIGNED           NOT NULL,
    `user_flags` BIGINT UNSIGNED            NOT NULL DEFAULT 0,
    `app_flags`  BIGINT UNSIGNED            NOT NULL DEFAULT 0,
    `app_opt_1`  VARCHAR(75)                NOT NULL DEFAULT '',
    `app_opt_2`  SET ("OPT_A", "OPT_B")     NOT NULL DEFAULT '',
    `app_opt_3`  ENUM ("RADIO_A","RADIO_B") NOT NULL DEFAULT 'RADIO_A',
    `updated_at` TIMESTAMP                  NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB;

ALTER TABLE `user_settings`
    ADD FOREIGN KEY (`lang_id`) REFERENCES `languages` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
    ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT;
