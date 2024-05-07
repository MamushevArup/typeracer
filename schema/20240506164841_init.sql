-- +goose Up
-- +goose StatementBegin
INSERT into admin (id, username, refresh_token) values (1,'big-brother', '');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE admin
-- +goose StatementEnd
