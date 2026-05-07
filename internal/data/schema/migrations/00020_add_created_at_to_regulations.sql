-- +goose Up
-- +goose StatementBegin
ALTER TABLE regulations ADD COLUMN created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE regulations DROP COLUMN IF EXISTS created_at;
-- +goose StatementEnd
