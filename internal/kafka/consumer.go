// Package kafka contains Kafka consumer logic for the Orders domain.
// It implements a robust, production-oriented message processing loop with:
//   - Structured validation of inbound messages (basic schema/business sanity checks)
//   - Dead-Letter Queue (DLQ) publishing for poison/unrecoverable messages
//   - Clear logging for observability and later diagnostics
//
// Design notes (for contributors):
//   - We prioritize "process what you can, isolate what you can't": invalid or permanently
//     failing messages are routed to the DLQ instead of blocking the partition.
//   - Validation here is intentionally lightweight. Deeper domain validation happens in
//     the service layer and at the DB via constraints. Keep this layer focused on
//     deserialization and basic structural checks.
//   - The consumer is cancellation-aware: context cancellation propagates to the underlying
//     Kafka reader, exiting cleanly.
//   - The DLQ preserves the original payload and adds minimal headers for post-mortem analysis.
//     Do not mutate the original message body when forwarding to DLQ.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/merkulovlad/wbtech-go/internal/logger"
	"github.com/merkulovlad/wbtech-go/internal/model"
	"github.com/merkulovlad/wbtech-go/internal/service/order"
	"github.com/segmentio/kafka-go"
)

// Consumer wraps a kafka-go Reader and an optional DLQ Writer and orchestrates
// message decoding, validation, and delegation to the order Service.
type Consumer struct {
	// reader consumes messages from the primary topic.
	reader *kafka.Reader
	// dlqWriter is an optional producer used to forward irrecoverable messages.
	dlqWriter *kafka.Writer

	// svc is the domain service used to persist/process orders.
	svc order.Service
	// log is the project-wide logger interface.
	log logger.InterfaceLogger

	// topic is the source topic name (for reference in DLQ headers).
	topic string
	// dlqTopic is the DLQ topic name (empty means DLQ disabled).
	dlqTopic string
}

// NewConsumer constructs a new Consumer.
//
// Parameters:
//   - brokers: Kafka bootstrap addresses.
//   - topic: source topic to consume from.
//   - groupID: consumer group identifier.
//   - dlqTopic: DLQ topic; if empty, DLQ publishing is disabled.
//   - svc: domain service to handle valid orders.
//   - log: logger implementation.
//
// Note: DLQ usage is recommended in production to avoid partition halts caused by poison messages.
func NewConsumer(brokers []string, topic, groupID, dlqTopic string, svc order.Service, log logger.InterfaceLogger) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	var w *kafka.Writer
	if dlqTopic != "" {
		w = &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    dlqTopic,
			Balancer: &kafka.LeastBytes{},
		}
	}
	return &Consumer{
		reader:    r,
		dlqWriter: w,
		svc:       svc,
		log:       log,
		topic:     topic,
		dlqTopic:  dlqTopic,
	}
}

// Run starts the consumer loop and blocks until the context is canceled or a fatal error occurs.
// The loop semantics are:
//  1. Read a message.
//  2. Decode JSON into model.Order.
//  3. Validate minimally (required fields, sensible ranges, non-empty items).
//  4. Invoke service.Create to perform domain processing/storage.
//  5. On unrecoverable failure: forward the original message to the DLQ with reason headers.
//     On success: ack via offset commit (handled internally by kafka-go Reader).
func (c *Consumer) Run(ctx context.Context) error {
	// Ensure resources are closed even on early returns.
	defer func() {
		if err := c.reader.Close(); err != nil {
			c.log.Errorf("kafka: reader close: %v", err)
		}
		if c.dlqWriter != nil {
			if err := c.dlqWriter.Close(); err != nil {
				c.log.Errorf("kafka: dlq writer close: %v", err)
			}
		}
	}()

	for {
		// ReadMessage blocks until a message arrives or the context is canceled.
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			// Returning the error exits the loop. For context cancellation, kafka-go returns ctx.Err().
			return err
		}

		// Decode payload into a strongly-typed Order.
		var o model.Order
		if err := json.Unmarshal(m.Value, &o); err != nil {
			c.log.Errorf("kafka: invalid JSON payload: %v", err)
			_ = c.sendToDLQ(ctx, m, "invalid_json", err)
			continue
		}

		// Minimal, defensive validation before entering domain logic.
		if err := validateOrder(&o); err != nil {
			c.log.Errorf("kafka: validation failed: %v", err)
			_ = c.sendToDLQ(ctx, m, "schema_validation", err)
			continue
		}

		// Delegate to domain service (idempotency and deeper validation happen there).
		if err := c.svc.Create(ctx, &o); err != nil {
			// Consider classifying transient vs permanent errors; for simplicity, DLQ everything here.
			c.log.Errorf("kafka: service create failed for order=%s: %v", o.OrderUID, err)
			_ = c.sendToDLQ(ctx, m, "business_error", err)
			continue
		}

		c.log.Infof("kafka: created order %s", o.OrderUID)
	}
}

// sendToDLQ forwards the original message to the DLQ topic, augmenting headers with diagnostics.
// If DLQ is disabled or the write fails, the error is logged and suppressed (best-effort policy).
func (c *Consumer) sendToDLQ(ctx context.Context, src kafka.Message, reason string, cause error) error {
	if c.dlqWriter == nil {
		// DLQ is optional; silently ignore if not configured.
		return nil
	}
	errText := reason
	if cause != nil {
		errText = fmt.Sprintf("%s: %v", reason, cause)
	}
	dlqMsg := kafka.Message{
		Key:   src.Key,   // preserve key for potential replay/partitioning affinity
		Value: src.Value, // preserve exact original payload
		Headers: append(src.Headers, []kafka.Header{
			{Key: "error", Value: []byte(errText)},
			{Key: "origin-topic", Value: []byte(c.topic)},
			{Key: "timestamp", Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
		}...),
	}
	if err := c.dlqWriter.WriteMessages(ctx, dlqMsg); err != nil {
		c.log.Errorf("kafka: DLQ write failed (topic=%s): %v", c.dlqTopic, err)
		return err
	}
	return nil
}

// validateOrder performs minimal structural validation of an Order.
// Keep this function free of business/stateful checks—only validate what is intrinsic
// to the single message (e.g., required fields, non-zero timestamps, basic ranges).
// Domain-level rules (e.g., state transitions, uniqueness) should live in the service layer.
func validateOrder(o *model.Order) error {
	if o == nil {
		return fmt.Errorf("order is nil")
	}

	var errs []string

	// Required non-empty identifiers/strings.
	if o.OrderUID == "" {
		errs = append(errs, "order_uid is required")
	}
	if o.TrackNumber == "" {
		errs = append(errs, "track_number is required")
	}
	if o.Entry == "" {
		errs = append(errs, "entry is required")
	}
	if o.CustomerID == "" {
		errs = append(errs, "customer_id is required")
	}
	if o.DeliveryService == "" {
		errs = append(errs, "delivery_service is required")
	}
	if o.ShardKey == "" {
		errs = append(errs, "shardkey is required")
	}
	if o.OofShard == "" {
		errs = append(errs, "oof_shard is required")
	}

	// Items must be present (empty order is not actionable).
	if len(o.Items) == 0 {
		errs = append(errs, "items must be non-empty")
	}

	// Timestamp should be set (zero time usually indicates producer bug).
	if o.DateCreated.IsZero() {
		errs = append(errs, "date_created must be set")
	}

	// Optional sanity: SmID should not be negative.
	if o.SmID < 0 {
		errs = append(errs, "sm_id must be >= 0")
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errs, "; "))
	}
	return nil
}
