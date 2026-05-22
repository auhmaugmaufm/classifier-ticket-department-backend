CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE tickets
    DROP COLUMN IF EXISTS status;

ALTER TABLE tickets
    RENAME COLUMN predict_status TO status;