-- +goose Up
-- +goose StatementBegin
CREATE TABLE regulations (
    id               UUID PRIMARY KEY,
    title            VARCHAR(255) NOT NULL,
    regulation_type  VARCHAR(50)  NOT NULL,
    issued_date      DATE,
    status           VARCHAR(50)  NOT NULL DEFAULT 'Active',
    category         VARCHAR(50)  NOT NULL DEFAULT 'External'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS regulations;
-- +goose StatementEnd
