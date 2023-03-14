CREATE TABLE `comments` (
    `id` CHAR(36) PRIMARY KEY,
    `commentable_type` VARCHAR(25),
    `commentable_id` CHAR(36),
    `user_id` CHAR(36),
    `body` TEXT NOT NULL,
    `resolved_at` TIMESTAMP NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NULL,

    CONSTRAINT uc_comments UNIQUE (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
)
ENGINE=InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;
