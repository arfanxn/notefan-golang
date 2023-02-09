CREATE TABLE pages (
  id VARCHAR(36) PRIMARY KEY,  
  title VARCHAR(50) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL
);