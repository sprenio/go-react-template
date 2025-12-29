package queue

import (
	"database/sql"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	db         *sql.DB
	rabbitConn *amqp.Connection
}

func NewConsumer(db *sql.DB, rabbitConn *amqp.Connection) *Consumer {
	return &Consumer{
		db:         db,
		rabbitConn: rabbitConn,
	}
}
