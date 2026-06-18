// Package consumer 实现审计模块的消息队列消费者。
package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	auditModel "github.com/YingmoY/PharmacyERP/internal/audit/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const operationLogQueue = "operation_log_queue"

// LogConsumer 操作日志 MQ 消费者，从 operation_log_queue 读取并持久化操作日志。
type LogConsumer struct {
	db       *gorm.DB
	mqClient *mq.Client
	log      *zap.Logger
}

// NewLogConsumer 创建操作日志消费者实例。
func NewLogConsumer(db *gorm.DB, mqClient *mq.Client, log *zap.Logger) *LogConsumer {
	return &LogConsumer{
		db:       db,
		mqClient: mqClient,
		log:      log,
	}
}

// Start 启动消费者，持续从 operation_log_queue 消费消息并写入数据库。
// 该方法会阻塞，应在独立 goroutine 中调用。
func (c *LogConsumer) Start(ctx context.Context) error {
	if c.mqClient == nil {
		c.log.Warn("MQ 客户端未初始化，操作日志消费者跳过启动")
		return nil
	}

	// 使用反射调用 channel，需要通过新连接方式消费
	// 由于 mq.Client 没有暴露 channel，需要建立独立消费连接
	conn, ch, err := c.mqClient.NewConsumerChannel()
	if err != nil {
		return fmt.Errorf("创建消费者 channel 失败: %w", err)
	}
	defer func() {
		_ = ch.Close()
		_ = conn.Close()
	}()

	// 声明队列（幂等操作）
	_, err = ch.QueueDeclare(
		operationLogQueue,
		true,  // 持久化
		false, // 不自动删除
		false, // 非排他
		false, // 不等待
		nil,
	)
	if err != nil {
		return fmt.Errorf("声明队列失败: %w", err)
	}

	// 设置预取数量，避免一次性取出过多消息
	if err := ch.Qos(10, 0, false); err != nil {
		return fmt.Errorf("设置 QoS 失败: %w", err)
	}

	msgs, err := ch.Consume(
		operationLogQueue,
		"audit-log-consumer", // 消费者标签
		false,                // 手动 ack
		false,                // 非排他
		false,                // no-local
		false,                // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("注册消费者失败: %w", err)
	}

	c.log.Info("操作日志消费者已启动", zap.String("queue", operationLogQueue))

	for {
		select {
		case <-ctx.Done():
			c.log.Info("操作日志消费者收到停止信号")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				c.log.Warn("操作日志队列 channel 已关闭")
				return nil
			}
			c.handleMessage(msg)
		}
	}
}

// handleMessage 处理单条操作日志消息。
func (c *LogConsumer) handleMessage(msg amqp.Delivery) {
	var event mq.LogEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		c.log.Error("解析操作日志消息失败",
			zap.Error(err),
			zap.ByteString("body", msg.Body),
		)
		_ = msg.Ack(false)
		return
	}

	// operator_id=0 violates the FK constraint; discard rather than infinite-requeue
	if event.OperatorID <= 0 {
		c.log.Warn("操作日志消息 operator_id 无效，丢弃",
			zap.Int64("operator_id", event.OperatorID),
			zap.String("action", event.Action),
		)
		_ = msg.Ack(false)
		return
	}

	// 将 mq.LogEvent 映射到 operation_log 表字段
	// BusinessType -> business_type / module
	// BusinessID   -> business_id / target_id
	// Action       -> action
	// OperatorID   -> operator_id
	// Detail       -> detail (JSONB)
	detailJSON, _ := json.Marshal(event.Detail)

	now := time.Now()
	opLog := &auditModel.OperationLog{
		Module:       event.BusinessType,
		BusinessType: event.BusinessType,
		BusinessID:   event.BusinessID,
		Action:       event.Action,
		OperatorID:   event.OperatorID,
		Detail:       detailJSON,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := c.db.Create(opLog).Error; err != nil {
		c.log.Error("写入操作日志失败",
			zap.String("action", event.Action),
			zap.Int64("operator_id", event.OperatorID),
			zap.Error(err),
		)
		// 写入失败则 nack，消息重新入队
		_ = msg.Nack(false, true)
		return
	}

	c.log.Debug("操作日志写入成功",
		zap.String("business_type", event.BusinessType),
		zap.String("action", event.Action),
		zap.Int64("operator_id", event.OperatorID),
	)
	_ = msg.Ack(false)
}
