ALTER TABLE `languages` ADD `i18n_code` CHAR(2) NULL AFTER `code`;
UPDATE `languages` SET `i18n_code` = code;
INSERT INTO `languages` (`code`, `i18n_code`, `name`)
VALUES ('de', 'de', 'Deutsch'),
       ('ua', 'uk', 'Українська');
