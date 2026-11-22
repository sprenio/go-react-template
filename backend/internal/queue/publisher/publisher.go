package publisher

import (
	"backend/internal/cache"
	"backend/internal/contexthelper"
	"backend/internal/queue"
	queueRepository "backend/internal/queue/repository"
	userRepository "backend/internal/user/repository"
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type QueuePublisher struct {
	queueRepository *queueRepository.QueueRepository
	rabbitConn      *amqp.Connection
	userRepository  *userRepository.UserRepository
	redisClient     *redis.Client
}

var ErrUnrecoverableError = apperror.New("fatal unrecoverable error")

func NewPublisher(ctx context.Context) *QueuePublisher {
	db := contexthelper.GetDb(ctx)
	rabbitConn := contexthelper.GetRabbitConn(ctx)
	queueRepo := queueRepository.NewQueueRepository(db)
	userRepo := userRepository.NewRepository(db)
	redisClient := contexthelper.GetRedisClient(ctx)
	return &QueuePublisher{queueRepository: queueRepo, rabbitConn: rabbitConn, userRepository: userRepo, redisClient: redisClient}
}

func (s *QueuePublisher) PublishEvent(ctx context.Context, queueId uint, generation uint, queuePrefix string, baseDelay time.Duration) error {
	event := queue.QueueEvent{QueueId: queueId, Generation: generation, BaseDelay: baseDelay, Retry: 0}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ch, err := s.rabbitConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := queue.SetupQueues(ch, queuePrefix, baseDelay); err != nil {
		return err
	}

	return ch.PublishWithContext(
		ctx,
		"",                                  // direct default exchange
		queue.GetMainQueueName(queuePrefix), // routing key (nazwa kolejki)
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (s *QueuePublisher) AcquireQueueLock(ctx context.Context, queueId uint) (*cache.RedisLock, bool, error) {
	queueKey := queue.GetQueueLockKey(queueId)
	lock := cache.NewLock(s.redisClient, queueKey)
	acquired, err := lock.Acquire(ctx, 5*time.Second)
	return lock, acquired, err
}

func PublishTask[T any](ctx context.Context, s *QueuePublisher, payload T, taskName queue.QueueTask, queuePrefix string, baseDelay time.Duration, maxRetries uint) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	queueId, err := s.queueRepository.Create(ctx, string(taskName), jsonData, maxRetries)
	if err != nil {
		return err
	}

	lock, acquired, err := s.AcquireQueueLock(ctx, queueId)
	if err != nil {
		return err
	}
	if !acquired {
		logger.ErrorCtx(ctx, "Failed to acquire lock for queue", "queueId", queueId)
		return apperror.New("Failed to acquire lock")
	}
	// Bezpieczne zwolnienie na koniec funkcji
	defer func() {
		if err := lock.Release(ctx); err != nil {
			logger.ErrorCtx(ctx, "Failed to release lock for queue", "queueId", queueId, "error", err)
		}
	}()
	err = s.PublishEvent(ctx, queueId, queue.InitialGeneration, queuePrefix, baseDelay)
	if err != nil {
		if dbErr := s.queueRepository.SetErrorStatus(ctx, queueId, err.Error(), nil); dbErr != nil {
			return apperror.Wrap(err, "Failed to set queue as failed", "dbErr", dbErr)
		}
		return err
	}

	if _, dbErr := s.queueRepository.SetQueuedStatus(ctx, queueId); dbErr != nil {
		logger.ErrorCtx(ctx, "Failed to set queue as queued", "queueId", queueId, "error", dbErr)
		// no error - queue will be processed anyway
	}

	return nil
}
