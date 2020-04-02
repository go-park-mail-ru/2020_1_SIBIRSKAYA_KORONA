DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS boards CASCADE;
DROP TABLE IF EXISTS board_admins CASCADE;
DROP TABLE IF EXISTS board_members CASCADE;

create UNLOGGED table users (
    id       bigserial primary key,
    name     TEXT not null,
    surname  TEXT not null,
    nickname TEXT not null unique,
    email    TEXT,
    avatar      TEXT,
    password TEXT not null
);
