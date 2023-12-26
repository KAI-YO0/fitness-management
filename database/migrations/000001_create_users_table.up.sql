CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  first_name VARCHAR (50) NOT NULL,
  last_name VARCHAR (50) NOT NULL,
  email VARCHAR (300) UNIQUE NOT NULL,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL
);
-- comments
COMMENT ON COLUMN users.id IS 'The user ID';
COMMENT ON COLUMN users.first_name IS 'The user first name';
COMMENT ON COLUMN users.last_name IS 'The user last name';
COMMENT ON COLUMN users.email IS 'The user email';
COMMENT ON COLUMN users.created_at IS 'Create time';
COMMENT ON COLUMN users.updated_at IS 'Update time';
COMMENT ON COLUMN users.deleted_at IS 'Delete time';
