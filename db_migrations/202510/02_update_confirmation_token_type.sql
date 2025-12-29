ALTER TABLE `confirmation_tokens`
    CHANGE `type` `type` ENUM ('register','email_change','password_change') NOT NULL;
