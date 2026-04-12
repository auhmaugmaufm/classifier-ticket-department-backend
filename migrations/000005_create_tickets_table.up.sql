CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE
    IF NOT EXISTS tickets (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        message TEXT NOT NULL,
        status TEXT NOT NULL,
        title TEXT NOT NULL,
        decription TEXT NOT NULL,
        department_id UUID NOT NULL,
        priority TEXT NOT NULL,
        submitted_date TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT fk_tickets_departments FOREIGN KEY (department_id) REFERENCES departments (id)
    );

CREATE INDEX IF NOT EXISTS idx_tickets_department_id ON tickets (department_id);