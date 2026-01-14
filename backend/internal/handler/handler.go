package handler

import (
	"time"
)

type Handler struct {
	startTime time.Time
}

func NewHandler() *Handler {
	return &Handler{startTime: time.Now()}
}
