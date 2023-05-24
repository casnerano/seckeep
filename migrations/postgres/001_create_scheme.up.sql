create extension if not exists "uuid-ossp";

create table if not exists users (
    uuid uuid primary key default uuid_generate_v4() not null,
    login character varying(100) not null,
    password character varying(100) not null,
    full_name character varying(100) not null,
    created_at timestamp default now() not null,
    constraint users_unique_login unique (login)
);

create type data_type as enum ('CREDENTIAL','TEXT', 'CARD', 'DOCUMENT');

create table if not exists data (
    uuid uuid primary key default uuid_generate_v4() not null,
    user_uuid uuid not null,
    type data_type not null,
    value bytea not null,
    created_at timestamp default now() not null,
    version timestamp default now() not null,
    constraint data_fk_user foreign key (user_uuid) references users (uuid)
)
