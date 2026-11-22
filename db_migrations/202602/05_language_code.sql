UPDATE `languages` SET `code` = `i18n_code`;
ALTER TABLE `languages` DROP `i18n_code`;
