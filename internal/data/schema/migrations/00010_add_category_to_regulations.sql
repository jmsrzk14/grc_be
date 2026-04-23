-- +goose Up
-- +goose StatementBegin
ALTER TABLE regulations ADD COLUMN category VARCHAR(50);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE regulations DROP COLUMN category;
-- +goose StatementEnd
