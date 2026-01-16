package queue

import (
	"backend/internal/contexthelper"
	"backend/pkg/logger"
	"context"
	"encoding/json"
)

type EmailTask struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Retries  int    `json:"retries"`
}

func (c *Consumer) StartEmailConsumer(ctx context.Context) error {
	rabbitConn := contexthelper.GetRabbitConn(ctx)
	ch, err := rabbitConn.Channel()
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to open channel: %v", err)
		return err
	}
	defer ch.Close()

	setupEmailQueues(ch)
	ch.Qos(1, 0, false)

	msgs, err := ch.Consume(emailMainQueue, "", false, false, false, false, nil)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to register email consumer: %v", err)
		return err
	}

	logger.InfoCtx(ctx, "📩 Email consumer started...")

	for {
		select {
		case <-ctx.Done():
			logger.InfoCtx(ctx, "⏹ Email consumer stopped by context")
			return nil
		case d, ok := <-msgs:
			if !ok {
				logger.InfoCtx(ctx, "📭 Email queue closed")
				return nil
			}

			var event QueueEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				logger.ErrorCtx(ctx, "❌ Invalid email task: %v", err)
				d.Ack(false)
				continue
			}
			logger.InfoCtx(ctx, "🔄 Handling email task: %s", event.Task)
			if err := c.HandleEmailTask(ctx, event.Task, event.Data); err != nil {
				logger.ErrorCtx(ctx, "Failed to handle email task: %v", err)

				event.Retries++
				if event.Retries > emailMaxRetries {
					logger.WarnCtx(ctx, "💀 Email max retries reached: %v", event)
					publishJSON(ch, emailDLQQueue, event)
				} else {
					logger.WarnCtx(ctx, "🔄 Email retry %d/%d", event.Retries, emailMaxRetries)
					publishJSON(ch, emailRetryQueue, event)
				}
				d.Ack(false)
				continue
			}

			d.Ack(false)
			logger.InfoCtx(ctx, "✅ Email task handled successfully: %s", event.Task)
		}
	}
}
