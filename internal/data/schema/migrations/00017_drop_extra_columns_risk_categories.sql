-- +goose Up
-- +goose StatementBegin
-- Menghapus kolom yang tidak sengaja ditambahkan ke risk_categories
ALTER TABLE risk_categories DROP COLUMN IF EXISTS appetite;
ALTER TABLE risk_categories DROP COLUMN IF EXISTS tolerance;
ALTER TABLE risk_categories DROP COLUMN IF EXISTS tenant_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- (Tidak perlu menambahkan kembali karena ini adalah proses pembersihan)
-- +goose StatementEnd
