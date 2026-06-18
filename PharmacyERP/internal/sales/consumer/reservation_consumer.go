// Package consumer 实现销售模块的消息队列消费者。
package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 预占状态常量
const (
	reservationStatusReserved = "RESERVED"
	reservationStatusExpired  = "EXPIRED"
)

// reservationRecord 预占记录（仅查询所需字段）。
type reservationRecord struct {
	ID       int64     `gorm:"column:id"`
	Status   string    `gorm:"column:status"`
	ExpireAt time.Time `gorm:"column:expire_at"`
}

// ReservationExpireConsumer 预占过期消费者，从 reservation_expire_queue 读取消息，检查并释放过期预占。
type ReservationExpireConsumer struct {
	db       *gorm.DB
	mqClient *mq.Client
	log      *zap.Logger
}

// NewReservationExpireConsumer 创建预占过期消费者实例。
func NewReservationExpireConsumer(db *gorm.DB, mqClient *mq.Client, log *zap.Logger) *ReservationExpireConsumer {
	return &ReservationExpireConsumer{
		db:       db,
		mqClient: mqClient,
		log:      log,
	}
}

// Start 启动预占过期消费者，持续处理过期预占释放逻辑。
// 该方法会阻塞，应在独立 goroutine 中调用。
func (c *ReservationExpireConsumer) Start(ctx context.Context) error {
	if c.mqClient == nil {
		c.log.Warn("MQ 客户端未初始化，预占过期消费者跳过启动")
		return nil
	}

	// 创建独立消费者 channel
	_, ch, err := c.mqClient.NewConsumerChannel()
	if err != nil {
		return fmt.Errorf("创建预占过期消费者 channel 失败: %w", err)
	}
	defer func() {
		_ = ch.Close()
	}()

	// 确保队列已声明
	_, err = ch.QueueDeclare(
		mq.ReservationExpireQueue,
		true,  // 持久化
		false, // 不自动删除
		false, // 非排他
		false, // 不等待
		nil,
	)
	if err != nil {
		return fmt.Errorf("声明预占过期队列失败: %w", err)
	}

	// 设置预取数量，避免并发处理过多消息
	if err := ch.Qos(5, 0, false); err != nil {
		return fmt.Errorf("设置 QoS 失败: %w", err)
	}

	msgs, err := ch.Consume(
		mq.ReservationExpireQueue,
		"reservation-expire-consumer", // 消费者标签
		false,                         // 手动 ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("注册预占过期消费者失败: %w", err)
	}

	c.log.Info("预占过期消费者已启动", zap.String("queue", mq.ReservationExpireQueue))

	for {
		select {
		case <-ctx.Done():
			c.log.Info("预占过期消费者收到停止信号")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				c.log.Warn("预占过期队列 channel 已关闭")
				return nil
			}
			c.handleMessage(ctx, msg)
		}
	}
}

// handleMessage 处理单条预占过期消息。
//
// 处理逻辑：
// 1. 解析消息体为 ReservationExpireEvent
// 2. 按 reservation_id 查找预占记录
// 3. 若记录不存在或已软删除，直接 ack
// 4. 若状态不为 RESERVED，直接 ack（已被消费或已释放）
// 5. 若当前时间 >= expire_at，加行锁更新状态为 EXPIRED，记录 released_at
// 6. 若当前时间 < expire_at，消息提前到达，直接 ack（调度精度问题可接受）
func (c *ReservationExpireConsumer) handleMessage(ctx context.Context, msg amqp.Delivery) {
	var event mq.ReservationExpireEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		c.log.Error("解析预占过期消息失败",
			zap.Error(err),
			zap.ByteString("body", msg.Body),
		)
		// 消息格式错误无法重试，直接 ack
		_ = msg.Ack(false)
		return
	}

	c.log.Debug("处理预占过期消息",
		zap.Int64("reservation_id", event.ReservationID),
		zap.String("reservation_no", event.ReservationNo),
		zap.String("trace_code", event.TraceCode),
		zap.Time("expire_at", event.ExpireAt),
	)

	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 步骤 1：使用行锁查询预占记录
		var record reservationRecord
		result := tx.
			Raw("SELECT id, status, expire_at FROM trace_reservation WHERE id = ? AND deleted_at IS NULL FOR UPDATE",
				event.ReservationID).
			Scan(&record)

		if result.Error != nil {
			return result.Error
		}

		// 步骤 2：记录不存在则跳过
		if result.RowsAffected == 0 {
			c.log.Debug("预占记录不存在或已删除，跳过处理",
				zap.Int64("reservation_id", event.ReservationID),
			)
			return nil
		}

		// 步骤 3：状态不为 RESERVED 则跳过（已被处理）
		if record.Status != reservationStatusReserved {
			c.log.Debug("预占记录状态非 RESERVED，跳过处理",
				zap.Int64("reservation_id", event.ReservationID),
				zap.String("status", record.Status),
			)
			return nil
		}

		// 步骤 4：检查是否已到期
		now := time.Now()
		if now.Before(event.ExpireAt) {
			c.log.Debug("预占尚未到期，跳过处理",
				zap.Int64("reservation_id", event.ReservationID),
				zap.Time("expire_at", event.ExpireAt),
			)
			return nil
		}

		// 步骤 5：更新预占状态为 EXPIRED
		if err := tx.Exec(
			"UPDATE trace_reservation SET status = ?, released_at = ?, updated_at = ? WHERE id = ?",
			reservationStatusExpired, now, now, event.ReservationID,
		).Error; err != nil {
			return err
		}

		c.log.Info("预占已标记为过期",
			zap.Int64("reservation_id", event.ReservationID),
			zap.String("reservation_no", event.ReservationNo),
			zap.String("trace_code", event.TraceCode),
		)

		// 步骤 6：记录操作日志
		c.log.Info("预占过期操作日志",
			zap.Int64("reservation_id", event.ReservationID),
			zap.String("action", "EXPIRE_RESERVATION"),
			zap.String("trace_code", event.TraceCode),
			zap.Int64("sales_order_id", event.SalesOrderID),
		)

		return nil
	})

	if err != nil {
		c.log.Error("处理预占过期消息失败",
			zap.Int64("reservation_id", event.ReservationID),
			zap.Error(err),
		)
		// 数据库错误，消息重新入队
		_ = msg.Nack(false, true)
		return
	}

	_ = msg.Ack(false)
}
