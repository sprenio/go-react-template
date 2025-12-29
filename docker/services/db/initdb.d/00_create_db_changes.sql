SELECT NOW() INTO @startDate;
CREATE TABLE IF NOT EXISTS `db_changes` (
    `id` int(4) UNSIGNED NOT NULL AUTO_INCREMENT,
    `date_dir` mediumint(3) UNSIGNED NOT NULL,
    `file_name` varchar(64) NOT NULL,
    `start_date` datetime NOT NULL,
    `complete_date` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO `db_changes` (`date_dir`, `file_name`, `start_date`, `complete_date`) VALUES (0, '000_init.sql', @startDate, NOW());
