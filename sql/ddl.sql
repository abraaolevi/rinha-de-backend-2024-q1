create table accounts (
    id serial primary key,
    account_limit bigint,
    balance bigint
);

create table transactions (
    id bigserial primary key,
    account_id serial references accounts (id),
    amount bigserial,
    operation char(1),
    description varchar(10),
    created_at timestamp default current_timestamp 
);