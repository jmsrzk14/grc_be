-- +goose Up
-- +goose StatementBegin
-- Tambah kolom is_active ke tabel tenant_regulations.
-- Ketika regulasi dicabut dari tenant, is_active di-set false (soft-deactivation).
-- Data historis tetap terjaga untuk keperluan audit.
ALTER TABLE tenant_regulations
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tenant_regulations DROP COLUMN IF EXISTS is_active;
-- +goose StatementEnd
