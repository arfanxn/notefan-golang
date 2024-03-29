CREATE TABLE pages (
  id CHAR(36) PRIMARY KEY,  
  space_id CHAR(36) NOT NULL,
  title VARCHAR(50) NOT NULL,
  `order` INT UNSIGNED NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL,

  CONSTRAINT uc_pages UNIQUE (id),
  FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE ON UPDATE CASCADE
) 
ENGINE=InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;
