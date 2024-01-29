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
  last_login timestamp not null,
  races int,
  avg_speed int,
  last_race_speed int,
  best_speed int,
  theme boolean not null
);

create table text(
     id uuid primary key,
     content text not null,
     author varchar(128) not null,
     occurrence int,
     accepted_at timestamp not null,
     length int not null,
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
    speed int not null,
    duration int not null,
    accuracy float not null,
    start_time timestamp not null,
    racer_id uuid references racer(id),
    text_id uuid references text(id)
);

create table race_history(
     single_id uuid references single(id),
     racer_id uuid references racer(id),
     text_id uuid references text(id),
     primary key (racer_id, single_id)
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
