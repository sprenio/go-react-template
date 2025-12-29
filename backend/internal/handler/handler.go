package handler

import (
	"database/sql"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	db *sql.DB
	rabbitConn *amqp.Connection
}

func NewHandler(db *sql.DB, rabbitConn *amqp.Connection) *Handler {
	return &Handler{db: db, rabbitConn: rabbitConn}
}
