package queue

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"medicare-gateway/internal/gateway"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName    = "medicare.gateway"
	dlxName         = "medicare.gateway.dlx"
	queueName       = "medicare.gateway.tasks"
	deadLetterQueue = "medicare.gateway.tasks.dlq"
	routingKey      = "medicare.task"
)

type Rabbit struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger *slog.Logger
}

func Connect(url string, logger *slog.Logger) (*Rabbit, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	r := &Rabbit{conn: conn, ch: ch, logger: logger}
	if err := r.declare(); err != nil {
		r.Close()
		return nil, err
	}
	return r, nil
}

func (r *Rabbit) declare() error {
	if err := r.ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		return err
	}
	if err := r.ch.ExchangeDeclare(dlxName, "direct", true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := r.ch.QueueDeclare(deadLetterQueue, true, false, false, false, nil); err != nil {
		return err
	}
	if err := r.ch.QueueBind(deadLetterQueue, routingKey, dlxName, false, nil); err != nil {
		return err
	}
	args := amqp.Table{
		"x-dead-letter-exchange":    dlxName,
		"x-dead-letter-routing-key": routingKey,
	}
	if _, err := r.ch.QueueDeclare(queueName, true, false, false, false, args); err != nil {
		return err
	}
	return r.ch.QueueBind(queueName, routingKey, exchangeName, false, nil)
}

func (r *Rabbit) PublishTask(ctx context.Context, task gateway.TaskMessage) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return r.ch.PublishWithContext(ctx, exchangeName, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		MessageId:    task.TaskID,
		Body:         body,
	})
}

func (r *Rabbit) Consume(ctx context.Context, handler func(context.Context, gateway.TaskMessage) error) error {
	if err := r.ch.Qos(4, 0, false); err != nil {
		return err
	}
	deliveries, err := r.ch.Consume(queueName, "medicare-gateway", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case delivery, ok := <-deliveries:
				if !ok {
					return
				}
				var task gateway.TaskMessage
				if err := json.Unmarshal(delivery.Body, &task); err != nil {
					r.logger.Error("invalid async task message", "error", err)
					_ = delivery.Nack(false, false)
					continue
				}
				callCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
				err := handler(callCtx, task)
				cancel()
				if err != nil {
					r.logger.Error("async task failed, moved to dlq", "task_id", task.TaskID, "infno", task.Infno, "error", err)
					_ = delivery.Nack(false, false)
					continue
				}
				_ = delivery.Ack(false)
			}
		}
	}()
	return nil
}

func (r *Rabbit) Close() error {
	if r.ch != nil {
		_ = r.ch.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
