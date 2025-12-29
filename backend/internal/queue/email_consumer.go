package queue

import (
	"backend/pkg/logger"
	"encoding/json"
	"context"
) 

type EmailTask struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Retries int `json:"retries"`
}

func (c *Consumer) StartEmailConsumer(ctx context.Context) error {
	ch, err := c.rabbitConn.Channel()
	if err != nil {
		logger.Error("Failed to open channel: %v", err)
		return err
	}
	defer ch.Close()

	setupEmailQueues(ch)
	ch.Qos(1, 0, false)

	msgs, err := ch.Consume(emailMainQueue, "", false, false, false, false, nil)
	if err != nil {
		logger.Error("Failed to register email consumer: %v", err)
		return err
	}

	logger.Info("📩 Email consumer started...")

	for {
		select {
		case <-ctx.Done():
			logger.Info("⏹ Email consumer stopped by context")
			return nil
		case d, ok := <-msgs:
			if !ok {
				logger.Info("📭 Email queue closed")
				return nil
			}

			var event QueueEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				logger.Error("❌ Invalid email task: %v", err)
				d.Ack(false)
				continue
			}

			if err := c.HandleEmailTask(ctx, event.Task, event.Data); err != nil {
				logger.Error("Failed to handle email task: %v", err)

				event.Retries++
				if event.Retries > emailMaxRetries {
					logger.Warn("💀 Email max retries reached: %v", event)
					publishJSON(ch, emailDLQQueue, event)
				} else {
					logger.Warn("🔄 Email retry %d/%d", event.Retries, emailMaxRetries)
					publishJSON(ch, emailRetryQueue, event)
				}
				d.Ack(false)
				continue
			}

			d.Ack(false)
			logger.Info("✅ Email task handled successfully: %s", event.Task)
		}
	}
}
