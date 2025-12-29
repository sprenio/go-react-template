CREATE TABLE `languages`
(
    `id`   TINYINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `code` CHAR(2)          NOT NULL UNIQUE,
    `name` VARCHAR(32)      NOT NULL
) ENGINE = InnoDB;

INSERT INTO `languages` (`code`, `name`)
VALUES ('pl', 'Polski'),
       ('en', 'English');
