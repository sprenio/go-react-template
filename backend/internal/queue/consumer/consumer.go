package consumer

import (
	"backend/internal/apierrors"
	"backend/internal/contexthelper"
	"backend/internal/queue"
	"backend/internal/queue/publisher"
	queueRepository "backend/internal/queue/repository"
	userRepository "backend/internal/user/repository"
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"context"
	"encoding/json"
	"time"

	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

const consumerShutdownTimeout = 30 * time.Second

type QueueConsumer struct {
	queuePublisher  *publisher.QueuePublisher
	queueRepository *queueRepository.QueueRepository
	emailRepository *queueRepository.EmailLogRepository
	rabbitConn      *amqp.Connection
	userRepository  *userRepository.UserRepository
}

func NewConsumer(ctx context.Context) *QueueConsumer {
	db := contexthelper.GetDb(ctx)
	rabbitConn := contexthelper.GetRabbitConn(ctx)
	queueRepo := queueRepository.NewQueueRepository(db)
	emailRepo := queueRepository.NewEmailLogRepository(db)
	userRepo := userRepository.NewRepository(db)
	queuePub := publisher.NewPublisher(ctx)
	return &QueueConsumer{queuePublisher: queuePub, queueRepository: queueRepo, rabbitConn: rabbitConn, userRepository: userRepo, emailRepository: emailRepo}
}

func StartConsumer[T any](
	ctx context.Context,
	s *QueueConsumer,
	queuePrefix string,
	baseDelay time.Duration,
	taskName queue.QueueTask,
	handler func(ctx context.Context, event queue.QueueEvent, payload T) error,
) error {
	ch, err := s.rabbitConn.Channel()
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to open channel for queue prefix", "queuePrefix", queuePrefix, "error", err)
		return err
	}
	defer ch.Close()

	if err := queue.SetupQueues(ch, queuePrefix, baseDelay); err != nil {
		return err
	}

	if err := ch.Qos(1, 0, false); err != nil {
		logger.ErrorCtx(ctx, "Failed to set QoS", "queuePrefix", queuePrefix, "error", err)
		return err
	}

	mainQueue := queue.GetMainQueueName(queuePrefix)
	consumerTag := string(taskName)

	msgs, err := ch.Consume(
		mainQueue,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to register consumer for queue", "queuePrefix", queuePrefix, "error", err)
		return err
	}

	logger.InfoCtx(ctx, "📩 Consumer started...", "taskName", taskName)

	for {
		select {
		case <-ctx.Done():
			logger.InfoCtx(ctx, "⏹ Stopping RabbitMQ consumer", "taskName", taskName)

			if err := ch.Cancel(consumerTag, false); err != nil {
				logger.ErrorCtx(ctx, "Failed to cancel RabbitMQ consumer", "taskName", taskName, "error", err)
				return err
			}

			logger.InfoCtx(ctx, "✅ RabbitMQ consumer cancelled", "taskName", taskName)
			return nil

		case d, ok := <-msgs:
			if !ok {
				logger.InfoCtx(ctx, "📭 Queue closed", "queueName", mainQueue)
				return nil
			}

			deliveryCtx, cancel := context.WithCancel(
				context.WithoutCancel(ctx),
			)

			deliveryDone := make(chan struct{})

			go func() {
				select {
				case <-ctx.Done():
					// SIGTERM: start countdown for task completion.
					timer := time.NewTimer(consumerShutdownTimeout)
					defer timer.Stop()

					select {
					case <-deliveryDone:
						// Task finished within graceful shutdown time.

					case <-timer.C:
						// Minął czas na dokończenie zadania.
						logger.WarnCtx(
							ctx,
							"Graceful shutdown timeout reached, cancelling active delivery",
							"taskName", taskName,
						)
						cancel()
					}

				case <-deliveryDone:
					// Wiadomość zakończyła się podczas normalnej pracy.
				}
			}()

			processDelivery(
				deliveryCtx,
				s,
				ch,
				d,
				taskName,
				queuePrefix,
				handler,
			)

			close(deliveryDone)
			cancel()

			// Jeśli aplikacja jest w trakcie zamykania,
			// nie pobieramy kolejnej wiadomości.
			if ctx.Err() != nil {
				logger.InfoCtx(ctx, "⏹ Active delivery finished during shutdown", "taskName", taskName)
				return nil
			}
		}
	}
}

// processDelivery to uniwersalny runner dla każdego konsumenta.
// T reprezentuje typ struktury payloadu (np. payload.RegisterPayload).
func processDelivery[T any](
	ctx context.Context,
	s *QueueConsumer,
	ch *amqp.Channel,
	d amqp.Delivery,
	taskName queue.QueueTask,
	queuePrefix string,
	handler func(ctx context.Context, event queue.QueueEvent, payload T) error,
) {

	var err error
	// 1. Walidacja JSON z kolejki (QueueEvent)
	var event queue.QueueEvent
	if err := json.Unmarshal(d.Body, &event); err != nil {
		logger.ErrorCtx(ctx, "❌ Invalid task event json", "error", err)
		_ = d.Ack(false) // Zły JSON - nie ponawiamy
		return
	}

	// 2. Pobranie rekordu z bazy
	queueEntry, err := s.queueRepository.GetById(ctx, event.QueueId)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get task from DB", "error", err)
		_ = d.Nack(false, false) // Leci do DLQ
		return
	}

	if err := ctx.Err(); err != nil {
		logger.WarnCtx(
			ctx,
			"Delivery interrupted before processing",
			"queueId", queueEntry.Id,
			"error", err,
		)
		_ = d.Nack(false, true)
		return
	}

	var nextRetryAt *time.Time
	defer func() {
		if err != nil {
			statusErr := s.queueRepository.SetErrorStatus(ctx, queueEntry.Id, err.Error(), nextRetryAt)
			if statusErr != nil {
				logger.ErrorCtx(ctx, "Failed to set error status for queue", "error", statusErr, "queueId", queueEntry.Id)
			}
		}
	}()

	if queueEntry.Status == queue.QueueStatusRetryPending {
		logger.InfoCtx(ctx, "⏭Skipping task. Task will be processed on next cycle", "taskId", queueEntry.Id, "status", queueEntry.Status)
		_ = d.Ack(false)
		return
	}
	// 3. Ochrona przed duplikatami z unieważnionego cyklu retry (starsza generacja)
	if event.Generation > 0 && event.Generation < queueEntry.Generation {
		logger.InfoCtx(ctx, "⏭Skipping task", "taskId", queueEntry.Id, "eventGeneration", event.Generation, "queueEntryGeneration", queueEntry.Generation)
		_ = d.Ack(false)
		return
	}

	// 4. Ochrona przed duplikatami w ramach tej samej generacji (at-least-once delivery)
	if queueEntry.Status == queue.QueueStatusDone {
		logger.InfoCtx(ctx, "⏭ Skipping task. Task already done", "taskId", queueEntry.Id)
		_ = d.Ack(false)
		return
	}

	// 5. Weryfikacja typu zadania
	if queueEntry.Task != taskName {
		logger.ErrorCtx(ctx, "Invalid task name", "expected", taskName, "got", queueEntry.Task)
		_ = d.Nack(false, false)
		return
	}

	// 6. Atomowa inkrementacja prób w bazie
	retries, err := s.queueRepository.IncrementRetries(ctx, queueEntry.Id)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to increment retries", "queueId", queueEntry.Id, "error", err)
		_ = d.Nack(false, false)
		return
	}

	// 7. Ochrona przed ponownym wykonaniem (gdy przekroczono retries przed uruchomieniem)
	if retries > queueEntry.MaxRetries {
		logger.WarnCtx(ctx, "Skipping execution", "queueId", queueEntry.Id, "retries", retries, "maxRetries", queueEntry.MaxRetries)
		_ = d.Nack(false, false)
		return
	}

	// 8. Oczekiwanie na zablokowanie kolejki
	if err = s.WaitUntilReadyToProcess(ctx, event.QueueId, 5*time.Second, 250*time.Millisecond); err != nil {
		logger.ErrorCtx(ctx, "Failed to wait for queue status", "error", err)
		if ctx.Err() != nil {
			logger.WarnCtx(
				ctx,
				"Delivery interrupted while waiting for queue",
				"queueId", queueEntry.Id,
				"error", ctx.Err(),
			)
			_ = d.Nack(false, true)
			return
		}

		retry, finalErr := s.handleProcessedQueue(ctx, err, retries, queueEntry, event, queuePrefix, ch)
		err = finalErr
		if retry {
			nextRetryAt = queue.GetNextRetryAt(event.BaseDelay, retries)
			_ = d.Ack(false)
		} else {
			_ = d.Nack(false, false)
		}
		return
	}

	// 9. Generyczne parsowanie payloadu z bazy do struktury T
	var payloadData T
	if err = json.Unmarshal(queueEntry.Payload, &payloadData); err != nil {
		logger.ErrorCtx(ctx, "Invalid payload json for task", "taskName", taskName, "queueId", queueEntry.Id, "error", err)
		_ = d.Nack(false, false)
		return
	}

	event.Retry = retries
	logger.InfoCtx(ctx, "🔄 Handling task", "taskName", taskName, "queueId", queueEntry.Id, "retry", retries, "maxRetries", queueEntry.MaxRetries, "generation", event.Generation)

	// 10. Wykonanie dedykowanej logiki biznesowej
	if err = handler(ctx, event, payloadData); err != nil {
		logger.ErrorCtx(ctx, "Failed to handle task", "taskName", taskName, "queueId", queueEntry.Id)
		if ctx.Err() != nil {
			logger.WarnCtx(ctx, "Delivery interrupted during handler", "queueId", queueEntry.Id, "error", ctx.Err())
			markDeliveryInterrupted(ctx, s, queueEntry.Id)
			_ = d.Nack(false, true)
			return
		}
		retry, finalErr := s.handleProcessedQueue(ctx, err, retries, queueEntry, event, queuePrefix, ch)
		err = finalErr
		if retry {
			nextRetryAt = queue.GetNextRetryAt(event.BaseDelay, retries)
			_ = d.Ack(false)
		} else {
			_ = d.Nack(false, false)
		}
		return
	}

	// 11. SUKCES
	if doneErr := s.queueRepository.SetDoneStatus(ctx, queueEntry.Id, retries); doneErr != nil {
		logger.ErrorCtx(ctx, "Failed to set queue as done", "error", doneErr)
	}

	_ = d.Ack(false)
	logger.InfoCtx(ctx, "✅ Task handled successfully", "taskName", taskName, "queueId", queueEntry.Id)
}

func (s *QueueConsumer) publishToRetry(ctx context.Context, ch *amqp.Channel, queuePrefix string, event queue.QueueEvent, retryLevel uint) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	targetRetryQueue := queue.GetRetryQueueName(queuePrefix, retryLevel)

	return ch.PublishWithContext(
		ctx,
		"",               // default exchange
		targetRetryQueue, // routing key = nazwa dedykowanej kolejki retry
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
func (s *QueueConsumer) WaitUntilReadyToProcess(
	ctx context.Context,
	queueId uint,
	timeout time.Duration,
	interval time.Duration,
) error {
	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		ready, err := s.queueRepository.SetProcessingStatus(ctx, queueId)
		if err != nil {
			return err
		}

		if ready {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeoutTimer.C:
			return errors.New("timeout waiting for queue status")
		case <-ticker.C:
		}
	}
}

func (s *QueueConsumer) handleProcessedQueue(ctx context.Context, err error, retries uint, queueEntry queue.Queue, event queue.QueueEvent, queuePrefix string, ch *amqp.Channel) (bool, error) {
	// Nieodwracalny błąd -> DLQ
	if apierrors.IsUnrecoverableError(err) {
		err = apperror.Wrap(err, "💀 Unrecoverable error, sending straight to DLQ")
		return false, err
	}

	// Przekroczono limit prób -> DLQ
	if retries >= queueEntry.MaxRetries {
		err = apperror.Wrap(err, "Max retries reached", "retries", retries, "maxRetries", queueEntry.MaxRetries, "queueId", queueEntry.Id)
		return false, err
	}

	// Ponowienie (Retry)
	logger.WarnCtx(ctx, "Publishing to retry bucket", "retries", retries, "maxRetries", queueEntry.MaxRetries, "queueId", queueEntry.Id)
	if publishErr := s.publishToRetry(ctx, ch, queuePrefix, event, retries); publishErr != nil {
		err = apperror.Wrap(err, "Failed to publish to retry queue", "error", publishErr)
		return false, err
	}
	return true, err
}

func markDeliveryInterrupted(
	ctx context.Context,
	s *QueueConsumer,
	queueID uint,
) {
	cleanupCtx, cancel := context.WithTimeout(
		context.WithoutCancel(ctx),
		5*time.Second,
	)
	defer cancel()

	if err := s.queueRepository.SetErrorStatus(
		cleanupCtx,
		queueID,
		"Processing interrupted during consumer shutdown",
		nil,
	); err != nil {
		logger.ErrorCtx(
			cleanupCtx,
			"Failed to mark interrupted queue as failed",
			"queueId", queueID,
			"error", err,
		)
	}
}
