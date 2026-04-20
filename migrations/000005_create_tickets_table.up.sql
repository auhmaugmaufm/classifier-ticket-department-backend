CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS tickets (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        message TEXT NOT NULL,
        status TEXT NOT NULL,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        form_id UUID NOT NULL,
        department_id UUID,
        priority TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT fk_tickets_departments FOREIGN KEY (department_id) REFERENCES departments (id),
        CONSTRAINT fk_tickets_forms FOREIGN KEY (form_id) REFERENCES forms (id)
    );

CREATE INDEX IF NOT EXISTS idx_tickets_department_id ON tickets (department_id);
CREATE INDEX IF NOT EXISTS idx_tickets_form_id ON tickets (form_id);