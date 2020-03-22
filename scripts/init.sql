DROP TABLE IF EXISTS users CASCADE;

create UNLOGGED table users (
    id       bigserial primary key,
    name     TEXT not null,
    surname  TEXT not null,
    nickname TEXT not null unique,
    email    TEXT,
    avatar      TEXT,
    password TEXT not null
);
