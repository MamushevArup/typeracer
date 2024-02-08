-- +goose Up
-- +goose StatementBegin
create table racer (
   id uuid primary key,
   email varchar(256) unique not null,
   password varchar(512) not null,
   username varchar(128) not null,
   avatar varchar(512) not null default '',
   created_at timestamp not null,
   last_login timestamp not null,
   races int default 0,
   avg_speed int default 0,
   last_race_speed int default 0,
   best_speed int default 0,
   refresh_token varchar(2056) not null,
   role varchar(64) not null,
   theme boolean not null default false
);
create table session(
    id serial primary key,
    user_id uuid references racer(id),
    last_login timestamp not null,
    role varchar(64) not null,
    refresh_token varchar(2056) not null,
    fingerprint varchar(200) not null
);
create table text(
     id uuid primary key,
     content text not null,
     author varchar(128) not null,
     occurrence int default 0,
     accepted_at timestamp not null,
     length int not null,
     avg_speed int default 0,
     contributor_id uuid references racer(id),
     likes int not null default 0,
     dislikes int not null default 0,
     source varchar(64) not null,
     source_title varchar(256) not null
);
create table contributor(
    user_id uuid references racer(id),
    sent_at timestamp not null,
    text_id uuid references text(id),
    username varchar(128)
);
create table single (
    id uuid primary key ,
    speed int not null default 0,
    duration int not null default 0,
    accuracy float not null default 0,
    start_time timestamp not null,
    racer_id uuid references racer(id),
    text_id uuid references text(id)
);
create table race_history(
     single_id uuid references single(id),
     racer_id uuid references racer(id),
     text_id uuid references text(id),
     race_type int not null default 0,
     primary key (racer_id, single_id)
);
create table random_text(
    text_id uuid references text(id)
);
create table moderation(
   racer_id uuid not null references racer(id),
   content text not null,
   author varchar(128) not null,
   length int not null,
   source varchar(64) not null,
   source_title varchar(128) not null,
   sent_at timestamp not null
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table moderation;
drop table race_history;
drop table single;
drop table session;
drop table contributor;
drop table random_text;
drop table text;
drop table racer;
-- +goose StatementEnd