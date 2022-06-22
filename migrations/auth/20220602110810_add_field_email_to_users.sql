-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN email varchar(100);
CREATE TYPE user_state AS ENUM ('disabled', 'enabled', 'notactivated');
ALTER TABLE users
ADD COLUMN "state" user_state;
ALTER TABLE users 
ADD CONSTRAINT unique_username
UNIQUE (username);
ALTER TABLE users 
ADD CONSTRAINT unique_email
UNIQUE (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users  
DROP COLUMN email;
ALTER TABLE users  
DROP COLUMN "state";
DROP TYPE user_state;
ALTER TABLE users 
DROP CONSTRAINT unique_username;
ALTER TABLE users 
DROP CONSTRAINT unique_email;
-- +goose StatementEnd
