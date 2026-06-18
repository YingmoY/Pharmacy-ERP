package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/YingmoY/PharmacyERP/internal/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	mandatory bool
}

type LogEvent struct {
	BusinessType string      `json:"business_type"`
	BusinessID   string      `json:"business_id"`
	Action       string      `json:"action"`
	OperatorID   int64       `json:"operator_id"`
	Detail       interface{} `json:"detail"`
}

func NewClient(cfg config.RabbitMQConfig) (*Client, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("connect rabbitmq failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open rabbitmq channel failed: %w", err)
	}

	queueName := cfg.LogQueue
	if queueName == "" {
		queueName = "operation_log_queue"
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare queue failed: %w", err)
	}

	if cfg.Prefetch > 0 {
		if err := ch.Qos(cfg.Prefetch, 0, false); err != nil {
			_ = ch.Close()
			_ = conn.Close()
			return nil, fmt.Errorf("set qos failed: %w", err)
		}
	}

	return &Client{conn: conn, channel: ch, queueName: queueName, mandatory: cfg.Mandatory}, nil
}

func (c *Client) PublishLogEvent(ctx context.Context, event LogEvent) error {
	if c == nil || c.channel == nil {
		return nil
	}
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal log event failed: %w", err)
	}

	return c.channel.PublishWithContext(
		ctx,
		"",
		c.queueName,
		c.mandatory,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// NewConsumerChannel 创建一个独立的消费者连接和 channel。
// 消费者应使用独立 channel 避免与发布者共享 channel 导致竞争问题。
// 调用方负责关闭返回的 channel 和 connection。
func (c *Client) NewConsumerChannel() (*amqp.Connection, *amqp.Channel, error) {
	if c == nil || c.conn == nil {
		return nil, nil, fmt.Errorf("mq client is not initialized")
	}
	// 复用现有 connection 创建新 channel（RabbitMQ 支持单连接多 channel）
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("open consumer channel failed: %w", err)
	}
	return c.conn, ch, nil
}

func (c *Client) Close() error {
	if c == nil {
		return nil
	}
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
