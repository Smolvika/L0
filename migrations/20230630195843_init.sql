-- +goose Up
CREATE TABLE IF NOT EXISTS orders(
    id SERIAL UNIQUE,
    order_uid  VARCHAR(30) NOT NULL,
    track_number  VARCHAR(30) NOT NULL,
    entry  VARCHAR(30) NOT NULL,
    locale  VARCHAR(2) NOT NULL,
    internal_signature  VARCHAR(50) NOT NULL,
    customer_id  VARCHAR(30) NOT NULL,
    delivery_service  VARCHAR(30) NOT NULL,
    shardkey  VARCHAR(30) NOT NULL,
    sm_id SERIAL NOT NULL,
    date_created VARCHAR(20) NOT NULL,
    oof_shard  VARCHAR(30) NOT NULL
    );

CREATE TABLE IF NOT EXISTS delivery(
    order_id SERIAL REFERENCES orders(id) ,
    name VARCHAR(30) NOT NULL,
    phone VARCHAR(11) NOT NULL,
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(30) NOT NULL,
    address VARCHAR(30) NOT NULL,
    region VARCHAR(30) NOT NULL,
    email VARCHAR(30) NOT NULL
    );

CREATE TABLE IF NOT EXISTS payment(
    order_id SERIAL REFERENCES orders(id),
    transaction VARCHAR(30) NOT NULL,
    request_id VARCHAR(30) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR(30) NOT NULL,
    amount NUMERIC NOT NULL,
    payment_dt BIGSERIAL NOT NULL,
    bank VARCHAR(30) NOT NULL,
    delivery_cost NUMERIC NOT NULL,
    goods_total SERIAL NOT NULL,
    custom_fee NUMERIC NOT NULL
    );

CREATE TABLE IF NOT EXISTS items(
    order_id SERIAL REFERENCES orders(id),
    chrt_id SERIAL NOT NULL,
    track_number VARCHAR(30) NOT NULL,
    price NUMERIC NOT NULL,
    rid VARCHAR(30) NOT NULL,
    "name" VARCHAR(30) NOT NULL,
    sale SERIAL CHECK(sale >=0 AND sale <= 100) NOT NULL,
    size VARCHAR(10) NOT NULL,
    total_price NUMERIC NOT NULL,
    nm_id BIGSERIAL NOT NULL,
    brand VARCHAR(50) NOT NULL,
    status INTEGER NOT NULL
    );


-- +goose Down
DROP TABLE delivery;
DROP TABLE payment;
DROP TABLE items;
DROP TABLE orders;

