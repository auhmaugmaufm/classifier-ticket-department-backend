CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE
    IF NOT EXISTS forms (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        form_id UUID NOT NULL,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT fk_forms_company_forms FOREIGN KEY (form_id) REFERENCES companies (id)
    );

CREATE INDEX IF NOT EXISTS idx_forms_form_id ON forms (form_id);