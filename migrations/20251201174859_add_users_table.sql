-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS auth_service.users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth_service.users;
-- +goose StatementEnd
