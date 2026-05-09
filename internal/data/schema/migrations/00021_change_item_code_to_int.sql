-- +goose Up
-- +goose StatementBegin
ALTER TABLE regulation_items ALTER COLUMN item_code TYPE integer USING (NULLIF(item_code, '')::integer);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE regulation_items ALTER COLUMN item_code TYPE varchar(100);
-- +goose StatementEnd
