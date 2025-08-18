package kafka

import (
	"context"
	"encoding/json"
	"github.com/merkulovlad/wbtech-go/internal/logger"
	"github.com/merkulovlad/wbtech-go/internal/model"
	"github.com/merkulovlad/wbtech-go/internal/service/order"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	svc    order.Service
	log    logger.InterfaceLogger
}

func NewConsumer(brokers []string, topic, groupID string, svc order.Service, log logger.InterfaceLogger) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{
		reader: r,
		svc:    svc,
		log:    log,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {

		}
	}(c.reader)

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err // exits if ctx is canceled
		}

		var o model.Order
		if err := json.Unmarshal(m.Value, &o); err != nil {
			c.log.Errorf("bad kafka message: %v", err)
			continue
		}

		if err := c.svc.Create(ctx, &o); err != nil {
			c.log.Errorf("failed to create order: %v", err)
			continue
		}

		c.log.Infof("created order %s", o.OrderUID)
	}
}
