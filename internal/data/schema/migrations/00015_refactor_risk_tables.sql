-- +goose Up
-- +goose StatementBegin
-- Menghapus unique constraint yang mencegah multiple risks dalam kategori yang sama untuk satu tenant
ALTER TABLE risk_category_tenants DROP CONSTRAINT IF EXISTS risk_category_tenants_risk_category_id_tenant_id_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_category_tenants ADD CONSTRAINT risk_category_tenants_risk_category_id_tenant_id_key UNIQUE (risk_category_id, tenant_id);
-- +goose StatementEnd
