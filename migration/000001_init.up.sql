CREATE TABLE Support (
   `id` VARCHAR(64) NOT NULL DEFAULT (UUID()),
   `name` VARCHAR(255) NOT NULL,
   `email` VARCHAR(255) NOT NULL,
   `message` TEXT NOT NULL,
   PRIMARY KEY (id),
   UNIQUE `support_email_idx`(`email` ASC)
);