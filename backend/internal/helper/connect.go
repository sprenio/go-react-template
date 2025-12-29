package helper

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/database"
	"backend/internal/queue"
	"backend/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectServicesWithRetry(ctx context.Context, dbDSN string, rabbitURL string) (*sql.DB, *amqp.Connection, error) {
	dbCh := make(chan error, 1)
	rabbitCh := make(chan struct {
		conn *amqp.Connection
		err  error
	}, 1)

	database.SetDSN(dbDSN)
	db, err := database.GetDB()
	if err != nil {
		logger.Fatal("Nie można połączyć się z bazą danych: %v", err)
	}

	// DB goroutine
	go func() {
		dbCh <- database.ConnectWithRetry(ctx, db)
	}()

	// RabbitMQ goroutine
	go func() {
		conn, err := queue.ConnectWithRetry(ctx, rabbitURL)
		rabbitCh <- struct {
			conn *amqp.Connection
			err  error
		}{conn, err}
	}()

	var dbErr error
	var rabbitConn *amqp.Connection
	var rabbitErr error

	for i := 0; i < 2; i++ {
		select {
		case dbErr = <-dbCh:
			if dbErr != nil {
				return nil, nil, fmt.Errorf("błąd połączenia z DB: %w", dbErr)
			}
		case r := <-rabbitCh:
			rabbitConn = r.conn
			rabbitErr = r.err
			if rabbitErr != nil {
				return nil, nil, fmt.Errorf("błąd połączenia z RabbitMQ: %w", rabbitErr)
			}
		case <-ctx.Done():
			return nil, nil, fmt.Errorf("timeout podczas łączenia z usługami: %w", ctx.Err())
		}
	}

	return db, rabbitConn, nil
}
