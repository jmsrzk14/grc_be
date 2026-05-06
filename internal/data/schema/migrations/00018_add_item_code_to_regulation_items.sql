-- +goose Up
-- +goose StatementBegin
ALTER TABLE regulation_items ADD COLUMN item_code VARCHAR(100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE regulation_items DROP COLUMN IF EXISTS item_code;
-- +goose StatementEnd
