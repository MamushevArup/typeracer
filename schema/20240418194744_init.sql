-- +goose Up
-- +goose StatementBegin
CREATE TABLE admin (
    id serial primary key ,
    username varchar(32) not null,
    refresh_token varchar(128) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE admin;
-- +goose StatementEnd
