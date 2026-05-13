-- +goose Up
-- +goose StatementBegin
-- Hapus kolom is_active dari tenant_regulations.
-- Kontrol akses per-tenant dilakukan via regulation_assesments.is_active (per sesi).
ALTER TABLE tenant_regulations DROP COLUMN IF EXISTS is_active;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tenant_regulations ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;
-- +goose StatementEnd
