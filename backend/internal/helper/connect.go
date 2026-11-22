package helper

import (
	"backend/internal/cache"
	"backend/internal/database"
	"backend/internal/queue"
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"context"
	"database/sql"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Services struct {
	DB          *sql.DB
	RabbitConn  *amqp.Connection
	RedisClient *redis.Client
}

func ConnectServicesWithRetry(ctx context.Context, dbDSN string, rabbitURL string, redisURL string) (*Services, error) {
	dbCh := make(chan error, 1)

	rabbitCh := make(chan struct {
		conn *amqp.Connection
		err  error
	}, 1)

	redisCh := make(chan struct {
		client *redis.Client
		err    error
	}, 1)

	database.SetDSN(dbDSN)
	db, err := database.GetDB()
	if err != nil {
		logger.Fatal("Nie można zainicjować bazy danych: %v", err)
	}

	// 1. DB goroutine
	go func() {
		dbCh <- database.ConnectWithRetry(ctx, db)
	}()

	// 2. RabbitMQ goroutine
	go func() {
		conn, err := queue.ConnectWithRetry(ctx, rabbitURL)
		rabbitCh <- struct {
			conn *amqp.Connection
			err  error
		}{conn, err}
	}()

	// 3. Redis goroutine
	go func() {
		client, err := cache.ConnectWithRetry(ctx, redisURL)
		redisCh <- struct {
			client *redis.Client
			err    error
		}{client, err}
	}()

	svc := &Services{
		DB: db,
	}

	// Oczekujemy na 3 usługi
	for i := 0; i < 3; i++ {
		select {
		case dbErr := <-dbCh:
			if dbErr != nil {
				return nil, apperror.Wrap(dbErr, "błąd połączenia z DB")
			}
		case r := <-rabbitCh:
			if r.err != nil {
				return nil, apperror.Wrap(r.err, "błąd połączenia z RabbitMQ")
			}
			svc.RabbitConn = r.conn
		case r := <-redisCh:
			if r.err != nil {
				return nil, apperror.Wrap(r.err, "błąd połączenia z Redisem")
			}
			svc.RedisClient = r.client
		case <-ctx.Done():
			return nil, apperror.Wrap(ctx.Err(), "timeout podczas łączenia z usługami")
		}
	}

	return svc, nil
}
