ALTER TABLE `confirmation_tokens`
    CHANGE `type` `type` ENUM ('register','email_change','reset_password') NOT NULL;
