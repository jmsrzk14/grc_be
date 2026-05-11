-- +goose Up
ALTER TABLE regulation_assesments ADD COLUMN is_active BOOLEAN DEFAULT TRUE;

-- +goose Down
ALTER TABLE regulation_assesments DROP COLUMN is_active;
