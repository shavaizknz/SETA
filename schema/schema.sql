CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE transactions (
    id uuid default uuid_generate_v4() primary key,
    transaction_id varchar(255) not null,
    account_id varchar(255) not null,
    amount numeric(10, 2) not null,
    status varchar(255) not null,
    type varchar(255) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    unique (transaction_id, account_id)
);
