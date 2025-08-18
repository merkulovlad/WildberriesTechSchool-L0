-- +goose Up
CREATE TABLE orders (
    order_uid          VARCHAR PRIMARY KEY,
    track_number       VARCHAR NOT NULL,
    entry              VARCHAR NOT NULL,
    locale             VARCHAR NOT NULL,
    date_created       TIMESTAMP NOT NULL,
    internal_signature VARCHAR,
    customer_id        VARCHAR NOT NULL,
    delivery_service   VARCHAR NOT NULL,
    shardkey           VARCHAR NOT NULL,
    sm_id              INTEGER NOT NULL,
    oof_shard          VARCHAR NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS orders;
