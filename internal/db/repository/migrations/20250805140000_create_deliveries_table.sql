-- +goose Up
CREATE TABLE deliveries (
                            id         SERIAL PRIMARY KEY,
                            order_uid  VARCHAR NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                            name       VARCHAR NOT NULL,
                            phone      VARCHAR NOT NULL,
                            address    VARCHAR NOT NULL,
                            city       VARCHAR NOT NULL,
                            region     VARCHAR NOT NULL,
                            zip        VARCHAR NOT NULL,
                            email      VARCHAR NOT NULL,
                            UNIQUE(order_uid)
);

-- +goose Down
DROP TABLE IF EXISTS deliveries;
