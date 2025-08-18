-- +goose Up
CREATE TABLE items (
                       id           SERIAL PRIMARY KEY,
                       order_uid    VARCHAR NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                       chrt_id      INT NOT NULL,
                       track_number VARCHAR NOT NULL,
                       price        INT NOT NULL,
                       rid          VARCHAR NOT NULL,
                       name         VARCHAR NOT NULL,
                       sale         INT NOT NULL,
                       size         VARCHAR NOT NULL,
                       total_price  INT NOT NULL,
                       nm_id        INT NOT NULL,
                       brand        VARCHAR NOT NULL,
                       status       INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS items;
