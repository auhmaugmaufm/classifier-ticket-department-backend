CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE tickets
    RENAME COLUMN status TO predict_status;

ALTER TABLE tickets
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'pending';