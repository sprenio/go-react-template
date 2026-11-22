ALTER TABLE `users` ADD `normalized_email` VARCHAR(100) NOT NULL AFTER `email`, ADD UNIQUE (`normalized_email`);
