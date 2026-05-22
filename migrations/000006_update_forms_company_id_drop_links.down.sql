CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    link TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_links_companies
        FOREIGN KEY (company_id)
        REFERENCES companies (id)
);

CREATE INDEX IF NOT EXISTS idx_links_company_id
ON links (company_id);

INSERT INTO links (company_id, link, created_at, updated_at, deleted_at)
SELECT DISTINCT
    forms.company_id,
    '/form/' || forms.company_id::text,
    NOW(),
    NOW(),
    NULL
FROM forms
WHERE forms.company_id IS NOT NULL;

ALTER TABLE forms
    ADD COLUMN IF NOT EXISTS link_id UUID;

UPDATE forms
SET link_id = links.id
FROM links
WHERE links.company_id = forms.company_id;

CREATE INDEX IF NOT EXISTS idx_forms_link_id
ON forms (link_id);

ALTER TABLE forms
    DROP CONSTRAINT IF EXISTS fk_forms_companies;

ALTER TABLE forms
    ADD CONSTRAINT fk_forms_links
        FOREIGN KEY (link_id)
        REFERENCES links (id);

ALTER TABLE forms
    DROP COLUMN IF EXISTS company_id;

DROP INDEX IF EXISTS idx_forms_company_id;