CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE
    IF NOT EXISTS departments (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        department_name TEXT NOT NULL,
        company_id UUID NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        CONSTRAINT fk_departments_companies FOREIGN KEY (company_id) REFERENCES companies (id)
    );

CREATE INDEX IF NOT EXISTS idx_departments_company_id ON departments (company_id);