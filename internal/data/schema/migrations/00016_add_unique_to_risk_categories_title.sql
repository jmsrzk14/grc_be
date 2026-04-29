-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_categories ADD CONSTRAINT risk_categories_title_key UNIQUE (title);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_categories DROP CONSTRAINT risk_categories_title_key;
-- +goose StatementEnd
