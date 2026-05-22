CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE forms
    ADD COLUMN IF NOT EXISTS company_id UUID;

UPDATE forms
SET company_id = links.company_id
FROM links
WHERE forms.link_id = links.id;

ALTER TABLE forms
    ALTER COLUMN company_id SET NOT NULL;

ALTER TABLE forms
    DROP CONSTRAINT IF EXISTS fk_forms_links;

ALTER TABLE forms
    ADD CONSTRAINT fk_forms_companies
        FOREIGN KEY (company_id)
        REFERENCES companies (id);

CREATE INDEX IF NOT EXISTS idx_forms_company_id
ON forms (company_id);

ALTER TABLE forms
    DROP COLUMN IF EXISTS link_id;

DROP TABLE IF EXISTS links;