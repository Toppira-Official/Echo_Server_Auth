-- +goose Up
CREATE TABLE IF NOT EXISTS credentials (
  id CHAR(27) PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  hashed_password VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_credentials_username
  ON credentials (username);

-- +goose Down
DROP INDEX IF EXISTS idx_credentials_username;
DROP TABLE IF EXISTS credentials;
