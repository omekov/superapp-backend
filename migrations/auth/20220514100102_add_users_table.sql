-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100),
    created_at timestamptz NULL DEFAULT now(),
	updated_at timestamptz NULL DEFAULT now()
);
-- +goose StatementEnd




-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
