-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS auth_service;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS auth_service;
-- +goose StatementEnd
