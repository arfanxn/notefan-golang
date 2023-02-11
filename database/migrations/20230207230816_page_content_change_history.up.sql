CREATE TABLE `page_content_change_history` (
  `before_page_content_id` CHAR(36),
  `after_page_content_id` CHAR(36),
  `user_id` CHAR(36),
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL,

  FOREIGN KEY (before_page_content_id) REFERENCES page_contents(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (after_page_content_id) REFERENCES page_contents(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
)
ENGINE=InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;
