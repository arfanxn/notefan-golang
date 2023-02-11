CREATE TABLE `notification_user` (
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

