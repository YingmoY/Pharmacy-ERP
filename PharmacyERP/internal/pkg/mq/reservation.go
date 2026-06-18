package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const ReservationExpireQueue = "reservation_expire_queue"

// ReservationExpireEvent 是预留过期消息结构
type ReservationExpireEvent struct {
	ReservationID  int64     `json:"reservation_id"`
	ReservationNo  string    `json:"reservation_no"`
	SalesOrderID   int64     `json:"sales_order_id"`
	TraceCode      string    `json:"trace_code"`
	ExpireAt       time.Time `json:"expire_at"`
}

// EnsureReservationExpireQueue 确保预留过期队列已声明（可在服务启动时调用）
func (c *Client) EnsureReservationExpireQueue() error {
	if c == nil || c.channel == nil {
		return nil
	}
	_, err := c.channel.QueueDeclare(
		ReservationExpireQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare reservation_expire_queue failed: %w", err)
	}
	return nil
}

// PublishReservationExpire 将预留过期消息发布到队列
// 注意：当前实现发布到 log_queue（因为没有 dead-letter 消费者），
// 未来接入 dead-letter 机制后可切换到 reservation_expire_queue
func (c *Client) PublishReservationExpire(ctx context.Context, event ReservationExpireEvent) error {
	if c == nil || c.channel == nil {
		return nil
	}
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal reservation expire event failed: %w", err)
	}

	// 发布到 log_queue（预留过期消费者将在后续实现）
	return c.channel.PublishWithContext(
		ctx,
		"",
		c.queueName, // 使用默认 log_queue
		c.mandatory,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Type:         "reservation_expire",
		},
	)
}
