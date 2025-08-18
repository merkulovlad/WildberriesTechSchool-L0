-- +goose Up
CREATE TABLE payments (
                          id             SERIAL PRIMARY KEY,
                          order_uid      VARCHAR NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                          transaction    VARCHAR NOT NULL,
                          request_id     VARCHAR NOT NULL,
                          currency       VARCHAR NOT NULL,
                          provider       VARCHAR NOT NULL,
                          amount         INT NOT NULL,
                          payment_dt     BIGINT NOT NULL,
                          bank           VARCHAR NOT NULL,
                          delivery_cost  INT NOT NULL,
                          goods_total    INT NOT NULL,
                          custom_fee     INT NOT NULL,
                          UNIQUE(order_uid)
);

-- +goose Down
DROP TABLE IF EXISTS payments;
