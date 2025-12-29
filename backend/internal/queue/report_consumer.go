package queue

import (
	"backend/pkg/logger"
	"context"
	"encoding/json"
)

type ReportTask struct {
	ReportID string `json:"report_id"`
	UserID   string `json:"user_id"`
	Retries  int    `json:"retries"`
}

const (
	reportMainQueue  = "report_tasks"
	reportRetryQueue = "report_tasks_retry"
	reportDLQQueue   = "report_tasks_dlq"

	reportMaxRetries = 5
	reportRetryDelay = 10000 // ms
)

func (c *Consumer) StartReportConsumer(ctx context.Context) error {
	ch, err := c.rabbitConn.Channel()
	if err != nil {
		logger.Error("Failed to open channel: %v", err)
		return err
	}
	defer ch.Close()

	setupQueues(ch, reportMainQueue, reportRetryQueue, reportDLQQueue, int32(reportRetryDelay))
	ch.Qos(1, 0, false)

	msgs, err := ch.Consume(reportMainQueue, "", false, false, false, false, nil)
	if err != nil {
		logger.Error("Failed to register report consumer: %v", err)
		return err
	}

	logger.Info("📄 Report consumer started...")

	for {
		select {
		case <-ctx.Done():
			logger.Info("📄 Report consumer stopped by context")
			return nil
		case d, ok := <-msgs:
			if !ok {
				logger.Info("📄 Report queue closed")
				return nil
			}
			var event QueueEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				logger.Error("❌ Invalid report task: %v", err)
				d.Ack(false)
				continue
			}

			if err := c.HandleReportTask(ctx, event.Task, event.Data); err != nil {
				logger.Error("Failed to handle report task: %v", err)

				event.Retries++
				if event.Retries > reportMaxRetries {
					logger.Warn("💀 Report max retries reached: %v", event)
					publishJSON(ch, reportDLQQueue, event)
				} else {
					logger.Warn("🔄 Report retry %d/%d", event.Retries, reportMaxRetries)
					publishJSON(ch, reportRetryQueue, event)
				}
				d.Ack(false)
				continue
			}
			// Successfully handled the report task
			d.Ack(false)
			logger.Info("✅ Report generated: %s", event.Data)
		}
	}
}
