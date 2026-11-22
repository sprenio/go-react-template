package queue

import (
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"context"
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueEvent struct {
	QueueId    uint          `json:"queue_id"`
	Generation uint          `json:"generation"`
	Retry      uint          `json:"retry"`
	BaseDelay  time.Duration `json:"base_delay"`
}

const (
	maxMqIntervalSeconds      = 30
	maxDelay                  = 12 * time.Hour
	maxRetriesQueue      uint = 5

	ResetPasswordQueuePrefix = "reset_password"
	ResetPasswordBaseDelay   = 5 * time.Second
	ResetPasswordMaxRetries  = 3

	RegisterQueuePrefix = "register_tasks"
	RegisterMaxRetries  = 5
	RegisterBaseDelay   = 5 * time.Second

	EmailQueuePrefix = "email_tasks"
	EmailMaxRetries  = 5
	EmailBaseDelay   = 2 * time.Second

	ConfirmationEmailType = "confirmation_email"
)

func ConnectWithRetry(ctx context.Context, url string) (*amqp.Connection, error) {
	var attempt int
	for {
		conn, err := amqp.Dial(url)
		if err == nil {
			logger.Info("Connected to RabbitMQ")
			return conn, nil
		}

		attempt++
		backoff := time.Duration(math.Min(math.Pow(2, float64(attempt)), maxMqIntervalSeconds)) * time.Second
		logger.WarnCtx(ctx, fmt.Sprintf("RabbitMQ connection error (attempt %d): %v. Retrying in %v", attempt, err, backoff))

		select {
		case <-ctx.Done():
			return nil, apperror.Wrap(err, "timeout while connecting to RabbitMQ")
		case <-time.After(backoff):
		}
	}
}

func SetupQueues(ch *amqp.Channel, queuePrefix string, baseDelay time.Duration) error {
	mainQueueName := GetMainQueueName(queuePrefix)
	dlqQueueName := GetDLQQueueName(queuePrefix)

	// 1. DLQ
	if _, err := ch.QueueDeclare(dlqQueueName, true, false, false, false, nil); err != nil {
		logger.Error("Failed to declare DLQ: %v", err)
		return err
	}

	// 2. Retry Queues (Multi-Bucket)
	for i := uint(1); i <= maxRetriesQueue; i++ {
		retryQueueName := GetRetryQueueName(queuePrefix, i)
		delay := calculateRetryDelay(baseDelay, i)

		_, err := ch.QueueDeclare(
			retryQueueName, true, false, false, false,
			amqp.Table{
				"x-dead-letter-exchange":    "",
				"x-dead-letter-routing-key": mainQueueName, // after TTL automatically returns to main
				"x-message-ttl":             int32(delay.Milliseconds()),
			},
		)
		if err != nil {
			logger.Error("Failed to declare Retry Queue %d: %v", i, err)
			return err
		}
	}

	// 3. Main Queue (after Nack(requeue=false) automatically goes to DLQ)
	_, err := ch.QueueDeclare(
		mainQueueName, true, false, false, false,
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": dlqQueueName,
		},
	)
	if err != nil {
		logger.Error("Failed to declare Main Queue: %v", err)
		return err
	}

	return nil
}

func GetNextRetryAt(baseDelay time.Duration, retryLevel uint) *time.Time {
	time := time.Now().Add(calculateRetryDelay(baseDelay, retryLevel))
	return &time
}

func calculateRetryDelay(baseDelay time.Duration, retryLevel uint) time.Duration {
	level := clampRetryLevel(retryLevel)

	// baseDelay * level * 5^(level - 1)
	multiplier := math.Pow(5, float64(level-1))
	delay := baseDelay * time.Duration(level) * time.Duration(multiplier)

	if delay > maxDelay || delay < 0 {
		return maxDelay
	}
	return delay
}

func GetRetryQueueName(queuePrefix string, retryLevel uint) string {
	level := clampRetryLevel(retryLevel)
	return fmt.Sprintf("%s_retry_%d", queuePrefix, level)
}

func GetMainQueueName(queuePrefix string) string {
	return queuePrefix
}

func GetDLQQueueName(queuePrefix string) string {
	return queuePrefix + "_dlq"
}

func GetQueueLockKey(queueId uint) string {
	return fmt.Sprintf("set_queue_lock_%d", queueId)
}

func clampRetryLevel(retryLevel uint) uint {
	if retryLevel > maxRetriesQueue {
		return maxRetriesQueue
	}
	if retryLevel < 1 {
		return 1
	}
	return retryLevel
}
