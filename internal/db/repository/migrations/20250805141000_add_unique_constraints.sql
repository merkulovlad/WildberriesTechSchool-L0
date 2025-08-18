-- +goose Up
-- Add unique constraints for ON CONFLICT support
ALTER TABLE deliveries ADD CONSTRAINT deliveries_order_uid_unique UNIQUE (order_uid);
ALTER TABLE payments ADD CONSTRAINT payments_order_uid_unique UNIQUE (order_uid);

-- +goose Down
-- Remove unique constraints
ALTER TABLE deliveries DROP CONSTRAINT IF EXISTS deliveries_order_uid_unique;
ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_order_uid_unique;
