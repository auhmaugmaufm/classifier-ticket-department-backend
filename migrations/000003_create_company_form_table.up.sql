CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE
    IF NOT EXISTS company_forms (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        company_id UUID NOT NULL,
        link_form TEXT NOT NULL
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT fk_company_forms_companies FOREIGN KEY (company_id) REFERENCES companies (id)
    );

CREATE INDEX IF NOT EXISTS idx_company_forms_company_id ON company_forms (company_id);