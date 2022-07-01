-- +goose Up
-- +goose StatementBegin
INSERT INTO users 
(
    id,
	username, 
	password,
	email,
	state
) VALUES (
    '4085b498-caf9-4599-9b7c-9993818a50a4',
    'superadmin',
    '$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq',
    'superadmin@gmail.com',
    'enabled'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE * FROM users WHERE id = '4085b498-caf9-4599-9b7c-9993818a50a4';
-- +goose StatementEnd
