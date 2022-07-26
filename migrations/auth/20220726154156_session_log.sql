-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_session_logs (
    id serial NOT NULL,
    session_id UUID NOT NULL,
    username VARCHAR(100) NOT NULL DEFAULT '',
    user_agent VARCHAR(255) NOT NULL,
    client_ip VARCHAR(255) NOT NULL,
    http_method VARCHAR(10) NOT NULL,
    http_path VARCHAR NOT NULL,
    http_status int4 NULL default 0,
    http_req_body jsonb not null default '{}'::jsonb,
    http_res_body jsonb not null default '{}'::jsonb,
    exprires_at timestamptz NULL DEFAULT now(),
	CONSTRAINT user_session_log_pk PRIMARY KEY (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_session_logs;
-- +goose StatementEnd
