-- +goose Up
-- +goose StatementBegin
create table racer (
  id uuid primary key,
  email varchar(256) not null,
  password varchar(512) not null,
  username varchar(128) not null,
  avatar varchar(512) not null,
  country varchar(128) not null,
  created_at timestamp not null,
  last_login timestamp,
  races int,
  avg_speed int,
  last_race_speed int,
  best_speed int,
  theme boolean
);

create table text(
     id uuid primary key,
     content text not null,
     author varchar(128),
     occurrence int,
     accepted_at timestamp,
     length int,
     avg_speed int,
     contributor_id uuid references racer(id)
);

create table contributor(
    user_id uuid references racer(id),
    sent_at timestamp not null,
    text_id uuid references text(id)
);

create table single (
    id uuid primary key ,
    speed int,
    duration int,
    accuracy float,
    start_time timestamp,
    racer_id uuid references racer(id),
    text_id uuid references text(id)
);

create table race_history(
     single_id uuid references single(id),
     racer_id uuid references racer(id),
     text_id uuid references text(id),
     primary key (racer_id, text_id)
);

create table random_text(
    text_id uuid references text(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table race_history;
drop table single;
drop table contributor;
drop table random_text;
drop table text;
drop table racer;
-- +goose StatementEnd
