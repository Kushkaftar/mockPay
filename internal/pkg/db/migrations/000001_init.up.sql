CREATE TABLE merchant 
(
    id serial not null unique,
    title varchar(255) not null unique,
    api_key varchar(64) not null unique
);

CREATE TABLE card (
    id serial primary key unique not null,
    pan bigint not null,
    card_holder varchar (255), 
    exp_month int not null,
    exp_year int not null,
    cvc int not null,
    hash_card varchar (32) not null,
    created_at timestamp default current_timestamp
);

CREATE TABLE transaction (
    id serial primary key unique not null,
    merchant_id int REFERENCES merchant ( id ) ON DELETE CASCADE not null,
    card_id int REFERENCES card ( id ) ON DELETE CASCADE not null,
    transaction_type smallint not null,
    transaction_status smallint not null,
    uuid varchar (40) unique not null,
    amount real not null,
    created_at timestamp default current_timestamp
);

CREATE TABLE refund (
    id serial primary key unique not null,
    target_transaction_id int REFERENCES transaction ( id ) ON DELETE CASCADE not null,
    refund_transaction_id int REFERENCES transaction ( id ) ON DELETE CASCADE not null,
    created_at timestamp default current_timestamp
);

CREATE TABLE card_balance (
    id serial primary key unique not null,
    card_id int REFERENCES card ( id ) ON DELETE CASCADE not null,
    card_balance real not null
);

CREATE TABLE merchant_balance (
    id serial primary key unique not null,
    merchant_id int REFERENCES merchant ( id ) ON DELETE CASCADE not null,
    merchant_balance real not null
);

CREATE TABLE balance_event (
    id serial primary key unique not null,
    customer_type smallint not null,
    transaction_id int REFERENCES transaction ( id ) ON DELETE CASCADE not null,
    old_balance real not null,
    new_balance real not null
);

CREATE UNIQUE INDEX api_key
    ON merchant ( api_key );

CREATE UNIQUE INDEX hash_card
    ON card ( hash_card );

CREATE UNIQUE INDEX uuid
    ON transaction ( uuid );

CREATE UNIQUE INDEX card_id
    ON card_balance ( card_id );