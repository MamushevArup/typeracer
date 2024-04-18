-- +goose Up
-- +goose StatementBegin
create table moderation(
    moderation_id uuid primary key,
    racer_id uuid not null references racer(id),
    content text not null,
    author varchar(128) not null,
    length int not null,
    source varchar(64) not null,
    source_title varchar(128) not null,
    sent_at timestamp not null,
    status int default 0 not null
);

create table rejected (
  moderation_id uuid references moderation(moderation_id),
  response text not null default ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rejected;
DROP TABLE moderation;
-- +goose StatementEnd
