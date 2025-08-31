package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/merkulovlad/wbtech-go/internal/model"
)

var ErrNotFound = errors.New("order not found")

const (
	qSelOrder = `
SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
       delivery_service, shardkey, sm_id, date_created, oof_shard
FROM orders WHERE order_uid = $1`

	qSelDelivery = `
SELECT name, phone, zip, city, address, region, email
FROM deliveries WHERE order_uid = $1`

	qSelPayment = `
SELECT transaction, request_id, currency, provider, amount, payment_dt, bank,
       delivery_cost, goods_total, custom_fee
FROM payments WHERE order_uid = $1`

	qSelItems = `
SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
FROM items WHERE order_uid = $1 ORDER BY id`
)

// GetOrder loads an order + delivery + payment + items.
// Cache can wrap this at a higher layer; repo only talks to DB.
func (o *OrderRepository) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	// keep tight timeouts to avoid hanging requests
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var ord model.Order
	err := o.db.QueryRowContext(ctx, qSelOrder, id).Scan(
		&ord.OrderUID, &ord.TrackNumber, &ord.Entry, &ord.Locale, &ord.InternalSignature,
		&ord.CustomerID, &ord.DeliveryService, &ord.ShardKey, &ord.SmID, &ord.DateCreated, &ord.OofShard,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("select orders: %w", err)
	}

	if err := o.db.QueryRowContext(ctx, qSelDelivery, id).Scan(
		&ord.Delivery.Name, &ord.Delivery.Phone, &ord.Delivery.Zip, &ord.Delivery.City,
		&ord.Delivery.Address, &ord.Delivery.Region, &ord.Delivery.Email,
	); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("select deliveries: %w", err)
	}

	if err := o.db.QueryRowContext(ctx, qSelPayment, id).Scan(
		&ord.Payment.Transaction, &ord.Payment.RequestID, &ord.Payment.Currency, &ord.Payment.Provider,
		&ord.Payment.Amount, &ord.Payment.PaymentDT, &ord.Payment.Bank,
		&ord.Payment.DeliveryCost, &ord.Payment.GoodsTotal, &ord.Payment.CustomFee,
	); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("select payments: %w", err)
	}

	rows, err := o.db.QueryContext(ctx, qSelItems, id)
	if err != nil {
		return nil, fmt.Errorf("select items: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	items := make([]model.Item, 0, 8)
	for rows.Next() {
		var it model.Item
		if err := rows.Scan(
			&it.ChrtID, &it.TrackNumber, &it.Price, &it.RID, &it.Name,
			&it.Sale, &it.Size, &it.TotalPrice, &it.NmID, &it.Brand, &it.Status,
		); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("items rows: %w", err)
	}
	ord.Items = items

	return &ord, nil
}

func (o *OrderRepository) UpsertOrder(ctx context.Context, ord *model.Order) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // no-op if already committed

	// orders
	if _, err := tx.ExecContext(ctx, `
INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id,
                    delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT (order_uid) DO UPDATE SET
  track_number=EXCLUDED.track_number,
  entry=EXCLUDED.entry,
  locale=EXCLUDED.locale,
  internal_signature=EXCLUDED.internal_signature,
  customer_id=EXCLUDED.customer_id,
  delivery_service=EXCLUDED.delivery_service,
  shardkey=EXCLUDED.shardkey,
  sm_id=EXCLUDED.sm_id,
  date_created=EXCLUDED.date_created,
  oof_shard=EXCLUDED.oof_shard
`,
		ord.OrderUID, ord.TrackNumber, ord.Entry, ord.Locale, ord.InternalSignature,
		ord.CustomerID, ord.DeliveryService, ord.ShardKey, ord.SmID, ord.DateCreated, ord.OofShard,
	); err != nil {
		return fmt.Errorf("upsert orders: %w", err)
	}

	// deliveries
	if _, err := tx.ExecContext(ctx, `
INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
ON CONFLICT (order_uid) DO UPDATE SET
  name=EXCLUDED.name, phone=EXCLUDED.phone, zip=EXCLUDED.zip, city=EXCLUDED.city,
  address=EXCLUDED.address, region=EXCLUDED.region, email=EXCLUDED.email
`,
		ord.OrderUID, ord.Delivery.Name, ord.Delivery.Phone, ord.Delivery.Zip, ord.Delivery.City,
		ord.Delivery.Address, ord.Delivery.Region, ord.Delivery.Email,
	); err != nil {
		return fmt.Errorf("upsert deliveries: %w", err)
	}

	// payments
	if _, err := tx.ExecContext(ctx, `
INSERT INTO payments (order_uid, transaction, request_id, currency, provider,
                      amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT (order_uid) DO UPDATE SET
  transaction=EXCLUDED.transaction, request_id=EXCLUDED.request_id,
  currency=EXCLUDED.currency, provider=EXCLUDED.provider, amount=EXCLUDED.amount,
  payment_dt=EXCLUDED.payment_dt, bank=EXCLUDED.bank,
  delivery_cost=EXCLUDED.delivery_cost, goods_total=EXCLUDED.goods_total, custom_fee=EXCLUDED.custom_fee
`,
		ord.OrderUID, ord.Payment.Transaction, ord.Payment.RequestID, ord.Payment.Currency,
		ord.Payment.Provider, ord.Payment.Amount, ord.Payment.PaymentDT, ord.Payment.Bank,
		ord.Payment.DeliveryCost, ord.Payment.GoodsTotal, ord.Payment.CustomFee,
	); err != nil {
		return fmt.Errorf("upsert payments: %w", err)
	}

	// items → replace all current items for this order
	if _, err := tx.ExecContext(ctx, `DELETE FROM items WHERE order_uid = $1`, ord.OrderUID); err != nil {
		return fmt.Errorf("delete items: %w", err)
	}
	if len(ord.Items) > 0 {
		stmt, err := tx.PrepareContext(ctx, `
INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
`)
		if err != nil {
			return fmt.Errorf("prepare items: %w", err)
		}
		defer func(stmt *sql.Stmt) {
			err := stmt.Close()
			if err != nil {
				log.Printf("failed to close statement: %v", err)
			}
		}(stmt)

		for _, it := range ord.Items {
			if _, err := stmt.ExecContext(ctx,
				ord.OrderUID, it.ChrtID, it.TrackNumber, it.Price, it.RID, it.Name,
				it.Sale, it.Size, it.TotalPrice, it.NmID, it.Brand, it.Status,
			); err != nil {
				return fmt.Errorf("insert item: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}
func (o *OrderRepository) GetRecent(ctx context.Context, limit int) ([]*model.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rows, err := o.db.QueryContext(ctx, `
        SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
               delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders
        ORDER BY date_created DESC
        LIMIT $1
    `, limit)
	if err != nil {
		return nil, fmt.Errorf("select recent orders: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	var orders []*model.Order
	for rows.Next() {
		var ord model.Order
		if err := rows.Scan(
			&ord.OrderUID, &ord.TrackNumber, &ord.Entry, &ord.Locale, &ord.InternalSignature,
			&ord.CustomerID, &ord.DeliveryService, &ord.ShardKey, &ord.SmID, &ord.DateCreated, &ord.OofShard,
		); err != nil {
			return nil, fmt.Errorf("scan recent order: %w", err)
		}
		orders = append(orders, &ord)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}
