-- +goose Up
-- +goose StatementBegin
-- 1. Tambahkan kolom tenant_id ke risk_categories
ALTER TABLE risk_categories ADD COLUMN tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;

-- 2. Tambahkan kolom appetite dan tolerance langsung ke risk_categories untuk simplifikasi
ALTER TABLE risk_categories ADD COLUMN appetite TEXT;
ALTER TABLE risk_categories ADD COLUMN tolerance TEXT;

-- 3. Update existing data (jika ada) - set ke tenant pertama atau biarkan null
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_categories DROP COLUMN IF EXISTS tolerance;
ALTER TABLE risk_categories DROP COLUMN IF EXISTS appetite;
ALTER TABLE risk_categories DROP COLUMN IF EXISTS tenant_id;
-- +goose StatementEnd
