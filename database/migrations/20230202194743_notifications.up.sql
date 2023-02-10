CREATE TABLE `notifications` (
  `id` CHAR(36) PRIMARY KEY,  
  `object_type` VARCHAR(25) NOT NULL,
  `object_id` CHAR(36) NOT NULL,
  `title` VARCHAR(50) NOT NULL,
  `type` VARCHAR(50) NOT NULL,
  `body` TEXT NOT NULL,
  `archived_at` TIMESTAMP NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL
)
ENGINE=InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;

CREATE TABLE `notification_user` (
  `id` CHAR(36) PRIMARY KEY,  
  `notification_id` CHAR(36) NOT NULL,
  `notifier_id` CHAR(36) NOT NULL,
  `notified_id` CHAR(36) NOT NULL,

  FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (notifier_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (notified_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
)
ENGINE=InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;

