package queue

import (
	"backend/pkg/logger"
	"encoding/json"
	"context"
	"fmt"
	"math"
	"time"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const maxMqIntervalSeconds = 30

type QueueEvent struct {
	Task    string          `json:"task"`
	Data    json.RawMessage `json:"data"`
	Retries int             `json:"retries"`
}


func ConnectWithRetry(ctx context.Context, url string) (*amqp.Connection, error) {
	var attempt int
	for {
		conn, err := amqp.Dial(url)
		if err == nil {
			logger.Info("Połączono z RabbitMQ")
			return conn, nil
		}

		attempt++
		backoff := time.Duration(math.Min(math.Pow(2, float64(attempt)), maxMqIntervalSeconds)) * time.Second
		logger.Warn("Błąd połączenia z RabbitMQ (próba %d): %v. Ponawiam za %v", attempt, err, backoff)

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout podczas łączenia z RabbitMQ: %w", err)
		case <-time.After(backoff):
		}
	}
}


func publishJSON(ch *amqp.Channel, queueName string, event QueueEvent) {
	body, err := json.Marshal(event)
	if err != nil {
		logger.Error("Marshal error: %v", err)
		os.Exit(1)
	}
	err = ch.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		logger.Error("Publish error: %v", err)
	}
}

func setupQueues(ch *amqp.Channel, mainQueue, retryQueue, dlqQueue string, retryDelay int32) error {
	// DLQ
	_, err := ch.QueueDeclare(
		dlqQueue, true, false, false, false, nil,
	)
	if err != nil {
		logger.Error("Failed to declare DLQ: %v", err)
		return err
	}

	// Retry Queue
	_, err = ch.QueueDeclare(
		retryQueue, true, false, false, false,
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": mainQueue,
			"x-message-ttl":             retryDelay,
		},
	)
	if err != nil {
		logger.Error("Failed to declare Retry Queue: %v", err)
		return err
	}

	// Main Queue
	_, err = ch.QueueDeclare(
		mainQueue, true, false, false, false,
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": dlqQueue,
		},
	)
	if err != nil {
		logger.Error("Failed to declare Main Queue: %v", err)
		return err
	}
	logger.Debug("Queues setup completed: %s, %s, %s", mainQueue, retryQueue, dlqQueue)
	return nil
}
