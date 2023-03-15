CREATE TABLE tokens (
  `id` CHAR(36) PRIMARY KEY,  
  `tokenable_type` VARCHAR(25) NOT NULL,
  `tokenable_id` CHAR(36) NOT NULL,
  `type` VARCHAR(25) NOT NULL,
  `body` VARCHAR(255) NOT NULL,
  `used_at` TIMESTAMP NULL,
  `expired_at` TIMESTAMP NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL,
  
  CONSTRAINT uc_tokens UNIQUE (id, body)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;
