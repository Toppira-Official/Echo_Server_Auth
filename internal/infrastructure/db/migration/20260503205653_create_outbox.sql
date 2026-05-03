-- +goose Up
CREATE TABLE IF NOT EXISTS outbox (
  id TEXT PRIMARY KEY,
  event_type TEXT NOT NULL,
  payload BYTEA NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  retry_count INT NOT NULL DEFAULT 0,
  max_retries INT NOT NULL DEFAULT 5,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  published_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_outbox_pending ON outbox (status, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_outbox_pending;

DROP TABLE IF EXISTS outbox;